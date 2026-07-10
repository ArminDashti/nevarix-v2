package main

import (
	"fmt"
	"os"

	"github.com/nevarix/nevarix-v2/apps/agent/cmd/nvx-agent/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
