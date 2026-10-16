package tmuxcli

import "context"

// list-sessions [-r] [-F format] [-f filter] [-O sort-order]
//               (alias: ls)
//         List all sessions managed by the server.  -F specifies the format of each
//         line and -f a filter.  Only sessions for which the filter is true are shown.
//         See the “FORMATS” section.  -O specifies the sort order: one of ‘index’,
//         ‘name’, ‘creation’ (time), or ‘activity’ (time).  -r reverses the sort
//         order.

// ListOption holds the flags shared by the list-* commands.
// -F is not exposed: each list method owns its format to parse the output.
type ListOption struct {
	// Filter is -f filter. Only items for which the filter is true are shown.
	Filter string
	// SortOrder is -O sort-order. Accepted values depend on the command.
	SortOrder string
	// Reverse is -r: reverses the sort order.
	Reverse bool
}

type ListSessionsOption struct {
	ListOption
}

type Session struct {
	Id       string // session_id, e.g. "$0"
	Name     string // session_name
	Windows  int    // session_windows
	Attached int    // session_attached: number of attached clients
	Created  int    // session_created: Unix time
	Activity int    // session_activity: Unix time
}

// ListSessions runs list-sessions.
// It fails with a [*CmdError] if no server is running.
func (s *ShellOut) ListSessions(ctx context.Context, opt ListSessionsOption) ([]Session, error) {
	vars := []string{
		"session_id",
		"session_name",
		"session_windows",
		"session_attached",
		"session_created",
		"session_activity",
	}
	args := opt.args(argv{"list-sessions"}).value("-F", format(vars))
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	rows, err := records(out, len(vars))
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, len(rows))
	for i, r := range rows {
		ss := &sessions[i]
		if err := parseFields(vars, r, &ss.Id, &ss.Name, &ss.Windows, &ss.Attached, &ss.Created, &ss.Activity); err != nil {
			return nil, err
		}
	}
	return sessions, nil
}

func (o ListOption) args(a argv) argv {
	return a.
		flag("-r", o.Reverse).
		value("-f", o.Filter).
		value("-O", o.SortOrder)
}
