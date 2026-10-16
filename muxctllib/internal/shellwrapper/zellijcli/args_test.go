package zellijcli

import (
	"context"
	"path/filepath"
	"slices"
	"testing"
)

func TestArgv(t *testing.T) {
	terminal1 := &PaneId{Id: 1}
	tests := []struct {
		name   string
		global GlobalOption
		stdout string
		call   func(ctx context.Context, s *ShellOut) error
		want   []string
	}{
		{
			name:   "global options precede the command",
			global: GlobalOption{Session: "s", ConfigFile: "c.kdl", ConfigDir: "/c", DataDir: "/d", Debug: true},
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionListClients(ctx)
				return err
			},
			want: []string{"--session=s", "--config=c.kdl", "--config-dir=/c", "--data-dir=/d", "--debug", "action", "list-clients"},
		},
		{
			name:   "version",
			stdout: "zellij 0.45.1\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.Version(ctx)
				return err
			},
			want: []string{"--version"},
		},
		{
			name:   "list-sessions",
			stdout: "a\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ListSessions(ctx, ListSessionsOption{Short: true, Reverse: true})
				return err
			},
			want: []string{"list-sessions", "--no-formatting", "--short", "--reverse"},
		},
		{
			name: "kill-all-sessions",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.KillAllSessions(ctx, KillAllSessionsOption{Yes: true})
			},
			want: []string{"kill-all-sessions", "--yes"},
		},
		{
			name: "kill-session",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.KillSession(ctx, "-s")
			},
			want: []string{"kill-session", "--", "-s"},
		},
		{
			name: "kill-session without target",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.KillSession(ctx, "")
			},
			want: []string{"kill-session"},
		},
		{
			name: "delete-session",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.DeleteSession(ctx, DeleteSessionOption{Force: true}, "s")
			},
			want: []string{"delete-session", "--force", "--", "s"},
		},
		{
			name:   "run with every flag",
			stdout: "terminal_3\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.Run(ctx, RunOption{
					Direction:             "right",
					Cwd:                   "/tmp",
					Floating:              true,
					InPlace:               true,
					CloseReplacedPane:     true,
					Name:                  "n",
					CloseOnExit:           true,
					StartSuspended:        true,
					X:                     "10%",
					Y:                     "1",
					Width:                 "50%",
					Height:                "20",
					Pinned:                new(true),
					Stacked:               true,
					Blocking:              true,
					BlockUntilExitSuccess: true,
					BlockUntilExitFailure: true,
					BlockUntilExit:        true,
					NearCurrentPane:       true,
					NoFocus:               true,
					Borderless:            new(false),
					TabId:                 new(2),
				}, "sh", "-c", "x")
				return err
			},
			want: []string{
				"run",
				"--direction=right", "--cwd=/tmp", "--floating", "--in-place", "--close-replaced-pane",
				"--name=n", "--close-on-exit", "--start-suspended",
				"--x=10%", "--y=1", "--width=50%", "--height=20", "--pinned=true",
				"--stacked", "--blocking", "--block-until-exit-success", "--block-until-exit-failure",
				"--block-until-exit", "--near-current-pane", "--no-focus", "--borderless=false", "--tab-id=2",
				"--", "sh", "-c", "x",
			},
		},
		{
			name:   "run without flags",
			stdout: "terminal_3\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.Run(ctx, RunOption{}, "-x")
				return err
			},
			want: []string{"run", "--", "-x"},
		},
		{
			name:   "new-pane with its own flags",
			stdout: "plugin_1\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionNewPane(ctx, ActionNewPaneOption{
					Floating:        true,
					Blocking:        true,
					Plugin:          "zellij:about",
					PaneId:          terminal1,
					Configuration:   "a=b",
					SkipPluginCache: true,
				})
				return err
			},
			want: []string{
				"action", "new-pane", "--floating", "--blocking",
				"--plugin=zellij:about", "--pane-id=terminal_1", "--configuration=a=b", "--skip-plugin-cache",
			},
		},
		{
			name:   "new-pane with a command",
			stdout: "terminal_2\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionNewPane(ctx, ActionNewPaneOption{}, "top")
				return err
			},
			want: []string{"action", "new-pane", "--", "top"},
		},
		{
			name:   "new-tab",
			stdout: "4\n",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionNewTab(ctx, ActionNewTabOption{
					Layout:                "l",
					LayoutString:          "layout {}",
					LayoutDir:             "/l",
					Name:                  "n",
					Cwd:                   "/tmp",
					InitialPlugin:         "p",
					CloseOnExit:           true,
					StartSuspended:        true,
					BlockUntilExitSuccess: true,
					BlockUntilExitFailure: true,
					BlockUntilExit:        true,
					NoFocus:               true,
				}, "sleep", "1")
				return err
			},
			want: []string{
				"action", "new-tab",
				"--layout=l", "--layout-string=layout {}", "--layout-dir=/l", "--name=n", "--cwd=/tmp",
				"--initial-plugin=p", "--close-on-exit", "--start-suspended",
				"--block-until-exit-success", "--block-until-exit-failure", "--block-until-exit", "--no-focus",
				"--", "sleep", "1",
			},
		},
		{
			name: "close-pane",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionClosePane(ctx, ActionClosePaneOption{PaneId: &PaneId{Plugin: true, Id: 2}})
			},
			want: []string{"action", "close-pane", "--pane-id=plugin_2"},
		},
		{
			name: "close-pane focused",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionClosePane(ctx, ActionClosePaneOption{})
			},
			want: []string{"action", "close-pane"},
		},
		{
			name: "close-tab-by-id",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionCloseTabById(ctx, 3)
			},
			want: []string{"action", "close-tab-by-id", "3"},
		},
		{
			name: "focus-pane-id",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionFocusPaneId(ctx, PaneId{Id: 4})
			},
			want: []string{"action", "focus-pane-id", "terminal_4"},
		},
		{
			name: "rename-pane",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionRenamePane(ctx, ActionRenamePaneOption{PaneId: terminal1}, "-name")
			},
			want: []string{"action", "rename-pane", "--pane-id=terminal_1", "--", "-name"},
		},
		{
			name: "go-to-tab-by-id",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionGoToTabById(ctx, 5)
			},
			want: []string{"action", "go-to-tab-by-id", "5"},
		},
		{
			name: "switch-session",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionSwitchSession(ctx, ActionSwitchSessionOption{
					TabPosition:  new(1),
					PaneId:       terminal1,
					Layout:       "l",
					LayoutString: "layout {}",
					LayoutDir:    "/l",
					Cwd:          "/tmp",
				}, "other")
			},
			want: []string{
				"action", "switch-session",
				"--tab-position=1", "--pane-id=terminal_1", "--layout=l", "--layout-string=layout {}",
				"--layout-dir=/l", "--cwd=/tmp", "--", "other",
			},
		},
		{
			name: "send-keys",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionSendKeys(ctx, ActionSendKeysOption{PaneId: terminal1}, "Ctrl c", "Enter")
			},
			want: []string{"action", "send-keys", "--pane-id=terminal_1", "--", "Ctrl c", "Enter"},
		},
		{
			name: "write-chars",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionWriteChars(ctx, ActionWriteCharsOption{}, "-x")
			},
			want: []string{"action", "write-chars", "--", "-x"},
		},
		{
			name:   "list-panes",
			stdout: "[]",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionListPanes(ctx)
				return err
			},
			want: []string{"action", "list-panes", "--all", "--json"},
		},
		{
			name:   "list-tabs",
			stdout: "[]",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionListTabs(ctx)
				return err
			},
			want: []string{"action", "list-tabs", "--all", "--json"},
		},
		{
			name:   "current-tab-info",
			stdout: "{}",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionCurrentTabInfo(ctx)
				return err
			},
			want: []string{"action", "current-tab-info", "--json"},
		},
		{
			name: "dump-screen",
			call: func(ctx context.Context, s *ShellOut) error {
				_, err := s.ActionDumpScreen(ctx, ActionDumpScreenOption{Path: "/tmp/x", Full: true, PaneId: terminal1, Ansi: true})
				return err
			},
			want: []string{"action", "dump-screen", "--path=/tmp/x", "--full", "--pane-id=terminal_1", "--ansi"},
		},
		{
			name: "toggle-fullscreen",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionToggleFullscreen(ctx, ActionToggleFullscreenOption{PaneId: terminal1})
			},
			want: []string{"action", "toggle-fullscreen", "--pane-id=terminal_1"},
		},
		{
			name: "resize",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionResize(ctx, ActionResizeOption{PaneId: terminal1}, "increase", "left")
			},
			want: []string{"action", "resize", "--pane-id=terminal_1", "--", "increase", "left"},
		},
		{
			name: "resize without direction",
			call: func(ctx context.Context, s *ShellOut) error {
				return s.ActionResize(ctx, ActionResizeOption{}, "decrease", "")
			},
			want: []string{"action", "resize", "--", "decrease"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, args := fakeZellij(t, tt.stdout, "", 0)
			s.Global = tt.global
			if err := tt.call(t.Context(), s); err != nil {
				t.Fatal(err)
			}
			if got := args(); !slices.Equal(got, tt.want) {
				t.Errorf("args\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestAttachArgs(t *testing.T) {
	s := &ShellOut{Global: GlobalOption{Session: "ignored"}}
	tests := []struct {
		name           string
		opt            AttachOption
		sessionName    string
		initialCommand []string
		want           []string
	}{
		{
			name:        "background",
			opt:         AttachOption{CreateBackground: true},
			sessionName: "s",
			want:        []string{"--session=ignored", "attach", "--create-background", "s"},
		},
		{
			name: "every flag with an initial command",
			opt: AttachOption{
				Create:           true,
				CreateBackground: true,
				Index:            new(0),
				ForceRunCommands: true,
				Token:            "t",
				Remember:         true,
				Forget:           true,
				CaCert:           "/ca.pem",
				Insecure:         true,
				CloseOnExit:      true,
				StartSuspended:   true,
			},
			sessionName:    "s",
			initialCommand: []string{"htop", "-d", "1"},
			want: []string{
				"--session=ignored", "attach",
				"--create", "--create-background", "--index=0", "--force-run-commands", "--token=t",
				"--remember", "--forget", "--ca-cert=/ca.pem", "--insecure", "--close-on-exit", "--start-suspended",
				"s", "--", "htop", "-d", "1",
			},
		},
		{
			name: "no session name",
			want: []string{"--session=ignored", "attach"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := s.Attach(t.Context(), tt.opt, tt.sessionName, tt.initialCommand...)
			if cmd.Path == "" || filepath.Base(cmd.Args[0]) != "zellij" {
				t.Errorf("command = %q, want zellij", cmd.Args[0])
			}
			if got := cmd.Args[1:]; !slices.Equal(got, tt.want) {
				t.Errorf("args\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}
