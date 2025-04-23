package main

import (
	"fmt"
	"os"

	"github.com/mirceanton/tesmartctl/cmd"
	"github.com/mirceanton/tesmartctl/internal/logging"
)

func main() {
	logging.Init()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
