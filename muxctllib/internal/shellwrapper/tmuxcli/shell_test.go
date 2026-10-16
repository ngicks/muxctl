package tmuxcli

import (
	"maps"
	"slices"
	"testing"
)

func TestFormat(t *testing.T) {
	got := format([]string{"pane_id", "@tag"})
	if want := "#{pane_id}\x1f#{@tag}"; got != want {
		t.Errorf("format = %q, want %q", got, want)
	}
}

func TestRecords(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		n       int
		want    [][]string
		wantErr bool
	}{
		{name: "empty", out: "", n: 2, want: nil},
		{name: "rows", out: "a\x1fb\nc\x1f\n", n: 2, want: [][]string{{"a", "b"}, {"c", ""}}},
		{name: "no trailing newline", out: "a\x1fb", n: 2, want: [][]string{{"a", "b"}}},
		{name: "blank lines skipped", out: "\na\x1fb\n\n", n: 2, want: [][]string{{"a", "b"}}},
		{name: "field count mismatch", out: "a\x1fb\x1fc\n", n: 2, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := records(tt.out, tt.n)
			if (err != nil) != tt.wantErr {
				t.Fatalf("records err = %v, wantErr %v", err, tt.wantErr)
			}
			if !slices.EqualFunc(got, tt.want, slices.Equal) {
				t.Errorf("records = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseFields(t *testing.T) {
	var (
		s    string
		n    int
		zero int
		b    bool
		f    bool
	)
	names := []string{"s", "n", "zero", "b", "f"}
	if err := parseFields(names, []string{"x", "42", "", "1", "0"}, &s, &n, &zero, &b, &f); err != nil {
		t.Fatal(err)
	}
	if s != "x" || n != 42 || zero != 0 || !b || f {
		t.Errorf("parseFields = %q %d %d %t %t", s, n, zero, b, f)
	}
	if err := parseFields([]string{"n"}, []string{"x"}, &n); err == nil {
		t.Error("parseFields with a non-number int field: want error")
	}
}

func TestUserOptionVars(t *testing.T) {
	got, err := userOptionVars([]string{"@a", "@b"})
	if err != nil || !slices.Equal(got, []string{"@a", "@b"}) {
		t.Errorf("userOptionVars = %q, %v", got, err)
	}
	if _, err := userOptionVars([]string{"a"}); err == nil {
		t.Error("userOptionVars without '@': want error")
	}
}

func TestUserOptionMap(t *testing.T) {
	if got := userOptionMap(nil, nil); got != nil {
		t.Errorf("userOptionMap(nil) = %v, want nil", got)
	}
	got := userOptionMap([]string{"@a", "@b"}, []string{"1", ""})
	if want := map[string]string{"@a": "1", "@b": ""}; !maps.Equal(got, want) {
		t.Errorf("userOptionMap = %v, want %v", got, want)
	}
}

func TestArgvTrailing(t *testing.T) {
	tests := []struct {
		name string
		vs   []string
		want []string
	}{
		{name: "none", vs: nil, want: []string{"cmd"}},
		{name: "some", vs: []string{"-x", "y"}, want: []string{"cmd", "--", "-x", "y"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := argv{"cmd"}.trailing(tt.vs)
			if !slices.Equal(a, tt.want) {
				t.Errorf("trailing = %q, want %q", a, tt.want)
			}
		})
	}
}
