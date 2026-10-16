# tmux / zellij CLI usage

This file lists every tmux / zellij invocation found in the two source repos and maps each one to its wrapper method in
`tmuxcli` / `zellijcli`. It also maps each tmux command to its zellij counterpart.

The source repos are:

- `github.com/ngicks/cmdman`, at commit `467af41`. Paths are relative to the repo root.
- `github.com/ngicks/run-in-tmux-popup`, at commit `9912705`. Paths are relative to the repo root.

## cmdman → tmuxcli

Nearly all calls go through the `executor.run` helper (`pkg/muxctl/tmux/exec.go:28`). It prefixes `-S <path>` when the
socket value contains a `/`, and `-L <name>` otherwise. The wrapper sets these through `GlobalOption.SocketPath` and
`GlobalOption.SocketName`.

`display-popup` in `cmdman/cli/tui.go` is the one call that skips `executor.run`, so it passes neither socket flag.

These `@cmdman_*` user options are read and written through set-option, show-options and `#{@...}` formats:

- window: `@cmdman_window`, `@cmdman_created`, `@cmdman_frame_def`, `@cmdman_scale`, `@cmdman_<key>`
- pane: `@cmdman_marker`, `@cmdman_leaf`, `@cmdman_frame`

| tmux command | flags / formats used | wrapper | call sites |
|---|---|---|---|
| `display-message` | `-p [-t T] <fmt>`. The formats are `#{session_name}`, `#{window_id}`, `#{pane_id}`, `#{pane_width}\t#{pane_height}`, `#{window_width}\t#{window_height}` and `#{window_id}\t#{window_name}\t#{window_panes}\t#{@cmdman_window}\t#{@cmdman_frame_def}` | `DisplayMessage` (raw format → string) | tmux.go:68, 92-98, 133; reuse.go:26; focus.go:38; apply.go:242, 256; frame.go:478 |
| `list-panes` | `-t <window> -F`, with fields drawn from `pane_id`, `pane_floating_flag`, `pane_dead`, `pane_width`, `pane_height`, `pane_title`, `@cmdman_frame`, `@cmdman_marker` and `@cmdman_leaf` | `ListPanes` → `[]Pane`. User options come through `ListPanesOption.UserOptions` | frame.go:60, 112; detach.go:144, 213, 250; leaf.go:133; stat.go:39 |
| `list-windows` | `-t <session> -F` or `-a -F`, with fields `window_id`, `session_name`, `window_name`, `@cmdman_window`, `@cmdman_frame_def`, `@cmdman_created` and `@cmdman_<key>` | `ListWindows` → `[]Window` | tmux.go:354; list.go:61-66 |
| `list-clients` | `-F #{client_activity}\t#{client_name} [-t =<session>]` | `ListClients` → `[]Client` | focus.go:72-76 |
| `show-options` | `-w -t <window>` (all options), `-gv default-shell`, `-w -t <window> -v @cmdman_<key>` | `ShowOptions` → `[]OptionEntry` | detach.go:127, 164; scale_state.go:49 |
| `set-option` | `-w [-u] -t <window> <opt> [value]` for `@cmdman_*`, `pane-border-status` and `remain-on-exit`. `-p [-u] -t <pane> <opt> [value]` for `@cmdman_frame`, `@cmdman_marker` and `@cmdman_leaf` | `SetOption` | tmux.go:209, 218, 231; frame.go:206, 294, 315, 456; detach.go:41, 43, 60, 61, 100, 102, 188, 199; leaf.go:26, 35, 75, 80-112; scale_state.go:73, 76 |
| `split-window` | `-h\|-v [-b] [-f] -d -l <size> -t <pane> -P -F #{pane_id}` | `SplitWindow` → `SpawnResult` | apply.go:165; frame.go:148, 430-440 |
| `new-window` | `-d -a -t =<s>:{end} -n <name> -P -F #{window_id} [-c dir]` | `NewWindow` → `SpawnResult` | tmux.go:388-396 |
| `new-session` | `-d -s <name> [-c dir]`. A "duplicate session" error is treated as success | `NewSession` | tmux.go:328-332 |
| `has-session` | `-t =<name>` | `HasSession` → `bool` | tmux.go:282, 325 |
| `respawn-pane` | `-k -t <pane> -- <argv...>` | `RespawnPane` | apply.go:234 |
| `kill-pane` | `-t <pane>` | `KillPane` | apply.go:224; frame.go:322 |
| `kill-window` | `-t <window>` | `KillWindow` | tmux.go:313 |
| `select-pane` | `-t <pane>`, `-t <pane> -T <title>` | `SelectPane` | apply.go:73; frame.go:311, 490; leaf.go:69; detach.go:71 |
| `select-window` | `-t <window>` | `SelectWindow` | focus.go:32 |
| `switch-client` | `-c <client> -t =<session>` | `SwitchClient` | focus.go:49 |
| `send-keys` | `-t <pane> <keys...>` | `SendKeys` | detach.go:194; leaf.go:31 |
| `attach-session` | `-t =<session>`. The argv is built but not run; the caller runs it later via `tea.ExecProcess` | `AttachSession` → `*exec.Cmd` | focus.go:60 |
| `display-popup` | `-E [-d cwd] [-w W] [-h H] [-x X] [-y Y] <cmd>`, with the terminal attached | `DisplayPopup` → `*exec.Cmd` | cmdman/cli/tui.go:360-403 |
| `capture-pane` | `-p -t <target>` | `CapturePane` | e2e only: e2e/cmdman/tui_launcher_landing_test.go:265 |
| `list-sessions` | (default format) | `ListSessions` → `[]Session` | e2e only: tui_launcher_landing_test.go:284 |
| `kill-server` | | `KillServer` | e2e only: e2e/cmdman/mux_test.go:62 |

