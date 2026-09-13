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

## Conventions

- Short lines, 80 to 110 characters, with H2 or larger headings.
- No bold and no em dashes.
- One plugin per bullet, with its full GitHub URL.
- Changes through a branch and a pull request.

## Releasing

- Semantic commit messages.
- Update this file when the map's process changes.
