package zellijcli

import "context"

// Open a new pane in the specified direction [right|down] If no direction is specified, will try to
// use the biggest available space. Returns: Created pane ID (format: terminal_<id> or plugin_<id>)
//
// Usage: zellij action new-pane [OPTIONS] [-- <COMMAND>...]
//
// Arguments:
//   [COMMAND]...
//
// Options:
//   -d, --direction <DIRECTION>          Direction to open the new pane in
//   -p, --plugin <PLUGIN>
//       --cwd <CWD>                      Change the working directory of the new pane
//   -f, --floating                       Open the new pane in floating mode
//   -i, --in-place                       Open the new pane in place of the current pane, temporarily
//                                        suspending it
//       --close-replaced-pane            Close the replaced pane instead of suspending it (only
//                                        effective with --in-place)
//       --pane-id <PANE_ID>              The pane to replace when opening in place, eg. terminal_1,
//                                        plugin_2 or 3 (only effective with --in-place; defaults to
//                                        the focused pane)
//   -n, --name <NAME>                    Name of the new pane
//   -c, --close-on-exit                  Close the pane immediately when its command exits
//   -s, --start-suspended                Start the command suspended, only running it after the you
//                                        first press ENTER
//       --configuration <CONFIGURATION>
//       --skip-plugin-cache
//   -x, --x <X>                          The x coordinates if the pane is floating as a bare integer
//                                        (eg. 1) or percent (eg. 10%)
//   -y, --y <Y>                          The y coordinates if the pane is floating as a bare integer
//                                        (eg. 1) or percent (eg. 10%)
//       --width <WIDTH>                  The width if the pane is floating as a bare integer (eg. 1)
//                                        or percent (eg. 10%)
//       --height <HEIGHT>                The height if the pane is floating as a bare integer (eg. 1)
//                                        or percent (eg. 10%)
//       --pinned <PINNED>                Whether to pin a floating pane so that it is always on top
//                                        [possible values: true, false]
//       --stacked
//   -b, --blocking                       Block until the command has finished and its pane has been
//                                        closed
//       --block-until-exit-success       Block until the command exits successfully (exit status 0) OR
//                                        its pane has been closed
//       --block-until-exit-failure       Block until the command exits with failure (non-zero exit
//                                        status) OR its pane has been closed
//       --block-until-exit               Block until the command exits (regardless of exit status) OR
//                                        its pane has been closed
//       --near-current-pane              if set, will open the pane near the current one rather than
//                                        following the user's focus
//       --no-focus                       if set, will open the pane without changing the focus of any
//                                        client, placing it relative to the pane the command was
//                                        issued from
//       --borderless <BORDERLESS>        start this pane without a border (warning: will make it
//                                        impossible to move with the mouse) [possible values: true,
//                                        false]
//       --tab-id <TAB_ID>                Target a specific tab by ID
//   -h, --help                           Print help

type ActionNewPaneOption struct {
	CommandPaneOption
	// Plugin is --plugin: a plugin URL or alias, opened instead of a command.
	Plugin string
	// PaneId is --pane-id: the pane InPlace replaces instead of the focused one.
	PaneId *PaneId
	// Configuration is --configuration.
	Configuration string
	// SkipPluginCache is --skip-plugin-cache.
	SkipPluginCache bool
}

// ActionNewPane runs action new-pane.
// The pane is nil when zellij prints no pane ID. zellij 0.45.1 does so for a Plugin
// that is already running, and usually for InPlace with CloseReplacedPane,
// where printing races with closing the replaced pane.
func (s *ShellOut) ActionNewPane(ctx context.Context, opt ActionNewPaneOption, command ...string) (*PaneId, error) {
	args := opt.args(argv{"action", "new-pane"}).trailing(command)
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return parseCreatedPaneId(out)
}

func (o ActionNewPaneOption) args(a argv) argv {
	return o.CommandPaneOption.args(a).
		value("--plugin", o.Plugin).
		paneId("--pane-id", o.PaneId).
		value("--configuration", o.Configuration).
		flag("--skip-plugin-cache", o.SkipPluginCache)
}
