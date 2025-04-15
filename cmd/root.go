package cmd

import (
	"net"
	"os"
	"runtime/debug"

	"os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "srv6-dynamic-sf-test",
	Short: "A proof of concept for SRv6 dynamic service function chaining",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("SRv6 dynamic service function chaining test")
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					log.Debug("Build revision", "revision", setting.Value)
				}
			}
		}

		if err := rlimit.RemoveMemlock(); err != nil {
			log.Fatal("Failed to remove memory lock", "error", err)
		}

		var objs counterObjects
		if err := loadCounterObjects(&objs, nil); err != nil {
			log.Fatal("Failed to load counter objects", "error", err)
		}
		defer objs.Close()

		ifname := os.Getenv("IF_NAME")
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			log.Fatal("Failed to get interface by name", "ifname", ifname, "error", err)
		}

		link, err := link.AttachXDP(link.XDPOptions{
			Program:   objs.CountPackets,
			Interface: iface.Index,
		})
		if err != nil {
			log.Fatal("Failed to attach XDP program", "error", err)
		}
		defer link.Close()

		log.Info("XDP program attached", "ifname", ifname)

		tick := time.Tick(1 * time.Second)
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		for {
			select {
			case <-tick:
				var count uint64
				if err := objs.PktCount.Lookup(uint32(0), &count); err != nil {
					log.Fatal("Failed to lookup packet count", "error", err)
				}
				log.Info("Packet count", "count", count)
			case <-stop:
				log.Info("Received interrupt signal, exiting")
				return
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
