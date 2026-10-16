package zellijcli

import "context"

// Switch to a different session
//
// Usage: zellij action switch-session [OPTIONS] <NAME>
//
// Arguments:
//   <NAME>  Name of the session to switch to
//
// Options:
//       --tab-position <TAB_POSITION>    Optional tab position to focus
//       --pane-id <PANE_ID>              Optional pane ID to focus (eg. "terminal_1" for terminal pane
//                                        with id 1, or "plugin_2" for plugin pane with id 2)
//   -l, --layout <LAYOUT>                Layout to apply when switching to the session (relative paths
//                                        start at layout-dir)
//       --layout-string <LAYOUT_STRING>  Raw KDL layout string to use directly
//       --layout-dir <LAYOUT_DIR>        Default folder to look for layouts
//   -c, --cwd <CWD>                      Change the working directory when switching
//   -h, --help                           Print help

type ActionSwitchSessionOption struct {
	// TabPosition is --tab-position.
	TabPosition *int
	// PaneId is --pane-id.
	PaneId *PaneId
	// Layout is --layout.
	Layout string
	// LayoutString is --layout-string: a KDL layout.
	LayoutString string
	// LayoutDir is --layout-dir.
	LayoutDir string
	// Cwd is --cwd.
	Cwd string
}

// ActionSwitchSession runs action switch-session.
func (s *ShellOut) ActionSwitchSession(ctx context.Context, opt ActionSwitchSessionOption, name string) error {
	args := opt.args(argv{"action", "switch-session"})
	args = append(args, "--", name)
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionSwitchSessionOption) args(a argv) argv {
	return a.
		int("--tab-position", o.TabPosition).
		paneId("--pane-id", o.PaneId).
		value("--layout", o.Layout).
		value("--layout-string", o.LayoutString).
		value("--layout-dir", o.LayoutDir).
		value("--cwd", o.Cwd)
}
