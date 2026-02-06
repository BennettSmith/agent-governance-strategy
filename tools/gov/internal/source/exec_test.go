package source

import "os/exec"

// Wrapped for test stubbing if needed later.
var execCommand = exec.Command
