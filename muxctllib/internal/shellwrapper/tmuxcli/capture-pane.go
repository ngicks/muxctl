package tmuxcli

import "context"

// capture-pane [-aeFHLpPqCJMN] [-b buffer-name] [-E end-line] [-S start-line] [-t
//         target-pane]
//               (alias: capturep)
//         Capture the contents of a pane.  If -p is given, the output goes to stdout,
//         otherwise to the buffer specified with -b or a new buffer if omitted.  If -a
//         is given, the alternate screen is used, and the history is not accessible.
//         If no alternate screen exists, an error will be returned unless -q is given.
//         Similarly, if the pane is in a mode, -M uses the screen for the mode.  If -e
//         is given, the output includes escape sequences for text and background
//         attributes.  -C also escapes non-printable characters as octal \xxx.  -T
//         ignores trailing positions that do not contain a character.  -N preserves
//         trailing spaces at each line's end and -J preserves trailing spaces and
//         joins any wrapped lines; -J implies -T.  -P captures only any output that
//         the pane has received that is the beginning of an as-yet incomplete escape
//         sequence.  -L includes the line number at the start of each line and -F
//         includes the flags (where ‘-’ is no flags, ‘D’ is an unused line, ‘O’ is a
//         line marked as output, ‘P’ is a line marked as a prompt, ‘X’ is a line
//         containing extended cells and ‘H’ is a line with hyperlinks).  With -H, only
//         hyperlinks in the specified lines are captured.  Multiple hyperlinks on the
//         same line are separated by spaces.
//
//         -S and -E specify the starting and ending line numbers, zero is the first
//         line of the visible pane and negative numbers are lines in the history.  ‘-’
//         to -S is the start of the history and to -E the end of the visible pane.
//         The default is to capture only the visible contents of the pane.

type CapturePaneOption struct {
	// Alternate is -a: capture the alternate screen.
	Alternate bool
	// EscapeNonPrintable is -C: also escape non-printable characters as octal \xxx.
	EscapeNonPrintable bool
	// EscapeSequences is -e: include escape sequences for text and background attributes.
	EscapeSequences bool
	// LineFlags is -F: prefix each line with its flags.
	LineFlags bool
	// Hyperlinks is -H: capture only the hyperlinks in the lines.
	Hyperlinks bool
	// JoinWrapped is -J: preserve trailing spaces and join wrapped lines.
	JoinWrapped bool
	// LineNumbers is -L: prefix each line with its line number.
	LineNumbers bool
	// Mode is -M: capture the screen of the mode the pane is in.
	Mode bool
	// PreserveTrailingSpaces is -N.
	PreserveTrailingSpaces bool
	// Print is -p: write to stdout instead of a buffer.
	Print bool
	// IncompleteEscape is -P: capture only the start of an incomplete escape sequence.
	IncompleteEscape bool
	// Quiet is -q: no error when there is no alternate screen.
	Quiet bool
	// TrimTrailing is -T: ignore trailing positions that do not contain a character.
	TrimTrailing bool
	// BufferName is -b buffer-name.
	BufferName string
	// EndLine is -E end-line: a line number, or "-" for the end of the visible pane.
	EndLine string
	// StartLine is -S start-line: a line number, or "-" for the start of the history.
	StartLine string
	// TargetPane is -t target-pane.
	TargetPane string
}

// CapturePane runs capture-pane and returns its stdout as is,
// which is empty unless Print is set.
func (s *ShellOut) CapturePane(ctx context.Context, opt CapturePaneOption) (string, error) {
	args := opt.args(argv{"capture-pane"})
	return s.run(ctx, args...)
}

func (o CapturePaneOption) args(a argv) argv {
	return a.
		flag("-a", o.Alternate).
		flag("-C", o.EscapeNonPrintable).
		flag("-e", o.EscapeSequences).
		flag("-F", o.LineFlags).
		flag("-H", o.Hyperlinks).
		flag("-J", o.JoinWrapped).
		flag("-L", o.LineNumbers).
		flag("-M", o.Mode).
		flag("-N", o.PreserveTrailingSpaces).
		flag("-p", o.Print).
		flag("-P", o.IncompleteEscape).
		flag("-q", o.Quiet).
		flag("-T", o.TrimTrailing).
		value("-b", o.BufferName).
		value("-E", o.EndLine).
		value("-S", o.StartLine).
		value("-t", o.TargetPane)
}
