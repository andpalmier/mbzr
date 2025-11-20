package main

import (
	"os"

	"github.com/andpalmier/mbzr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
