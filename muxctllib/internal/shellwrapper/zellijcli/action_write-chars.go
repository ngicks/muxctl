package zellijcli

import "context"

// Write characters to the terminal
//
// Usage: zellij action write-chars [OPTIONS] <CHARS>
//
// Arguments:
//   <CHARS>
//
// Options:
//   -p, --pane-id <PANE_ID>  The pane_id of the pane, eg. terminal_1, plugin_2 or 3 (equivalent to
//                            terminal_3)
//   -h, --help               Print help

type ActionWriteCharsOption struct {
	// PaneId is --pane-id. Nil writes to the focused pane.
	PaneId *PaneId
}

// ActionWriteChars runs action write-chars.
func (s *ShellOut) ActionWriteChars(ctx context.Context, opt ActionWriteCharsOption, chars string) error {
	args := opt.args(argv{"action", "write-chars"})
	args = append(args, "--", chars)
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionWriteCharsOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
