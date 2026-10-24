// Package types holds the values muxctllib passes between its packages and
// the terminal multiplexers.
package types

import "context"

// Pane is a pane in a window.
type Pane struct {
	// Id identifies the pane to the multiplexer.
	Id string
	// Width is the pane's width in columns.
	Width int
	// Height is the pane's height in lines.
	Height int
	// Floating reports whether the pane floats above the tiled panes.
	Floating bool
	// Focused reports whether the pane is the window's active pane.
	Focused bool
	// Options holds the pane's own options, such as tmux's per-pane user
	// options. It is nil when the multiplexer has no per-pane options, e.g.
	// zellij.
	Options map[string]string
}

// PaneLister lists the panes of a window.
type PaneLister interface {
	Panes(ctx context.Context) ([]Pane, error)
}

// PaneResizer resizes a pane along one axis.
type PaneResizer interface {
	// Resize makes pane cells wide when horizontal is true and cells high
	// otherwise.
	Resize(ctx context.Context, pane string, horizontal bool, cells int) error
}

// PaneFocuser focuses a pane.
type PaneFocuser interface {
	// Focus makes pane the window's active pane.
	Focus(ctx context.Context, pane string) error
}
