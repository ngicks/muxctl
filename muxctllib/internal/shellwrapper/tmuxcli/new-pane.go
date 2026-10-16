package tmuxcli

import "context"

// new-pane [-bdefhIkPvZ] [-c start-directory] [-e environment] [-F format] [-l size]
//         [-m message] [-p percentage] [-R inactive-border-style] [-s style] [-S
//         active-border-style] [-t target-pane] [shell-command [argument ...]]
//               (alias: newp)
//         Create a new pane.  The new pane is created by splitting target-pane.  If -d
//         is given, the session does not make the new pane the current pane.  -Z zooms
//         if the window is not zoomed, or keeps it zoomed if already zoomed.  -s sets
//         the style for the pane content.  -S sets the border style when the pane is
//         active and -R sets the border style when the pane is inactive (see
//         “STYLES”).
//
//         -h does a horizontal split and -v a vertical split; if neither is specified,
//         -v is assumed.  The -l option specifies the size of the new pane in lines
//         (for vertical split) or in columns (for horizontal split); size may be
//         followed by ‘%’ to specify a percentage of the available space.  -p is a
//         shorthand option for this.  The -b option causes the new pane to be created
//         to the left of or above target-pane.  The -f option creates a new pane
//         spanning the full window height (with -h) or full window width (with -v),
//         instead of splitting the active pane.
//
//         -k keeps the pane open after the optional shell-command exits and waits for
//         a key to be pressed before closing it.  The message shown is controlled by
//         the remain-on-exit-format option.  -m message is equivalent to -k but also
//         sets the remain-on-exit-format option for this pane to message.
//
//         -E, or an empty shell-command, ('') will create an empty pane with no
//         command running in it; display-message -I can write to an empty pane.  The
//         -I flag will create an empty pane and forward any output from stdin to it.
//         For example:
//
//               $ make 2>&1|tmux splitw -dI &
//
//         All other options have the same meaning as for the new-window command.

// PaneSplitOption holds the flags split-window shares with new-pane.
type PaneSplitOption struct {
	SpawnOption
	// Before is -b: create the pane to the left of or above TargetPane.
	Before bool
	// Empty is -E: create a pane with no command running in it.
	// It is missing from the synopsis and from tmux list-commands, but described
	// in the text and accepted by tmux 3.7c.
	Empty bool
	// Full is -f: span the full window height (with -h) or width (with -v).
	Full bool
	// Horizontal is -h: split horizontally.
	Horizontal bool
	// ForwardStdin is -I: create an empty pane and forward stdin to it.
	ForwardStdin bool
	// KeepOnExit is -k: keep the pane open after the command exits until a key is pressed.
	KeepOnExit bool
	// Size is -l size: lines or columns, or a percentage followed by '%'.
	Size string
	// Message is -m message: like -k, also setting remain-on-exit-format for the pane.
	Message string
	// Percentage is -p percentage.
	Percentage *int
	// InactiveBorderStyle is -R inactive-border-style.
	InactiveBorderStyle string
	// Style is -s style: the style of the pane content.
	Style string
	// ActiveBorderStyle is -S active-border-style.
	ActiveBorderStyle string
	// TargetPane is -t target-pane.
	TargetPane string
	// Vertical is -v: split vertically, the default.
	Vertical bool
	// Zoom is -Z: zoom the window, or keep it zoomed.
	Zoom bool
}

// NewPaneOption holds the flags of new-pane.
// -P and -F are not exposed: the method owns them to return a [SpawnResult].
type NewPaneOption struct {
	PaneSplitOption

	// -x, -y, -X and -Y are missing from the man page but listed by
	// tmux list-commands and accepted by tmux 3.7c; split-window rejects them.

	// Width is -x width of the floating pane: columns, or a percentage followed by '%'.
	Width string
	// Height is -y height of the floating pane: lines, or a percentage followed by '%'.
	Height string
	// X is -X x-position of the floating pane.
	X string
	// Y is -Y y-position of the floating pane.
	Y string
}

// NewPane runs new-pane.
// tmux 3.7c makes the new pane floating; [ShellOut.SplitWindow] makes a tiled one.
func (s *ShellOut) NewPane(ctx context.Context, opt NewPaneOption, shellCommand ...string) (SpawnResult, error) {
	args := opt.args(argv{"new-pane"})
	return s.spawn(ctx, args, shellCommand)
}

func (o PaneSplitOption) args(a argv) argv {
	return o.SpawnOption.args(a).
		flag("-b", o.Before).
		flag("-E", o.Empty).
		flag("-f", o.Full).
		flag("-h", o.Horizontal).
		flag("-I", o.ForwardStdin).
		flag("-k", o.KeepOnExit).
		value("-l", o.Size).
		value("-m", o.Message).
		int("-p", o.Percentage).
		value("-R", o.InactiveBorderStyle).
		value("-s", o.Style).
		value("-S", o.ActiveBorderStyle).
		value("-t", o.TargetPane).
		flag("-v", o.Vertical).
		flag("-Z", o.Zoom)
}

func (o NewPaneOption) args(a argv) argv {
	return o.PaneSplitOption.args(a).
		value("-x", o.Width).
		value("-y", o.Height).
		value("-X", o.X).
		value("-Y", o.Y)
}
