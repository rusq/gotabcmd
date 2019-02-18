package gotabcmd

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

const (
	// name of the tabcmd executable.  Must be in the path.
	cTabcmd = "tabcmd"
	// CCommandTimeout is a default command execution timeout
	CCommandTimeout = 15 * time.Second
)

// tabcmd is a wrapper around "tabcmd"
type tabcmd struct {
	timeout        time.Duration
	commandContext func(context.Context, string, ...string) *exec.Cmd
}

var commandContext = exec.CommandContext

// newTabcmd returns new Tabcmd executor.  If commandTimeout is 0, it is
// set to the default value.
func newTabcmd(commandTimeout time.Duration) *tabcmd {
	if commandTimeout == 0 {
		// set the default value
		commandTimeout = CCommandTimeout
	}
	return &tabcmd{
		timeout:        commandTimeout,
		commandContext: exec.CommandContext}
}

// run satisfies the TabExecutor interface, executes the cTabcmd with given
// action and args, returns output and error.  Output is stdout if no error
// occurred and stderr of the program otherwise.
func (t *tabcmd) run(action string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	cmd := t.commandContext(ctx, cTabcmd, append([]string{action}, args...)...)

	var bufout, buferr bytes.Buffer
	cmd.Stdout = &bufout
	cmd.Stderr = &buferr

	err := cmd.Run()
	if err != nil {
		return buferr.String(), err
	}

	return bufout.String(), nil
}
