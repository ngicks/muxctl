package zellijcli

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// List active sessions
//
// Usage: zellij list-sessions [OPTIONS]
//
// Options:
//   -n, --no-formatting  Do not add colors and formatting to the list (useful for parsing)
//   -s, --short          Print just the session name
//   -r, --reverse        List the sessions in reverse order (default is ascending order)
//   -h, --help           Print help

// ListSessionsOption holds the flags of list-sessions.
// --no-formatting is not exposed: ListSessions always passes it to parse the output.
type ListSessionsOption struct {
	// Short is --short. Only names are printed, so the other fields of [Session] stay zero.
	Short bool
	// Reverse is --reverse.
	Reverse bool
}

type Session struct {
	Name string
	// Created is the age zellij prints, e.g. "1m 22s".
	Created string
	// Current reports the session named by $ZELLIJ_SESSION_NAME.
	Current bool
	// Exited reports a dead session that attaching resurrects.
	Exited bool
}

// ListSessions runs list-sessions.
// It returns no sessions and a nil error when there are none; zellij itself exits 1 then.
func (s *ShellOut) ListSessions(ctx context.Context, opt ListSessionsOption) ([]Session, error) {
	args := opt.args(argv{"list-sessions", "--no-formatting"})
	out, err := s.run(ctx, args...)
	if err != nil {
		if cmdErr, ok := errors.AsType[*CmdError](err); ok &&
			cmdErr.ExitCode() == 1 &&
			strings.Contains(cmdErr.Stderr, "No active zellij sessions found") {
			return nil, nil
		}
		return nil, err
	}
	return parseSessions(out, opt.Short)
}

func (o ListSessionsOption) args(a argv) argv {
	return a.
		flag("--short", o.Short).
		flag("--reverse", o.Reverse)
}

// parseSessions parses lines like "name [Created 1m 22s ago] (current)".
func parseSessions(out string, short bool) ([]Session, error) {
	var sessions []Session
	for line := range strings.Lines(out) {
		line = strings.TrimSuffix(line, "\n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		if short {
			sessions = append(sessions, Session{Name: line})
			continue
		}
		// The last marker is searched because a session name may contain it, a duration may not.
		name, rest, ok := strings.CutLast(line, " [Created ")
		if !ok {
			return nil, fmt.Errorf("zellij: unexpected list-sessions line %q", line)
		}
		created, suffix, ok := strings.Cut(rest, " ago]")
		if !ok {
			return nil, fmt.Errorf("zellij: unexpected list-sessions line %q", line)
		}
		ss := Session{Name: name, Created: created}
		switch strings.TrimSpace(suffix) {
		case "":
		case "(current)":
			ss.Current = true
		case "(EXITED - attach to resurrect)":
			ss.Exited = true
		default:
			return nil, fmt.Errorf("zellij: unexpected list-sessions line %q", line)
		}
		sessions = append(sessions, ss)
	}
	return sessions, nil
}
