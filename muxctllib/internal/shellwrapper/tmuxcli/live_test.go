package tmuxcli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

func liveShell(t *testing.T) *ShellOut {
	t.Helper()
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not found:", err)
	}
	var env []string
	for _, kv := range os.Environ() {
		// An inherited $TMUX would point commands without -t at the enclosing session.
		if strings.HasPrefix(kv, "TMUX=") || strings.HasPrefix(kv, "TMUX_PANE=") {
			continue
		}
		env = append(env, kv)
	}
	s := &ShellOut{
		Global: GlobalOption{
			SocketName: fmt.Sprintf("tmuxcli-test-%d", os.Getpid()),
			// The user's configuration could change indices, sizes and defaults.
			ConfigFile: "/dev/null",
		},
		Env: env,
	}
	t.Cleanup(func() {
		ctx := context.Background()
		sock, err := s.DisplayMessage(ctx, DisplayMessageOption{Print: true}, "#{socket_path}")
		if err != nil {
			return
		}
		_ = s.KillServer(ctx)
		// kill-server leaves the socket file behind.
		_ = os.Remove(sock)
	})
	return s
}

func paneById(t *testing.T, s *ShellOut, opt ListPanesOption, id string) Pane {
	t.Helper()
	panes, err := s.ListPanes(t.Context(), opt)
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(panes, func(p Pane) bool { return p.Id == id })
	if i < 0 {
		t.Fatalf("pane %s not in %+v", id, panes)
	}
	return panes[i]
}

func windowById(t *testing.T, s *ShellOut, id string) Window {
	t.Helper()
	windows, err := s.ListWindows(t.Context(), ListWindowsOption{All: true})
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(windows, func(w Window) bool { return w.Id == id })
	if i < 0 {
		t.Fatalf("window %s not in %+v", id, windows)
	}
	return windows[i]
}

