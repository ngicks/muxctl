package zellijcli

import "context"

// Renames the focused pane
//
// Usage: zellij action rename-pane [OPTIONS] <NAME>
//
// Arguments:
//   <NAME>
//
// Options:
//   -p, --pane-id <PANE_ID>  Target a specific pane by ID (eg. terminal_1, plugin_2, or 3)
//   -h, --help               Print help

type ActionRenamePaneOption struct {
	// PaneId is --pane-id. Nil renames the focused pane.
	PaneId *PaneId
}

// ActionRenamePane runs action rename-pane.
func (s *ShellOut) ActionRenamePane(ctx context.Context, opt ActionRenamePaneOption, name string) error {
	args := opt.args(argv{"action", "rename-pane"})
	args = append(args, "--", name)
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionRenamePaneOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
