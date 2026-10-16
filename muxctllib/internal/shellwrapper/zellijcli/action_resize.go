package zellijcli

import "context"

// [increase|decrease] the focused panes area at the [left|down|up|right] border
//
// Usage: zellij action resize [OPTIONS] <RESIZE> [DIRECTION]
//
// Arguments:
//   <RESIZE>
//   [DIRECTION]
//
// Options:
//   -p, --pane-id <PANE_ID>  Target a specific pane by ID (eg. terminal_1, plugin_2, or 3)
//   -h, --help               Print help

type ActionResizeOption struct {
	// PaneId is --pane-id. Nil resizes the focused pane.
	PaneId *PaneId
}

// ActionResize runs action resize.
// resize is "increase" or "decrease"; direction is "left", "down", "up", "right", or empty for all borders.
func (s *ShellOut) ActionResize(ctx context.Context, opt ActionResizeOption, resize string, direction string) error {
	args := opt.args(argv{"action", "resize"})
	args = append(args, "--", resize)
	if direction != "" {
		args = append(args, direction)
	}
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionResizeOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
