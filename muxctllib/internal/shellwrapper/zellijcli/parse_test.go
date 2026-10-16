package zellijcli

import (
	"reflect"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in   string
		want Version
	}{
		{in: "zellij 0.45.1\n", want: Version{Raw: "zellij 0.45.1", Major: 0, Minor: 45, Patch: 1}},
		{in: "zellij 1.2.3-rc.1\n", want: Version{Raw: "zellij 1.2.3-rc.1", Major: 1, Minor: 2, Patch: 3}},
		{in: "zellij 1.2.3+build\n", want: Version{Raw: "zellij 1.2.3+build", Major: 1, Minor: 2, Patch: 3}},
		{in: "zellij 1.2\n", want: Version{Raw: "zellij 1.2"}},
		{in: "zellij x.y.z\n", want: Version{Raw: "zellij x.y.z"}},
		{in: "", want: Version{}},
	}
	for _, tt := range tests {
		if got := parseVersion(tt.in); got != tt.want {
			t.Errorf("parseVersion(%q) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestParseSessions(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		short   bool
		want    []Session
		wantErr bool
	}{
		{
			name: "suffixes",
			in: "plain [Created 1m 22s ago] \n" +
				"cur [Created 0s ago] (current)\n" +
				"dead [Created 3h 4m 5s ago] (EXITED - attach to resurrect)\n",
			want: []Session{
				{Name: "plain", Created: "1m 22s"},
				{Name: "cur", Created: "0s", Current: true},
				{Name: "dead", Created: "3h 4m 5s", Exited: true},
			},
		},
		{
			name: "name containing the marker",
			in:   "a [Created b [Created 9s ago] \n",
			want: []Session{{Name: "a [Created b", Created: "9s"}},
		},
		{
			name:  "short",
			in:    "a\nb c\n",
			short: true,
			want:  []Session{{Name: "a"}, {Name: "b c"}},
		},
		{
			name: "empty",
			in:   "",
		},
		{
			name:    "no created marker",
			in:      "nonsense\n",
			wantErr: true,
		},
		{
			name:    "unknown suffix",
			in:      "a [Created 3s ago] (weird)\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSessions(tt.in, tt.short)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %t", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestListSessionsNoSessions(t *testing.T) {
	tests := []struct {
		name    string
		stderr  string
		exit    int
		wantErr bool
	}{
		{name: "none", stderr: "No active zellij sessions found.\n", exit: 1},
		{name: "other failure with the same status", stderr: "something else\n", exit: 1, wantErr: true},
		{name: "same message with another status", stderr: "No active zellij sessions found.\n", exit: 2, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := fakeZellij(t, "", tt.stderr, tt.exit)
			got, err := s.ListSessions(t.Context(), ListSessionsOption{})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %t", err, tt.wantErr)
			}
			if got != nil {
				t.Errorf("sessions = %+v, want nil", got)
			}
		})
	}
}

func TestParseClients(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []Client
		wantErr bool
	}{
		{
			name: "padded columns",
			in: "CLIENT_ID ZELLIJ_PANE_ID RUNNING_COMMAND\n" +
				"1         terminal_12    sleep 300      \n" +
				"2         plugin_3       zellij:about\n",
			want: []Client{
				{Id: 1, PaneId: "terminal_12", RunningCommand: "sleep 300"},
				{Id: 2, PaneId: "plugin_3", RunningCommand: "zellij:about"},
			},
		},
		{
			name: "header only",
			in:   "CLIENT_ID ZELLIJ_PANE_ID RUNNING_COMMAND\n",
		},
		{
			name:    "bad client id",
			in:      "CLIENT_ID ZELLIJ_PANE_ID RUNNING_COMMAND\nx terminal_1 sh\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseClients(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %t", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseCreatedPaneId(t *testing.T) {
	tests := []struct {
		in      string
		want    *PaneId
		wantErr bool
	}{
		{in: "terminal_1\n", want: &PaneId{Id: 1}},
		{in: "plugin_10\n", want: &PaneId{Plugin: true, Id: 10}},
		{in: "", want: nil},
		{in: "\n", want: nil},
		{in: "junk\n", wantErr: true},
	}
	for _, tt := range tests {
		got, err := parseCreatedPaneId(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseCreatedPaneID(%q) err = %v, wantErr %t", tt.in, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseCreatedPaneID(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestActionNewTabId(t *testing.T) {
	s, _ := fakeZellij(t, "7\n", "", 0)
	id, err := s.ActionNewTab(t.Context(), ActionNewTabOption{})
	if err != nil || id != 7 {
		t.Errorf("ActionNewTab() = %d, %v; want 7, nil", id, err)
	}

	s, _ = fakeZellij(t, "", "", 0)
	if _, err := s.ActionNewTab(t.Context(), ActionNewTabOption{}); err == nil {
		t.Error("ActionNewTab() with empty stdout: want error")
	}
}

func TestActionNewPaneNoId(t *testing.T) {
	s, _ := fakeZellij(t, "", "", 0)
	id, err := s.ActionNewPane(t.Context(), ActionNewPaneOption{
		InPlace:           true,
		CloseReplacedPane: true,
		PaneId:            &PaneId{Id: 1},
	}, "sh")
	if err != nil || id != nil {
		t.Errorf("ActionNewPane() = %v, %v; want nil, nil", id, err)
	}
}
