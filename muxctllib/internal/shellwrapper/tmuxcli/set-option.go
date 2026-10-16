package tmuxcli

import "context"

// set-option [-aFgopqsuUw] [-t target-pane] option [value]
//               (alias: set)
//         Set a pane option with -p, a window option with -w, a server option with -s,
//         otherwise a session option.  If the option is not a user option, -w or -s
//         may be unnecessary - tmux will infer the scope from the option name,
//         assuming -w for pane options.  If -g is given, the global session or window
//         option is set.
//
//         -F expands formats in the option value.  The -u flag unsets an option, so a
//         session inherits the option from the global options (or with -g, restores a
//         global option to the default).  -U unsets an option (like -u) but if the
//         option is a pane option also unsets the option on any panes in the window.
//         value depends on the option and its type and can be omitted for flag and
//         choice options to toggle its value (choice options toggle between the first
//         two choices).
//
//         The -o flag prevents setting an option that is already set and the -q flag
//         suppresses errors about unknown or ambiguous options.
//
//         With -a, and if the option expects a string or a style, value is appended to
//         the existing setting.  For example:
//
//               set -g status-left "foo"
//               set -ag status-left "bar"
//
//         Will result in ‘foobar’.  And:
//
//               set -g status-style "bg=red"
//               set -ag status-style "fg=blue"
//
//         Will result in a red background and blue foreground.  Without -a, the result
//         would be the default background and a blue foreground.

// OptionScope selects the options set-option and show-options act on.
// With none set, tmux uses session options or infers the scope from the option name.
type OptionScope struct {
	// Pane is -p: pane options.
	Pane bool
	// Window is -w: window options.
	Window bool
	// Server is -s: server options.
	Server bool
	// Global is -g: global session or window options.
	Global bool
}

type SetOptionOption struct {
	OptionScope
	// Append is -a: append value to a string or style option.
	Append bool
	// ExpandFormats is -F: expand formats in value.
	ExpandFormats bool
	// OnlyIfUnset is -o: do not set an option that is already set.
	OnlyIfUnset bool
	// Quiet is -q: suppress errors about unknown or ambiguous options.
	Quiet bool
	// Unset is -u: unset the option so it is inherited again.
	Unset bool
	// UnsetPanes is -U: like Unset, also unsetting a pane option on every pane in the window.
	UnsetPanes bool
	// TargetPane is -t target-pane: the pane, window or session, depending on the scope.
	TargetPane string
}

// SetOption runs set-option.
// A nil value is omitted, to toggle a flag or choice option or to unset one.
func (s *ShellOut) SetOption(ctx context.Context, opt SetOptionOption, option string, value *string) error {
	args := opt.args(argv{"set-option"})
	args = append(args, option)
	if value != nil {
		// Not argv.value: an empty value is meaningful here.
		args = append(args, *value)
	}
	_, err := s.run(ctx, args...)
	return err
}

func (o OptionScope) args(a argv) argv {
	return a.
		flag("-p", o.Pane).
		flag("-w", o.Window).
		flag("-s", o.Server).
		flag("-g", o.Global)
}

func (o SetOptionOption) args(a argv) argv {
	return o.OptionScope.args(a).
		flag("-a", o.Append).
		flag("-F", o.ExpandFormats).
		flag("-o", o.OnlyIfUnset).
		flag("-q", o.Quiet).
		flag("-u", o.Unset).
		flag("-U", o.UnsetPanes).
		value("-t", o.TargetPane)
}
