package zellijcli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func lookZellij(t *testing.T) string {
	t.Helper()
	path, err := exec.LookPath("zellij")
	if err != nil {
		t.Skip("zellij not found in PATH")
	}
	return path
}

// liveEnv is the current environment without the variables that would tie
// zellij to the session this test may run in, and with a plain shell for new panes.
func liveEnv(extra ...string) []string {
	var env []string
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "ZELLIJ") || strings.HasPrefix(kv, "SHELL=") {
			continue
		}
		env = append(env, kv)
	}
	env = append(env, "SHELL=/bin/sh")
	return append(env, extra...)
}

// waitFor polls cond, which zellij's asynchronous actions need, until it holds or 5 seconds pass.
func waitFor(t *testing.T, what string, cond func() (bool, error)) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		ok, err := cond()
		if ok {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s (last error: %v)", what, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func findPane(ctx context.Context, s *ShellOut, id PaneId) (*Pane, error) {
	panes, err := s.ActionListPanes(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range panes {
		if p.PaneId() == id {
			return &p, nil
		}
	}
	return nil, nil
}

func TestLiveListSessionsNone(t *testing.T) {
	path := lookZellij(t)
	s := &ShellOut{Path: path, Env: liveEnv("ZELLIJ_SOCKET_DIR="+t.TempDir(), "XDG_CACHE_HOME="+t.TempDir())}
	sessions, err := s.ListSessions(t.Context(), ListSessionsOption{})
	if err != nil || sessions != nil {
		t.Fatalf("ListSessions() = %+v, %v; want nil, nil", sessions, err)
	}
}

func TestLive(t *testing.T) {
	path := lookZellij(t)
	ctx := t.Context()

	// An empty config keeps zellij from opening its first-run setup pane.
	configDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(configDir, "config.kdl"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	env := liveEnv("ZELLIJ_CONFIG_DIR=" + configDir)
	name := fmt.Sprintf("zellijcli-test-%d", os.Getpid())

	s := &ShellOut{Path: path, Env: env}
	t.Cleanup(func() {
		// kill-session already removes a short-lived session, so delete-session failing is expected.
		ctx := context.Background()
		if err := s.KillSession(ctx, name); err != nil {
			t.Logf("kill-session: %v", err)
		}
		if err := s.DeleteSession(ctx, DeleteSessionOption{Force: true}, name); err != nil {
			t.Logf("delete-session: %v", err)
		}
	})

	v, err := s.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(v.Raw, "zellij ") || (v.Major == 0 && v.Minor == 0) {
		t.Errorf("Version() = %+v", v)
	}

	if out, err := s.Attach(ctx, AttachOption{CreateBackground: true}, name).CombinedOutput(); err != nil {
		t.Fatalf("attach: %v: %s", err, out)
	}
	waitFor(t, "session "+name, func() (bool, error) {
		sessions, err := s.ListSessions(ctx, ListSessionsOption{})
		return slices.ContainsFunc(sessions, func(ss Session) bool { return ss.Name == name && !ss.Exited && ss.Created != "" }), err
	})
	short, err := s.ListSessions(ctx, ListSessionsOption{Short: true})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(short, Session{Name: name}) {
		t.Errorf("ListSessions(Short) = %+v, want %s", short, name)
	}

	z := &ShellOut{Path: path, Env: env, Global: GlobalOption{Session: name}}
	waitFor(t, "the initial terminal pane", func() (bool, error) {
		panes, err := z.ActionListPanes(ctx)
		return slices.ContainsFunc(panes, func(p Pane) bool { return !p.IsPlugin }), err
	})

	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	catId, err := z.ActionNewPane(ctx, ActionNewPaneOption{Cwd: dir, Name: "cat"}, "cat")
	if err != nil {
		t.Fatal(err)
	}
	if catId == nil || catId.Plugin {
		t.Fatalf("ActionNewPane() = %v, want a terminal pane", catId)
	}
	waitFor(t, "the cat pane", func() (bool, error) {
		p, err := findPane(ctx, z, *catId)
		return p != nil && p.Title == "cat" && strings.Contains(p.PaneCommand, "cat") && p.PaneCwd == dir, err
	})

	if err := z.ActionRenamePane(ctx, ActionRenamePaneOption{PaneId: catId}, "renamed"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "rename-pane", func() (bool, error) {
		p, err := findPane(ctx, z, *catId)
		return p != nil && p.Title == "renamed", err
	})

	if err := z.ActionFocusPaneId(ctx, *catId); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "focus-pane-id", func() (bool, error) {
		p, err := findPane(ctx, z, *catId)
		return p != nil && p.IsFocused, err
	})

	for _, want := range []bool{true, false} {
		if err := z.ActionToggleFullscreen(ctx, ActionToggleFullscreenOption{PaneId: catId}); err != nil {
			t.Fatal(err)
		}
		waitFor(t, fmt.Sprintf("toggle-fullscreen to %t", want), func() (bool, error) {
			p, err := findPane(ctx, z, *catId)
			return p != nil && p.IsFullscreen == want, err
		})
	}

	if err := z.ActionResize(ctx, ActionResizeOption{PaneId: catId}, "increase", "left"); err != nil {
		t.Fatal(err)
	}

	const marker = "zellijcli-marker"
	if err := z.ActionWriteChars(ctx, ActionWriteCharsOption{PaneId: catId}, marker); err != nil {
		t.Fatal(err)
	}
	if err := z.ActionSendKeys(ctx, ActionSendKeysOption{PaneId: catId}, "Enter"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "dump-screen to show the written chars", func() (bool, error) {
		dump, err := z.ActionDumpScreen(ctx, ActionDumpScreenOption{PaneId: catId})
		return strings.Contains(dump, marker), err
	})
	dumpPath := filepath.Join(t.TempDir(), "dump.txt")
	if _, err := z.ActionDumpScreen(ctx, ActionDumpScreenOption{Path: dumpPath, Full: true, PaneId: catId}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "dump-screen --path", func() (bool, error) {
		b, err := os.ReadFile(dumpPath)
		return strings.Contains(string(b), marker), err
	})

	respawned, err := z.ActionNewPane(ctx, ActionNewPaneOption{
		InPlace:           true,
		CloseReplacedPane: true,
		PaneId:            catId,
	}, "sleep", "300")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("in-place new-pane replacing %v returned %v", catId, respawned)
	waitFor(t, "the replaced pane to close", func() (bool, error) {
		p, err := findPane(ctx, z, *catId)
		return p == nil, err
	})

	popup, err := z.Run(ctx, RunOption{
		Floating:    true,
		CloseOnExit: true,
		Name:        "popup",
		Pinned:      new(true),
		X:           "10%",
		Y:           "10%",
		Width:       "50%",
		Height:      "50%",
	}, "sleep", "300")
	if err != nil {
		t.Fatal(err)
	}
	if popup == nil {
		t.Fatal("Run() returned no pane")
	}
	waitFor(t, "the floating pane", func() (bool, error) {
		p, err := findPane(ctx, z, *popup)
		return p != nil && p.IsFloating && p.Title == "popup", err
	})
	if err := z.ActionClosePane(ctx, ActionClosePaneOption{PaneId: popup}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "close-pane", func() (bool, error) {
		p, err := findPane(ctx, z, *popup)
		return p == nil, err
	})

	tabId, err := z.ActionNewTab(ctx, ActionNewTabOption{Name: "zellijcli-tab"})
	if err != nil {
		t.Fatal(err)
	}
	hasTab := func() (bool, error) {
		tabs, err := z.ActionListTabs(ctx)
		return slices.ContainsFunc(tabs, func(tab Tab) bool { return tab.TabId == tabId && tab.Name == "zellijcli-tab" }), err
	}
	waitFor(t, "new-tab", hasTab)
	if err := z.ActionGoToTabById(ctx, tabId); err != nil {
		t.Fatal(err)
	}
	if err := z.ActionCloseTabById(ctx, tabId); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "close-tab-by-id", func() (bool, error) {
		ok, err := hasTab()
		return !ok && err == nil, err
	})

	clients, err := z.ActionListClients(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(clients) != 0 {
		t.Errorf("ActionListClients() = %+v, want none", clients)
	}

	_, err = z.ActionCurrentTabInfo(ctx)
	if cmdErr, ok := errors.AsType[*CmdError](err); !ok || cmdErr.ExitCode() != 2 {
		t.Errorf("ActionCurrentTabInfo() err = %v, want *CmdError with exit status 2 (no attached client)", err)
	}
}
