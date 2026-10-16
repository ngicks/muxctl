package zellijcli

import "context"

// Kill all sessions
//
// Usage: zellij kill-all-sessions [OPTIONS]
//
// Options:
//   -y, --yes   Automatic yes to prompts
//   -h, --help  Print help

type KillAllSessionsOption struct {
	// Yes is --yes: answers yes to the confirmation prompt.
	Yes bool
}

// KillAllSessions runs kill-all-sessions.
func (s *ShellOut) KillAllSessions(ctx context.Context, opt KillAllSessionsOption) error {
	args := opt.args(argv{"kill-all-sessions"})
	_, err := s.run(ctx, args...)
	return err
}

func (o KillAllSessionsOption) args(a argv) argv {
	return a.flag("--yes", o.Yes)
}
