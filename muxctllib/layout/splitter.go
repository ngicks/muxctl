package layout

import (
	"context"
	"errors"

	"github.com/ngicks/muxctl/muxctllib/types"
)

// PaneSplitter is what [SplitterApplier] needs from a terminal multiplexer.
type PaneSplitter interface {
	types.BasicMuxer
	// PaneSplit splits pane along dir and returns the new pane. The new pane is
	// cells wide for a horizontal split and cells high for a vertical one,
	// placed after pane, and runs leaf's command. PaneSplit sets options on the
	// new pane, so a later PaneList reports them in [types.Pane.Options]; a
	// multiplexer without per-pane options ignores them.
	PaneSplit(
		ctx context.Context,
		pane string,
		dir Dir,
		cells int,
		leaf LeafSpec,
		options map[string]string,
	) (string, error)
}

// PathOption is the pane option that records which layout leaf a pane was
// created for. Its value is the index of each pane from the root split down,
// e.g. "0.2.1". It lives in the pane's options rather than in [types.Pane]
// because no multiplexer knows about layout leaves natively.
const PathOption = "muxctl_layout_path"

// SplitterApplier is an [Applier] that builds a layout by splitting panes
// through Splitter.
type SplitterApplier struct {
	Splitter PaneSplitter
}

var _ Applier = SplitterApplier{}

// Apply is not implemented yet.
func (a SplitterApplier) Apply(ctx context.Context, l Layout, opt *ApplyOption) error {
	return errors.New("layout: SplitterApplier.Apply: not implemented")
}
