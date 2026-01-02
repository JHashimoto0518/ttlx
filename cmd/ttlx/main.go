// Package main is the entry point for the ttlx CLI tool.
package main

import (
	"os"

	"github.com/JHashimoto0518/ttlx/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
