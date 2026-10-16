package tmuxcli

import "context"

// kill-pane [-a] [-t target-pane]
//               (alias: killp)
//         Destroy the given pane.  If no panes remain in the containing window, it is
//         also destroyed.  The -a option kills all but the pane given with -t.

type KillPaneOption struct {
	// AllOthers is -a: kill all panes but TargetPane.
	AllOthers bool
	// TargetPane is -t target-pane.
	TargetPane string
}

// KillPane runs kill-pane.
func (s *ShellOut) KillPane(ctx context.Context, opt KillPaneOption) error {
	args := opt.args(argv{"kill-pane"})
	_, err := s.run(ctx, args...)
	return err
}

func (o KillPaneOption) args(a argv) argv {
	return a.
		flag("-a", o.AllOthers).
		value("-t", o.TargetPane)
}
