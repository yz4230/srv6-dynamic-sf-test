package main

import (
	"srv6-dynamic-sf-test/cmd"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
}
