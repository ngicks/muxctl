package zellijcli

import "context"

// Dumps the viewport and optionally scrollback of a pane to a file or STDOUT
//
// Usage: zellij action dump-screen [OPTIONS]
//
// Options:
//       --path <PATH>        File path to dump the pane content to. If omitted, prints to STDOUT
//   -f, --full               Dump the pane with full scrollback
//   -p, --pane-id <PANE_ID>  The pane_id of the pane, eg. terminal_1, plugin_2 or 3 (equivalent to
//                            terminal_3). If not specified, dumps the focused pane
//   -a, --ansi               Preserve ANSI styling in the dump output
//   -h, --help               Print help

type ActionDumpScreenOption struct {
	// Path is --path. When set, zellij writes the dump there and prints nothing.
	Path string
	// Full is --full.
	Full bool
	// PaneId is --pane-id. Nil dumps the focused pane.
	PaneId *PaneId
	// Ansi is --ansi.
	Ansi bool
}

// ActionDumpScreen runs action dump-screen and returns its stdout as is.
// zellij reports success with an empty dump if the pane does not exist.
func (s *ShellOut) ActionDumpScreen(ctx context.Context, opt ActionDumpScreenOption) (string, error) {
	args := opt.args(argv{"action", "dump-screen"})
	return s.run(ctx, args...)
}

func (o ActionDumpScreenOption) args(a argv) argv {
	return a.
		value("--path", o.Path).
		flag("--full", o.Full).
		paneId("--pane-id", o.PaneId).
		flag("--ansi", o.Ansi)
}
