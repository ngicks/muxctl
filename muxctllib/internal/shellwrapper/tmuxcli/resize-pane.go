package tmuxcli

import (
	"context"
	"strconv"
)

// resize-pane [-DLMRTUZ] [-t target-pane] [-x width] [-y height] [adjustment]
//               (alias: resizep)
//         Resize a pane, up, down, left or right by adjustment with -U, -D, -L or -R,
//         or to an absolute size with -x or -y.  The adjustment is given in lines or
//         columns (the default is 1); -x and -y may be a given as a number of lines or
//         columns or followed by ‘%’ for a percentage of the window size (for example
//         ‘-x 10%’).  With -Z, the active pane is toggled between zoomed (occupying
//         the whole of the window) and unzoomed (its normal position in the layout).
//
//         -M begins mouse resizing (only valid if bound to a mouse key binding, see
//         “MOUSE SUPPORT”).
//
//         -T trims all lines below the current cursor position and moves lines out of
//         the history to replace them.

type ResizePaneOption struct {
	// Down is -D: resize downward by the adjustment.
	Down bool
	// Left is -L: resize to the left by the adjustment.
	Left bool
	// Mouse is -M: begin mouse resizing.
	Mouse bool
	// Right is -R: resize to the right by the adjustment.
	Right bool
	// Trim is -T: trim lines below the cursor, replacing them from the history.
	Trim bool
	// Up is -U: resize upward by the adjustment.
	Up bool
	// Zoom is -Z: toggle the zoom of the active pane.
	Zoom bool
	// Width is -x width: columns, or a percentage followed by '%'.
	Width string
	// Height is -y height: lines, or a percentage followed by '%'.
	Height string
	// TargetPane is -t target-pane.
	TargetPane string
}

// ResizePane runs resize-pane.
// A zero adjustment is omitted, which tmux takes as 1.
func (s *ShellOut) ResizePane(ctx context.Context, opt ResizePaneOption, adjustment int) error {
	args := opt.args(argv{"resize-pane"})
	if adjustment != 0 {
		args = append(args, strconv.Itoa(adjustment))
	}
	_, err := s.run(ctx, args...)
	return err
}

func (o ResizePaneOption) args(a argv) argv {
	return a.
		flag("-D", o.Down).
		flag("-L", o.Left).
		flag("-M", o.Mouse).
		flag("-R", o.Right).
		flag("-T", o.Trim).
		flag("-U", o.Up).
		flag("-Z", o.Zoom).
		value("-x", o.Width).
		value("-y", o.Height).
		value("-t", o.TargetPane)
}
