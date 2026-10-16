package tmuxcli

import (
	"context"
	"errors"
)

// has-session [-t target-session]
//               (alias: has)
//         Report an error and exit with 1 if the specified session does not exist.  If
//         it does exist, exit with 0.

type HasSessionOption struct {
	// TargetSession is -t target-session.
	TargetSession string
}

// HasSession runs has-session and reports whether the session exists.
// No running server also reports false, since tmux exits with 1 for both.
func (s *ShellOut) HasSession(ctx context.Context, opt HasSessionOption) (bool, error) {
	args := opt.args(argv{"has-session"})
	_, err := s.run(ctx, args...)
	if err == nil {
		return true, nil
	}
	if cmdErr, ok := errors.AsType[*CmdError](err); ok && cmdErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}

func (o HasSessionOption) args(a argv) argv {
	return a.value("-t", o.TargetSession)
}
