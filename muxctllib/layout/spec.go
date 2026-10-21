// Package layout describes tiled pane layouts for terminal multiplexers as a
// YAML or JSON document.
package layout

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Spec is the YAML document: an ordered list of named layouts under one
// "layouts" key. A list keeps the order the document wrote, which a mapping
// would lose; the order is what a consumer cycles through.
//
//	layouts:
//	  - name: dev
//	    desc: description of layout
//	    layout:
//	      dir: horizontal
//	      sizes: [2, 5c, 20%]
//	      panes:
//	        - command: [echo, foo]
//	        - command:
//	            title: title to show
//	            args: [echo, bar]
//	            envs: [FOO=BAR]
//	        - dir: v
//	          sizes: [3, 1]
//	          panes:
//	            - command: [echo, baz]
//	            - command: [echo, qux]
//	              focus: true
//
// Spec encodes and decodes through go.yaml.in/yaml/v4 and encoding/json/v2.
// Decoding canonicalizes every [Dir] and [Size], so a document does not round
// trip byte for byte. A decoded Spec encodes back to an equal one.
type Spec struct {
	// Layouts lists the layouts in document order.
	Layouts []Layout `yaml:"layouts" json:"layouts"`
}

// Layout is one named layout.
type Layout struct {
	// Name identifies the layout. It is required and unique within Layouts.
	Name string `yaml:"name" json:"name"`
	// Desc describes the layout.
	Desc string `yaml:"desc,omitempty" json:"desc,omitzero"`
	// Layout is the root split.
	Layout SplitSpec `yaml:"layout" json:"layout"`
}

// PaneSpec is one entry of [SplitSpec.Panes]: a split or a leaf. The fields
// are unexported so only [SplitPane], [LeafPane] and the unmarshalers build
// one, and each of them sets exactly one side. A document mapping that sets
// panes together with command, or neither of them, fails to decode with
// [ErrSplitXorLeaf]. The zero PaneSpec is neither and fails to marshal.
type PaneSpec struct {
	split *SplitSpec
	leaf  *LeafSpec
}

// SplitPane returns a PaneSpec holding s.
func SplitPane(s SplitSpec) PaneSpec { return PaneSpec{split: &s} }

// LeafPane returns a PaneSpec holding l.
func LeafPane(l LeafSpec) PaneSpec { return PaneSpec{leaf: &l} }

// Split returns the split and true when p holds one.
func (p PaneSpec) Split() (SplitSpec, bool) {
	if p.split == nil {
		return SplitSpec{}, false
	}
	return *p.split, true
}

// Leaf returns the leaf and true when p holds one.
func (p PaneSpec) Leaf() (LeafSpec, bool) {
	if p.leaf == nil {
		return LeafSpec{}, false
	}
	return *p.leaf, true
}

// paneKeys holds only the keys that tell a split from a leaf. The values are
// decoded into any just to see whether each key is present; the pane itself is
// decoded afterwards into [SplitSpec] or [LeafSpec] through their own tags.
type paneKeys struct {
	Panes   any `yaml:"panes" json:"panes"`
	Command any `yaml:"command" json:"command"`
}

// isSplit reports whether the keys mark a split. Having both keys, or neither,
// is [ErrSplitXorLeaf].
func (k paneKeys) isSplit() (bool, error) {
	hasPanes, hasCommand := k.Panes != nil, k.Command != nil
	if hasPanes == hasCommand {
		return false, ErrSplitXorLeaf
	}
	return hasPanes, nil
}

// MarshalYAML writes the held side as one mapping.
func (p PaneSpec) MarshalYAML() (any, error) {
	switch {
	case p.split != nil:
		return p.split, nil
	case p.leaf != nil:
		return p.leaf, nil
	}
	return nil, ErrSplitXorLeaf
}

// MarshalJSONTo writes the held side as one object.
func (p PaneSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	switch {
	case p.split != nil:
		return json.MarshalEncode(enc, p.split)
	case p.leaf != nil:
		return json.MarshalEncode(enc, p.leaf)
	}
	return ErrSplitXorLeaf
}

