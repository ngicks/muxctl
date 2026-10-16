package zellijcli

import "context"

// Focus a specific pane by its ID
//
// Usage: zellij action focus-pane-id <PANE_ID>
//
// Arguments:
//   <PANE_ID>  The pane_id of the pane, eg. terminal_1, plugin_2 or 3
//
// Options:
//   -h, --help  Print help

// ActionFocusPaneId runs action focus-pane-id.
func (s *ShellOut) ActionFocusPaneId(ctx context.Context, id PaneId) error {
	_, err := s.run(ctx, "action", "focus-pane-id", id.String())
	return err
}
