---
name: local-muxctllib-shellwrapper-sync
description: "Update muxctllib/internal/shellwrapper (tmuxcli, zellijcli) and keep its USAGE.md in sync: re-collect tmux/zellij CLI usage, add or change command wrappers, refresh the tmux ↔ zellij counterpart map. Use when wrapping a new tmux/zellij command or flag, when a new tmux/zellij call site appears (in muxctllib, an app, cmdman or run-in-tmux-popup), when a tmux/zellij upgrade changes behaviour, or when asked to update, sync or audit USAGE.md."
---

# local-muxctllib-shellwrapper-sync

Keep three things consistent:

- the wrappers in `muxctllib/internal/shellwrapper/tmuxcli` and `.../zellijcli`;
- the CLI usage they serve;
- `muxctllib/internal/shellwrapper/USAGE.md`.

The code conventions are in the `local.muxctllib.shellwrapper` instruction. Follow them for any
code you write here.

Paths below are relative to the worktree root. `SKILL_DIR` is this skill's directory.

## 1. Collect usage

List every tmux / zellij invocation in the sources `USAGE.md` covers:

- this repo: `grep -rn 'shellwrapper/\(tmuxcli\|zellijcli\)' --include='*.go' .`, then read each
  method call;
- the source repos named at the top of `USAGE.md`, at their current commit.

For each source repo:

1. Find the places that run the binary, e.g. `exec.Command*` with `"tmux"`/`"zellij"`, or an
   executor such as cmdman's `executor.run`.
2. Trace every argv those places receive, including argv slices built up with `append`.
3. For each invocation, record:
   - the subcommand, and every distinct flag set, in the order the code passes them;
   - `-F`/`-p` format strings, verbatim;
   - positional arguments, and whether `--` is used;
   - whether stdout is consumed, and how it is parsed;
   - `file:line`.

This is a broad read-only sweep. Delegate it to an explorer subagent per repo when one is available,
and ask for tables in exactly that shape.

## 2. Diff against USAGE.md and the wrappers

Classify every finding as one of:

- **covered**: the method exists and the flags exist as option fields. Update `file:line` and flag
  details in `USAGE.md` if they moved.
- **missing flag**: add the field to the option struct (step 3), then update `USAGE.md`.
- **missing command**: add a new command file (step 3). Also look for the zellij counterpart, or
  record that there is none.
- **gone**: the call site no longer exists. Remove its row. Keep the wrapper unless the user asks
  to drop it.

## 3. Add or change wrappers

Get the authoritative text and flags from the installed binaries. Never write them from memory.

```sh
# tmux: render the man page without hyphenation or justification, then extract one section
MANWIDTH=90 MANPAGER=cat man --nh --nj tmux | sed 's/.\x08//g' > /tmp/tmux.man
deno run --allow-read "$SKILL_DIR/scripts/mansec.ts" /tmp/tmux.man <command>   # prints `//` comment lines
tmux -L probe list-commands | grep '^<command> '                    # exact flag template

# zellij: the help text is the "man page"
zellij <subcommand> --help
zellij action <name> --help
```

- `mansec.ts` cannot find `kill-server`, because its man entry has no synopsis. Copy that one by
  hand.
- Where the man page and `list-commands` disagree, probe the real binary. Examples: tmux 3.7c
  rejects the documented `display-popup -K`, and `split-window` rejects `-x/-y/-X/-Y`.
- Probe output formats and exit codes live before writing a parser:
  - tmux: `tmux -L <private-socket> ...`, then `kill-server` and remove the socket file;
  - zellij: `ZELLIJ_CONFIG_DIR=<tmp with empty config.kdl> zellij attach -b -c <name>`, then
    `zellij --session=<name> action ...`, then `kill-session`.

Then write the file in the mandatory order: free-floating comment, then option, result, method,
helpers. Add table tests for the new `args()` and any parser, and extend `live_test.go`.

## 4. Update USAGE.md

- **Source repos:** update each repo's commit in the header (`git -C <repo> rev-parse --short HEAD`).
- **Per-repo tables:** one row per subcommand, with the flags and formats used, the wrapper method,
  and the call sites.
- **Counterpart table:** keep the model tmux server ↔ all zellij sessions, session ↔ session,
  window ↔ tab, pane ↔ pane.
  - Record behaviour differences that were verified live in the notes column.
  - Mark a missing counterpart with a dash (—) and a reason. Never invent a command.
- Every backticked identifier must exist. Check them:

  ```sh
  cd muxctllib/internal/shellwrapper
  for id in $(grep -oE '`[A-Z][A-Za-z]+`' USAGE.md | tr -d '`' | sort -u); do
    grep -rqE "\b$id\b" --include='*.go' . || echo "missing: $id"
  done
  ```

## 5. Verify

Run from the worktree root:

```sh
go build ./muxctllib/... && go vet ./muxctllib/... && test -z "$(gofmt -l muxctllib)"
go test -count=1 -race ./muxctllib/internal/shellwrapper/...
deno run --allow-read "$SKILL_DIR/scripts/check-layout.ts" muxctllib/internal/shellwrapper
```

`check-layout.ts` exits non-zero and names each file that breaks the section order, whose man
comment is not followed by a blank line, or whose option struct is not named `<Method>Option`.

In the final report, list:

- what changed in the code;
- what changed in `USAGE.md`;
- any CLI fact you verified live that shaped the API.
