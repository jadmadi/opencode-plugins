# opencode-mimo-ports

Planning home for nine OpenCode V2 plugins ported from MiMoCode features.

This repo holds the program plan and one spec per feature. Each plugin gets its
own repo when it is built, using the same conventions as
[opencode-model-switcher](https://github.com/jadmadi/opencode-model-switcher)
and [opencode-compose-next](https://github.com/jadmadi/opencode-compose-next).

## Status

| # | Feature          | Spec                                       | Future repo                     | Status  |
| - | ---------------- | ------------------------------------------ | ------------------------------- | ------- |
| 1 | Task tool        | docs/compose/spec/task-tool.md             | opencode-task-tool              | delivered |
| 2 | Memory           | docs/compose/spec/memory.md                | opencode-memory                 | delivered |
| 3 | Goal             | docs/compose/spec/goal.md                  | opencode-goal                   | planned |
| 4 | Workflow runner  | docs/compose/spec/workflow-runner.md       | opencode-workflows              | planned |
| 5 | Loop             | docs/compose/spec/loop.md                  | opencode-loop                   | planned |
| 6 | Context limit    | docs/compose/spec/context-limit.md         | opencode-context-limit          | planned |
| 7 | Skip permissions | docs/compose/spec/skip-permissions.md      | opencode-skip-permissions       | planned |
| 8 | Max mode         | docs/compose/spec/max-mode.md              | opencode-max-mode               | planned |
| 9 | Distill          | docs/compose/spec/distill.md               | opencode-distill                | planned |

## Waves and order

| Wave | Features                  | Why                                        |
| ---- | ------------------------- | ------------------------------------------ |
| 1    | task tool, memory         | Foundation. compose-next assumes a tracker |
| 2    | goal, workflow runner     | Autonomy. Both build on the task tool      |
| 3    | loop, context-limit, skip-permissions, max mode, distill | Small, independent wins |

Build order: task tool, memory, goal, workflow runner, loop, context-limit,
skip-permissions, max mode, distill.

## Conventions for each port

- One repo per plugin, MIT licensed, with README, AGENTS.md, CONTRIBUTING.md,
  LICENSE, and a NOTICE when content is adapted from MiMoCode.
- Single-file, dependency-free plugin. Export a plain `{ id, setup }` object and
  use Bun globals. Do not import `@opencode/plugin`.
- Tests with `bun test`.
- Work through the compose-next workflow: spec, workspace, implement, verify,
  review with the `reviewer` subagent, finalize, then a pull request.

## Source and license

The features are inspired by MiMoCode (https://github.com/XiaomiMiMo/MiMo-Code),
MIT, Copyright (c) 2026 MiMo Code, Xiaomi Corporation. MiMoCode also ships a
`USE_RESTRICTIONS.md` and a trademark policy. Check both before copying any
content, and keep attribution in a NOTICE file. This note is a flag, not legal
advice.
