package layout

import "context"

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
