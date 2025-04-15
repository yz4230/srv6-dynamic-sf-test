package cmd

import (
	"runtime/debug"

	"github.com/charmbracelet/log"
)

func printCommitSHA() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				log.Debug("Build revision", "revision", setting.Value)
			}
		}
	}
}
