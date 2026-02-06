package main

import (
	"fmt"
	"os"

	"agent-governance-strategy/tools/gov/internal/cli"
)

func main() {
	exitCode := cli.Run(os.Args, os.Stdout, os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	fmt.Fprintln(os.Stdout)
}

