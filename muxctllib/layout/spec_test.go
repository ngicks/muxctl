package layout

import (
	"encoding/json/v2"
	"errors"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestDir_Canonical(t *testing.T) {
	tests := []struct {
		in      Dir
		want    Dir
		wantErr error
	}{
		{"", DirHorizontal, nil},
		{DirHorizontal, DirHorizontal, nil},
		{DirH, DirHorizontal, nil},
		{DirVertical, DirVertical, nil},
		{DirV, DirVertical, nil},
		{"Horizontal", "", ErrInvalidDir},
		{"left", "", ErrInvalidDir},
	}
	for _, tt := range tests {
		got, err := tt.in.Canonicalize()
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Dir(%q).Canonical() err = %v, want %v", tt.in, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("Dir(%q).Canonical() = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// errDecoder stands for "the decoder's own error": non-nil and none of this
// package's sentinels. Those sentinels mean a string that failed to parse, so
// a value of the wrong shape must not carry them.
var errDecoder = errors.New("decoder error")

func checkErr(t *testing.T, label string, err, want error) {
	t.Helper()
	if want == errDecoder {
		switch {
		case err == nil:
			t.Errorf("%s: err = nil, want a decoder error", label)
		case errors.Is(err, ErrInvalidDir), errors.Is(err, ErrInvalidSize):
			t.Errorf("%s: err = %v, want a decoder error without a layout sentinel", label, err)
		}
		return
	}
	if !errors.Is(err, want) {
		t.Errorf("%s: err = %v, want %v", label, err, want)
	}
}

func TestDir_Unmarshal(t *testing.T) {
	tests := []struct {
		yamlIn, jsonIn       string
		want                 Dir
		wantYAMLErr, wantErr error
	}{
		{`h`, `"h"`, DirHorizontal, nil, nil},
		{`horizontal`, `"horizontal"`, DirHorizontal, nil, nil},
		{`v`, `"v"`, DirVertical, nil, nil},
		{`vertical`, `"vertical"`, DirVertical, nil, nil},
		{`""`, `""`, DirHorizontal, nil, nil},
		{`sideways`, `"sideways"`, "", ErrInvalidDir, ErrInvalidDir},
		{`[h]`, `["h"]`, "", errDecoder, errDecoder},
		// A bare YAML 1 is still a scalar and reads as the string "1".
		{`1`, `1`, "", ErrInvalidDir, errDecoder},
	}
	for _, tt := range tests {
		var y Dir
		err := yaml.Unmarshal([]byte(tt.yamlIn), &y)
		checkErr(t, "yaml "+tt.yamlIn, err, tt.wantYAMLErr)
		if err == nil && y != tt.want {
			t.Errorf("yaml %q = %q, want %q", tt.yamlIn, y, tt.want)
		}

		var j Dir
		err = json.Unmarshal([]byte(tt.jsonIn), &j)
		checkErr(t, "json "+tt.jsonIn, err, tt.wantErr)
		if err == nil && j != tt.want {
			t.Errorf("json %s = %q, want %q", tt.jsonIn, j, tt.want)
		}
	}
}

func TestSize_Canonical(t *testing.T) {
	tests := []struct {
		in      Size
		want    Size
		wantErr error
	}{
		{"", "1", nil},
		{"2", "2", nil},
		{"5c", "5c", nil},
		{"20%", "20%", nil},
		{"100%", "100%", nil},
		{"007", "7", nil},
		{"+3c", "3c", nil},
		{"0", "", ErrInvalidSize},
		{"-3", "", ErrInvalidSize},
		{"0c", "", ErrInvalidSize},
		{"101%", "", ErrInvalidSize},
		{"c", "", ErrInvalidSize},
		{"5cm", "", ErrInvalidSize},
		{"1.5", "", ErrInvalidSize},
		{"  3", "", ErrInvalidSize},
	}
	for _, tt := range tests {
		got, err := tt.in.Canonicalize()
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Size(%q).Canonical() err = %v, want %v", tt.in, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("Size(%q).Canonical() = %q, want %q", tt.in, got, tt.want)
		}
		if err != nil {
			continue
		}
		again, err := got.Canonicalize()
		if err != nil || again != got {
			t.Errorf("Size(%q).Canonical() = %q, %v; want it unchanged", got, again, err)
		}
	}
}

func TestSize_Unmarshal(t *testing.T) {
	tests := []struct {
		yamlIn, jsonIn       string
		want                 Size
		wantYAMLErr, wantErr error
	}{
		{`2`, `2`, "2", nil, nil},
		{`"2"`, `"2"`, "2", nil, nil},
		{`5c`, `"5c"`, "5c", nil, nil},
		{`20%`, `"20%"`, "20%", nil, nil},
		{`""`, `""`, "1", nil, nil},
		{`0`, `0`, "", ErrInvalidSize, ErrInvalidSize},
		{`1.5`, `1.5`, "", ErrInvalidSize, ErrInvalidSize},
		{`101%`, `"101%"`, "", ErrInvalidSize, ErrInvalidSize},
		{`[2]`, `[2]`, "", errDecoder, errDecoder},
		{`{a: 1}`, `{"a":1}`, "", errDecoder, errDecoder},
		// A bare YAML true is still a scalar and reads as the string "true".
		{`true`, `true`, "", ErrInvalidSize, errDecoder},
	}
	for _, tt := range tests {
		var y Size
		err := yaml.Unmarshal([]byte(tt.yamlIn), &y)
		checkErr(t, "yaml "+tt.yamlIn, err, tt.wantYAMLErr)
		if err == nil && y != tt.want {
			t.Errorf("yaml %q = %q, want %q", tt.yamlIn, y, tt.want)
		}

		var j Size
		err = json.Unmarshal([]byte(tt.jsonIn), &j)
		checkErr(t, "json "+tt.jsonIn, err, tt.wantErr)
		if err == nil && j != tt.want {
			t.Errorf("json %s = %q, want %q", tt.jsonIn, j, tt.want)
		}
	}
}

func TestCommandSpec_Unmarshal(t *testing.T) {
	tests := []struct {
		yamlIn, jsonIn string
		want           CommandSpec
		wantErr        error
	}{
		{
			`[echo, foo]`, `["echo","foo"]`,
			CommandSpec{Args: []string{"echo", "foo"}}, nil,
		},
		{
			"title: t\nargs: [echo, bar]\nenvs: [FOO=BAR]",
			`{"title":"t","args":["echo","bar"],"envs":["FOO=BAR"]}`,
			CommandSpec{Args: []string{"echo", "bar"}, Title: "t", Envs: []string{"FOO=BAR"}}, nil,
		},
		{`{}`, `{}`, CommandSpec{}, nil},
		{`echo`, `"echo"`, CommandSpec{}, ErrInvalidCommand},
		{`1`, `1`, CommandSpec{}, ErrInvalidCommand},
	}
	for _, tt := range tests {
		var y CommandSpec
		err := yaml.Unmarshal([]byte(tt.yamlIn), &y)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("yaml %q: err = %v, want %v", tt.yamlIn, err, tt.wantErr)
		}
		if err == nil && !reflect.DeepEqual(y, tt.want) {
			t.Errorf("yaml %q = %+v, want %+v", tt.yamlIn, y, tt.want)
		}

		var j CommandSpec
		err = json.Unmarshal([]byte(tt.jsonIn), &j)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("json %s: err = %v, want %v", tt.jsonIn, err, tt.wantErr)
		}
		if err == nil && !reflect.DeepEqual(j, tt.want) {
			t.Errorf("json %s = %+v, want %+v", tt.jsonIn, j, tt.want)
		}
	}
}

