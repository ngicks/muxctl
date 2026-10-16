package zellijcli

import (
	"context"
	"encoding/json/v2"
	"fmt"
)

// List all tabs with their information
//
// Returns: Tab information in table or JSON format
//
// Usage: zellij action list-tabs [OPTIONS]
//
// Options:
//   -s, --state
//           Include state information (active, fullscreen, sync, floating visibility)
//
//   -d, --dimensions
//           Include dimension information (viewport, display area)
//
//   -p, --panes
//           Include pane counts
//
//   -l, --layout
//           Include layout information (swap layout name and dirty state)
//
//   -a, --all
//           Include all available fields
//
//   -j, --json
//           Output as JSON
//
//   -h, --help
//           Print help (see a summary with '-h')

// Tab is one entry of action list-tabs --json.
// zellij's help says current-tab-info --json prints the same TabInfo.
type Tab struct {
	Position                int    `json:"position"`
	Name                    string `json:"name"`
	Active                  bool   `json:"active"`
	PanesToHide             int    `json:"panes_to_hide"`
	IsFullscreenActive      bool   `json:"is_fullscreen_active"`
	IsSyncPanesActive       bool   `json:"is_sync_panes_active"`
	AreFloatingPanesVisible bool   `json:"are_floating_panes_visible"`
	OtherFocusedClients     []int  `json:"other_focused_clients"`
	// ActiveSwapLayoutName is nil for a tab that no client has rendered yet.
	ActiveSwapLayoutName         *string `json:"active_swap_layout_name"`
	IsSwapLayoutDirty            bool    `json:"is_swap_layout_dirty"`
	ViewportRows                 int     `json:"viewport_rows"`
	ViewportColumns              int     `json:"viewport_columns"`
	DisplayAreaRows              int     `json:"display_area_rows"`
	DisplayAreaColumns           int     `json:"display_area_columns"`
	SelectableTiledPanesCount    int     `json:"selectable_tiled_panes_count"`
	SelectableFloatingPanesCount int     `json:"selectable_floating_panes_count"`
	TabId                        int     `json:"tab_id"`
	HasBellNotification          bool    `json:"has_bell_notification"`
	IsFlashingBell               bool    `json:"is_flashing_bell"`
}

// ActionListTabs runs action list-tabs --all --json.
// The other flags only pick columns of the table output, so they are not exposed.
func (s *ShellOut) ActionListTabs(ctx context.Context) ([]Tab, error) {
	out, err := s.run(ctx, "action", "list-tabs", "--all", "--json")
	if err != nil {
		return nil, err
	}
	var tabs []Tab
	if err := json.Unmarshal([]byte(out), &tabs); err != nil {
		return nil, fmt.Errorf("zellij: decode list-tabs: %w", err)
	}
	return tabs, nil
}
