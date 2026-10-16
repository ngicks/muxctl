package zellijcli

import (
	"context"
	"os/exec"
)

// Attach to a session
//
// Usage: zellij attach [OPTIONS] [SESSION_NAME] [-- <INITIAL_COMMAND>...] [COMMAND]
//
// Commands:
//   options  Change the behaviour of zellij
//   help     Print this message or the help of the given subcommand(s)
//
// Arguments:
//   [SESSION_NAME]        Name of the session to attach to
//   [INITIAL_COMMAND]...  Command to run in the first pane of the session, if it is created
//
// Options:
//   -c, --create              Create a session if one does not exist
//   -b, --create-background   Create a detached session in the background if one does not exist
//       --index <INDEX>       Number of the session index in the active sessions ordered creation date
//   -f, --force-run-commands  If resurrecting a dead session, immediately run all its commands on
//                             startup
//   -t, --token <TOKEN>       Authentication token for remote sessions
//   -r, --remember            Save session for automatic re-authentication (4 weeks)
//       --forget              Delete saved session before connecting
//       --ca-cert <FILE>      Path to a custom CA certificate (PEM format) for verifying the remote
//                             server
//       --insecure            Skip TLS certificate validation (DANGEROUS — development only)
//       --close-on-exit       Close the initial command's pane immediately when it exits
//       --start-suspended     Start the initial command suspended, only running it after you first
//                             press ENTER
//   -h, --help                Print help

type AttachOption struct {
	// Create is --create.
	Create bool
	// CreateBackground is --create-background. The command then exits once the session exists.
	CreateBackground bool
	// Index is --index.
	Index *int
	// ForceRunCommands is --force-run-commands.
	ForceRunCommands bool
	// Token is --token.
	Token string
	// Remember is --remember.
	Remember bool
	// Forget is --forget.
	Forget bool
	// CaCert is --ca-cert.
	CaCert string
	// Insecure is --insecure.
	Insecure bool
	// CloseOnExit is --close-on-exit.
	CloseOnExit bool
	// StartSuspended is --start-suspended.
	StartSuspended bool
}

// Attach returns an unstarted zellij attach.
// An attached client runs until it detaches, so the caller owns the process
// lifecycle; with CreateBackground just Run it.
// The "options" subcommand is not supported.
func (s *ShellOut) Attach(ctx context.Context, opt AttachOption, sessionName string, initialCommand ...string) *exec.Cmd {
	args := opt.args(argv{"attach"})
	if sessionName != "" {
		args = append(args, sessionName)
	}
	args = args.trailing(initialCommand)
	return s.command(ctx, args...)
}

func (o AttachOption) args(a argv) argv {
	return a.
		flag("--create", o.Create).
		flag("--create-background", o.CreateBackground).
		int("--index", o.Index).
		flag("--force-run-commands", o.ForceRunCommands).
		value("--token", o.Token).
		flag("--remember", o.Remember).
		flag("--forget", o.Forget).
		value("--ca-cert", o.CaCert).
		flag("--insecure", o.Insecure).
		flag("--close-on-exit", o.CloseOnExit).
		flag("--start-suspended", o.StartSuspended)
}
