package main

import (
	"fmt"
	"os"

	"github.com/nevarix/nevarix-v2/apps/hub/cmd/nvx-hub/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
