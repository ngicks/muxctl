package tmuxcli

import "context"

// new-session [-AdDEPX] [-c start-directory] [-e environment] [-f flags] [-F format]
//         [-n window-name] [-s session-name] [-t group-name] [-x width] [-y height]
//         [shell-command [argument ...]]
//               (alias: new)
//         Create a new session with name session-name.
//
//         The new session is attached to the current terminal unless -d is given.
//         window-name and shell-command are the name of and shell command to execute
//         in the initial window.  With -d, the initial size comes from the global
//         default-size option; -x and -y can be used to specify a different size.  ‘-’
//         uses the size of the current client if any.  If -x or -y is given, the
//         default-size option is set for the session.  -f sets a comma-separated list
//         of client flags (see attach-session).
//
//         If run from a terminal, any termios(4) special characters are saved and used
//         for new windows in the new session.
//
//         The -A flag makes new-session behave like attach-session if session-name
//         already exists; if -A is given, -D behaves like -d to attach-session, and -X
//         behaves like -x to attach-session.
//
//         If -t is given, it specifies a session group.  Sessions in the same group
//         share the same set of windows - new windows are linked to all sessions in
//         the group and any windows closed removed from all sessions.  The current and
//         previous window and any session options remain independent and any session
//         in a group may be killed without affecting the others.  The group-name
//         argument may be:
//
//         1.      the name of an existing group, in which case the new session is
//                 added to that group;
//
//         2.      the name of an existing session - the new session is added to the
//                 same group as that session, creating a new group if necessary;
//
//         3.      the name for a new group containing only the new session.
//
//         -n and shell-command are invalid if -t is used.
//
//         The -P option prints information about the new session after it has been
//         created.  By default, it uses the format ‘#{session_name}:’ but a different
//         format may be specified with -F.
//
//         If -E is used, the update-environment option will not be applied.  -e takes
//         the form ‘VARIABLE=value’ and sets an environment variable for the newly
//         created session; it may be specified multiple times.

// NewSessionOption holds the flags of new-session.
// -P and -F are not exposed: the method owns them to return a [SpawnResult].
type NewSessionOption struct {
	// AttachExisting is -A: attach to SessionName instead if it already exists.
	AttachExisting bool
	// Detached is -d: do not attach the new session to the current terminal.
	Detached bool
	// DetachOthers is -D: with -A, detach other clients attached to the session.
	DetachOthers bool
	// NoUpdateEnvironment is -E: do not apply the update-environment option.
	NoUpdateEnvironment bool
	// HangupOthers is -X: like DetachOthers, also sending SIGHUP to their parent processes.
	HangupOthers bool
	// StartDirectory is -c start-directory.
	StartDirectory string
	// Environment is -e environment, each "VARIABLE=value".
	Environment []string
	// ClientFlags is -f flags: a comma-separated list of client flags.
	ClientFlags string
	// WindowName is -n window-name.
	WindowName string
	// SessionName is -s session-name.
	SessionName string
	// GroupName is -t group-name: the session group to join or create.
	GroupName string
	// Width is -x width.
	Width string
	// Height is -y height.
	Height string
}

// NewSession runs new-session.
// Without Detached, or with AttachExisting when the session exists, tmux
// attaches to the terminal, which the captured tmux process lacks;
// use [ShellOut.AttachSession] for that.
func (s *ShellOut) NewSession(ctx context.Context, opt NewSessionOption, shellCommand ...string) (SpawnResult, error) {
	args := opt.args(argv{"new-session"})
	return s.spawn(ctx, args, shellCommand)
}

func (o NewSessionOption) args(a argv) argv {
	return a.
		flag("-A", o.AttachExisting).
		flag("-d", o.Detached).
		flag("-D", o.DetachOthers).
		flag("-E", o.NoUpdateEnvironment).
		flag("-X", o.HangupOthers).
		value("-c", o.StartDirectory).
		values("-e", o.Environment).
		value("-f", o.ClientFlags).
		value("-n", o.WindowName).
		value("-s", o.SessionName).
		value("-t", o.GroupName).
		value("-x", o.Width).
		value("-y", o.Height)
}
