package tmuxcli

import "context"

// switch-client [-ElnprZ] [-c target-client] [-O sort-order] [-t target-session] [-T
//         key-table]
//               (alias: switchc)
//         Switch the current session for client target-client to target-session.  As a
//         special case, -t may refer to a pane (a target that contains ‘:’, ‘.’ or
//         ‘%’), to change session, window and pane.  In that case, -Z keeps the window
//         zoomed if it was zoomed.  If -l, -n or -p is used, the client is moved to
//         the last, next or previous session respectively.  -O may be used with -n and
//         -p to specify the field to sort on: one of ‘name’, ‘size’, ‘creation’
//         (time), or ‘activity’ (time).  -r toggles the client read-only and
//         ignore-size flags (see the attach-session command).
//
//         If -E is used, update-environment option will not be applied.
//
//         -T sets the client's key table; the next key will be looked up using
//         key-table.  After that key, the client is returned to its default key table
//         (normally root).  This may be used to configure multiple prefix keys, or to
//         bind commands to sequences of keys.  For example, to make typing ‘abc’ run
//         the list-keys command:
//
//               bind-key -Ttable2 c list-keys
//               bind-key -Ttable1 b switch-client -Ttable2
//               bind-key -Troot   a switch-client -Ttable1

type SwitchClientOption struct {
	// NoUpdateEnvironment is -E: do not apply the update-environment option.
	NoUpdateEnvironment bool
	// Last is -l: move to the last session.
	Last bool
	// Next is -n: move to the next session.
	Next bool
	// Previous is -p: move to the previous session.
	Previous bool
	// ToggleReadOnly is -r: toggle the client's read-only and ignore-size flags.
	ToggleReadOnly bool
	// KeepZoom is -Z: keep the window zoomed when TargetSession is a pane.
	KeepZoom bool
	// TargetClient is -c target-client.
	TargetClient string
	// SortOrder is -O sort-order: the order of sessions for Next and Previous.
	SortOrder string
	// TargetSession is -t target-session, which may also be a pane.
	TargetSession string
	// KeyTable is -T key-table.
	KeyTable string
}

// SwitchClient runs switch-client.
func (s *ShellOut) SwitchClient(ctx context.Context, opt SwitchClientOption) error {
	args := opt.args(argv{"switch-client"})
	_, err := s.run(ctx, args...)
	return err
}

func (o SwitchClientOption) args(a argv) argv {
	return a.
		flag("-E", o.NoUpdateEnvironment).
		flag("-l", o.Last).
		flag("-n", o.Next).
		flag("-p", o.Previous).
		flag("-r", o.ToggleReadOnly).
		flag("-Z", o.KeepZoom).
		value("-c", o.TargetClient).
		value("-O", o.SortOrder).
		value("-t", o.TargetSession).
		value("-T", o.KeyTable)
}
