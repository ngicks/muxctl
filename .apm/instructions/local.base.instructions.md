---
description: "muxctl repository overview, layout, and naming rules for local agent primitives"
applyTo: "*"
---

### General

muxctl is a mono repo that holds a terminal-multiplexer library and TUI apps built on it.

- `muxctllib` is the library. It drives tmux and zellij through their command line interfaces.
- The TUI apps are built with bubbletea on top of `muxctllib`. Each app is its own Go module.

### Tech stack

- Go 1.27 (`go 1.27.0` in every `go.mod` and in `go.work`).
- bubbletea for the TUI apps.
- tmux (3.7c or later) and zellij (0.45 or later). They are driven only through their CLIs.
- Deno (TypeScript) for helper scripts, e.g. the scripts bundled with skills. Do not use Python.
  Run them with `deno run` and the narrowest `--allow-*` permissions they need.

### Structure overview

```
.                          (the git worktree, e.g. main/)
├── go.work                ties the modules together; add each new module with `go work use`
├── muxctllib/             library module github.com/ngicks/muxctl/muxctllib
│   └── internal/
│       └── shellwrapper/  typed wrappers over the tmux / zellij CLIs
│           ├── USAGE.md   inventory of CLI usage → wrapper method, plus tmux ↔ zellij map
│           ├── tmuxcli/   one file per tmux command
│           └── zellijcli/ one file per zellij (sub)command
├── <app>/                 a bubbletea TUI app, as its own module (none yet)
└── .apm/                  local APM package (these instructions and skills)
```

Run `go` commands from the worktree root. `go.work` resolves the modules, so for example
`go test ./muxctllib/...` works there.

### Toolchain caveat

The golangci-lint and ngcheckers PostToolUse hooks may be built with an older Go. If so, they fail on
this module with "package requires newer Go version" or "Go language version ... is lower than the
targeted Go version". That failure is about the tools, not the code. Verify with
`go build`, `go vet`, `gofmt -l` and `go test` instead, and say so in your report.

### Naming local agent primitives

APM deploys instructions and skills from several packages into the same directories.
When two of them share a name, the local one wins, and otherwise the first one found wins.
A bare name such as `base` can therefore shadow, or be shadowed by, a dependency's file without
any warning. To avoid that, name every primitive in reverse-DNS style, from the owner down to the
topic:

- Instructions in this repo: `.apm/instructions/local.<area>.<topic>.instructions.md`.
  - They always start with `local.`, which marks them as owned by this repo.
  - `<area>` is the module or top-level directory the rule is about (`muxctllib`, an app name),
    or `base` for rules that cover the whole repo.
  - Examples: `local.base.instructions.md`, `local.muxctllib.shellwrapper.instructions.md`.
- Skills in this repo: `.apm/skills/local-<area>-<topic>/SKILL.md`, and the frontmatter `name`
  must equal the directory name.
  - The Agent Skills spec only allows `[a-z0-9-]` in skill names, so skills use the same
    reverse-DNS order as instructions, with the dots replaced by hyphens.
  - Example: `local-muxctllib-shellwrapper-sync`.
- Dependencies keep their own prefixes (e.g. `ngicks.go.basics`). Never reuse a dependency's
  prefix for a local primitive.

### Installing / updating APM packages

Do not run `apm install`, `apm update` or `apm compile` directly. Use the shell functions
`apm-install` and `apm-bump` instead, and run them from the bare-repository root: the directory
that contains `.bare/` and the worktrees (e.g. `main/`). Never run them from inside a worktree.

- `apm-install` installs from `apm.yml` / `apm.lock.yaml`. Use it after adding, renaming or
  editing a local primitive under `.apm/`, or after changing `apm.yml`.
- `apm-bump` does the same with `--update`, so dependencies move to their latest refs. Run it
  only when asked to update dependencies.

Why not plain `apm`: this repo is a bare repo with worktrees, and every worktree shares the
agent config deployed at the bare root. The wrappers:

- clear the deployed `.agents/`, `.claude/` and `.codex/` at the root;
- install from the default-branch worktree with `--root ..`, so files land at the bare root rather
  than inside one worktree;
- copy `apm.yml`, `apm.lock.yaml` and `.apm/` so the root and the worktree stay in sync;
- compile `AGENTS.md` for codex.

A plain `apm install` inside a worktree does none of this. It writes a second, private
`.claude/` there and leaves the root copy stale.

The wrappers are defined in the user's interactive zsh profile. If `apm-install` is not found,
stop and ask the user to run it, e.g. `! cd <bare root> && apm-install`. Do not fall back to
plain `apm install`.
