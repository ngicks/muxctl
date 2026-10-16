package tmuxcli

import "context"

// respawn-pane [-k] [-c start-directory] [-e environment] [-t target-pane]
//         [shell-command [argument ...]]
//               (alias: respawnp)
//         Reactivate a pane in which the command has exited (see the remain-on-exit
//         window option).  If shell-command is not given, the command used when the
//         pane was created or last respawned is executed.  The pane must be already
//         inactive, unless -k is given, in which case any existing command is killed.
//         -c specifies a new working directory for the pane.  The -e option has the
//         same meaning as for the new-window command.

type RespawnPaneOption struct {
	// Kill is -k: kill the command still running in the pane.
	Kill bool
	// StartDirectory is -c start-directory.
	StartDirectory string
	// Environment is -e environment, each "VARIABLE=value".
	Environment []string
	// TargetPane is -t target-pane.
	TargetPane string
}

// RespawnPane runs respawn-pane.
// An empty shellCommand reruns the command the pane was created or last respawned with.
func (s *ShellOut) RespawnPane(ctx context.Context, opt RespawnPaneOption, shellCommand ...string) error {
	args := opt.args(argv{"respawn-pane"}).trailing(shellCommand)
	_, err := s.run(ctx, args...)
	return err
}

func (o RespawnPaneOption) args(a argv) argv {
	return a.
		flag("-k", o.Kill).
		value("-c", o.StartDirectory).
		values("-e", o.Environment).
		value("-t", o.TargetPane)
}
