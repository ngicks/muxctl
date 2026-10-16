package zellijcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Create a new tab, optionally with a specified tab layout and name
//
// Returns: The created tab's ID as a single number on stdout
//
// Usage: zellij action new-tab [OPTIONS] [-- <INITIAL_COMMAND>...]
//
// Arguments:
//   [INITIAL_COMMAND]...
//           Optional initial command to run in the new tab
//
// Options:
//   -l, --layout <LAYOUT>
//           Layout to use for the new tab
//
//       --layout-string <LAYOUT_STRING>
//           Raw KDL layout string to use directly (instead of a layout file path)
//
//       --layout-dir <LAYOUT_DIR>
//           Default folder to look for layouts
//
//   -n, --name <NAME>
//           Name of the new tab
//
//   -c, --cwd <CWD>
//           Change the working directory of the new tab
//
//       --initial-plugin <INITIAL_PLUGIN>
//           Initial plugin to load in the new tab
//
//       --close-on-exit
//           Close the pane immediately when its command exits
//
//       --start-suspended
//           Start the command suspended, only running it after you first press ENTER
//
//       --block-until-exit-success
//           Block until the command exits successfully (exit status 0) OR its pane has been closed
//
//       --block-until-exit-failure
//           Block until the command exits with failure (non-zero exit status) OR its pane has been
//           closed
//
//       --block-until-exit
//           Block until the command exits (regardless of exit status) OR its pane has been closed
//
//       --no-focus
//           if set, will create the tab without changing the focus of any client
//
//   -h, --help
//           Print help (see a summary with '-h')

type ActionNewTabOption struct {
	// Layout is --layout.
	Layout string
	// LayoutString is --layout-string: a KDL layout.
	LayoutString string
	// LayoutDir is --layout-dir.
	LayoutDir string
	// Name is --name.
	Name string
	// Cwd is --cwd.
	Cwd string
	// InitialPlugin is --initial-plugin.
	InitialPlugin string
	// CloseOnExit is --close-on-exit.
	CloseOnExit bool
	// StartSuspended is --start-suspended.
	StartSuspended bool
	// BlockUntilExitSuccess is --block-until-exit-success.
	BlockUntilExitSuccess bool
	// BlockUntilExitFailure is --block-until-exit-failure.
	BlockUntilExitFailure bool
	// BlockUntilExit is --block-until-exit.
	BlockUntilExit bool
	// NoFocus is --no-focus.
	NoFocus bool
}

// ActionNewTab runs action new-tab and returns the ID of the created tab.
func (s *ShellOut) ActionNewTab(ctx context.Context, opt ActionNewTabOption, initialCommand ...string) (int, error) {
	args := opt.args(argv{"action", "new-tab"}).trailing(initialCommand)
	out, err := s.run(ctx, args...)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("zellij: parse tab id: %w", err)
	}
	return id, nil
}

func (o ActionNewTabOption) args(a argv) argv {
	return a.
		value("--layout", o.Layout).
		value("--layout-string", o.LayoutString).
		value("--layout-dir", o.LayoutDir).
		value("--name", o.Name).
		value("--cwd", o.Cwd).
		value("--initial-plugin", o.InitialPlugin).
		flag("--close-on-exit", o.CloseOnExit).
		flag("--start-suspended", o.StartSuspended).
		flag("--block-until-exit-success", o.BlockUntilExitSuccess).
		flag("--block-until-exit-failure", o.BlockUntilExitFailure).
		flag("--block-until-exit", o.BlockUntilExit).
		flag("--no-focus", o.NoFocus)
}
