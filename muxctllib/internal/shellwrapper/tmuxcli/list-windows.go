package tmuxcli

import "context"

// list-windows [-ar] [-F format] [-f filter] [-O sort-order] [-t target-session]
//               (alias: lsw)
//         If -a is given, list all windows on the server.  Otherwise, list windows in
//         the current session or in target-session.  -F specifies the format of each
//         line and -f a filter.  Only windows for which the filter is true are shown.
//         See the “FORMATS” section.  -O specifies the sort order: one of ‘index’,
//         ‘name’, ‘size’ (area), ‘creation’ (time), or ‘activity’ (time).  -r reverses
//         the sort order.

type ListWindowsOption struct {
	ListOption
	// All is -a: list all windows on the server, ignoring TargetSession.
	All bool
	// TargetSession is -t target-session.
	TargetSession string
	// UserOptions are user option names such as "@foo" reported in [Window.UserOptions].
	UserOptions []string
}

type Window struct {
	Id          string // window_id, e.g. "@0"
	Index       int    // window_index
	Name        string // window_name
	SessionId   string // session_id
	SessionName string // session_name
	Panes       int    // window_panes: number of panes
	Active      bool   // window_active
	Zoomed      bool   // window_zoomed_flag
	Width       int    // window_width
	Height      int    // window_height
	// UserOptions maps each name in [ListWindowsOption.UserOptions] to its value, "" if unset.
	UserOptions map[string]string
}

// ListWindows runs list-windows.
// It fails with a [*CmdError] if no server is running.
func (s *ShellOut) ListWindows(ctx context.Context, opt ListWindowsOption) ([]Window, error) {
	userVars, err := userOptionVars(opt.UserOptions)
	if err != nil {
		return nil, err
	}
	vars := append([]string{
		"window_id",
		"window_index",
		"window_name",
		"session_id",
		"session_name",
		"window_panes",
		"window_active",
		"window_zoomed_flag",
		"window_width",
		"window_height",
	}, userVars...)
	args := opt.args(argv{"list-windows"}).value("-F", format(vars))
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	rows, err := records(out, len(vars))
	if err != nil {
		return nil, err
	}
	nBuiltin := len(vars) - len(userVars)
	windows := make([]Window, len(rows))
	for i, r := range rows {
		w := &windows[i]
		if err := parseFields(vars, r,
			&w.Id, &w.Index, &w.Name, &w.SessionId, &w.SessionName,
			&w.Panes, &w.Active, &w.Zoomed, &w.Width, &w.Height,
		); err != nil {
			return nil, err
		}
		w.UserOptions = userOptionMap(opt.UserOptions, r[nBuiltin:])
	}
	return windows, nil
}

func (o ListWindowsOption) args(a argv) argv {
	return o.ListOption.args(a).
		flag("-a", o.All).
		value("-t", o.TargetSession)
}