func leafOf(args ...string) PaneSpec {
	return LeafPane(LeafSpec{Command: CommandSpec{Args: args}})
}

func TestPaneSpec_Accessors(t *testing.T) {
	split := SplitSpec{Dir: DirVertical, Sizes: Sizes{"1"}, Panes: []PaneSpec{leafOf("a")}}
	leaf := LeafSpec{Command: CommandSpec{Args: []string{"a"}}, Focus: true}

	p := SplitPane(split)
	if got, ok := p.Split(); !ok || !reflect.DeepEqual(got, split) {
		t.Errorf("SplitPane(s).Split() = %+v, %v", got, ok)
	}
	if _, ok := p.Leaf(); ok {
		t.Error("SplitPane(s).Leaf() reported a leaf")
	}

	p = LeafPane(leaf)
	if got, ok := p.Leaf(); !ok || !reflect.DeepEqual(got, leaf) {
		t.Errorf("LeafPane(l).Leaf() = %+v, %v", got, ok)
	}
	if _, ok := p.Split(); ok {
		t.Error("LeafPane(l).Split() reported a split")
	}

	var zero PaneSpec
	if _, ok := zero.Split(); ok {
		t.Error("zero.Split() reported a split")
	}
	if _, ok := zero.Leaf(); ok {
		t.Error("zero.Leaf() reported a leaf")
	}
}

