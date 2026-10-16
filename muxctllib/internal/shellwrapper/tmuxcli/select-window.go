package tmuxcli

import "context"

// select-window [-lnpT] [-t target-window]
//               (alias: selectw)
//         Select the window at target-window.  -l, -n and -p are equivalent to the
//         last-window, next-window and previous-window commands.  If -T is given and
//         the selected window is already the current window, the command behaves like
//         last-window.

type SelectWindowOption struct {
	// Last is -l: same as the last-window command.
	Last bool
	// Next is -n: same as the next-window command.
	Next bool
	// Previous is -p: same as the previous-window command.
	Previous bool
	// Toggle is -T: behave like last-window if TargetWindow is already current.
	Toggle bool
	// TargetWindow is -t target-window.
	TargetWindow string
}

// SelectWindow runs select-window.
func (s *ShellOut) SelectWindow(ctx context.Context, opt SelectWindowOption) error {
	args := opt.args(argv{"select-window"})
	_, err := s.run(ctx, args...)
	return err
}

func (o SelectWindowOption) args(a argv) argv {
	return a.
		flag("-l", o.Last).
		flag("-n", o.Next).
		flag("-p", o.Previous).
		flag("-T", o.Toggle).
		value("-t", o.TargetWindow)
}
