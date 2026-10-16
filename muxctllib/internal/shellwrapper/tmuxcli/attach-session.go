package tmuxcli

import (
	"context"
	"os/exec"
)

// attach-session [-dErx] [-c working-directory] [-f flags] [-t target-session]
//               (alias: attach)
//         If run from outside tmux, attach to target-session in the current terminal.
//         target-session must already exist - to create a new session, see the
//         new-session command (with -A to create or attach).  If used from inside,
//         switch the currently attached session to target-session.  If -d is
//         specified, any other clients attached to the session are detached.  If -x is
//         given, send SIGHUP to the parent process of the client as well as detaching
//         the client, typically causing it to exit.  -f sets a comma-separated list of
//         client flags.  The flags are:
//
//         active-pane
//                 the client has an independent active pane
//
//         ignore-size
//                 the client does not affect the size of other clients
//
//         no-detach-on-destroy
//                 do not detach the client when the session it is attached to is
//                 destroyed if there are any other sessions
//
//         no-output
//                 the client does not receive pane output in control mode
//
//         pause-after=seconds
//                 output is paused once the pane is seconds behind in control mode
//
//         read-only
//                 the client is read-only
//
//         wait-exit
//                 wait for an empty line input before exiting in control mode
//
//         A leading ‘!’ turns a flag off if the client is already attached.  -r is an
//         alias for -f read-only,ignore-size.  When a client is read-only, only keys
//         bound to the detach-client or switch-client commands have any effect.  A
//         client with the active-pane flag allows the active pane to be selected
//         independently of the window's active pane used by clients without the flag.
//         This only affects the cursor position and commands issued from the client;
//         other features such as hooks and styles continue to use the window's active
//         pane.
//
//         If no server is started, attach-session will attempt to start it; this will
//         fail unless sessions are created in the configuration file.
//
//         The target-session rules for attach-session are slightly adjusted: if tmux
//         needs to select the most recently used session, it will prefer the most
//         recently used unattached session.
//
//         -c will set the session working directory (used for new windows) to
//         working-directory.
//
//         If -E is used, the update-environment option will not be applied.

type AttachSessionOption struct {
	// DetachOthers is -d: detach other clients attached to the session.
	DetachOthers bool
	// NoUpdateEnvironment is -E: do not apply the update-environment option.
	NoUpdateEnvironment bool
	// ReadOnly is -r: same as ClientFlags "read-only,ignore-size".
	ReadOnly bool
	// HangupOthers is -x: like DetachOthers, also sending SIGHUP to their parent processes.
	HangupOthers bool
	// WorkingDirectory is -c working-directory: the session working directory for new windows.
	WorkingDirectory string
	// ClientFlags is -f flags: a comma-separated list of client flags.
	ClientFlags string
	// TargetSession is -t target-session.
	TargetSession string
}

// AttachSession returns an unstarted attach-session command.
// The client takes over the terminal on the command's stdin, stdout and stderr
// and runs until it detaches, so the caller owns the process lifecycle.
func (s *ShellOut) AttachSession(ctx context.Context, opt AttachSessionOption) *exec.Cmd {
	args := opt.args(argv{"attach-session"})
	return s.command(ctx, args...)
}

func (o AttachSessionOption) args(a argv) argv {
	return a.
		flag("-d", o.DetachOthers).
		flag("-E", o.NoUpdateEnvironment).
		flag("-r", o.ReadOnly).
		flag("-x", o.HangupOthers).
		value("-c", o.WorkingDirectory).
		value("-f", o.ClientFlags).
		value("-t", o.TargetSession)
}
