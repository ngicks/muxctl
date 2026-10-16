---
description: "Conventions for muxctllib/internal/shellwrapper (tmuxcli, zellijcli) and keeping its USAGE.md in sync"
applyTo: "**/muxctllib/internal/shellwrapper/**"
---

### shellwrapper

`muxctllib/internal/shellwrapper/{tmuxcli,zellijcli}` wrap the tmux and zellij CLIs.
`USAGE.md` next to them records which wrapper method serves each real CLI call, and which zellij
command is the counterpart of each tmux command.

#### Keep code and USAGE.md in sync

A change to any of the following is incomplete until `USAGE.md` matches it:

- adding, removing or renaming a command file, method, option field or result type;
- a new call site of a tmux / zellij command in muxctllib, an app, or one of the source repos
  listed in `USAGE.md`;
- a tmux ↔ zellij counterpart that is added, found to be wrong, or confirmed to have no match.

Use the `local-muxctllib-shellwrapper-sync` skill for the full procedure.

#### Package layout

- `shell.go` holds `ShellOut` (`Path`, `Global`, `Env`, `Dir`), `GlobalOption`, `CmdError`, and the
  package-wide helpers (`argv`, format / record parsing, `PaneId`).
  - The wrapper never reads `$TMUX`, `$TMUX_PANE` or `$ZELLIJ*` itself. The caller passes the
    environment in through `Env`.
- There is one file per command, and one exported method on `*ShellOut` per file.
  - Files are named after the command, with hyphens kept: `list-panes.go`, `action_new-pane.go`.
  - zellij method names keep the full command path: `ActionNewPane`.
- Every command file is laid out in exactly this order:
  1. a free-floating comment that copies the command's man page section (tmux) or `--help`
     output (zellij) verbatim. A blank line must follow it, so it does not become a doc comment;
  2. the option type(s): `<Method>Option`, plus any shared option struct this command is the
     origin of;
  3. the result type(s);
  4. the method;
  5. unexported helpers.
- All flags go in `<Method>Option`. Positional arguments are method parameters.
  - Every flag that the man page, `tmux list-commands` or `--help` lists gets a field.
  - A flag the man page leaves out but the binary accepts gets a one-line comment saying why it
    is there.
  - Omit the option struct entirely when the command has no flags, or when the method forces
    all of them (e.g. `--all --json`).
- Flags that the docs define on another command are shared by embedding, not duplicated:
  - tmux: `ListOption`, `SpawnOption` (new-window) → `PaneSplitOption` (new-pane, split-window),
    and `OptionScope` (set-option, show-options);
  - zellij: `CommandPaneOption` (run, action new-pane).
- Typed results:
  - tmux list and spawn methods own `-F`/`-P` and parse fields separated by `\x1f` into structs.
    User options (`@name`) come through `UserOptions`.
  - zellij list methods force `--all --json` and decode the output.
  - Interactive or long-running commands (`attach-session`, `display-popup`, zellij `attach`)
    return an unstarted `*exec.Cmd`.
- Arguments are built value-style, like `append`: `func (a argv) flag(...) argv` and
  `func (o XOption) args(a argv) argv`.
  - Write one method chain with a direct `return`. Don't repeat `a = a.flag(...)`.
  - Never call a helper without using its result: the flag is silently lost.
- Abbreviations are not all-caps: `Id`, `PaneId`, `Pid`, `Tty`, `Ansi`.

#### Verifying behaviour

Check CLI facts against the real binaries before encoding them: output formats, exit codes,
quoting, and flags that are rejected or undocumented. tmux: `tmux -L <private-socket> ...`.
zellij: a background session (`zellij attach -b -c <name>`) with a temporary `ZELLIJ_CONFIG_DIR`
that holds an empty `config.kdl`.

Tests:

- table tests cover each option's `args()` and each parser;
- live tests (`live_test.go`) skip when the binary is missing, and use a private tmux socket or a
  uniquely named zellij session that they always clean up.
