package tmuxcli

import (
	"context"
	"strings"
)

// show-options [-AgHpqsvw] [-t target-pane] [option]
//               (alias: show)
//         Show the pane options (or a single option if option is provided) with -p,
//         the window options with -w, the server options with -s, otherwise the
//         session options.  If the option is not a user option, -w or -s may be
//         unnecessary - tmux will infer the scope from the option name, assuming -w
//         for pane options.  Global session or window options are listed if -g is
//         used.  -v shows only the option value, not the name.  If -q is set, no error
//         will be returned if option is unset.  -H includes hooks (omitted by
//         default).  -A includes options inherited from a parent set of options, such
//         options are marked with an asterisk.

type ShowOptionsOption struct {
	OptionScope
	// IncludeInherited is -A: include options inherited from a parent set of options.
	IncludeInherited bool
	// IncludeHooks is -H: include hooks.
	IncludeHooks bool
	// Quiet is -q: no error if the option is unset.
	Quiet bool
	// ValueOnly is -v: show only values; every [OptionEntry.Name] is then the option argument.
	ValueOnly bool
	// TargetPane is -t target-pane: the pane, window or session, depending on the scope.
	TargetPane string
}

// OptionEntry is one line of show-options output.
type OptionEntry struct {
	// Name is the option name, with "[index]" for an array item.
	Name string
	// Value is the option value with tmux's quoting removed.
	Value string
	// Inherited reports that the option comes from a parent set of options,
	// which tmux marks with '*' after the name when IncludeInherited is set.
	Inherited bool
}

// ShowOptions runs show-options. An empty option shows all options.
func (s *ShellOut) ShowOptions(ctx context.Context, opt ShowOptionsOption, option string) ([]OptionEntry, error) {
	args := opt.args(argv{"show-options"})
	if option != "" {
		args = append(args, option)
	}
	out, err := s.run(ctx, args...)
	if err != nil {
		return nil, err
	}
	return parseOptionEntries(out, option, opt.ValueOnly), nil
}

func (o ShowOptionsOption) args(a argv) argv {
	return o.OptionScope.args(a).
		flag("-A", o.IncludeInherited).
		flag("-H", o.IncludeHooks).
		flag("-q", o.Quiet).
		flag("-v", o.ValueOnly).
		value("-t", o.TargetPane)
}

func parseOptionEntries(out, option string, valueOnly bool) []OptionEntry {
	var entries []OptionEntry
	for line := range strings.Lines(out) {
		line = strings.TrimSuffix(line, "\n")
		if valueOnly {
			// -v prints values raw, so an empty line is an empty value.
			entries = append(entries, OptionEntry{Name: option, Value: line})
			continue
		}
		if line == "" {
			continue
		}
		name, value, _ := strings.Cut(line, " ")
		name, inherited := strings.CutSuffix(name, "*")
		entries = append(entries, OptionEntry{Name: name, Value: unquoteOptionValue(value), Inherited: inherited})
	}
	return entries
}

// unquoteOptionValue reverses the quoting show-options applies to a value:
// tmux escapes it vis(3)-style (\t, \\, \" and octal \ooo among others),
// then wraps it in single or double quotes if it contains special characters.
// Escapes are decoded inside single quotes too, since tmux escapes before quoting.
func unquoteOptionValue(v string) string {
	if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
		v = v[1 : len(v)-1]
	}
	if !strings.Contains(v, `\`) {
		return v
	}
	var b strings.Builder
	for i := 0; i < len(v); i++ {
		if v[i] != '\\' || i+1 == len(v) {
			b.WriteByte(v[i])
			continue
		}
		i++
		if i+2 < len(v) && '0' <= v[i] && v[i] <= '3' && isOctal(v[i+1]) && isOctal(v[i+2]) {
			b.WriteByte((v[i]-'0')<<6 | (v[i+1]-'0')<<3 | (v[i+2] - '0'))
			i += 2
			continue
		}
		b.WriteByte(unescapeByte(v[i]))
	}
	return b.String()
}

func isOctal(c byte) bool {
	return '0' <= c && c <= '7'
}

// unescapeByte decodes the character after a backslash as tmux's
// configuration parser does.
func unescapeByte(c byte) byte {
	switch c {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 'e':
		return 0x1b
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 's':
		return ' '
	case 't':
		return '\t'
	case 'v':
		return '\v'
	}
	return c
}
