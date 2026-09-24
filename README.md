# opencode-plugins

A home for OpenCode V2 plugins and tools: ports from other agent CLIs, and
original tools. Each plugin lives in its own repo. This repo is the map: what
exists, where it lives, and the conventions every plugin follows. Design specs
for the ports and sila-prime live under `docs/compose/spec/`; each plugin repo
carries the delivered copy with its report. The machine-readable index is
`plugins.json`. See `AGENTS.md` for how to keep this map current.

## OpenCode

These plugins run on OpenCode. New accounts through my referral link get $5 in
usage credits, and I get $5 too:

https://opencode.ai/go?ref=N9H3ZEP22A

## Models

- [model-switcher](https://github.com/jadmadi/opencode-model-switcher): slash
  commands that switch provider and model mid-session. The command set comes
  from a user JSON file, so anyone can add their own.
- [context-limit](https://github.com/jadmadi/opencode-context-limit): a working
  context budget per model. It lowers the model window through a model
  transform, so compaction fires earlier. It can only lower a window.
- [max-mode](https://github.com/jadmadi/opencode-max-mode): a `best_of_n` tool
  that runs 2 to 8 candidate answers in parallel and a judge picks the winner.

## Sessions and memory

- [memory](https://github.com/jadmadi/opencode-memory): project memory,
  checkpoints, and notes, stored per project and injected at the start of a
  session. Registers read, append, and search tools plus three commands.
- [sila-prime](https://github.com/jadmadi/opencode-sila-prime): runs
  `sila prime` on the first prompt and injects the cross-tool briefing, so the
  sila knowledge store primes OpenCode without the CLI. Registers a `/sila`
  command for sila subcommands.
- [task-tool](https://github.com/jadmadi/opencode-task-tool): a tree `task`
  tool (T1, T1.1) so multi-step work survives long turns and compaction.

## Autonomy and workflow

- [goal](https://github.com/jadmadi/opencode-goal): a judged stopping condition
  per session. When a turn ends, a judge model decides if the goal is met, and
  the plugin nudges and continues until it is, or gives up at a cap.
- [workflows](https://github.com/jadmadi/opencode-workflows): deterministic
  multi-agent workflows. Ships deep-research: brief, plan, research, reflect,
  write, and review. Phases write artifacts and a run is resumable.
- [loop](https://github.com/jadmadi/opencode-loop): a prompt on a fixed cadence.
  A tick is skipped while the previous run is still active.
- [compose-next](https://github.com/jadmadi/opencode-compose-next): the
  spec to ship workflow as a skill (grill, spec, workspace, implement, verify,
  review, finalize, finish), with a read-only reviewer subagent.

## Safety

- [skip-permissions](https://github.com/jadmadi/opencode-skip-permissions):
  auto-approve permission prompts for one session. A permission hook upgrades an
  ask to an allow; a deny is final and never reaches the hook.
- [distill](https://github.com/jadmadi/opencode-distill): find repeated
  workflows in session history and propose skills, commands, or subagents.
  Nothing is written until an explicit apply, and nothing is overwritten.

## Commands and tools

| Name              | Plugin          | Purpose                                             |
| ----------------- | --------------- | --------------------------------------------------- |
| `/ds-go` `/ds` `/zai` `/oc-zen` `/oc-thinking` | model-switcher | Switch model, or cycle a group |
| `/context-limit`  | context-limit   | Show or set a working context budget                |
| `/max`            | max-mode        | Set the default best-of-N count                     |
| `best_of_n`       | max-mode        | Run candidates and return the judged winner         |
| `/remember` `/memory` `/checkpoint` | memory | Write or read the memory files           |
| `memory_read` `memory_append` `memory_search` | memory | Tool access to memory     |
| `/sila`           | sila-prime      | Run a sila subcommand (default: prime)              |
| `task`            | task-tool       | Add, update, list, or clear tree tasks              |
| `/goal`           | goal            | Set, show, or clear the stopping condition          |
| `/workflow`       | workflows       | Run a workflow, or list them                        |
| `/loop`           | loop            | Register, list, or stop a recurring prompt          |
| `/compose-next`   | compose-next    | Start the spec to ship workflow                     |
| `/skip-permissions` | skip-permissions | Toggle auto-approve for this session             |
| `/distill`        | distill         | Propose, list, or apply reusable artifacts          |

## Limits and overrides

- Goal, max mode, and distill call `ctx.generate.text`, which fails on OpenCode
  Go with `Request is missing x-opencode-session`. Point them at a working model
  with `GOAL_MODEL`, `MAX_MODE_MODEL`, and `DISTILL_MODEL`.
- The workflow runner uses child sessions and inherits the invoking session's
  model, so it needs no override.
- One loop per session, because the turn-end event names the session, not the
  run.
- `skip-permissions` never overrides a deny. `context-limit` never raises a
  window. `distill` writes only on an explicit apply and never overwrites.
- `sila-prime` calls the `sila` CLI. `SILA_BIN` points at the binary,
  `SILA_PRIME=off` disables the first-prompt briefing, and `SILA_PRIME_BUDGET`
  caps it. The session is marked after one attempt, so `/sila` refreshes on
  demand.
- sila's MCP server is configured in `~/.config/opencode/opencode.json` under
  `mcp.servers.sila`, with a permission block that allows seven tools
  (`prime_context`, `search_knowledge`, `search_messages`, `omni_search`,
  `get_latest_handoff`, `save_handoff`, `save_memo`) and denies the rest, so 41
  tools do not load into every prompt.

## Conventions

- One repo per plugin, MIT licensed, with README, AGENTS.md, CONTRIBUTING.md,
  LICENSE, and a NOTICE when content is adapted from another project.
- Single-file, dependency-free plugin. Export a plain `{ id, setup }` object and
  use Bun globals. Do not import `@opencode/plugin`.
- Tests with `bun test`.
- Built through the compose-next workflow: spec, workspace, implement, verify,
  review with the `reviewer` subagent, finalize, then a pull request.
- Adding, renaming, or changing a plugin means updating this map. The checklist
  is in `AGENTS.md`.
- Version in `package.json` plus an exported `VERSION` constant, kept equal by
  a test. Releases use CalVer `YYYY.MM.MICRO` and are tagged `v<version>` on
  main after merge, with a GitHub release. See `AGENTS.md`.
- Runtime API baseline in `runtime-baseline.json`, checked by
  `go run script/check-runtime-api.go` from the repo root. Run it before
  changing any plugin's `ctx` usage.

## Sources and license

Some plugins here are ports; each names its source and keeps the license notice
in its `NOTICE` file. The first program ported from MiMoCode
(https://github.com/XiaomiMiMo/MiMo-Code), MIT, Copyright (c) 2026 MiMo Code,
Xiaomi Corporation. MiMoCode also ships a `USE_RESTRICTIONS.md` and a trademark
policy. Check those before copying content. This note is a flag, not legal
advice.
