package main

import "testing"

func TestRun_Help(t *testing.T) {
	code := run([]string{"agent-gov", "help"})
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
}

func TestMainWithExit_PropagatesExitCode(t *testing.T) {
	var got int = -1
	mainWithExit(func(code int) { got = code }, []string{"agent-gov"})
	if got != 2 {
		t.Fatalf("expected exit code 2, got %d", got)
	}
}

