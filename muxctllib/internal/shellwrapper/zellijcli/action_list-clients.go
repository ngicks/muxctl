package zellijcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Usage: zellij action list-clients
//
// Options:
//   -h, --help  Print help

type Client struct {
	Id int
	// PaneId is the focused pane as zellij prints it, e.g. "terminal_1".
	// It is kept unparsed since the column is free text.
	PaneId         string
	RunningCommand string
}

// ActionListClients runs action list-clients.
func (s *ShellOut) ActionListClients(ctx context.Context) ([]Client, error) {
	out, err := s.run(ctx, "action", "list-clients")
	if err != nil {
		return nil, err
	}
	return parseClients(out)
}

// parseClients parses the table under the header "CLIENT_ID ZELLIJ_PANE_ID RUNNING_COMMAND".
// Columns are padded with spaces and the running command may contain spaces itself.
func parseClients(out string) ([]Client, error) {
	var clients []Client
	header := true
	for line := range strings.Lines(out) {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if header {
			header = false
			continue
		}
		id, rest := cutField(line)
		paneId, rest := cutField(rest)
		n, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("zellij: parse list-clients line %q: %w", line, err)
		}
		clients = append(clients, Client{Id: n, PaneId: paneId, RunningCommand: strings.TrimSpace(rest)})
	}
	return clients, nil
}

func cutField(s string) (field, rest string) {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if i := strings.IndexFunc(s, unicode.IsSpace); i >= 0 {
		return s[:i], s[i:]
	}
	return s, ""
}
