package tmuxcli

import (
	"os/exec"
	"slices"
	"testing"
)

func TestOptionArgs(t *testing.T) {
	list := ListOption{Filter: "#{pane_active}", SortOrder: "index", Reverse: true}
	listArgs := []string{"-r", "-f", "#{pane_active}", "-O", "index"}
	spawn := SpawnOption{Detached: true, StartDirectory: "/tmp", Environment: []string{"A=1", "B=2"}}
	spawnArgs := []string{"-d", "-c", "/tmp", "-e", "A=1", "-e", "B=2"}
	paneSplit := PaneSplitOption{
		SpawnOption:         spawn,
		Before:              true,
		Empty:               true,
		Full:                true,
		Horizontal:          true,
		ForwardStdin:        true,
		KeepOnExit:          true,
		Size:                "30%",
		Message:             "bye",
		Percentage:          new(40),
		InactiveBorderStyle: "fg=red",
		Style:               "bg=blue",
		ActiveBorderStyle:   "fg=green",
		TargetPane:          "%1",
		Vertical:            true,
		Zoom:                true,
	}
	paneSplitArgs := slices.Concat(spawnArgs, []string{
		"-b", "-E", "-f", "-h", "-I", "-k", "-l", "30%", "-m", "bye", "-p", "40",
		"-R", "fg=red", "-s", "bg=blue", "-S", "fg=green", "-t", "%1", "-v", "-Z",
	})
	scope := OptionScope{Pane: true, Window: true, Server: true, Global: true}
	scopeArgs := []string{"-p", "-w", "-s", "-g"}

	tests := []struct {
		name string
		got  []string
		want []string
	}{
		{"zero option emits nothing", argsOf(NewPaneOption{}), nil},
		{"ListSessionsOption", argsOf(ListSessionsOption{ListOption: list}), listArgs},
		{
			"ListPanesOption",
			argsOf(ListPanesOption{ListOption: list, All: true, Session: true, Target: "main", UserOptions: []string{"@x"}}),
			slices.Concat(listArgs, []string{"-a", "-s", "-t", "main"}),
		},
		{
			"ListWindowsOption",
			argsOf(ListWindowsOption{ListOption: list, All: true, TargetSession: "main", UserOptions: []string{"@x"}}),
			slices.Concat(listArgs, []string{"-a", "-t", "main"}),
		},
		{
			"ListClientsOption",
			argsOf(ListClientsOption{ListOption: list, TargetSession: "main"}),
			slices.Concat(listArgs, []string{"-t", "main"}),
		},
		{
			"NewWindowOption",
			argsOf(NewWindowOption{
				SpawnOption:    spawn,
				After:          true,
				Before:         true,
				Kill:           true,
				SelectExisting: true,
				WindowName:     "w",
				TargetWindow:   "main:1",
			}),
			slices.Concat(spawnArgs, []string{"-a", "-b", "-k", "-S", "-n", "w", "-t", "main:1"}),
		},
		{
			"NewPaneOption",
			argsOf(NewPaneOption{PaneSplitOption: paneSplit, Width: "80%", Height: "10", X: "5", Y: "3"}),
			slices.Concat(paneSplitArgs, []string{"-x", "80%", "-y", "10", "-X", "5", "-Y", "3"}),
		},
		{"SplitWindowOption", argsOf(SplitWindowOption{PaneSplitOption: paneSplit}), paneSplitArgs},
		{
			"NewSessionOption",
			argsOf(NewSessionOption{
				AttachExisting:      true,
				Detached:            true,
				DetachOthers:        true,
				NoUpdateEnvironment: true,
				HangupOthers:        true,
				StartDirectory:      "/tmp",
				Environment:         []string{"A=1"},
				ClientFlags:         "read-only",
				WindowName:          "w",
				SessionName:         "s",
				GroupName:           "g",
				Width:               "80",
				Height:              "-",
			}),
			[]string{
				"-A", "-d", "-D", "-E", "-X", "-c", "/tmp", "-e", "A=1", "-f", "read-only",
				"-n", "w", "-s", "s", "-t", "g", "-x", "80", "-y", "-",
			},
		},
		{
			"AttachSessionOption",
			argsOf(AttachSessionOption{
				DetachOthers:        true,
				NoUpdateEnvironment: true,
				ReadOnly:            true,
				HangupOthers:        true,
				WorkingDirectory:    "/tmp",
				ClientFlags:         "active-pane",
				TargetSession:       "main",
			}),
			[]string{"-d", "-E", "-r", "-x", "-c", "/tmp", "-f", "active-pane", "-t", "main"},
		},
		{
			"CapturePaneOption",
			argsOf(CapturePaneOption{
				Alternate:              true,
				EscapeNonPrintable:     true,
				EscapeSequences:        true,
				LineFlags:              true,
				Hyperlinks:             true,
				JoinWrapped:            true,
				LineNumbers:            true,
				Mode:                   true,
				PreserveTrailingSpaces: true,
				Print:                  true,
				IncompleteEscape:       true,
				Quiet:                  true,
				TrimTrailing:           true,
				BufferName:             "buf",
				EndLine:                "-",
				StartLine:              "-10",
				TargetPane:             "%1",
			}),
			[]string{
				"-a", "-C", "-e", "-F", "-H", "-J", "-L", "-M", "-N", "-p", "-P", "-q", "-T",
				"-b", "buf", "-E", "-", "-S", "-10", "-t", "%1",
			},
		},
		{
			"DisplayMessageOption",
			argsOf(DisplayMessageOption{
				ListVariables: true,
				KeepUpdating:  true,
				ForwardStdin:  true,
				Literal:       true,
				IgnoreKeys:    true,
				Print:         true,
				Verbose:       true,
				TargetClient:  "/dev/pts/1",
				Delay:         new(0),
				Format:        "#{pane_id}",
				TargetPane:    "%1",
			}),
			[]string{
				"-a", "-C", "-I", "-l", "-N", "-p", "-v",
				"-c", "/dev/pts/1", "-d", "0", "-F", "#{pane_id}", "-t", "%1",
			},
		},
		{
			"DisplayPopupOption",
			argsOf(DisplayPopupOption{
				NoBorder:       true,
				Close:          true,
				CloseOnExit:    true,
				AnyKey:         true,
				DisableClose:   true,
				BorderLines:    "rounded",
				TargetClient:   "/dev/pts/1",
				StartDirectory: "/tmp",
				Environment:    []string{"A=1"},
				Height:         "50%",
				Style:          "bg=red",
				BorderStyle:    "fg=red",
				TargetPane:     "%1",
				Title:          "#{pane_id}",
				Width:          "80",
				X:              "C",
				Y:              "S",
			}),
			[]string{
				"-B", "-C", "-E", "-k", "-N", "-b", "rounded", "-c", "/dev/pts/1", "-d", "/tmp",
				"-e", "A=1", "-h", "50%", "-s", "bg=red", "-S", "fg=red", "-t", "%1",
				"-T", "#{pane_id}", "-w", "80", "-x", "C", "-y", "S",
			},
		},
		{"DisplayPopupOption CloseOnSuccess", argsOf(DisplayPopupOption{CloseOnSuccess: true}), []string{"-E", "-E"}},
		{
			"DisplayPopupOption CloseOnSuccess with CloseOnExit",
			argsOf(DisplayPopupOption{CloseOnExit: true, CloseOnSuccess: true}),
			[]string{"-E", "-E"},
		},
		{"HasSessionOption", argsOf(HasSessionOption{TargetSession: "main"}), []string{"-t", "main"}},
		{"KillPaneOption", argsOf(KillPaneOption{AllOthers: true, TargetPane: "%1"}), []string{"-a", "-t", "%1"}},
		{"KillWindowOption", argsOf(KillWindowOption{AllOthers: true, TargetWindow: "@1"}), []string{"-a", "-t", "@1"}},
		{
			"ResizePaneOption",
			argsOf(ResizePaneOption{
				Down:       true,
				Left:       true,
				Mouse:      true,
				Right:      true,
				Trim:       true,
				Up:         true,
				Zoom:       true,
				Width:      "10%",
				Height:     "5",
				TargetPane: "%1",
			}),
			[]string{"-D", "-L", "-M", "-R", "-T", "-U", "-Z", "-x", "10%", "-y", "5", "-t", "%1"},
		},
		{
			"RespawnPaneOption",
			argsOf(RespawnPaneOption{Kill: true, StartDirectory: "/tmp", Environment: []string{"A=1"}, TargetPane: "%1"}),
			[]string{"-k", "-c", "/tmp", "-e", "A=1", "-t", "%1"},
		},
		{
			"SelectPaneOption",
			argsOf(SelectPaneOption{
				Down:         true,
				DisableInput: true,
				EnableInput:  true,
				Left:         true,
				Last:         true,
				ClearMarked:  true,
				Mark:         true,
				Right:        true,
				Up:           true,
				KeepZoom:     true,
				Title:        "t",
				TargetPane:   "%1",
			}),
			[]string{"-D", "-d", "-e", "-L", "-l", "-M", "-m", "-R", "-U", "-Z", "-T", "t", "-t", "%1"},
		},
		{
			"SelectWindowOption",
			argsOf(SelectWindowOption{Last: true, Next: true, Previous: true, Toggle: true, TargetWindow: "@1"}),
			[]string{"-l", "-n", "-p", "-T", "-t", "@1"},
		},
		{
			"SendKeysOption",
			argsOf(SendKeysOption{
				ExpandFormats:   true,
				Hex:             true,
				ToClient:        true,
				Literal:         true,
				Mouse:           true,
				Reset:           true,
				CopyModeCommand: true,
				TargetClient:    "/dev/pts/1",
				RepeatCount:     new(3),
				TargetPane:      "%1",
			}),
			[]string{"-F", "-H", "-K", "-l", "-M", "-R", "-X", "-c", "/dev/pts/1", "-N", "3", "-t", "%1"},
		},
		{
			"SetOptionOption",
			argsOf(SetOptionOption{
				OptionScope:   scope,
				Append:        true,
				ExpandFormats: true,
				OnlyIfUnset:   true,
				Quiet:         true,
				Unset:         true,
				UnsetPanes:    true,
				TargetPane:    "%1",
			}),
			slices.Concat(scopeArgs, []string{"-a", "-F", "-o", "-q", "-u", "-U", "-t", "%1"}),
		},
		{
			"ShowOptionsOption",
			argsOf(ShowOptionsOption{
				OptionScope:      scope,
				IncludeInherited: true,
				IncludeHooks:     true,
				Quiet:            true,
				ValueOnly:        true,
				TargetPane:       "%1",
			}),
			slices.Concat(scopeArgs, []string{"-A", "-H", "-q", "-v", "-t", "%1"}),
		},
		{
			"SwitchClientOption",
			argsOf(SwitchClientOption{
				NoUpdateEnvironment: true,
				Last:                true,
				Next:                true,
				Previous:            true,
				ToggleReadOnly:      true,
				KeepZoom:            true,
				TargetClient:        "/dev/pts/1",
				SortOrder:           "name",
				TargetSession:       "main",
				KeyTable:            "root",
			}),
			[]string{
				"-E", "-l", "-n", "-p", "-r", "-Z",
				"-c", "/dev/pts/1", "-O", "name", "-t", "main", "-T", "root",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.Equal(tt.got, tt.want) {
				t.Errorf("args =\n%q\nwant\n%q", tt.got, tt.want)
			}
		})
	}
}

