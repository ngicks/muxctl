package tmuxcli

import "context"

// list-clients [-r] [-F format] [-f filter] [-O sort-order] [-t target-session]
//               (alias: lsc)
//         List all clients attached to the server.  -F specifies the format of each
//         line and -f a filter.  Only clients for which the filter is true are shown.
//         See the “FORMATS” section.  -O specifies the sort order: one of ‘name’,
//         ‘size’, ‘creation’ (time), or ‘activity’ (time).  -r reverses the sort
//         order.  If target-session is specified, list only clients connected to that
//         session.

type ListClientsOption struct {
	ListOption
	// TargetSession is -t target-session: list only clients attached to it.
	TargetSession string
}

type Client struct {
	Name        string // client_name
	Tty         string // client_tty
	Pid         int    // client_pid
	SessionName string // client_session
	Activity    int    // client_activity: Unix time
	Created     int    // client_created: Unix time
	Width       int    // client_width
	Height      int    // client_height
	Termname    string // client_termname
	ReadOnly    bool   // client_readonly
}

// ListClients runs list-clients.
// It fails with a [*CmdError] if no server is running.
func (s *ShellOut) ListClients(ctx context.Context, opt ListClientsOption) ([]Client, error) {
	vars := []string{
		"client_name",
		"client_tty",
		"client_pid",
		"client_session",
		"client_activity",
		"client_created",
		"client_width",
		"client_height",
		"client_termname",
		"client_readonly",
	}
	args := opt.args(argv{"list-clients"}).value("-F", format(vars))
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	rows, err := records(out, len(vars))
	if err != nil {
		return nil, err
	}
	clients := make([]Client, len(rows))
	for i, r := range rows {
		c := &clients[i]
		if err := parseFields(vars, r, &c.Name, &c.Tty, &c.Pid, &c.SessionName, &c.Activity, &c.Created, &c.Width, &c.Height, &c.Termname, &c.ReadOnly); err != nil {
			return nil, err
		}
	}
	return clients, nil
}

func (o ListClientsOption) args(a argv) argv {
	return o.ListOption.args(a).value("-t", o.TargetSession)
}