// UnmarshalYAML reads one mapping into a split or a leaf, chosen by whether
// the panes key or the command key is present.
func (p *PaneSpec) UnmarshalYAML(node *yaml.Node) error {
	var keys paneKeys
	if err := node.Decode(&keys); err != nil {
		return err
	}
	isSplit, err := keys.isSplit()
	if err != nil {
		return fmt.Errorf("%w (line %d)", err, node.Line)
	}
	if isSplit {
		var s SplitSpec
		if err := node.Decode(&s); err != nil {
			return err
		}
		*p = SplitPane(s)
		return nil
	}
	var l LeafSpec
	if err := node.Decode(&l); err != nil {
		return err
	}
	*p = LeafPane(l)
	return nil
}

// UnmarshalJSONFrom reads one object into a split or a leaf, chosen by whether
// the panes member or the command member is present.
func (p *PaneSpec) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	v, err := dec.ReadValue()
	if err != nil {
		return err
	}
	var keys paneKeys
	if err := json.Unmarshal(v, &keys); err != nil {
		return err
	}
	isSplit, err := keys.isSplit()
	if err != nil {
		return err
	}
	if isSplit {
		var s SplitSpec
		if err := json.Unmarshal(v, &s); err != nil {
			return err
		}
		*p = SplitPane(s)
		return nil
	}
	var l LeafSpec
	if err := json.Unmarshal(v, &l); err != nil {
		return err
	}
	*p = LeafPane(l)
	return nil
}

// SplitSpec divides its area along Dir and gives each pane the part its entry
// in Sizes asks for.
type SplitSpec struct {
	// Dir is the split axis. An absent or empty dir is horizontal.
	Dir Dir `yaml:"dir,omitempty" json:"dir,omitzero"`
	// Sizes has exactly one entry per pane, index-parallel to Panes. Decoding
	// fails with [ErrSizesMismatch] otherwise.
	Sizes Sizes `yaml:"sizes,omitempty" json:"sizes,omitzero"`
	// Panes are the split's children in order. The presence of the panes key
	// is what marks a [PaneSpec] as a split.
	Panes []PaneSpec `yaml:"panes,omitempty" json:"panes,omitzero"`
}

// UnmarshalYAML decodes the split and validates its sizes.
func (s *SplitSpec) UnmarshalYAML(node *yaml.Node) error {
	// plain has SplitSpec's fields but not its methods, so Decode fills the
	// fields instead of recursing into this method.
	type plain SplitSpec
	var p plain
	if err := node.Decode(&p); err != nil {
		return err
	}
	if err := SplitSpec(p).validate(); err != nil {
		return fmt.Errorf("%w (line %d)", err, node.Line)
	}
	*s = SplitSpec(p)
	return nil
}

// UnmarshalJSONFrom decodes the split and validates its sizes.
func (s *SplitSpec) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// See UnmarshalYAML for why plain exists.
	type plain SplitSpec
	var p plain
	if err := json.UnmarshalDecode(dec, &p); err != nil {
		return err
	}
	if err := SplitSpec(p).validate(); err != nil {
		return err
	}
	*s = SplitSpec(p)
	return nil
}

// validate reports [ErrSizesMismatch] unless Sizes has one entry per pane.
func (s SplitSpec) validate() error {
	if len(s.Sizes) != len(s.Panes) {
		return fmt.Errorf("%w: %d sizes for %d panes", ErrSizesMismatch, len(s.Sizes), len(s.Panes))
	}
	return nil
}

// LeafSpec is a pane that runs a command.
type LeafSpec struct {
	// Command marks the mapping as a leaf.
	Command CommandSpec `yaml:"command" json:"command"`
	// Focus asks for this pane to be focused once the layout is built. At most
	// one leaf per layout sets it. With none set, the first leaf in document
	// order is focused.
	Focus bool `yaml:"focus,omitempty" json:"focus,omitzero"`
}

