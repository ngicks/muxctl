package zellijcli

import "context"

// Close the focused pane
//
// Usage: zellij action close-pane [OPTIONS]
//
// Options:
//   -p, --pane-id <PANE_ID>  Target a specific pane by ID (eg. terminal_1, plugin_2, or 3)
//   -h, --help               Print help

type ActionClosePaneOption struct {
	// PaneId is --pane-id. Nil closes the focused pane.
	PaneId *PaneId
}

// ActionClosePane runs action close-pane.
// zellij reports success even if the pane does not exist.
func (s *ShellOut) ActionClosePane(ctx context.Context, opt ActionClosePaneOption) error {
	args := opt.args(argv{"action", "close-pane"})
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionClosePaneOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