func TestLive(t *testing.T) {
	s := liveShell(t)
	ctx := t.Context()

	v, err := s.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v.Raw == "" {
		t.Errorf("Version = %+v, want Raw set", v)
	}

	if has, err := s.HasSession(ctx, HasSessionOption{TargetSession: "main"}); err != nil || has {
		t.Fatalf("HasSession without server = %t, %v; want false, nil", has, err)
	}

	sess, err := s.NewSession(ctx, NewSessionOption{
		Detached:    true,
		SessionName: "main",
		WindowName:  "w0",
		Width:       "200",
		Height:      "50",
	}, "cat")
	if err != nil {
		t.Fatal(err)
	}
	if sess.SessionName != "main" || !strings.HasPrefix(sess.SessionId, "$") ||
		!strings.HasPrefix(sess.WindowId, "@") || sess.WindowIndex != 0 || !strings.HasPrefix(sess.PaneId, "%") {
		t.Fatalf("NewSession = %+v", sess)
	}

	if has, err := s.HasSession(ctx, HasSessionOption{TargetSession: "main"}); err != nil || !has {
		t.Errorf("HasSession(main) = %t, %v; want true, nil", has, err)
	}
	if has, err := s.HasSession(ctx, HasSessionOption{TargetSession: "nope"}); err != nil || has {
		t.Errorf("HasSession(nope) = %t, %v; want false, nil", has, err)
	}

	win, err := s.NewWindow(ctx, NewWindowOption{
		Detached:     true,
		Environment:  []string{"TMUXCLI_TEST=1"},
		WindowName:   "w1",
		TargetWindow: "main:",
	}, "cat")
	if err != nil {
		t.Fatal(err)
	}
	if win.SessionId != sess.SessionId || win.WindowIndex != 1 || win.WindowId == sess.WindowId {
		t.Fatalf("NewWindow = %+v", win)
	}

	existing, err := s.NewWindow(ctx, NewWindowOption{
		Detached:       true,
		SelectExisting: true,
		WindowName:     "w1",
		TargetWindow:   "main",
	})
	if err != nil || existing != (SpawnResult{}) {
		t.Errorf("NewWindow selecting existing = %+v, %v; want zero, nil", existing, err)
	}

	split, err := s.SplitWindow(ctx, SplitWindowOption{
		Detached:   true,
		Horizontal: true,
		TargetPane: win.PaneId,
	}, "cat")
	if err != nil {
		t.Fatal(err)
	}
	if split.WindowId != win.WindowId || split.PaneId == win.PaneId {
		t.Fatalf("SplitWindow = %+v", split)
	}

	float, err := s.NewPane(ctx, NewPaneOption{
		Detached:   true,
		TargetPane: win.PaneId,
		Width:      "20",
		Height:     "10",
		X:          "5",
		Y:          "3",
	}, "cat")
	if err != nil {
		t.Fatal(err)
	}
	if float.WindowId != win.WindowId {
		t.Fatalf("NewPane = %+v", float)
	}

	paneScope := SetOptionOption{Pane: true, TargetPane: float.PaneId}
	tag := `a "b" \c`
	if err := s.SetOption(ctx, paneScope, "@tag", &tag); err != nil {
		t.Fatal(err)
	}
	if err := s.SetOption(ctx, paneScope, "@empty", new("")); err != nil {
		t.Fatal(err)
	}
	showScope := ShowOptionsOption{Pane: true, TargetPane: float.PaneId}
	entries, err := s.ShowOptions(ctx, showScope, "")
	if want := []OptionEntry{{Name: "@empty"}, {Name: "@tag", Value: tag}}; err != nil || !slices.Equal(entries, want) {
		t.Errorf("ShowOptions = %+v, %v; want %+v", entries, err, want)
	}
	valueOnly := showScope
	valueOnly.ValueOnly = true
	entries, err = s.ShowOptions(ctx, valueOnly, "@tag")
	if want := []OptionEntry{{Name: "@tag", Value: tag}}; err != nil || !slices.Equal(entries, want) {
		t.Errorf("ShowOptions -v = %+v, %v; want %+v", entries, err, want)
	}
	unset := paneScope
	unset.Unset = true
	if err := s.SetOption(ctx, unset, "@empty", nil); err != nil {
		t.Fatal(err)
	}
	quiet := showScope
	quiet.Quiet = true
	if entries, err := s.ShowOptions(ctx, quiet, "@empty"); err != nil || len(entries) != 0 {
		t.Errorf("ShowOptions of unset option = %+v, %v; want none, nil", entries, err)
	}

	inWindow := ListPanesOption{Target: win.WindowId, UserOptions: []string{"@tag"}}
	fp := paneById(t, s, inWindow, float.PaneId)
	if !fp.Floating || fp.Width != 20 || fp.Height != 10 || fp.Left != 5 || fp.Top != 3 ||
		fp.SessionName != "main" || fp.WindowId != win.WindowId || fp.WindowIndex != 1 ||
		fp.UserOptions["@tag"] != tag {
		t.Errorf("floating pane = %+v", fp)
	}
	sp := paneById(t, s, inWindow, split.PaneId)
	if sp.Floating || sp.Dead || sp.Pid == 0 || sp.UserOptions["@tag"] != "" {
		t.Errorf("split pane = %+v", sp)
	}
	if _, err := s.ListPanes(ctx, ListPanesOption{UserOptions: []string{"tag"}}); err == nil {
		t.Error("ListPanes with a user option lacking '@': want error")
	}
	for _, opt := range []ListPanesOption{{Session: true, Target: "main"}, {All: true}} {
		if panes, err := s.ListPanes(ctx, opt); err != nil || len(panes) != 4 {
			t.Errorf("ListPanes(%+v) = %d panes, %v; want 4", opt, len(panes), err)
		}
	}

	windows, err := s.ListWindows(ctx, ListWindowsOption{TargetSession: "main"})
	if err != nil {
		t.Fatal(err)
	}
	if len(windows) != 2 || windows[0].Id != sess.WindowId || windows[0].Name != "w0" ||
		windows[1].Id != win.WindowId || windows[1].Name != "w1" || windows[1].SessionName != "main" ||
		windows[1].SessionId != sess.SessionId || windows[1].Width != 200 {
		t.Errorf("ListWindows = %+v", windows)
	}

	sessions, err := s.ListSessions(ctx, ListSessionsOption{})
	if err != nil || len(sessions) != 1 || sessions[0].Name != "main" || sessions[0].Windows != 2 {
		t.Errorf("ListSessions = %+v, %v", sessions, err)
	}

	if clients, err := s.ListClients(ctx, ListClientsOption{}); err != nil || len(clients) != 0 {
		t.Errorf("ListClients = %+v, %v; want none, nil", clients, err)
	}

	msg, err := s.DisplayMessage(ctx, DisplayMessageOption{Print: true, TargetPane: split.PaneId}, "#{pane_id} #{window_id}")
	if want := split.PaneId + " " + win.WindowId; err != nil || msg != want {
		t.Errorf("DisplayMessage = %q, %v; want %q", msg, err, want)
	}
	msg, err = s.DisplayMessage(ctx, DisplayMessageOption{Print: true, Format: "#{session_name}", TargetPane: split.PaneId}, "")
	if err != nil || msg != "main" {
		t.Errorf("DisplayMessage -F = %q, %v; want %q", msg, err, "main")
	}
	msg, err = s.DisplayMessage(ctx, DisplayMessageOption{Print: true, Literal: true, TargetPane: split.PaneId}, "-x #{y}")
	if err != nil || msg != "-x #{y}" {
		t.Errorf("DisplayMessage -l = %q, %v; want %q", msg, err, "-x #{y}")
	}

	if err := s.SelectPane(ctx, SelectPaneOption{Title: "hello", TargetPane: split.PaneId}); err != nil {
		t.Fatal(err)
	}
	if got := paneById(t, s, inWindow, split.PaneId).Title; got != "hello" {
		t.Errorf("pane title = %q, want %q", got, "hello")
	}

	if err := s.SendKeys(ctx, SendKeysOption{TargetPane: split.PaneId}, "-x-marker", "Enter"); err != nil {
		t.Fatal(err)
	}
	var captured string
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		captured, err = s.CapturePane(ctx, CapturePaneOption{Print: true, TargetPane: split.PaneId})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(captured, "-x-marker") {
			break
		}
	}
	if !strings.Contains(captured, "-x-marker") {
		t.Errorf("CapturePane = %q, want it to contain the sent keys", captured)
	}

	if err := s.RespawnPane(ctx, RespawnPaneOption{Kill: true, TargetPane: split.PaneId}, "cat"); err != nil {
		t.Fatal(err)
	}
	if got := paneById(t, s, inWindow, split.PaneId).Pid; got == sp.Pid {
		t.Errorf("RespawnPane kept pid %d", got)
	}

	if err := s.ResizePane(ctx, ResizePaneOption{Zoom: true, TargetPane: split.PaneId}, 0); err != nil {
		t.Fatal(err)
	}
	if !windowById(t, s, win.WindowId).Zoomed {
		t.Error("window not zoomed after ResizePane -Z")
	}
	if err := s.ResizePane(ctx, ResizePaneOption{Zoom: true, TargetPane: split.PaneId}, 0); err != nil {
		t.Fatal(err)
	}
	if windowById(t, s, win.WindowId).Zoomed {
		t.Error("window still zoomed after second ResizePane -Z")
	}
	if err := s.ResizePane(ctx, ResizePaneOption{Left: true, TargetPane: split.PaneId}, 3); err != nil {
		t.Fatal(err)
	}
	if got := paneById(t, s, inWindow, split.PaneId).Width; got != sp.Width+3 && got != sp.Width-3 {
		t.Errorf("pane width after ResizePane -L 3 = %d, want %d changed by 3", got, sp.Width)
	}

	if err := s.SelectWindow(ctx, SelectWindowOption{TargetWindow: win.WindowId}); err != nil {
		t.Fatal(err)
	}
	if !windowById(t, s, win.WindowId).Active {
		t.Error("window not active after SelectWindow")
	}

	if err := s.KillPane(ctx, KillPaneOption{TargetPane: float.PaneId}); err != nil {
		t.Fatal(err)
	}
	if panes, err := s.ListPanes(ctx, ListPanesOption{Target: win.WindowId}); err != nil || len(panes) != 2 {
		t.Errorf("ListPanes after KillPane = %+v, %v; want 2 panes", panes, err)
	}
	err = s.KillPane(ctx, KillPaneOption{TargetPane: "%9999"})
	if cmdErr, ok := errors.AsType[*CmdError](err); !ok || cmdErr.ExitCode() != 1 || cmdErr.Stderr == "" {
		t.Errorf("KillPane of a missing pane = %v; want *CmdError with exit status 1 and stderr", err)
	}

	if err := s.KillWindow(ctx, KillWindowOption{TargetWindow: win.WindowId}); err != nil {
		t.Fatal(err)
	}
	if windows, err := s.ListWindows(ctx, ListWindowsOption{TargetSession: "main"}); err != nil || len(windows) != 1 {
		t.Errorf("ListWindows after KillWindow = %+v, %v; want 1 window", windows, err)
	}

	sock, err := s.DisplayMessage(ctx, DisplayMessageOption{Print: true, TargetPane: sess.PaneId}, "#{socket_path}")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.KillServer(ctx); err != nil {
		t.Fatal(err)
	}
	// kill-server leaves the socket file behind.
	if err := os.Remove(sock); err != nil {
		t.Error(err)
	}
	if has, err := s.HasSession(ctx, HasSessionOption{TargetSession: "main"}); err != nil || has {
		t.Errorf("HasSession after KillServer = %t, %v; want false, nil", has, err)
	}
}
