package zellijcli

import (
	"context"
	"strconv"
)

// Go to tab with stable ID
//
// Usage: zellij action go-to-tab-by-id <ID>
//
// Arguments:
//   <ID>
//
// Options:
//   -h, --help  Print help

// ActionGoToTabById runs action go-to-tab-by-id.
func (s *ShellOut) ActionGoToTabById(ctx context.Context, id int) error {
	_, err := s.run(ctx, "action", "go-to-tab-by-id", strconv.Itoa(id))
	return err
}