// CommandSpec is what a leaf runs. It has two forms: a sequence, which is Args
// alone, and a mapping with the fields below.
type CommandSpec struct {
	// Args is the argv run in the pane, executed directly without a shell.
	// Empty leaves the multiplexer's default shell.
	Args []string `yaml:"args,omitempty" json:"args,omitzero"`
	// Title is the pane title shown by the multiplexer. Empty leaves it to the
	// command.
	Title string `yaml:"title,omitempty" json:"title,omitzero"`
	// Envs is extra environment for the command, each entry "KEY=VALUE".
	Envs []string `yaml:"envs,omitempty" json:"envs,omitzero"`
}

// UnmarshalYAML accepts a sequence as Args or a mapping with the struct's
// fields.
func (c *CommandSpec) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		var args []string
		if err := node.Decode(&args); err != nil {
			return err
		}
		*c = CommandSpec{Args: args}
		return nil
	case yaml.MappingNode:
		type plain CommandSpec
		var fields plain
		if err := node.Decode(&fields); err != nil {
			return err
		}
		*c = CommandSpec(fields)
		return nil
	}
	return fmt.Errorf("%w: got %v at line %d", ErrInvalidCommand, node.Kind, node.Line)
}

// UnmarshalJSONFrom accepts an array as Args or an object with the struct's
// fields.
func (c *CommandSpec) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case '[':
		var args []string
		if err := json.UnmarshalDecode(dec, &args); err != nil {
			return err
		}
		*c = CommandSpec{Args: args}
		return nil
	case '{':
		type plain CommandSpec
		var fields plain
		if err := json.UnmarshalDecode(dec, &fields); err != nil {
			return err
		}
		*c = CommandSpec(fields)
		return nil
	}
	return fmt.Errorf("%w: got %v", ErrInvalidCommand, dec.PeekKind())
}

// Dir is the axis along which a [SplitSpec] lays out its panes. The words follow
// tmux: a horizontal split puts panes side by side. zellij names the same split
// "vertical", so the words are mapped, never passed through.
//
// Decoding accepts the short forms and the empty string and canonicalizes to
// the long form, so a decoded Dir is always one of the two long constants.
type Dir string

// Spellings of Dir. The long forms are canonical.
const (
	// DirHorizontal places panes side by side, left to right.
	// tmux: split-window -h. zellij KDL: split_direction="vertical".
	DirHorizontal Dir = "horizontal"
	DirH          Dir = "h"
	// DirVertical stacks panes top to bottom.
	// tmux: split-window -v. zellij KDL: split_direction="horizontal".
	DirVertical Dir = "vertical"
	DirV        Dir = "v"
)

// Canonicalize returns the long form of d. An empty Dir is [DirHorizontal].
func (d Dir) Canonicalize() (Dir, error) {
	switch d {
	case "", DirHorizontal, DirH:
		return DirHorizontal, nil
	case DirVertical, DirV:
		return DirVertical, nil
	}
	return "", fmt.Errorf("%w: %q", ErrInvalidDir, string(d))
}

// UnmarshalYAML decodes a scalar through [Dir.Canonicalize]. A node that is not a
// scalar fails with the decoder's own error, not with [ErrInvalidDir]: that
// error means a string that is not a direction.
func (d *Dir) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	return d.set(s)
}

// UnmarshalJSON decodes a string through [Dir.Canonicalize]. A value that is not
// a string fails with the decoder's own error, not with [ErrInvalidDir].
func (d *Dir) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.set(s)
}

func (d *Dir) set(s string) error {
	dir, err := Dir(s).Canonicalize()
	if err != nil {
		return err
	}
	*d = dir
	return nil
}

// Sizes is the [SplitSpec.Sizes] list, index-parallel to the split's panes.
type Sizes []Size