In the table above, paths without a directory are under `pkg/muxctl/tmux/`.

## run-in-tmux-popup → tmuxcli / zellijcli

The tmux binary path and the zellij binary path each come from field 2 of `PINENTRY_USER_DATA`. When that field is
empty, the default is `tmux` or `zellij`. When `$TMUX` is empty, `TMUX=<SessionMeta>` is injected into the child
environment; the wrapper takes this through `ShellOut.Env`.

| command | flags used | wrapper | call sites |
|---|---|---|---|
| tmux `display-popup` | `[-c client] [-T title] [-x X] [-y Ybottom] [-w W] [-h H] [-e K=V...] -E <cmdline>`. It is started in the background and its exit status is the payload's | `DisplayPopup` → `*exec.Cmd` | runinpopup/internal/tmux/popup.go:40-56 |
| tmux `display-popup -C` | `-C [-c client]` | `DisplayPopup` with `Close` | popup.go:77-86 |
| tmux `new-pane` (floating) | `[-t session] -P -F #{pane_id} [-X x] [-Y y] [-x w] [-y h] [-e K=V...] -- <cmdline>` | `NewPane` → `SpawnResult` | popup.go:140-151 |
| tmux `kill-pane` | `-t <pane>` | `KillPane` | popup.go:201-206 |
| tmux `-V` | output parsed as `tmux <maj>.<min><suffix>`; the check is `< 3.7c` | `Version` | runinpopup/internal/tmux/version.go:12-14 |
| tmux `display-message` | `-p [-t session] #{window_zoomed_flag}:#{pane_id}` | `DisplayMessage` | runinpopup/internal/tmux/zoom.go:11-25 |
| tmux `resize-pane` | `-Z -t <pane>` | `ResizePane` | zoom.go:28-31 |
| zellij `run` | `[--session=S] run [--name=N] [--x=X] [--y=Y] [--width=W] [--height=H] --floating --close-on-exit --pinned=true -- <payload...>`, which prints `terminal_<n>` | `Run` → `*PaneId` (with `GlobalOption.Session`) | runinpopup/internal/zellij/run.go:59-77 |
| zellij `action close-pane` | `[--session=S] action close-pane --pane-id=terminal_<n>` | `ActionClosePane` | run.go:141-150 |

