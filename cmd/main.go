package main

import (
	"flag"
	"fmt"
	"os"

	"k8spreview/pkg/ui"
	"k8spreview/pkg/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("k8spreview %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.Date)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: k8spreview <yaml-file>")
		fmt.Println("       k8spreview -version")
		os.Exit(1)
	}

	yamlPath := args[0]
	if err := ui.Run(yamlPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
