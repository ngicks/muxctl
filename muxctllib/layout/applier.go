package layout

import (
	"context"
	"errors"

	"github.com/ngicks/muxctl/muxctllib/types"
)

// Applier builds a [Layout] in a terminal multiplexer.
type Applier interface {
	// Apply arranges the panes of a single window into layout. It does not create
	// or switch windows, tabs or sessions.
	Apply(ctx context.Context, layout Layout, opt *ApplyOption) error
}

// ApplyOption holds the options of [Applier.Apply].
type ApplyOption struct {
	// KeepFocus moves focus back to the pane that had it before Apply. When a
	// leaf sets [LeafSpec.Focus], Apply still focuses that leaf once before
	// moving back, because the app in it may be waiting for a focus event.
	KeepFocus bool
	// UnmanagedPane says what Apply does with the panes the window already has.
	UnmanagedPane UnmanagedPane
}

// UnmanagedPane is how [Applier.Apply] treats a pane that is already in the
// window and is not part of the layout.
type UnmanagedPane string

const (
	// UnmanagedPaneCloseAll closes every unmanaged pane.
	UnmanagedPaneCloseAll UnmanagedPane = "close_all"
	// UnmanagedPaneCloseTiled closes unmanaged tiled panes and leaves unmanaged
	// floating panes as they are. The empty UnmanagedPane means the same.
	UnmanagedPaneCloseTiled UnmanagedPane = "close_tiled"
	// UnmanagedPaneFail fails on any unmanaged pane.
	UnmanagedPaneFail UnmanagedPane = "fail"
)

// Splitter is what [SplitterApplier] needs from a terminal multiplexer.
type Splitter interface {
	types.PaneLister
	types.PaneResizer
	types.PaneFocuser
	// Split splits pane along dir and returns the new pane. The new pane is
	// cells wide for a horizontal split and cells high for a vertical one,
	// placed after pane, and runs leaf's command. Split sets options on the
	// new pane, so a later Panes reports them in [types.Pane.Options]; a
	// multiplexer without per-pane options ignores them.
	Split(
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
	Splitter Splitter
}

var _ Applier = SplitterApplier{}

// Apply is not implemented yet.
func (a SplitterApplier) Apply(ctx context.Context, l Layout, opt *ApplyOption) error {
	return errors.New("layout: SplitterApplier.Apply: not implemented")
}
