package zellijcli

import "context"

// Delete a specific session
//
// Usage: zellij delete-session [OPTIONS] [TARGET_SESSION]
//
// Arguments:
//   [TARGET_SESSION]  Name of target session
//
// Options:
//   -f, --force  Kill the session if it's running before deleting it
//   -h, --help   Print help

type DeleteSessionOption struct {
	// Force is --force: a running session is killed, then deleted.
	Force bool
}

// DeleteSession runs delete-session. An empty targetSession is omitted.
func (s *ShellOut) DeleteSession(ctx context.Context, opt DeleteSessionOption, targetSession string) error {
	args := opt.args(argv{"delete-session"})
	if targetSession != "" {
		args = append(args, "--", targetSession)
	}
	_, err := s.run(ctx, args...)
	return err
}

func (o DeleteSessionOption) args(a argv) argv {
	return a.flag("--force", o.Force)
}
