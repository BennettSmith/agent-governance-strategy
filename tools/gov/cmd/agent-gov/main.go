package main

import (
	"os"

	"agent-governance-strategy/tools/gov/internal/cli"
)

func main() {
	mainWithExit(os.Exit, os.Args)
}

func run(args []string) int {
	return cli.Run(args, os.Stdout, os.Stderr)
}

func mainWithExit(exit func(int), args []string) {
	exit(run(args))
}
