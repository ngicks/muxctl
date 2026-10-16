// Package tmuxcli wraps the tmux(1) command line interface.
//
// Every tmux command wrapped here lives in its own file named after the
// command (e.g. list-panes.go) and is exposed as one method on [ShellOut].
// Flags of a command are gathered in <Command>Option, positional arguments are
// method parameters, and output is parsed into typed results where tmux prints
// something meaningful.
package tmuxcli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ShellOut runs tmux as a child process.
// The zero value runs "tmux" from PATH against the default server,
// inheriting the environment and working directory of the current process.
type ShellOut struct {
	// Path is the tmux executable. Empty means "tmux".
	Path string
	// Global holds flags placed before the command name.
	Global GlobalOption
	// Env is the environment of the tmux process. Nil inherits the current one.
	// tmux reads $TMUX from it to pick the default server and session.
	Env []string
	// Dir is the working directory of the tmux process. Empty inherits the current one.
	Dir string
}

// GlobalOption holds the tmux flags that precede the command.
type GlobalOption struct {
	// SocketName is -L socket-name: a named server in tmux's socket directory.
	SocketName string
	// SocketPath is -S socket-path: a full path to the server socket.
	// tmux ignores -L when -S is given.
	SocketPath string
	// ConfigFile is -f file: an alternative configuration file.
	ConfigFile string
}

func (o GlobalOption) args(a argv) argv {
	return a.
		value("-L", o.SocketName).
		value("-S", o.SocketPath).
		value("-f", o.ConfigFile)
}

// CmdError is returned when tmux exits unsuccessfully or cannot be started.
type CmdError struct {
	// Args is the full argument list tmux was invoked with, global flags included.
	Args []string
	// Stderr is what tmux wrote to stderr, with surrounding white space trimmed.
	Stderr string
	Err    error
}

func (e *CmdError) Error() string {
	if e.Stderr == "" {
		return fmt.Sprintf("tmux %s: %v", strings.Join(e.Args, " "), e.Err)
	}
	return fmt.Sprintf("tmux %s: %v: %s", strings.Join(e.Args, " "), e.Err, e.Stderr)
}

func (e *CmdError) Unwrap() error {
	return e.Err
}

// ExitCode returns tmux's exit status, or -1 if tmux did not exit normally.
func (e *CmdError) ExitCode() int {
	if exitErr, ok := errors.AsType[*exec.ExitError](e.Err); ok {
		return exitErr.ExitCode()
	}
	return -1
}

func (s *ShellOut) command(ctx context.Context, args ...string) *exec.Cmd {
	path := s.Path
	if path == "" {
		path = "tmux"
	}
	// Sized for every global flag so neither the flags nor args are copied twice.
	full := make(argv, 0, 6+len(args))
	full = s.Global.args(full)
	full = append(full, args...)
	cmd := exec.CommandContext(ctx, path, full...)
	cmd.Env = s.Env
	cmd.Dir = s.Dir
	return cmd
}

// run runs tmux to completion and returns its stdout as is.
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
// Empty values and false flags are omitted, which matches tmux's own defaults.
type argv []string

func (a argv) flag(name string, on bool) argv {
	if on {
		a = append(a, name)
	}
	return a
}

func (a argv) value(name, v string) argv {
	if v != "" {
		a = append(a, name, v)
	}
	return a
}

func (a argv) values(name string, vs []string) argv {
	for _, v := range vs {
		a = append(a, name, v)
	}
	return a
}

func (a argv) int(name string, v *int) argv {
	if v != nil {
		a = append(a, name, strconv.Itoa(*v))
	}
	return a
}

// trailing appends trailing positional arguments after "--",
// so that tmux does not parse one starting with '-' as a flag.
func (a argv) trailing(vs []string) argv {
	if len(vs) > 0 {
		a = append(a, "--")
		a = append(a, vs...)
	}
	return a
}

// fieldSep separates fields in the formats this package builds for -F.
// tmux prints it raw, and unlike a tab it cannot occur in a pane title or an
// option value that tmux's parser accepted from a user.
const fieldSep = "\x1f"

// format joins format variables such as "pane_id" into a -F format string.
func format(vars []string) string {
	var b strings.Builder
	for i, v := range vars {
		if i > 0 {
			b.WriteString(fieldSep)
		}
		b.WriteString("#{")
		b.WriteString(v)
		b.WriteString("}")
	}
	return b.String()
}

// records splits output produced by a [format] of n variables into rows of n fields.
func records(out string, n int) ([][]string, error) {
	var rows [][]string
	for line := range strings.Lines(out) {
		line = strings.TrimSuffix(line, "\n")
		if line == "" {
			continue
		}
		fields := strings.Split(line, fieldSep)
		if len(fields) != n {
			return nil, fmt.Errorf("tmux: want %d fields, got %d in %q", n, len(fields), line)
		}
		rows = append(rows, fields)
	}
	return rows, nil
}

// userOptionVars returns format variables for user options ("@name") so they
// can be appended to a [format].
func userOptionVars(names []string) ([]string, error) {
	for _, n := range names {
		if !strings.HasPrefix(n, "@") {
			return nil, fmt.Errorf("tmux: user option name must start with '@': %q", n)
		}
	}
	return names, nil
}

// userOptionMap pairs names with the trailing fields of a record.
func userOptionMap(names, fields []string) map[string]string {
	if len(names) == 0 {
		return nil
	}
	m := make(map[string]string, len(names))
	for i, n := range names {
		m[n] = fields[i]
	}
	return m
}

func atoi(field, v string) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("tmux: parse %s: %w", field, err)
	}
	return n, nil
}

// parseFields assigns record fields to typed destinations in order.
// Each destination is *string, *int, or *bool ("1" is true).
func parseFields(names []string, fields []string, dsts ...any) error {
	for i, dst := range dsts {
		switch d := dst.(type) {
		case *string:
			*d = fields[i]
		case *bool:
			*d = fields[i] == "1"
		case *int:
			if fields[i] == "" {
				*d = 0
				continue
			}
			n, err := atoi(names[i], fields[i])
			if err != nil {
				return err
			}
			*d = n
		default:
			panic(fmt.Sprintf("tmuxcli: unsupported destination %T", dst))
		}
	}
	return nil
}
