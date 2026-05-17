package main

import (
	"fmt"
	"os"

	"github.com/Kukuruzoed/config-analyzer/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
