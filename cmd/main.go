package main

import (
	"flag"
	"fmt"
	"os"

	"k8spreview/pkg/k8s"
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
	var resources []k8s.Resource
	if len(args) == 0 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			fmt.Printf("Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		// If data is being piped in, read from stdin
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			resources, err = k8s.Parse(os.Stdin)
			if err != nil {
				fmt.Printf("Error parsing YAML from stdin: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Usage: k8spreview <file1.yaml> [file2.yaml ...]")
			fmt.Println("       cat file.yaml | k8spreview -")
			os.Exit(1)
		}
	} else {
		for _, path := range args {
			var rs []k8s.Resource
			var err error
			if path == "-" {
				rs, err = k8s.Parse(os.Stdin)
			} else {
				rs, err = k8s.ParseFromFile(path)
			}
			if err != nil {
				fmt.Printf("Error parsing YAML from %s: %v\n", path, err)
				os.Exit(1)
			}
			resources = append(resources, rs...)
		}
	}

	if err := ui.RunWithResources(resources); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
