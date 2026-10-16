package tmuxcli

import "context"

// send-keys [-FHKlMRX] [-c target-client] [-N repeat-count] [-t target-pane] [key ...]
//               (alias: send)
//         Send a key or keys to a window or client.  Each argument key is the name of
//         the key (such as ‘C-a’ or ‘NPage’) to send; if the string is not recognised
//         as a key, it is sent as a series of characters.  If -K is given, keys are
//         sent to target-client, so they are looked up in the client's key table,
//         rather than to target-pane.  All arguments are sent sequentially from first
//         to last.  If no keys are given and the command is bound to a key, then that
//         key is used.
//
//         The -l flag disables key name lookup and processes the keys as literal UTF-8
//         characters.  The -H flag expects each key to be a hexadecimal number for an
//         ASCII character.
//
//         The -R flag causes the terminal state to be reset.
//
//         -M passes through a mouse event (only valid if bound to a mouse key binding,
//         see “MOUSE SUPPORT”).
//
//         -X is used to send a command into copy mode - see the “WINDOWS AND PANES”
//         section.  -N specifies a repeat count and -F expands formats in arguments
//         where appropriate.

type SendKeysOption struct {
	// ExpandFormats is -F: expand formats in the keys.
	ExpandFormats bool
	// Hex is -H: each key is the hexadecimal number of an ASCII character.
	Hex bool
	// ToClient is -K: send the keys to TargetClient's key table instead of to TargetPane.
	ToClient bool
	// Literal is -l: send the keys as literal UTF-8 characters, without key name lookup.
	Literal bool
	// Mouse is -M: pass through a mouse event.
	Mouse bool
	// Reset is -R: reset the terminal state.
	Reset bool
	// CopyModeCommand is -X: send a command into copy mode.
	CopyModeCommand bool
	// TargetClient is -c target-client.
	TargetClient string
	// RepeatCount is -N repeat-count.
	RepeatCount *int
	// TargetPane is -t target-pane.
	TargetPane string
}

// SendKeys runs send-keys.
func (s *ShellOut) SendKeys(ctx context.Context, opt SendKeysOption, keys ...string) error {
	args := opt.args(argv{"send-keys"}).trailing(keys)
	_, err := s.run(ctx, args...)
	return err
}

func (o SendKeysOption) args(a argv) argv {
	return a.
		flag("-F", o.ExpandFormats).
		flag("-H", o.Hex).
		flag("-K", o.ToClient).
		flag("-l", o.Literal).
		flag("-M", o.Mouse).
		flag("-R", o.Reset).
		flag("-X", o.CopyModeCommand).
		value("-c", o.TargetClient).
		int("-N", o.RepeatCount).
		value("-t", o.TargetPane)
}
