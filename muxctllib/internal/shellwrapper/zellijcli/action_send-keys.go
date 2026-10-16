package zellijcli

import "context"

// Send one or more keys to the terminal (e.g., "Ctrl a", "F1", "Alt Shift b")
//
// Usage: zellij action send-keys [OPTIONS] <KEYS>...
//
// Arguments:
//   <KEYS>...  Keys to send as space-separated strings
//
// Options:
//   -p, --pane-id <PANE_ID>  The pane_id of the pane, eg. terminal_1, plugin_2 or 3 (equivalent to
//                            terminal_3)
//   -h, --help               Print help

type ActionSendKeysOption struct {
	// PaneId is --pane-id. Nil sends to the focused pane.
	PaneId *PaneId
}

// ActionSendKeys runs action send-keys. Each key is one argument, e.g. "Ctrl a" or "Enter".
func (s *ShellOut) ActionSendKeys(ctx context.Context, opt ActionSendKeysOption, keys ...string) error {
	args := opt.args(argv{"action", "send-keys"})
	args = append(args, "--")
	args = append(args, keys...)
	_, err := s.run(ctx, args...)
	return err
}

func (o ActionSendKeysOption) args(a argv) argv {
	return a.paneId("--pane-id", o.PaneId)
}
