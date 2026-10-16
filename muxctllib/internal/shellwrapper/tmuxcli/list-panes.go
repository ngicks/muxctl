package tmuxcli

import "context"

// list-panes [-ars] [-F format] [-f filter] [-O sort-order] [-t target]
//               (alias: lsp)
//         If -a is given, target is ignored and all panes on the server are listed.
//         If -s is given, target is a session (or the current session).  If neither is
//         given, target is a window (or the current window).  -F specifies the format
//         of each line and -f a filter.  Only panes for which the filter is true are
//         shown.  See the “FORMATS” section.  -O specifies the sort order: one of
//         ‘name’ (title), ‘index’, ‘size’ (area), ‘creation’ (time), or ‘activity’
//         (time).  -r reverses the sort order.

type ListPanesOption struct {
	ListOption
	// All is -a: list all panes on the server, ignoring Target.
	All bool
	// Session is -s: Target is a session instead of a window.
	Session bool
	// Target is -t target: a window, or a session with Session.
	Target string
	// UserOptions are user option names such as "@foo" reported in [Pane.UserOptions].
	UserOptions []string
}

type Pane struct {
	Id             string // pane_id, e.g. "%0"
	Index          int    // pane_index
	SessionName    string // session_name
	WindowId       string // window_id
	WindowIndex    int    // window_index
	Title          string // pane_title
	Active         bool   // pane_active
	Dead           bool   // pane_dead
	DeadStatus     int    // pane_dead_status: exit status of the command in a dead pane
	Floating       bool   // pane_floating_flag
	Width          int    // pane_width
	Height         int    // pane_height
	Left           int    // pane_left
	Top            int    // pane_top
	Pid            int    // pane_pid
	CurrentCommand string // pane_current_command
	CurrentPath    string // pane_current_path
	// UserOptions maps each name in [ListPanesOption.UserOptions] to its value, "" if unset.
	UserOptions map[string]string
}

// ListPanes runs list-panes.
// It fails with a [*CmdError] if no server is running.
func (s *ShellOut) ListPanes(ctx context.Context, opt ListPanesOption) ([]Pane, error) {
	userVars, err := userOptionVars(opt.UserOptions)
	if err != nil {
		return nil, err
	}
	vars := append([]string{
		"pane_id",
		"pane_index",
		"session_name",
		"window_id",
		"window_index",
		"pane_title",
		"pane_active",
		"pane_dead",
		"pane_dead_status",
		"pane_floating_flag",
		"pane_width",
		"pane_height",
		"pane_left",
		"pane_top",
		"pane_pid",
		"pane_current_command",
		"pane_current_path",
	}, userVars...)
	args := opt.args(argv{"list-panes"}).value("-F", format(vars))
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	rows, err := records(out, len(vars))
	if err != nil {
		return nil, err
	}
	nBuiltin := len(vars) - len(userVars)
	panes := make([]Pane, len(rows))
	for i, r := range rows {
		p := &panes[i]
		if err := parseFields(vars, r,
			&p.Id, &p.Index, &p.SessionName, &p.WindowId, &p.WindowIndex, &p.Title,
			&p.Active, &p.Dead, &p.DeadStatus, &p.Floating,
			&p.Width, &p.Height, &p.Left, &p.Top, &p.Pid, &p.CurrentCommand, &p.CurrentPath,
		); err != nil {
			return nil, err
		}
		p.UserOptions = userOptionMap(opt.UserOptions, r[nBuiltin:])
	}
	return panes, nil
}

func (o ListPanesOption) args(a argv) argv {
	return o.ListOption.args(a).
		flag("-a", o.All).
		flag("-s", o.Session).
		value("-t", o.Target)
}
