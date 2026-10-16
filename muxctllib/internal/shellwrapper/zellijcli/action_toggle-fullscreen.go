package zellijcli

import "context"

// Toggle between fullscreen focus pane and normal layout
//
// Usage: zellij action toggle-fullscreen [OPTIONS]
//
// Options:
//   -p, --pane-id <PANE_ID>  Target a specific pane by ID (eg. terminal_1, plugin_2, or 3)
//   -h, --help               Print help

type ActionToggleFullscreenOption struct {
	// PaneId is --pane-id. Nil toggles the focused pane.
	PaneId *PaneId
}

// ActionToggleFullscreen runs action toggle-fullscreen.
func (s *ShellOut) ActionToggleFullscreen(ctx context.Context, opt ActionToggleFullscreenOption) error {
	args := opt.args(argv{"action", "toggle-fullscreen"})
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionToggleFullscreenOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
