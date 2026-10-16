package tmuxcli

import "context"

// kill-window [-a] [-t target-window]
//               (alias: killw)
//         Kill the current window or the window at target-window, removing it from any
//         sessions to which it is linked.  The -a option kills all but the window
//         given with -t.

type KillWindowOption struct {
	// AllOthers is -a: kill all windows but TargetWindow.
	AllOthers bool
	// TargetWindow is -t target-window.
	TargetWindow string
}

// KillWindow runs kill-window.
func (s *ShellOut) KillWindow(ctx context.Context, opt KillWindowOption) error {
	args := opt.args(argv{"kill-window"})
	_, err := s.run(ctx, args...)
	return err
}

func (o KillWindowOption) args(a argv) argv {
	return a.
		flag("-a", o.AllOthers).
		value("-t", o.TargetWindow)
}
