package zellijcli

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// List all panes in the current session
//
// Returns: Formatted list of panes (table or JSON) to stdout
//
// Usage: zellij action list-panes [OPTIONS]
//
// Options:
//   -t, --tab
//           Include tab information (name, position, ID)
//
//   -c, --command
//           Include running command information
//
//   -s, --state
//           Include pane state (focused, floating, exited, etc.)
//
//   -g, --geometry
//           Include geometry (position, size)
//
//   -a, --all
//           Include all available fields
//
//   -j, --json
//           Output as JSON
//
//   -h, --help
//           Print help (see a summary with '-h')

// Pane is one entry of action list-panes --json.
type Pane struct {
	Id                      int     `json:"id"`
	IsPlugin                bool    `json:"is_plugin"`
	IsFocused               bool    `json:"is_focused"`
	IsFullscreen            bool    `json:"is_fullscreen"`
	IsFloating              bool    `json:"is_floating"`
	IsSuppressed            bool    `json:"is_suppressed"`
	Title                   string  `json:"title"`
	Exited                  bool    `json:"exited"`
	ExitStatus              *int    `json:"exit_status"`
	IsHeld                  bool    `json:"is_held"`
	PaneX                   int     `json:"pane_x"`
	PaneContentX            int     `json:"pane_content_x"`
	PaneY                   int     `json:"pane_y"`
	PaneContentY            int     `json:"pane_content_y"`
	PaneRows                int     `json:"pane_rows"`
	PaneContentRows         int     `json:"pane_content_rows"`
	PaneColumns             int     `json:"pane_columns"`
	PaneContentColumns      int     `json:"pane_content_columns"`
	CursorCoordinatesInPane *[2]int `json:"cursor_coordinates_in_pane"`
	TerminalCommand         *string `json:"terminal_command"`
	PluginUrl               *string `json:"plugin_url"`
	IsSelectable            bool    `json:"is_selectable"`
	// IndexInPaneGroup is kept raw: its shape is undocumented and observed only as {}.
	IndexInPaneGroup jsontext.Value `json:"index_in_pane_group"`
	DefaultFg        *string        `json:"default_fg"`
	DefaultBg        *string        `json:"default_bg"`
	TabId            int            `json:"tab_id"`
	TabPosition      int            `json:"tab_position"`
	TabName          string         `json:"tab_name"`
	// PaneCommand and PaneCwd are absent for plugins and for terminals whose process exited.
	PaneCommand string `json:"pane_command,omitzero"`
	PaneCwd     string `json:"pane_cwd,omitzero"`
}

// ActionListPanes runs action list-panes --all --json.
// The other flags only pick columns of the table output, so they are not exposed.
func (s *ShellOut) ActionListPanes(ctx context.Context) ([]Pane, error) {
	out, err := s.run(ctx, "action", "list-panes", "--all", "--json")
	if err != nil {
		return nil, err
	}
	var panes []Pane
	if err := json.Unmarshal([]byte(out), &panes); err != nil {
		return nil, fmt.Errorf("zellij: decode list-panes: %w", err)
	}
	return panes, nil
}

func (p Pane) PaneId() PaneId {
	return PaneId{Plugin: p.IsPlugin, Id: p.Id}
}
