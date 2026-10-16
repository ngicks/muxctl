package tmuxcli

import (
	"slices"
	"testing"
)

func TestUnquoteOptionValue(t *testing.T) {
	// Quoted forms are what tmux 3.7c show-options printed for each value.
	tests := []struct {
		quoted string
		want   string
	}{
		{`\~foo`, "~foo"},
		{`"a\$b"`, "a$b"},
		{`"#x"`, "#x"},
		{`"a;b"`, "a;b"},
		{`"it's"`, "it's"},
		{`tab\tx`, "tab\tx"},
		{`é`, "é"},
		{`a\\b`, `a\b`},
		{`"{x}"`, "{x}"},
		{`'a"b'`, `a"b`},
		{`''`, ""},
		{`"a b \"q\" \\t"`, `a b "q" \t`},
		{`e\033x`, "e\x1bx"},
		{`c\001x`, "c\x01x"},
		{`d\177x`, "d\x7fx"},
		{`'a"b\\c'`, `a"b\c`},
		{`'a"\tb'`, "a\"\tb"},
		{`nl\nx`, "nl\nx"},
		{`"a%b"`, "a%b"},
		{`\~`, "~"},
		{`x~`, "x~"},
		{`" lead"`, " lead"},
		{`-dash`, "-dash"},
		{"a`b", "a`b"},
		{`"a'b\"c"`, `a'b"c`},
		{`\\~`, `\~`},
		{``, ``},
		{`on`, `on`},
		{`trailing\`, `trailing\`},
	}
	for _, tt := range tests {
		if got := unquoteOptionValue(tt.quoted); got != tt.want {
			t.Errorf("unquoteOptionValue(%s) = %q, want %q", tt.quoted, got, tt.want)
		}
	}
}

func TestParseOptionEntries(t *testing.T) {
	tests := []struct {
		name      string
		out       string
		option    string
		valueOnly bool
		want      []OptionEntry
	}{
		{name: "empty", out: "", want: nil},
		{
			name: "names and values",
			out: "@tag \"a b\"\n" +
				"cursor-colour* none\n" +
				"after-set-hook\n" +
				"update-environment[0] DISPLAY\n" +
				"@empty ''\n",
			want: []OptionEntry{
				{Name: "@tag", Value: "a b"},
				{Name: "cursor-colour", Value: "none", Inherited: true},
				{Name: "after-set-hook"},
				{Name: "update-environment[0]", Value: "DISPLAY"},
				{Name: "@empty"},
			},
		},
		{
			name:      "value only is raw",
			out:       "a b\tc\n\n",
			option:    "@v",
			valueOnly: true,
			want:      []OptionEntry{{Name: "@v", Value: "a b\tc"}, {Name: "@v"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOptionEntries(tt.out, tt.option, tt.valueOnly)
			if !slices.Equal(got, tt.want) {
				t.Errorf("parseOptionEntries = %+v, want %+v", got, tt.want)
			}
		})
	}
}
