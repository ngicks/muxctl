package tmuxcli

import "context"

// split-window [-bdefhIkPvZ] [-c start-directory] [-e environment] [-F format] [-l
//         size] [-m message] [-p percentage] [-R inactive-border-style] [-s style] [-S
//         active-border-style] [-t target-pane] [shell-command [argument ...]]
//               (alias: splitw)
//         Creates a new pane by splitting target-pane.  Shares behavior with new-pane.
//
//         See new-pane for more details.

// SplitWindowOption holds the flags of split-window.
// -P and -F are not exposed: the method owns them to return a [SpawnResult].
type SplitWindowOption struct {
	PaneSplitOption
}

// SplitWindow runs split-window.
func (s *ShellOut) SplitWindow(ctx context.Context, opt SplitWindowOption, shellCommand ...string) (SpawnResult, error) {
	args := opt.args(argv{"split-window"})
	return s.spawn(ctx, args, shellCommand)
}
