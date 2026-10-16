package zellijcli

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// fakeZellij returns a ShellOut whose zellij records its arguments, one per line,
// and answers with the given stdout, stderr and exit status.
func fakeZellij(t *testing.T, stdout, stderr string, exit int) (*ShellOut, func() []string) {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "zellij")
	body := "#!/bin/sh\n" +
		"printf '%s\\n' \"$@\" > \"$FAKE_ARGS\"\n" +
		"printf '%s' \"$FAKE_STDOUT\"\n" +
		"printf '%s' \"$FAKE_STDERR\" >&2\n" +
		"exit \"$FAKE_EXIT\"\n"
	if err := os.WriteFile(script, []byte(body), 0o700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	s := &ShellOut{
		Path: script,
		Env: []string{
			"FAKE_ARGS=" + argsFile,
			"FAKE_STDOUT=" + stdout,
			"FAKE_STDERR=" + stderr,
			"FAKE_EXIT=" + strconv.Itoa(exit),
		},
	}
	return s, func() []string {
		t.Helper()
		b, err := os.ReadFile(argsFile)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	}
}

func TestCmdError(t *testing.T) {
	s, _ := fakeZellij(t, "", "boom\n", 3)
	_, err := s.ActionListClients(t.Context())
	cmdErr, ok := errors.AsType[*CmdError](err)
	if !ok {
		t.Fatalf("err = %v, want *CmdError", err)
	}
	if cmdErr.ExitCode() != 3 {
		t.Errorf("ExitCode() = %d, want 3", cmdErr.ExitCode())
	}
	if cmdErr.Stderr != "boom" {
		t.Errorf("Stderr = %q, want %q", cmdErr.Stderr, "boom")
	}
	if want := "zellij action list-clients: exit status 3: boom"; cmdErr.Error() != want {
		t.Errorf("Error() = %q, want %q", cmdErr.Error(), want)
	}
}

func TestParsePaneId(t *testing.T) {
	tests := []struct {
		in      string
		want    PaneId
		wantErr bool
	}{
		{in: "terminal_1", want: PaneId{Id: 1}},
		{in: "plugin_2", want: PaneId{Plugin: true, Id: 2}},
		{in: "3", want: PaneId{Id: 3}},
		{in: "terminal_0", want: PaneId{Id: 0}},
		{in: "", wantErr: true},
		{in: "terminal_", wantErr: true},
		{in: "pane_1", wantErr: true},
		{in: "-1", wantErr: true},
		{in: "+1", wantErr: true},
		{in: "plugin_x", wantErr: true},
		{in: "terminal_1 ", wantErr: true},
	}
	for _, tt := range tests {
		got, err := ParsePaneId(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParsePaneID(%q) err = %v, wantErr %t", tt.in, err, tt.wantErr)
			continue
		}
		if err != nil {
			continue
		}
		if got != tt.want {
			t.Errorf("ParsePaneID(%q) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestPaneIdString(t *testing.T) {
	for _, id := range []PaneId{{Id: 0}, {Id: 12}, {Plugin: true, Id: 3}} {
		back, err := ParsePaneId(id.String())
		if err != nil || back != id {
			t.Errorf("ParsePaneID(%q) = %+v, %v; want %+v", id.String(), back, err, id)
		}
	}
	if got := (PaneId{Plugin: true, Id: 3}).String(); got != "plugin_3" {
		t.Errorf("String() = %q, want plugin_3", got)
	}
	if got := (PaneId{Id: 7}).String(); got != "terminal_7" {
		t.Errorf("String() = %q, want terminal_7", got)
	}
}