## tmux ↔ zellij counterparts

The two tools line up like this:

| tmux | zellij |
|---|---|
| server | the set of all sessions |
| session | session |
| window | tab |
| pane | pane |

A dash (—) in the table below means zellij has no counterpart.

| tmux (`tmuxcli`) | zellij (`zellijcli`) | notes |
|---|---|---|
| `-V` `Version` | `--version` `Version` | |
| `list-sessions` `ListSessions` | `list-sessions` `ListSessions` | |
| `has-session` `HasSession` | `list-sessions` `ListSessions` | zellij has no has-session; look for the name in the list |
| `new-session -d` `NewSession` | `attach --create-background` `Attach` | |
| `attach-session` `AttachSession` | `attach` `Attach` | both return `*exec.Cmd` |
| `kill-server` `KillServer` | `kill-all-sessions --yes` `KillAllSessions` | |
| — | `kill-session` `KillSession`, `delete-session` `DeleteSession` | a killed zellij session may linger as "EXITED" until it is deleted |
| `display-popup` `DisplayPopup` | `run --floating --close-on-exit` `Run` | zellij returns a pane id and does not block |
| `new-pane` `NewPane` | `action new-pane --floating` `ActionNewPane` | in tmux 3.7c `new-pane` always creates a floating pane (verified live). Its `-X/-Y/-x/-y` map to zellij's `--x/--y/--width/--height` |
| `split-window` `SplitWindow` | `action new-pane --direction` `ActionNewPane` | this is the tiled split. split-window rejects `-x/-y/-X/-Y` |
| `respawn-pane -k` `RespawnPane` | `action new-pane --in-place --close-replaced-pane --pane-id` `ActionNewPane` | this replaces the pane, and zellij gives the replacement a new pane id. It may print no id (a race seen live), so `ActionNewPane` returns `nil` in that case |
| `new-window` `NewWindow` | `action new-tab` `ActionNewTab` | |
| `kill-pane` `KillPane` | `action close-pane` `ActionClosePane` | |
| `kill-window` `KillWindow` | `action close-tab-by-id` `ActionCloseTabById` | |
| `select-pane -t` `SelectPane` | `action focus-pane-id` `ActionFocusPaneId` | |
| `select-pane -T` `SelectPane` | `action rename-pane` `ActionRenamePane` | |
| `select-window` `SelectWindow` | `action go-to-tab-by-id` `ActionGoToTabById` | |
| `switch-client` `SwitchClient` | `action switch-session` `ActionSwitchSession` | |
| `send-keys` `SendKeys` | `action send-keys` `ActionSendKeys` / `action write-chars` `ActionWriteChars` | `write-chars` corresponds to `send-keys -l` |
| `list-panes` `ListPanes` | `action list-panes --all --json` `ActionListPanes` | |
| `list-windows` `ListWindows` | `action list-tabs --all --json` `ActionListTabs` | |
| `list-clients` `ListClients` | `action list-clients` `ActionListClients` | |
| `display-message -p` `DisplayMessage` | `action current-tab-info --json` `ActionCurrentTabInfo` | this covers the current tab only; with no attached client, zellij reports "No active tab found" |
| `capture-pane -p` `CapturePane` | `action dump-screen` `ActionDumpScreen` | |
| `resize-pane -Z` `ResizePane` | `action toggle-fullscreen` `ActionToggleFullscreen` | |
| `resize-pane -U/-D/-L/-R` `ResizePane` | `action resize` `ActionResize` | |
| `set-option` / `show-options` | — | zellij has no user options, so cmdman's `@cmdman_*` state needs another store |
