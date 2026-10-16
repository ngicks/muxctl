package tmuxcli

import (
	"context"
	"strings"
)

// display-message [-aCIlNpv] [-c target-client] [-d delay] [-t target-pane] [message]
//               (alias: display)
//         Display a message.  If -p is given, the output is printed to stdout,
//         otherwise it is displayed in the target-client status line for up to delay
//         milliseconds.  If delay is not given, the display-time option is used; a
//         delay of zero waits for a key press.  -N ignores key presses and closes only
//         after the delay expires.  If -C is given, the pane will continue to be
//         updated while the message is displayed.  If -l is given, message is printed
//         unchanged.  Otherwise, the format of message is described in the “FORMATS”
//         section; information is taken from target-pane if -t is given, otherwise the
//         active pane.
//
//         -v prints verbose logging as the format is parsed and -a lists the format
//         variables and their values.
//
//         -I forwards any input read from stdin to the empty pane given by
//         target-pane.

type DisplayMessageOption struct {
	// ListVariables is -a: list the format variables and their values.
	ListVariables bool
	// KeepUpdating is -C: keep updating the pane while the message is displayed.
	KeepUpdating bool
	// ForwardStdin is -I: forward stdin to the empty pane TargetPane.
	ForwardStdin bool
	// Literal is -l: print the message unchanged instead of expanding it as a format.
	Literal bool
	// IgnoreKeys is -N: ignore key presses and close only after Delay.
	IgnoreKeys bool
	// Print is -p: print to stdout instead of the status line.
	Print bool
	// Verbose is -v: log verbosely as the format is parsed.
	Verbose bool
	// TargetClient is -c target-client.
	TargetClient string
	// Delay is -d delay in milliseconds; zero waits for a key press.
	Delay *int
	// Format is -F format: the message given as a flag instead of an argument.
	// It is missing from the man page but listed by tmux list-commands.
	Format string
	// TargetPane is -t target-pane.
	TargetPane string
}

// DisplayMessage runs display-message and returns its stdout without the
// trailing newline. An empty message is omitted.
func (s *ShellOut) DisplayMessage(ctx context.Context, opt DisplayMessageOption, message string) (string, error) {
	args := opt.args(argv{"display-message"})
	if message != "" {
		args = args.trailing([]string{message})
	}
	out, err := s.run(ctx, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

func (o DisplayMessageOption) args(a argv) argv {
	return a.
		flag("-a", o.ListVariables).
		flag("-C", o.KeepUpdating).
		flag("-I", o.ForwardStdin).
		flag("-l", o.Literal).
		flag("-N", o.IgnoreKeys).
		flag("-p", o.Print).
		flag("-v", o.Verbose).
		value("-c", o.TargetClient).
		int("-d", o.Delay).
		value("-F", o.Format).
		value("-t", o.TargetPane)
}