func TestPaneSpec_Unmarshal(t *testing.T) {
	tests := []struct {
		name           string
		yamlIn, jsonIn string
		want           PaneSpec
		wantErr        error
	}{
		{
			"leaf short",
			`command: [echo, foo]`, `{"command":["echo","foo"]}`,
			leafOf("echo", "foo"), nil,
		},
		{
			"leaf long with focus",
			"command:\n  args: [echo, bar]\n  title: t\nfocus: true",
			`{"command":{"args":["echo","bar"],"title":"t"},"focus":true}`,
			LeafPane(LeafSpec{
				Command: CommandSpec{Args: []string{"echo", "bar"}, Title: "t"},
				Focus:   true,
			}), nil,
		},
		{
			"split with defaults",
			"sizes: [1, 1]\npanes:\n  - command: [a]\n  - command: [b]",
			`{"sizes":[1,1],"panes":[{"command":["a"]},{"command":["b"]}]}`,
			SplitPane(
				SplitSpec{Sizes: Sizes{"1", "1"}, Panes: []PaneSpec{leafOf("a"), leafOf("b")}},
			),
			nil,
		},
		{
			"split canonicalizes dir and sizes",
			"dir: v\nsizes: [3, 1c, 50%]\npanes:\n  - command: [a]\n  - command: [b]\n  - command: [c]",
			`{"dir":"v","sizes":[3,"1c","50%"],` +
				`"panes":[{"command":["a"]},{"command":["b"]},{"command":["c"]}]}`,
			SplitPane(SplitSpec{
				Dir:   DirVertical,
				Sizes: []Size{"3", "1c", "50%"},
				Panes: []PaneSpec{leafOf("a"), leafOf("b"), leafOf("c")},
			}),
			nil,
		},
		{
			"nested",
			"sizes: [1, 1]\npanes:\n  - sizes: [1]\n    panes:\n      - command: [a]\n  - command: [b]",
			`{"sizes":[1,1],"panes":[{"sizes":[1],"panes":[{"command":["a"]}]},{"command":["b"]}]}`,
			SplitPane(SplitSpec{Sizes: Sizes{"1", "1"}, Panes: []PaneSpec{
				SplitPane(SplitSpec{Sizes: Sizes{"1"}, Panes: []PaneSpec{leafOf("a")}}),
				leafOf("b"),
			}}),
			nil,
		},
		{
			"sizes mismatch",
			"sizes: [1]\npanes:\n  - command: [a]\n  - command: [b]",
			`{"sizes":[1],"panes":[{"command":["a"]},{"command":["b"]}]}`,
			PaneSpec{}, ErrSizesMismatch,
		},
		{
			"both",
			"command: [a]\npanes:\n  - command: [b]",
			`{"command":["a"],"panes":[{"command":["b"]}]}`,
			PaneSpec{}, ErrSplitXorLeaf,
		},
		{
			"neither",
			`focus: true`, `{"focus":true}`,
			PaneSpec{}, ErrSplitXorLeaf,
		},
		{
			"dir alone is not a split",
			`dir: h`, `{"dir":"h"}`,
			PaneSpec{}, ErrSplitXorLeaf,
		},
		{
			"present empty panes is a split",
			`panes: []`, `{"panes":[]}`,
			SplitPane(SplitSpec{Panes: []PaneSpec{}}), nil,
		},
		{
			"null command counts as absent",
			`command: null`, `{"command":null}`,
			PaneSpec{}, ErrSplitXorLeaf,
		},
		{
			"nested error surfaces",
			"sizes: [1]\npanes:\n  - dir: up\n    sizes: [1]\n    panes: [{command: [a]}]",
			`{"sizes":[1],"panes":[{"dir":"up","sizes":[1],"panes":[{"command":["a"]}]}]}`,
			PaneSpec{}, ErrInvalidDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var y PaneSpec
			err := yaml.Unmarshal([]byte(tt.yamlIn), &y)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("yaml: err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(y, tt.want) {
				t.Errorf("yaml = %s, want %s", dump(y), dump(tt.want))
			}

			var j PaneSpec
			err = json.Unmarshal([]byte(tt.jsonIn), &j)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("json: err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(j, tt.want) {
				t.Errorf("json = %s, want %s", dump(j), dump(tt.want))
			}
		})
	}
}