// ComputeSplits divides extent cells among the panes along the split axis, one
// result per entry of s. extent is the parent's width for a horizontal split
// and its height for a vertical one.
//
// One cell per separator between panes is set aside first. Cell and percent
// panes then take what they ask for in order, each capped at what is left, so
// the result never exceeds extent. A percent is of the whole extent. Weighted
// panes share the rest in proportion to their weights, and the rounding
// remainder goes to the last weighted pane. A pane that ends up with 0 cells
// has no room; its separator stays set aside.
//
// ComputeSplits fails with [ErrInvalidSize] when an entry does not parse.
func (s Sizes) ComputeSplits(extent int) ([]int, error) {
	n := len(s)
	cells := make([]int, n)
	if n == 0 {
		return cells, nil
	}
	extent = max(extent, 0)
	room := max(extent-(n-1), 0)

	weights := make([]int, n)
	sumWeight := 0
	for i, size := range s {
		v, suffix, err := size.parse()
		if err != nil {
			return nil, err
		}
		switch suffix {
		case "c":
			cells[i] = min(v, room)
		case "%":
			cells[i] = min(extent*v/100, room)
		default:
			weights[i] = v
			sumWeight += v
			continue
		}
		room -= cells[i]
	}

	last, assigned := -1, 0
	for i, w := range weights {
		if w == 0 {
			continue
		}
		cells[i] = room * w / sumWeight
		assigned += cells[i]
		last = i
	}
	if last >= 0 {
		cells[last] += room - assigned
	}
	return cells, nil
}

// Size is a pane's part of its parent along the split axis:
//
//   - "5c" is an exact number of cells: columns for a horizontal split, lines
//     for a vertical one.
//   - "20%" is a percentage of the parent's extent, 1..100.
//   - "2" is a weight. Weighted panes share what is left after the cell and
//     percent panes, in proportion to their weights.
//
// Decoding canonicalizes through [Size.Canonicalize], so a decoded Size is always
// in one of the three forms above with no sign or leading zeros.
type Size string

// Canonicalize returns the canonical spelling of s. An empty Size is a weight of
// 1. N must be a positive integer, and a percent must be in 1..100.
func (s Size) Canonicalize() (Size, error) {
	n, suffix, err := s.parse()
	if err != nil {
		return "", err
	}
	return Size(strconv.Itoa(n) + suffix), nil
}

// parse splits s into its magnitude and its suffix: "c", "%", or "" for a
// weight. An empty Size is a weight of 1.
func (s Size) parse() (int, string, error) {
	if s == "" {
		return 1, "", nil
	}
	body, suffix := string(s), ""
	switch {
	case strings.HasSuffix(body, "c"):
		body, suffix = strings.TrimSuffix(body, "c"), "c"
	case strings.HasSuffix(body, "%"):
		body, suffix = strings.TrimSuffix(body, "%"), "%"
	}
	n, err := strconv.Atoi(body)
	if err != nil {
		return 0, "", fmt.Errorf("%w: %q is not a number", ErrInvalidSize, string(s))
	}
	if n <= 0 {
		return 0, "", fmt.Errorf("%w: %q must be > 0", ErrInvalidSize, string(s))
	}
	if suffix == "%" && n > 100 {
		return 0, "", fmt.Errorf("%w: %q must be in 1..100", ErrInvalidSize, string(s))
	}
	return n, suffix, nil
}

// UnmarshalYAML decodes a scalar through [Size.Canonicalize]. A bare number is a
// weight. A node that is not a scalar fails with the decoder's own error, not
// with [ErrInvalidSize]: that error means a string that is not a size.
func (s *Size) UnmarshalYAML(node *yaml.Node) error {
	var str string
	if err := node.Decode(&str); err != nil {
		return err
	}
	return s.set(str)
}

// UnmarshalJSONFrom decodes a string or a number through [Size.Canonicalize]. A
// number is a weight. Any other value fails with the decoder's own error, not
// with [ErrInvalidSize].
func (s *Size) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == '0' {
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		return s.set(tok.String())
	}
	var str string
	if err := json.UnmarshalDecode(dec, &str); err != nil {
		return err
	}
	return s.set(str)
}

func (s *Size) set(str string) error {
	size, err := Size(str).Canonicalize()
	if err != nil {
		return err
	}
	*s = size
	return nil
}
