package tmuxcli

import "context"

// select-pane [-DdeLlMmRUZ] [-T title] [-t target-pane]
//               (alias: selectp)
//         Make pane target-pane the active pane in its window.  If one of -D, -L, -R,
//         or -U is used, respectively the pane below, to the left, to the right, or
//         above the target pane is used.  -Z keeps the window zoomed if it was zoomed.
//         -l is the same as using the last-pane command.  -e enables or -d disables
//         input to the pane.  -T sets the pane title.
//
//         -m and -M are used to set and clear the marked pane.  There is one marked
//         pane at a time, setting a new marked pane clears the last.  The marked pane
//         is the default target for -s to join-pane, move-pane, swap-pane and
//         swap-window.

type SelectPaneOption struct {
	// Down is -D: select the pane below TargetPane.
	Down bool
	// DisableInput is -d: disable input to the pane.
	DisableInput bool
	// EnableInput is -e: enable input to the pane.
	EnableInput bool
	// Left is -L: select the pane to the left of TargetPane.
	Left bool
	// Last is -l: same as the last-pane command.
	Last bool
	// ClearMarked is -M: clear the marked pane.
	ClearMarked bool
	// Mark is -m: make the pane the marked pane.
	Mark bool
	// Right is -R: select the pane to the right of TargetPane.
	Right bool
	// Up is -U: select the pane above TargetPane.
	Up bool
	// KeepZoom is -Z: keep the window zoomed if it was zoomed.
	KeepZoom bool
	// Title is -T title.
	Title string
	// TargetPane is -t target-pane.
	TargetPane string
}

// SelectPane runs select-pane.
func (s *ShellOut) SelectPane(ctx context.Context, opt SelectPaneOption) error {
	args := opt.args(argv{"select-pane"})
	_, err := s.run(ctx, args...)
	return err
}

func (o SelectPaneOption) args(a argv) argv {
	return a.
		flag("-D", o.Down).
		flag("-d", o.DisableInput).
		flag("-e", o.EnableInput).
		flag("-L", o.Left).
		flag("-l", o.Last).
		flag("-M", o.ClearMarked).
		flag("-m", o.Mark).
		flag("-R", o.Right).
		flag("-U", o.Up).
		flag("-Z", o.KeepZoom).
		value("-T", o.Title).
		value("-t", o.TargetPane)
}
