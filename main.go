package main

import (
	"github.com/andpalmier/mbzr/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
