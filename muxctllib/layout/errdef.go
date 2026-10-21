package layout

import "errors"

// Errors returned by this package.
var (
	ErrInvalidDir     = errors.New("layout: invalid dir")
	ErrInvalidSize    = errors.New("layout: invalid size")
	ErrInvalidCommand = errors.New("layout: command must be a sequence or a mapping")
	ErrSplitXorLeaf   = errors.New("layout: a pane sets either panes or command")
	ErrSizesMismatch  = errors.New("layout: sizes must have one entry per pane")
)