func TestPaneSpec_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		in       PaneSpec
		wantYAML string
		wantJSON string
	}{
		{
			"leaf",
			LeafPane(LeafSpec{Command: CommandSpec{Args: []string{"echo", "a"}}, Focus: true}),
			"command:\n    args:\n        - echo\n        - a\nfocus: true\n",
			`{"command":{"args":["echo","a"]},"focus":true}`,
		},
		{
			"split",
			SplitPane(
				SplitSpec{
					Dir:   DirVertical,
					Sizes: []Size{"1", "2c"},
					Panes: []PaneSpec{leafOf("a"), leafOf("b")},
				},
			),
			"dir: vertical\nsizes:\n    - \"1\"\n    - 2c\n" +
				"panes:\n    - command:\n        args:\n            - a\n" +
				"    - command:\n        args:\n            - b\n",
			`{"dir":"vertical","sizes":["1","2c"],` +
				`"panes":[{"command":{"args":["a"]}},{"command":{"args":["b"]}}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y, err := yaml.Marshal(tt.in)
			if err != nil {
				t.Fatalf("yaml: %v", err)
			}
			if string(y) != tt.wantYAML {
				t.Errorf("yaml = %q, want %q", y, tt.wantYAML)
			}
			j, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("json: %v", err)
			}
			if string(j) != tt.wantJSON {
				t.Errorf("json = %s, want %s", j, tt.wantJSON)
			}
		})
	}

	var zero PaneSpec
	if _, err := yaml.Marshal(zero); !errors.Is(err, ErrSplitXorLeaf) {
		t.Errorf("yaml.Marshal(zero) err = %v, want %v", err, ErrSplitXorLeaf)
	}
	if _, err := json.Marshal(zero); !errors.Is(err, ErrSplitXorLeaf) {
		t.Errorf("json.Marshal(zero) err = %v, want %v", err, ErrSplitXorLeaf)
	}
}

const specDoc = `
layouts:
  - name: dev
    desc: description of layout
    layout:
      dir: horizontal
      sizes: [2, 5c, 20%]
      panes:
        - command: [echo, foo]
        - command:
            title: title to show
            args: [echo, bar]
            envs: [FOO=BAR]
        - dir: v
          sizes: [3, 1]
          panes:
            - command: [echo, baz]
            - command: [echo, qux]
              focus: true
  - name: single
    layout:
      sizes: [1]
      panes:
        - command: [sh]
`

var specWant = Spec{Layouts: []Layout{
	{
		Name: "dev",
		Desc: "description of layout",
		Layout: SplitSpec{
			Dir:   DirHorizontal,
			Sizes: []Size{"2", "5c", "20%"},
			Panes: []PaneSpec{
				leafOf("echo", "foo"),
				LeafPane(LeafSpec{Command: CommandSpec{
					Args:  []string{"echo", "bar"},
					Title: "title to show",
					Envs:  []string{"FOO=BAR"},
				}}),
				SplitPane(SplitSpec{
					Dir:   DirVertical,
					Sizes: []Size{"3", "1"},
					Panes: []PaneSpec{
						leafOf("echo", "baz"),
						LeafPane(
							LeafSpec{
								Command: CommandSpec{Args: []string{"echo", "qux"}},
								Focus:   true,
							},
						),
					},
				}),
			},
		},
	},
	{
		Name:   "single",
		Layout: SplitSpec{Sizes: Sizes{"1"}, Panes: []PaneSpec{leafOf("sh")}},
	},
}}

func TestSpec_Unmarshal(t *testing.T) {
	var got Spec
	if err := yaml.Unmarshal([]byte(specDoc), &got); err != nil {
		t.Fatalf("yaml: %v", err)
	}
	if !reflect.DeepEqual(got, specWant) {
		t.Errorf("yaml = %s\nwant %s", dump(got), dump(specWant))
	}
}

func TestSpec_RoundTrip(t *testing.T) {
	y, err := yaml.Marshal(specWant)
	if err != nil {
		t.Fatalf("yaml.Marshal: %v", err)
	}
	var fromYAML Spec
	if err := yaml.Unmarshal(y, &fromYAML); err != nil {
		t.Fatalf("yaml.Unmarshal: %v\n%s", err, y)
	}
	if !reflect.DeepEqual(fromYAML, specWant) {
		t.Errorf("yaml round trip = %s\nwant %s", dump(fromYAML), dump(specWant))
	}

	j, err := json.Marshal(specWant)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var fromJSON Spec
	if err := json.Unmarshal(j, &fromJSON); err != nil {
		t.Fatalf("json.Unmarshal: %v\n%s", err, j)
	}
	if !reflect.DeepEqual(fromJSON, specWant) {
		t.Errorf("json round trip = %s\nwant %s", dump(fromJSON), dump(specWant))
	}
}

// dump renders a spec for a failure message. The unexported pointers in
// PaneSpec print as addresses otherwise.
func dump(v any) string {
	b, err := json.Marshal(v, json.Deterministic(true))
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func TestSizes_ComputeSplits(t *testing.T) {
	tests := []struct {
		name    string
		sizes   Sizes
		extent  int
		want    []int
		wantErr error
	}{
		{"no panes", nil, 80, []int{}, nil},
		{"one pane takes all", Sizes{"1"}, 80, []int{80}, nil},
		{"equal weights, remainder to last", Sizes{"1", "1"}, 80, []int{39, 40}, nil},
		{"three equal weights", Sizes{"1", "1", "1"}, 12, []int{3, 3, 4}, nil},
		{"weight ratio", Sizes{"3", "1"}, 41, []int{30, 10}, nil},
		{"empty entry is weight 1", Sizes{"2", ""}, 31, []int{20, 10}, nil},
		{"mixed units", Sizes{"2", "5c", "20%"}, 100, []int{73, 5, 20}, nil},
		{"weights share what percent leaves", Sizes{"2", "1", "50%"}, 100, []int{32, 16, 50}, nil},
		{"non-canonical entries parse", Sizes{"+2", "02c"}, 10, []int{7, 2}, nil},
		{"cells capped by what is left", Sizes{"50c", "50c"}, 60, []int{50, 9}, nil},
		{"percent capped by what is left", Sizes{"50%", "50%"}, 100, []int{50, 49}, nil},
		{"weight left with no room", Sizes{"10c", "1"}, 10, []int{9, 0}, nil},
		{"extent below separators", Sizes{"1", "1", "1"}, 1, []int{0, 0, 0}, nil},
		{"negative extent", Sizes{"1", "1"}, -5, []int{0, 0}, nil},
		{"malformed entry", Sizes{"x"}, 80, nil, ErrInvalidSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sizes.ComputeSplits(tt.extent)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeSplits(%d) = %v, want %v", tt.extent, got, tt.want)
			}
		})
	}
}
