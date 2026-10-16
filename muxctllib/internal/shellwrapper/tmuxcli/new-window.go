package tmuxcli

import "context"

// new-window [-abdkPS] [-c start-directory] [-e environment] [-F format] [-n
//         window-name] [-t target-window] [shell-command [argument ...]]
//               (alias: neww)
//         Create a new window.  With -a or -b, the new window is inserted at the next
//         index after or before the specified target-window, moving windows up if
//         necessary; otherwise target-window is the new window location.
//
//         If -d is given, the session does not make the new window the current window.
//         target-window represents the window to be created; if the target already
//         exists an error is shown, unless the -k flag is used, in which case it is
//         destroyed.  If -S is given and a window named window-name already exists, it
//         is selected (unless -d is also given in which case the command does
//         nothing).
//
//         shell-command is the command to execute.  If shell-command is not specified,
//         the value of the default-command option is used.  -c specifies the working
//         directory in which the new window is created.
//
//         When the shell command completes, the window closes.  See the remain-on-exit
//         option to change this behaviour.
//
//         -e takes the form ‘VARIABLE=value’ and sets an environment variable for the
//         newly created window; it may be specified multiple times.
//
//         The TERM environment variable must be set to ‘screen’ or ‘tmux’ for all
//         programs running inside tmux.  New windows will automatically have
//         ‘TERM=screen’ added to their environment, but care must be taken not to
//         reset this in shell startup files or by the -e option.
//
//         The -P option prints information about the new window after it has been
//         created.  By default, it uses the format ‘#{session_name}:#{window_index}’
//         but a different format may be specified with -F.

// SpawnOption holds the flags that new-pane and split-window take
// with the same meaning as for new-window.
type SpawnOption struct {
	// Detached is -d: the session does not make the new window or pane current.
	Detached bool
	// StartDirectory is -c start-directory.
	StartDirectory string
	// Environment is -e environment, each "VARIABLE=value".
	Environment []string
}

// NewWindowOption holds the flags of new-window.
// -P and -F are not exposed: the method owns them to return a [SpawnResult].
type NewWindowOption struct {
	SpawnOption
	// After is -a: insert at the next index after TargetWindow.
	After bool
	// Before is -b: insert at the index before TargetWindow.
	Before bool
	// Kill is -k: destroy the window at TargetWindow if it already exists.
	Kill bool
	// SelectExisting is -S: select the window named WindowName if it already exists.
	SelectExisting bool
	// WindowName is -n window-name.
	WindowName string
	// TargetWindow is -t target-window.
	TargetWindow string
}

// SpawnResult identifies the session, window and pane a command created.
// It is the zero value when tmux created nothing and so printed nothing,
// as when new-window -S selects an existing window.
type SpawnResult struct {
	SessionId   string // session_id, e.g. "$0"
	SessionName string // session_name
	WindowId    string // window_id, e.g. "@0"
	WindowIndex int    // window_index
	PaneId      string // pane_id, e.g. "%0"
}

// NewWindow runs new-window.
func (s *ShellOut) NewWindow(ctx context.Context, opt NewWindowOption, shellCommand ...string) (SpawnResult, error) {
	args := opt.args(argv{"new-window"})
	return s.spawn(ctx, args, shellCommand)
}

func (o SpawnOption) args(a argv) argv {
	return a.
		flag("-d", o.Detached).
		value("-c", o.StartDirectory).
		values("-e", o.Environment)
}

func (o NewWindowOption) args(a argv) argv {
	return o.SpawnOption.args(a).
		flag("-a", o.After).
		flag("-b", o.Before).
		flag("-k", o.Kill).
		flag("-S", o.SelectExisting).
		value("-n", o.WindowName).
		value("-t", o.TargetWindow)
}

var spawnVars = []string{
	"session_id",
	"session_name",
	"window_id",
	"window_index",
	"pane_id",
}

// spawn runs a command that creates a session, window or pane,
// appending -P -F and then shellCommand to args.
func (s *ShellOut) spawn(ctx context.Context, args argv, shellCommand []string) (SpawnResult, error) {
	args = args.
		flag("-P", true).
		value("-F", format(spawnVars)).
		trailing(shellCommand)
	out, err := s.run(ctx, args...)
	if err != nil {
		return SpawnResult{}, err
	}
	return parseSpawnResult(out)
}

func parseSpawnResult(out string) (SpawnResult, error) {
	rows, err := records(out, len(spawnVars))
	if err != nil || len(rows) == 0 {
		return SpawnResult{}, err
	}
	var r SpawnResult
	err = parseFields(spawnVars, rows[0], &r.SessionId, &r.SessionName, &r.WindowId, &r.WindowIndex, &r.PaneId)
	return r, err
}