func TestUnstartedCommand(t *testing.T) {
	s := &ShellOut{
		Path:   "/opt/tmux",
		Global: GlobalOption{SocketName: "sock", SocketPath: "/tmp/sock", ConfigFile: "/dev/null"},
		Env:    []string{"A=1"},
		Dir:    "/tmp",
	}
	global := []string{"/opt/tmux", "-L", "sock", "-S", "/tmp/sock", "-f", "/dev/null"}
	ctx := t.Context()
	tests := []struct {
		name string
		cmd  *exec.Cmd
		want []string
	}{
		{
			"AttachSession",
			s.AttachSession(ctx, AttachSessionOption{ReadOnly: true, TargetSession: "main"}),
			slices.Concat(global, []string{"attach-session", "-r", "-t", "main"}),
		},
		{
			"DisplayPopup",
			s.DisplayPopup(ctx, DisplayPopupOption{CloseOnExit: true}, "-x", "y"),
			slices.Concat(global, []string{"display-popup", "-E", "--", "-x", "y"}),
		},
		{
			"DisplayPopup close",
			s.DisplayPopup(ctx, DisplayPopupOption{Close: true}),
			slices.Concat(global, []string{"display-popup", "-C"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.Equal(tt.cmd.Args, tt.want) {
				t.Errorf("Args =\n%q\nwant\n%q", tt.cmd.Args, tt.want)
			}
			if tt.cmd.Path != s.Path || !slices.Equal(tt.cmd.Env, s.Env) || tt.cmd.Dir != s.Dir {
				t.Errorf("Path, Env, Dir = %q, %q, %q", tt.cmd.Path, tt.cmd.Env, tt.cmd.Dir)
			}
		})
	}
}

func argsOf(o interface{ args(argv) argv }) []string {
	return o.args(nil)
}
