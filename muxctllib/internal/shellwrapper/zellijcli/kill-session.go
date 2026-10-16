package zellijcli

import "context"

// Kill a specific session
//
// Usage: zellij kill-session [TARGET_SESSION]
//
// Arguments:
//   [TARGET_SESSION]  Name of target session
//
// Options:
//   -h, --help  Print help

// KillSession runs kill-session. An empty targetSession is omitted.
func (s *ShellOut) KillSession(ctx context.Context, targetSession string) error {
	args := argv{"kill-session"}
	if targetSession != "" {
		args = append(args, "--", targetSession)
	}
	_, err := s.run(ctx, args...)
	return err
}
