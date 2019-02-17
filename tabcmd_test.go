package gotabcmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNewTabcmd(t *testing.T) {
	type args struct {
		commandTimeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *tabcmd
	}{
		{"timeout setting", args{5 * time.Hour}, &tabcmd{timeout: 5 * time.Hour, commandContext: exec.CommandContext}},
		{"default timeout", args{0}, &tabcmd{timeout: CCommandTimeout, commandContext: exec.CommandContext}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newTabcmd(tt.args.commandTimeout)
			// compare timeout
			if got.timeout != tt.want.timeout {
				t.Errorf("NewTabcmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

// from https://github.com/golang/go/blob/master/src/os/exec/exec_test.go
func helperCommandContext(ctx context.Context, name string, args ...string) (cmd *exec.Cmd) {
	// testenv.MustHaveExec(t)

	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, append([]string{name}, args...)...)
	if ctx != nil {
		cmd = exec.CommandContext(ctx, os.Args[0], cs...)
	} else {
		cmd = exec.Command(os.Args[0], cs...)
	}
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// from https://github.com/golang/go/blob/master/src/os/exec/exec_test.go
func helperCommand(t *testing.T, name string, args ...string) *exec.Cmd {
	return helperCommandContext(nil, name, args...)
}

// func TestTabcmd_Run(t *testing.T) {
// 	tc := NewTabcmd(0)
// 	tc.commandContext = helperCommandContext

// 	out, err := tc.Run("success_to_stdout", "arg1", "arg2")
// 	if

// }

func TestTabcmd_Run(t *testing.T) {
	type fields struct {
		timeout        time.Duration
		commandContext func(context.Context, string, ...string) *exec.Cmd
	}
	type args struct {
		action string
		args   []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"success", fields{CCommandTimeout, helperCommandContext},
			args{"succeed", []string{"param1", "param2", "param3"}},
			"Success: [param1 param2 param3]", false},
		{"stderr", fields{CCommandTimeout, helperCommandContext},
			args{"fail", []string{"fail1", "fail2"}},
			"Failure: [fail1 fail2]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &tabcmd{
				timeout:        tt.fields.timeout,
				commandContext: tt.fields.commandContext,
			}
			got, err := tc.run(tt.args.action, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tabcmd.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tabcmd.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

// from https://github.com/golang/go/blob/dca707b2a040642bb46aa4da4fb4eb6188cc2502/src/os/exec/exec_test.go#L724
// TestHelperProcess isn't a real test. It's used as a helper process
// for TestTabcmd_Run.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, action, args := args[0], args[1], args[2:]
	if cmd != cTabcmd {
		os.Stderr.WriteString("Invalid tabcmd executable")
		os.Exit(1)
	}
	switch action {
	case "succeed":
		os.Stdout.WriteString(fmt.Sprintf("Success: %v", args))
		os.Exit(0)
	case "fail":
		os.Stderr.WriteString(fmt.Sprintf("Failure: %s", args))
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}

}
