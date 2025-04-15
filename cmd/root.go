package cmd

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"os/signal"

	"github.com/charmbracelet/log"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

type Config struct {
	Prefix net.IPNet
}

var config Config

var rootCmd = &cobra.Command{
	Use:   "srv6-dynamic-sf-test",
	Short: "A proof of concept for SRv6 dynamic service function chaining",
	Run: func(cmd *cobra.Command, args []string) {
		printCommitSHA()

		if err := rlimit.RemoveMemlock(); err != nil {
			log.Fatal("Failed to remove memory lock", "error", err)
		}

		var objs testObjects
		if err := loadTestObjects(&objs, nil); err != nil {
			log.Fatal("Failed to load test objects", "error", err)
		}
		defer objs.Close()

		bpfEncap := &netlink.BpfEncap{}
		bpfEncap.SetProg(nl.LWT_BPF_XMIT, objs.Test.FD(), "lwt_xmit/test")
		route := &netlink.Route{
			Dst:   &config.Prefix,
			Gw:    config.Prefix.IP,
			Encap: bpfEncap,
		}

		if err := netlink.RouteAdd(route); err != nil {
			log.Fatal("Failed to add route", "error", err)
		}
		defer func() {
			if err := netlink.RouteDel(route); err != nil {
				log.Fatal("Failed to delete route", "error", err)
			}
		}()
		log.Info("Route added", "route", route)

		stop := make(chan os.Signal, 1)
		stopped := false
		signal.Notify(stop, os.Interrupt)

		wg := &sync.WaitGroup{}

		{
			r, err := perf.NewReader(objs.Logs, os.Getpagesize())
			if err != nil {
				log.Fatal("Failed to create perf reader", "error", err)
			}
			defer r.Close()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for !stopped {
					r.SetDeadline(time.Now().Add(100 * time.Millisecond))
					record, err := r.Read()
					if errors.Is(err, os.ErrDeadlineExceeded) {
						continue
					}
					if err != nil {
						log.Fatal("Failed to read perf record", "error", err)
					}

					log.WithPrefix("ebpf").Info(string(record.RawSample))
				}
			}()
		}

		<-stop
		stopped = true
		log.Info("Received interrupt signal, exiting...")
		wg.Wait()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()

	rootCmd.Flags().IPNetVarP(&config.Prefix, "prefix", "p", net.IPNet{}, "SRv6 prefix")
	rootCmd.MarkFlagRequired("prefix")

	viper.BindPFlags(rootCmd.Flags())
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) {
			rootCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
