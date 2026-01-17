package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fisher0x/hashtocrack/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch strings.ToLower(command) {
	case "help", "-h", "--help":
		cli.PrintHelp(Version)
	case "version", "-v", "--version":
		fmt.Printf("Cracky v%s\n", Version)
	default:
		opts := cli.ParseArgs(os.Args[1:])
		cli.Run(opts)
	}
}
