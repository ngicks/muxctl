package zellijcli

import (
	"context"
	"encoding/json/v2"
	"fmt"
)

// Get information about the currently active tab
//
// Returns: Tab name and ID by default, or full info in JSON
//
// Usage: zellij action current-tab-info [OPTIONS]
//
// Options:
//   -j, --json
//           Output as JSON with full TabInfo
//
//   -h, --help
//           Print help (see a summary with '-h')

// ActionCurrentTabInfo runs action current-tab-info --json.
// zellij resolves the tab through the client issuing the action, so outside
// an attached client it fails with a [*CmdError] ("No active tab found for current client").
func (s *ShellOut) ActionCurrentTabInfo(ctx context.Context) (Tab, error) {
	out, err := s.run(ctx, "action", "current-tab-info", "--json")
	if err != nil {
		return Tab{}, err
	}
	var tab Tab
	if err := json.Unmarshal([]byte(out), &tab); err != nil {
		return Tab{}, fmt.Errorf("zellij: decode current-tab-info: %w", err)
	}
	return tab, nil
}
