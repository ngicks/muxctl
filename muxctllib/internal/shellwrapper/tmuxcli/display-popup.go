package tmuxcli

import (
	"context"
	"os/exec"
)

// display-popup [-BCEkN] [-b border-lines] [-c target-client] [-d start-directory] [-e
//         environment] [-h height] [-s style] [-S border-style] [-t target-pane] [-T
//         title] [-w width] [-x position] [-y position] [shell-command [argument ...]]
//               (alias: popup)
//         Display a popup running shell-command (or default-command when omitted) on
//         target-client.  A popup is a rectangular box drawn over the top of any
//         panes.  If the command is run inside an existing popup, that popup is
//         modified.  Only the -b, -B, -C, -E, -EE, -K, -N, -s, and -S options are
//         accepted in this case; all others are ignored.
//
//         -E closes the popup automatically when shell-command exits.  Two -E closes
//         the popup only if shell-command exited with success.  -k allows any key to
//         dismiss the popup instead of only ‘Escape’ or ‘C-c’.
//
//         -x and -y give the position of the popup, they have the same meaning as for
//         the display-menu command.  -w and -h give the width and height - both may be
//         a percentage (followed by ‘%’).  If omitted, half of the terminal size is
//         used.
//
//         -B does not surround the popup by a border.
//
//         -b sets the type of characters used for drawing popup borders.  When -B is
//         specified, the -b option is ignored.  See popup-border-lines for possible
//         values for border-lines.
//
//         -s sets the style for the popup and -S sets the style for the popup border
//         (see “STYLES”).
//
//         -e takes the form ‘VARIABLE=value’ and sets an environment variable for the
//         popup; it may be specified multiple times.
//
//         -T is a format for the popup title (see “FORMATS”).
//
//         The -C flag closes any popup on the client.
//
//         -N disables any previously specified -E, -EE, or -k option.

type DisplayPopupOption struct {
	// NoBorder is -B: do not surround the popup by a border.
	NoBorder bool
	// Close is -C: close any popup on the client.
	Close bool
	// CloseOnExit is -E: close the popup when the command exits.
	CloseOnExit bool
	// CloseOnSuccess is -EE: close the popup only when the command exits with success.
	// It takes precedence over CloseOnExit.
	CloseOnSuccess bool
	// AnyKey is -k: any key dismisses the popup, not only Escape or C-c.
	AnyKey bool
	// DisableClose is -N: disable any previously given -E, -EE or -k.
	DisableClose bool
	// BorderLines is -b border-lines.
	BorderLines string
	// TargetClient is -c target-client.
	TargetClient string
	// StartDirectory is -d start-directory.
	StartDirectory string
	// Environment is -e environment, each "VARIABLE=value".
	Environment []string
	// Height is -h height: lines, or a percentage followed by '%'.
	Height string
	// Style is -s style.
	Style string
	// BorderStyle is -S border-style.
	BorderStyle string
	// TargetPane is -t target-pane.
	TargetPane string
	// Title is -T title, a format.
	Title string
	// Width is -w width: columns, or a percentage followed by '%'.
	Width string
	// X is -x position.
	X string
	// Y is -y position.
	Y string
}

// DisplayPopup returns an unstarted display-popup command.
// -E/popups block until closed, so the caller owns the process lifecycle.
func (s *ShellOut) DisplayPopup(ctx context.Context, opt DisplayPopupOption, shellCommand ...string) *exec.Cmd {
	args := opt.args(argv{"display-popup"}).trailing(shellCommand)
	return s.command(ctx, args...)
}

func (o DisplayPopupOption) args(a argv) argv {
	return a.
		flag("-B", o.NoBorder).
		flag("-C", o.Close).
		flag("-E", o.CloseOnExit || o.CloseOnSuccess).
		flag("-E", o.CloseOnSuccess).
		flag("-k", o.AnyKey).
		flag("-N", o.DisableClose).
		value("-b", o.BorderLines).
		value("-c", o.TargetClient).
		value("-d", o.StartDirectory).
		values("-e", o.Environment).
		value("-h", o.Height).
		value("-s", o.Style).
		value("-S", o.BorderStyle).
		value("-t", o.TargetPane).
		value("-T", o.Title).
		value("-w", o.Width).
		value("-x", o.X).
		value("-y", o.Y)
}
