package zellijcli

import (
	"encoding/json/jsontext"
	"errors"
	"reflect"
	"testing"
)

// panesJson is trimmed from real "zellij action list-panes --all --json" output (zellij 0.45.1).
const panesJson = `[
  {
    "id": 1, "is_plugin": true, "is_focused": false, "is_fullscreen": false, "is_floating": false,
    "is_suppressed": false, "title": "tab-bar", "exited": false, "exit_status": null, "is_held": false,
    "pane_x": 0, "pane_content_x": 0, "pane_y": 0, "pane_content_y": 0,
    "pane_rows": 1, "pane_content_rows": 1, "pane_columns": 50, "pane_content_columns": 50,
    "cursor_coordinates_in_pane": null, "terminal_command": null, "plugin_url": "tab-bar",
    "is_selectable": false, "index_in_pane_group": {}, "default_fg": null, "default_bg": null,
    "tab_id": 0, "tab_position": 0, "tab_name": "Tab #1"
  },
  {
    "id": 13, "is_plugin": false, "is_focused": true, "is_fullscreen": false, "is_floating": false,
    "is_suppressed": false, "title": "sleep 300", "exited": false, "exit_status": null, "is_held": false,
    "pane_x": 25, "pane_content_x": 25, "pane_y": 17, "pane_content_y": 18,
    "pane_rows": 16, "pane_content_rows": 14, "pane_columns": 25, "pane_content_columns": 25,
    "cursor_coordinates_in_pane": [0, 1], "terminal_command": "sleep 300", "plugin_url": null,
    "is_selectable": true, "index_in_pane_group": {}, "default_fg": "#ff0000", "default_bg": "#00ff00",
    "tab_id": 0, "tab_position": 0, "tab_name": "Tab #1",
    "pane_command": "sleep 300", "pane_cwd": "/home/user"
  },
  {
    "id": 14, "is_plugin": false, "is_focused": false, "is_fullscreen": false, "is_floating": false,
    "is_suppressed": false, "title": "sh -c exit 3", "exited": true, "exit_status": 3, "is_held": true,
    "pane_x": 25, "pane_content_x": 25, "pane_y": 33, "pane_content_y": 34,
    "pane_rows": 16, "pane_content_rows": 15, "pane_columns": 25, "pane_content_columns": 25,
    "cursor_coordinates_in_pane": [0, 1], "terminal_command": "sh -c exit 3", "plugin_url": null,
    "is_selectable": true, "index_in_pane_group": {}, "default_fg": null, "default_bg": null,
    "tab_id": 0, "tab_position": 0, "tab_name": "Tab #1"
  }
]`

func TestActionListPanesDecode(t *testing.T) {
	s, _ := fakeZellij(t, panesJson, "", 0)
	got, err := s.ActionListPanes(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	want := []Pane{
		{
			Id: 1, IsPlugin: true, Title: "tab-bar",
			PaneRows: 1, PaneContentRows: 1, PaneColumns: 50, PaneContentColumns: 50,
			PluginUrl: new("tab-bar"), IndexInPaneGroup: jsontext.Value("{}"), TabName: "Tab #1",
		},
		{
			Id: 13, IsFocused: true, Title: "sleep 300",
			PaneX: 25, PaneContentX: 25, PaneY: 17, PaneContentY: 18,
			PaneRows: 16, PaneContentRows: 14, PaneColumns: 25, PaneContentColumns: 25,
			CursorCoordinatesInPane: &[2]int{0, 1}, TerminalCommand: new("sleep 300"), IsSelectable: true,
			IndexInPaneGroup: jsontext.Value("{}"), DefaultFg: new("#ff0000"), DefaultBg: new("#00ff00"),
			TabName: "Tab #1", PaneCommand: "sleep 300", PaneCwd: "/home/user",
		},
		{
			Id: 14, Title: "sh -c exit 3", Exited: true, ExitStatus: new(3), IsHeld: true,
			PaneX: 25, PaneContentX: 25, PaneY: 33, PaneContentY: 34,
			PaneRows: 16, PaneContentRows: 15, PaneColumns: 25, PaneContentColumns: 25,
			CursorCoordinatesInPane: &[2]int{0, 1}, TerminalCommand: new("sh -c exit 3"), IsSelectable: true,
			IndexInPaneGroup: jsontext.Value("{}"), TabName: "Tab #1",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
	if id := got[0].PaneId(); id != (PaneId{Plugin: true, Id: 1}) {
		t.Errorf("PaneID() = %+v", id)
	}
	if id := got[1].PaneId(); id != (PaneId{Id: 13}) {
		t.Errorf("PaneID() = %+v", id)
	}
}

// tabJson is one entry of real "zellij action list-tabs --all --json" output (zellij 0.45.1).
const tabJson = `{
  "position": 1, "name": "t2", "active": false, "panes_to_hide": 0,
  "is_fullscreen_active": false, "is_sync_panes_active": false, "are_floating_panes_visible": true,
  "other_focused_clients": [2], "active_swap_layout_name": "BASE", "is_swap_layout_dirty": true,
  "viewport_rows": 48, "viewport_columns": 50, "display_area_rows": 50, "display_area_columns": 50,
  "selectable_tiled_panes_count": 1, "selectable_floating_panes_count": 1, "tab_id": 1,
  "has_bell_notification": false, "is_flashing_bell": false
}`

var wantTab = Tab{
	Position: 1, Name: "t2", AreFloatingPanesVisible: true,
	OtherFocusedClients: []int{2}, ActiveSwapLayoutName: new("BASE"), IsSwapLayoutDirty: true,
	ViewportRows: 48, ViewportColumns: 50, DisplayAreaRows: 50, DisplayAreaColumns: 50,
	SelectableTiledPanesCount: 1, SelectableFloatingPanesCount: 1, TabId: 1,
}

func TestActionListTabsDecode(t *testing.T) {
	s, _ := fakeZellij(t, "["+tabJson+"]", "", 0)
	got, err := s.ActionListTabs(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if want := []Tab{wantTab}; !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestActionCurrentTabInfoDecode(t *testing.T) {
	s, _ := fakeZellij(t, tabJson, "", 0)
	got, err := s.ActionCurrentTabInfo(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, wantTab) {
		t.Errorf("got  %+v\nwant %+v", got, wantTab)
	}
}

func TestDecodeErrorIsNotCmdError(t *testing.T) {
	s, _ := fakeZellij(t, "not json", "", 0)
	_, err := s.ActionListPanes(t.Context())
	if err == nil {
		t.Fatal("want error")
	}
	if _, ok := errors.AsType[*CmdError](err); ok {
		t.Errorf("err = %v, want a decode error", err)
	}
}
