package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "srv6-dynamic-sf-test",
	Short: "A proof of concept for SRv6 dynamic service function chaining",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("SRv6 dynamic service function chaining test")
		return nil
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
