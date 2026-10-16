// Package zellijcli wraps the zellij command line interface.
//
// Every zellij command wrapped here lives in its own file named after the
// command path (e.g. action_new-pane.go for "zellij action new-pane") and is
// exposed as one method on [ShellOut] whose name keeps the whole path
// (ActionNewPane). Flags of a command are gathered in <Method>Option,
// positional arguments are method parameters, and output is parsed into typed
// results where zellij prints something meaningful.
package zellijcli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ShellOut runs zellij as a child process.
// The zero value runs "zellij" from PATH, inheriting the environment and
// working directory of the current process.
type ShellOut struct {
	// Path is the zellij executable. Empty means "zellij".
	Path string
	// Global holds flags placed before the command name.
	Global GlobalOption
	// Env is the environment of the zellij process. Nil inherits the current one.
	// zellij reads $ZELLIJ_SESSION_NAME from it to pick the session actions are
	// sent to when Global.Session is empty.
	Env []string
	// Dir is the working directory of the zellij process. Empty inherits the current one.
	Dir string
}

// GlobalOption holds the zellij flags that precede the command.
type GlobalOption struct {
	// Session is --session. zellij's help only mentions naming a new session,
	// but action and run are sent to this session too.
	Session string
	// ConfigFile is --config: an alternative configuration file.
	ConfigFile string
	// ConfigDir is --config-dir: an alternative configuration directory.
	ConfigDir string
	// DataDir is --data-dir: where zellij looks for plugins.
	DataDir string
	// Debug is --debug.
	Debug bool
}

func (o GlobalOption) args(a argv) argv {
	return a.
		value("--session", o.Session).
		value("--config", o.ConfigFile).
		value("--config-dir", o.ConfigDir).
		value("--data-dir", o.DataDir).
		flag("--debug", o.Debug)
}

// CmdError is returned when zellij exits unsuccessfully or cannot be started.
type CmdError struct {
	// Args is the full argument list zellij was invoked with, global flags included.
	Args []string
	// Stderr is what zellij wrote to stderr, with surrounding white space trimmed.
	Stderr string
	Err    error
}

func (e *CmdError) Error() string {
	if e.Stderr == "" {
		return fmt.Sprintf("zellij %s: %v", strings.Join(e.Args, " "), e.Err)
	}
	return fmt.Sprintf("zellij %s: %v: %s", strings.Join(e.Args, " "), e.Err, e.Stderr)
}

func (e *CmdError) Unwrap() error {
	return e.Err
}

// ExitCode returns zellij's exit status, or -1 if zellij did not exit normally.
func (e *CmdError) ExitCode() int {
	if exitErr, ok := errors.AsType[*exec.ExitError](e.Err); ok {
		return exitErr.ExitCode()
	}
	return -1
}

// PaneId identifies a pane within a session.
// zellij numbers terminal and plugin panes independently, so both parts are needed.
type PaneId struct {
	Plugin bool
	Id     int
}

// String returns the form zellij prints and accepts, "terminal_N" or "plugin_N".
func (p PaneId) String() string {
	if p.Plugin {
		return "plugin_" + strconv.Itoa(p.Id)
	}
	return "terminal_" + strconv.Itoa(p.Id)
}

// ParsePaneId parses "terminal_N", "plugin_N", or a bare "N", which zellij takes as a terminal.
func ParsePaneId(s string) (PaneId, error) {
	var p PaneId
	num := s
	if n, ok := strings.CutPrefix(s, "terminal_"); ok {
		num = n
	} else if n, ok := strings.CutPrefix(s, "plugin_"); ok {
		num = n
		p.Plugin = true
	}
	id, err := strconv.ParseUint(num, 10, 32)
	if err != nil {
		return PaneId{}, fmt.Errorf("zellij: parse pane id %q: %w", s, err)
	}
	p.Id = int(id)
	return p, nil
}

func (s *ShellOut) command(ctx context.Context, args ...string) *exec.Cmd {
	path := s.Path
	if path == "" {
		path = "zellij"
	}
	// Sized for every global flag so neither the flags nor args are copied twice.
	full := make(argv, 0, 9+len(args))
	full = s.Global.args(full)
	full = append(full, args...)
	cmd := exec.CommandContext(ctx, path, full...)
	cmd.Env = s.Env
	cmd.Dir = s.Dir
	return cmd
}

// run runs zellij to completion and returns its stdout as is.
func (s *ShellOut) run(ctx context.Context, args ...string) (string, error) {
	cmd := s.command(ctx, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", &CmdError{
			Args:   cmd.Args[1:],
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}
	return stdout.String(), nil
}

// argv accumulates command line arguments.
// Options are always spelled "--name=value": short letters differ between
// subcommands (-b is --borderless for run but --blocking for action new-pane),
// and a value starting with '-' cannot be mistaken for a flag.
// Empty values, false flags and nil pointers are omitted, which matches zellij's own defaults.
type argv []string

func (a argv) flag(name string, on bool) argv {
	if on {
		a = append(a, name)
	}
	return a
}

func (a argv) value(name, v string) argv {
	if v != "" {
		a = append(a, name+"="+v)
	}
	return a
}

func (a argv) int(name string, v *int) argv {
	if v != nil {
		a = append(a, name+"="+strconv.Itoa(*v))
	}
	return a
}

func (a argv) bool(name string, v *bool) argv {
	if v != nil {
		a = append(a, name+"="+strconv.FormatBool(*v))
	}
	return a
}

func (a argv) paneId(name string, v *PaneId) argv {
	if v != nil {
		a = append(a, name+"="+v.String())
	}
	return a
}

// trailing appends an optional trailing command after "--",
// so its arguments are not parsed as zellij flags.
func (a argv) trailing(command []string) argv {
	if len(command) > 0 {
		a = append(a, "--")
		a = append(a, command...)
	}
	return a
}
