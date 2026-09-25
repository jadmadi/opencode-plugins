# AGENTS.md

Guidance for agents working in this repository.

## What this is

The map for the OpenCode plugin ecosystem: what exists, where it lives, the
commands and tools each plugin adds, and the conventions they follow. Each
plugin lives in its own repo under `jadmadi/opencode-*`. This repo has no code.

## Keep the map current

Update this repo in the same piece of work whenever a plugin is added, renamed,
or changed in a way the map states. A plugin is not done until the map matches
it.

Checklist for a new or renamed plugin:

- Add or update its bullet under the right section in `README.md`, with a link
  to its repo and one line on what it does.
- Add its commands and tools to the "Commands and tools" table.
- Add its limits and environment overrides to "Limits and overrides".
- Add its design spec under `docs/compose/spec/<plugin>.md`, or say why there is
  none.
- Verify the repo link resolves with `gh repo list jadmadi --limit 200`.
- Ship the change on a branch and open a pull request.

For a behavior fix in a plugin, update the map only when the fix changes
something the map states: a command, a tool, an override, or a limit. Routine
internal fixes do not need a map change.

## Specs

The specs here are the designs. Each plugin repo carries the delivered copy,
with its report and journey log. Keep the two aligned when behavior changes; do
not let them drift. model-switcher and compose-next are original tools with no
design spec.

## Runtime API baseline

The plugins target the shipped OpenCode plugin API, not the website docs. The
docs move ahead of releases, so `runtime-baseline.json` records what the
plugins build against and `script/check-runtime-api.go` verifies it. Run the
check before any change to a plugin's `ctx` usage, from this repo root:

```sh
go run script/check-runtime-api.go
```

- HOLD means the baseline still matches. Keep the plain `{ id, setup }` object
  and the domains the baseline names.
- MIGRATE means the runtime or the published package moved. Re-probe with the
  steps below, then update the baseline and the plugins together.
- UNKNOWN means the check could not reach the runtime or the registry. Do not
  treat it as HOLD.

Baseline as of 2026-09-25, OpenCode 2.0.16 (re-probed):

- `ctx.model` carries model reads and transforms (members: `default`, `list`,
  `reload`, `transform`), and `ctx.provider` carries provider reads and
  transforms (members: `get`, `list`, `reload`, `transform`). Both take
  `transform(editor)` callbacks.
- `ctx.catalog` does not exist. Migrate `ctx.catalog.model.*` to
  `ctx.model.*` and `ctx.catalog.provider.*` to `ctx.provider.*`.
- `import { Plugin } from "@opencode/plugin"` does not resolve for local file
  plugins. Re-probed on 2.0.16: a directory plugin importing it fails to load
  with `Cannot find package '@opencode/plugin'`. The build docs use that import
  in their example, so the docs are ahead of the runtime. Keep the plain
  `{ id, setup }` object. Only `@opencode/plugin/tui` is aliased, for CLI
  plugins.

Probe recipe when the check says MIGRATE:

1. Drop a temporary plugin in `~/.config/opencode/plugins/` that writes
   `typeof ctx.catalog`, `typeof ctx.model`, `typeof ctx.provider`, and the
   member keys of each domain to a file. Touch it, read the file, then remove
   it.
2. Drop a directory plugin importing `{ Plugin } from "@opencode/plugin"` and
   see whether it loads.
3. Record the results in `runtime-baseline.json`. If the import resolves and
   `ctx.model` exists, migrate the plugins to `Plugin.define` and the
   model/provider domains, then bump versions, update the changelog, tag, and
   release.

Known stale text: some plugin AGENTS.md files say npm only publishes dev
snapshots of `@opencode/plugin`. That is no longer true. Fix the line when
that repo is next changed.

## Conventions

- Short lines, 80 to 110 characters, with H2 or larger headings.
- No bold and no em dashes.
- One plugin per bullet, with its full GitHub URL.
- Changes through a branch and a pull request.
- Plugin repos carry no agent provisioning: no `.cursor/`, `.devin/`,
  `CLAUDE.md`, `goals/`, or private files.
- `main` is the published channel for `github:` installs: every commit reaches
  users. Merge only when the tests pass.
- Plugins carry no dependencies. The test run writes `aube-lock.yaml`; ignore
  it, do not commit it.

## Versioning and releases

- `package.json.version` is the source of truth. Each plugin exports `VERSION`
  from its `.ts` file, and a test asserts the two match. Drift fails
  `bun test`.
- Each release records the OpenCode version it was built and tested under, as
  `builtUnder` in `package.json` plus a `builtUnder-<version>` tag, for example
  `builtUnder-2.0.16`.
- A status output ends with a line `<plugin> <version>` when the plugin has a
  status surface: model-switcher, context-limit, loop, goal, max-mode,
  skip-permissions, workflows, and distill. One-line status messages get the
  same line after a newline.
- Versions use CalVer `YYYY.MM.MICRO`: year, month, and a release counter that
  starts at 0, for example `2026.9.0`. Do not zero-pad the month or the counter,
  so the value stays valid semver.
- Bump the micro for each release within a month. The year and month follow the
  calendar.
- `CHANGELOG.md` records each release. Add an entry under the new version when
  you bump it.
- After a change merges to main, tag the repo `v<version>` (for example
  `v2026.9.0`) and publish a GitHub release with short notes.
- Install docs keep the `main` URL and add one line: pin a release by replacing
  `main` with a tag such as `v2026.9.0`.
- The first CalVer release is `2026.9.0` for every plugin.

## Releasing

- Semantic commit messages.
- Update this file when the map's process changes.
- After a plugin change merges, bump its version, update `CHANGELOG.md`, tag
  `v<version>`, and publish a release. The versioning rules are above.
