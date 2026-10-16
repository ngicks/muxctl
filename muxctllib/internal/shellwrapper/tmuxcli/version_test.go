package tmuxcli

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		out  string
		want Version
	}{
		{"tmux 3.7c\n", Version{Raw: "3.7c", Major: 3, Minor: 7, Suffix: "c"}},
		{"tmux 3.4\n", Version{Raw: "3.4", Major: 3, Minor: 4}},
		{"tmux 3.4-rc\n", Version{Raw: "3.4-rc", Major: 3, Minor: 4, Suffix: "-rc"}},
		{"tmux next-3.8\n", Version{Raw: "next-3.8", Major: 3, Minor: 8}},
		{"tmux openbsd-7.5\n", Version{Raw: "openbsd-7.5"}},
		{"tmux master\n", Version{Raw: "master"}},
		{"tmux 3.x\n", Version{Raw: "3.x"}},
	}
	for _, tt := range tests {
		if got := parseVersion(tt.out); got != tt.want {
			t.Errorf("parseVersion(%q) = %+v, want %+v", tt.out, got, tt.want)
		}
	}
}

func TestVersionLess(t *testing.T) {
	v := func(major, minor int, suffix string) Version {
		return Version{Major: major, Minor: minor, Suffix: suffix}
	}
	tests := []struct {
		a, b Version
		want bool
	}{
		{v(3, 7, "b"), v(3, 7, "c"), true},
		{v(3, 7, ""), v(3, 7, "a"), true},
		{v(3, 7, "c"), v(3, 7, "c"), false},
		{v(3, 7, "c"), v(3, 7, "b"), false},
		{v(3, 6, "z"), v(3, 7, ""), true},
		{v(2, 9, "a"), v(3, 0, ""), true},
		{v(3, 8, ""), v(3, 7, "c"), false},
	}
	for _, tt := range tests {
		if got := tt.a.Less(tt.b); got != tt.want {
			t.Errorf("%+v.Less(%+v) = %t, want %t", tt.a, tt.b, got, tt.want)
		}
	}
}
