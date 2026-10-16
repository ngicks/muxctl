package zellijcli

import (
	"context"
	"strings"
)

// Run a command in a new pane Returns: Created pane ID (format: terminal_<id>)
//
// Usage: zellij run [OPTIONS] -- <COMMAND>...
//
// Arguments:
//   <COMMAND>...  Command to run
//
// Options:
//   -d, --direction <DIRECTION>     Direction to open the new pane in
//       --cwd <CWD>                 Change the working directory of the new pane
//   -f, --floating                  Open the new pane in floating mode
//   -i, --in-place                  Open the new pane in place of the current pane, temporarily
//                                   suspending it
//       --close-replaced-pane       Close the replaced pane instead of suspending it (only effective
//                                   with --in-place)
//   -n, --name <NAME>               Name of the new pane
//   -c, --close-on-exit             Close the pane immediately when its command exits
//   -s, --start-suspended           Start the command suspended, only running after you first presses
//                                   ENTER
//   -x, --x <X>                     The x coordinates if the pane is floating as a bare integer (eg.
//                                   1) or percent (eg. 10%)
//   -y, --y <Y>                     The y coordinates if the pane is floating as a bare integer (eg.
//                                   1) or percent (eg. 10%)
//       --width <WIDTH>             The width if the pane is floating as a bare integer (eg. 1) or
//                                   percent (eg. 10%)
//       --height <HEIGHT>           The height if the pane is floating as a bare integer (eg. 1) or
//                                   percent (eg. 10%)
//       --pinned <PINNED>           Whether to pin a floating pane so that it is always on top
//                                   [possible values: true, false]
//       --stacked
//       --blocking                  Block until the command has finished and its pane has been closed
//       --block-until-exit-success  Block until the command exits successfully (exit status 0) OR its
//                                   pane has been closed
//       --block-until-exit-failure  Block until the command exits with failure (non-zero exit status)
//                                   OR its pane has been closed
//       --block-until-exit          Block until the command exits (regardless of exit status) OR its
//                                   pane has been closed
//       --near-current-pane         if set, will open the pane near the current one rather than
//                                   following the user's focus
//       --no-focus                  if set, will open the pane without changing the focus of any
//                                   client, placing it relative to the pane the command was issued
//                                   from
//   -b, --borderless <BORDERLESS>   start this pane without a border (warning: will make it impossible
//                                   to move with the mouse) [possible values: true, false]
//       --tab-id <TAB_ID>           Target a specific tab by ID
//   -h, --help                      Print help

// CommandPaneOption holds the flags that run and action new-pane share.
type CommandPaneOption struct {
	// Direction is --direction: "right" or "down". Empty lets zellij pick the biggest free space.
	Direction string
	// Cwd is --cwd.
	Cwd string
	// Floating is --floating.
	Floating bool
	// InPlace is --in-place.
	InPlace bool
	// CloseReplacedPane is --close-replaced-pane. It only takes effect with InPlace.
	CloseReplacedPane bool
	// Name is --name.
	Name string
	// CloseOnExit is --close-on-exit.
	CloseOnExit bool
	// StartSuspended is --start-suspended.
	StartSuspended bool
	// X is --x. X, Y, Width and Height are bare integers ("1") or percentages ("10%").
	X string
	// Y is --y.
	Y string
	// Width is --width.
	Width string
	// Height is --height.
	Height string
	// Pinned is --pinned.
	Pinned *bool
	// Stacked is --stacked.
	Stacked bool
	// Blocking is --blocking. zellij does not exit until the pane is closed.
	Blocking bool
	// BlockUntilExitSuccess is --block-until-exit-success.
	BlockUntilExitSuccess bool
	// BlockUntilExitFailure is --block-until-exit-failure.
	BlockUntilExitFailure bool
	// BlockUntilExit is --block-until-exit.
	BlockUntilExit bool
	// NearCurrentPane is --near-current-pane. zellij rejects it together with TabId.
	NearCurrentPane bool
	// NoFocus is --no-focus.
	NoFocus bool
	// Borderless is --borderless.
	Borderless *bool
	// TabId is --tab-id.
	TabId *int
}

type RunOption struct {
	CommandPaneOption
}

// Run runs run.
// The pane is nil when zellij prints no pane ID; see [ShellOut.ActionNewPane].
func (s *ShellOut) Run(ctx context.Context, opt RunOption, command ...string) (*PaneId, error) {
	args := opt.args(argv{"run"})
	// Unlike new-pane, run requires the separator even though the command is mandatory.
	args = append(args, "--")
	args = append(args, command...)
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return parseCreatedPaneId(out)
}

func (o CommandPaneOption) args(a argv) argv {
	return a.
		value("--direction", o.Direction).
		value("--cwd", o.Cwd).
		flag("--floating", o.Floating).
		flag("--in-place", o.InPlace).
		flag("--close-replaced-pane", o.CloseReplacedPane).
		value("--name", o.Name).
		flag("--close-on-exit", o.CloseOnExit).
		flag("--start-suspended", o.StartSuspended).
		value("--x", o.X).
		value("--y", o.Y).
		value("--width", o.Width).
		value("--height", o.Height).
		bool("--pinned", o.Pinned).
		flag("--stacked", o.Stacked).
		flag("--blocking", o.Blocking).
		flag("--block-until-exit-success", o.BlockUntilExitSuccess).
		flag("--block-until-exit-failure", o.BlockUntilExitFailure).
		flag("--block-until-exit", o.BlockUntilExit).
		flag("--near-current-pane", o.NearCurrentPane).
		flag("--no-focus", o.NoFocus).
		bool("--borderless", o.Borderless).
		int("--tab-id", o.TabId)
}

// parseCreatedPaneId parses the pane ID that run and action new-pane print.
func parseCreatedPaneId(out string) (*PaneId, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, nil
	}
	id, err := ParsePaneId(fields[0])
	if err != nil {
		return nil, err
	}
	return &id, nil
}
