package main

import (
	"fmt"
	"os"

	"k8spreview/pkg/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: k8spreview <yaml-file>")
		os.Exit(1)
	}

	yamlPath := os.Args[1]
	if err := ui.Run(yamlPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
