package zellijcli

import (
	"context"
	"strconv"
)

// Close tab with stable ID
//
// Usage: zellij action close-tab-by-id <ID>
//
// Arguments:
//   <ID>
//
// Options:
//   -h, --help  Print help

// ActionCloseTabById runs action close-tab-by-id.
// zellij reports success even if the tab does not exist.
func (s *ShellOut) ActionCloseTabById(ctx context.Context, id int) error {
	_, err := s.run(ctx, "action", "close-tab-by-id", strconv.Itoa(id))
	return err
}
