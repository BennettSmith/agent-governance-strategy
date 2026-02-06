package cli

import (
	"context"
	"os/exec"
)

var execCommandContext = exec.CommandContext

var _ = context.Background
