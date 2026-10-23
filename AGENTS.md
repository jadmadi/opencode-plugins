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

- HOLD means the baseline still matches. Keep the model/provider domains and
  export `Plugin.define({ id, setup })` from the host `@opencode/plugin`
  module.
- MIGRATE means the runtime or the published package moved. Re-probe with the
  steps below, then update the baseline and the plugins together.
- UNKNOWN means the check could not reach the runtime or the registry. Do not
  treat it as HOLD.

Baseline as of 2026-10-23, OpenCode 2.1.3 (re-probed):

- `ctx.model` carries model reads and transforms (members: `default`, `list`,
  `reload`, `transform`), and `ctx.provider` carries provider reads and
  transforms (members: `get`, `list`, `reload`, `transform`). Both take
  `transform(editor)` callbacks.
- `ctx.catalog` does not exist. Migrate `ctx.catalog.model.*` to
  `ctx.model.*` and `ctx.catalog.provider.*` to `ctx.provider.*`.
- `import { Plugin } from "@opencode/plugin"` resolves for local file plugins
  on 2.1.3. Export `Plugin.define({ id, setup })`. Each server plugin repo
  stubs the module in `test-preload.ts` so `bun test` runs outside OpenCode.
- `ctx.credential` is available.
- `ctx.session` has `list`, `remove`, and `compact`.
- `ctx.permission.assert` is available.

Probe recipe when the check says MIGRATE:

1. Drop a temporary plugin in `~/.config/opencode/plugins/` that writes
   `typeof ctx.catalog`, `typeof ctx.model`, `typeof ctx.provider`,
   `typeof ctx.credential`, `typeof ctx.session`, and the member keys of each
   domain to a file. Touch it, read the file, then remove it.
2. Drop a directory plugin importing `{ Plugin } from "@opencode/plugin"` and
   see whether it loads. Call `Plugin.define({ id, setup })` and check it
   returns the same object.
3. Record the results in `runtime-baseline.json`. If the import resolves and
   `ctx.model` exists, keep `Plugin.define` and the model/provider domains,
   then bump versions, update the changelog, tag, and release.

Known stale text: `opencode-sila` is a multi-file TUI plugin and is not one of
the 12 server plugins migrated to `Plugin.define`. Its server entry still uses
a plain `{ id, setup }` object, and its AGENTS.md says the runtime does not
resolve `@opencode/plugin`. Update that repo when it is next changed.

## Conventions

- Short lines, 80 to 110 characters, with H2 or larger headings.
- No bold and no em dashes.
- One plugin per bullet, with its full GitHub URL.
- Changes through a branch and a pull request.
- Plugin repos carry no agent provisioning: no `.cursor/`, `.devin/`,
  `CLAUDE.md`, `goals/`, or private files.
- `main` is the published channel for `github:` installs: every commit reaches
  users. Merge only when the tests pass.
- No GitHub Actions. Each repo installs the sila pre-push gate, so the tests run
  locally before the push, and `main` requires a pull request.
- Every plugin repo is licensed AGPL-3.0-only.
- Plugins carry no dependencies. The test run writes `aube-lock.yaml`; ignore
  it, do not commit it.

## Versioning and releases

- `package.json.version` is the source of truth. Each plugin exports `VERSION`
  from its `.ts` file, and a test asserts the two match. Drift fails
  `bun test`.
- Each release records the OpenCode version it was built and tested under, as
  `builtUnder` in `package.json` plus a `builtUnder-<version>` tag, for example
  `builtUnder-2.1.3`.
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
  `v2026.10.0`) and publish a GitHub release with short notes.
- Install docs keep the `main` URL and add one line: pin a release by replacing
  `main` with a tag such as `v2026.10.0`.
- The first CalVer release is `2026.9.0` for every plugin.

## Releasing

- Semantic commit messages.
- Update this file when the map's process changes.
- After a plugin change merges, bump its version, update `CHANGELOG.md`, tag
  `v<version>`, and publish a release. The versioning rules are above.

<!-- project-sync:begin -->
### Continuous Project Metadata Maintenance (`project.jsonc`)
This repository maintains an authoritative `project.jsonc` at its root describing its title, tagline, version, capabilities, features, installation commands, tech stack, and FAQs for jadmadi.net synchronization.
Whenever you implement a new feature, add a CLI command or flag, bump the version, or modify installation instructions:
- You MUST update `project.jsonc` at the repository root in the same commit.
- Update `lastEdit` to the current local timestamp (`YY-MM-DD-HH-MM`).
- Keep `version`, `features`, `stats`, and `install_command` strictly synchronized with the code.
<!-- project-sync:end -->

<!-- sila:begin (managed by sila; content inside is regenerated by `sila install` / `sila agents sync`) -->
## 0. Safety Rails (NEVER Do This)
- No unprompted remote push: Local commits are autonomous; remote push requires explicit user consent.
- No unauthorized repo creation or namespace re-homing: NEVER run 'gh repo create' or push external/organizational specifications to personal namespaces; fail closed on push/auth hurdles.

## 1. Working Relationship & Discipline
- Simplified, concise, and straight to the point.
- Verified Repo Runbook Primacy: Never guess commands; verify exit code 0 locally.
- Closed-Loop Goal Attribution (6 Steps): When completing a goal package under `goals/<slug>/`:
  1. Mark every completed task `- [x]` in `goals/<slug>/plan.md` (unchecked tasks block the `implemented` transition).
  2. Write verification evidence with an explicit standalone `Verdict: PASS` declaration when using markdown evidence.
  3. Commit the implementation with `(goals/<slug>/goal.md)` or `[goal: <slug>]` under a `feat:`, `fix:`, `refactor:`, `perf:`, or `chore:` prefix (`docs:` and `test:` commits are ignored by the goal matcher).
  4. Certify the goal with `sila goals certify <slug> --evidence=goals/<slug>/evidence.md` or `sila goals audit <slug>`.
  5. Run `sila goals` to reconcile `goals/README.md` and commit the updated registry and certification artifacts.
  6. Verify the stop condition: `sila goals --json` reports `status: "implemented"`, `verification_status: "verified"`, and `unverified_count: 0` for delivered goals, with a clean `git status`.

---

## Sila Trigger Contract

- First command in a repository: `sila salam` prints this loop with live next steps; run `sila --help` to inspect all available tools.
- Task start: run `sila heal` (idempotent; sila's rules, plugin, and MCP entry only, no project git hooks; preview with `--dry-run`), commit refreshed managed files, then `sila prime` and `sila goals ready`.
- Goal completion (6 steps): check every `plan.md` task -> write evidence with `Verdict: PASS` -> commit with `(goals/<slug>/goal.md)` under `feat|fix|refactor|perf|chore` -> `sila goals certify <slug> --evidence=...` or `sila goals audit <slug>` -> run `sila goals` to sync the registry -> verify `sila goals --json` reports `unverified_count: 0` and `git status` is clean.
- Stuck on an error, a package, or an unfamiliar API: run `sila recall "<error or symptom>"` before web search or a workaround.
- Problem solved: record it with `sila lesson "<symptom>" --cause=... --fix=... --files=...`.
- Decision made: record it with `sila memo "<directive>"`.
- Session end: `sila handoff "<summary>" --open="..."` and `sila sync` (or the session-end hook).
- Workspace search (`fffind`, `ffgrep`, ripgrep, rtk) finds code in the current tree; `sila` is cross-harness continuity: past sessions, decisions, and lessons. Use both.
- Sila itself gets in the way or misses something: run `sila feedback "<what happened>" --type=issue|tip|idea` and keep working.
<!-- sila:end -->

<!-- sila:begin (principles) -->
# Sila Universal Principles Pack & Engineering Constitution

## 1. Layer 1: Radical Simplicity (KISS & YAGNI - The Baseline Default)

**Philosophy**: Write only what is strictly necessary today. Reject speculative abstractions, hollow interfaces, premature generalization, and configurable frameworks before real demand exists.

- Write simple, pure functions before designing multi-layer abstractions or class hierarchies.
- Do not add configuration options, extension points, or plugins for hypothetical future requirements.
- Delete dead code, obsolete parameters, and speculative stubs immediately.
- Keep the dependency footprint minimal; prefer standard library solutions.
- *Anti-Pattern to Avoid*: Speculative enterprise OOP hierarchies, hollow wrapper structs, premature generalization, configurable frameworks with single implementations.

## 2. Layer 2: Human & Agent Ergonomics (CUPID - Surface & Interfaces)

**Philosophy**: Surfaces, APIs, and CLI interfaces must optimize for agent and developer joy: Composable, Unix-like, Predictable, Idiomatic, and Domain-based.

- Composable: Small, independent pieces that fit together cleanly with clear input/output contracts.
- Unix Philosophy: Do one thing well; support standard stdio streams, composable pipelines, and JSON output.
- Predictable: Consistent error returns, deterministic behavior, zero hidden side effects, idempotent mutations.
- Idiomatic: Follow the natural conventions, formatting, and idioms of the target language.
- Domain-Based: Model business reality directly in terms of the problem domain without artificial design-pattern jargon.
- *Anti-Pattern to Avoid*: Inconsistent return envelopes, hidden global state mutations, non-standard naming conventions, generic catch-all manager/helper classes.

## 3. Layer 3: Pragmatic Structural Seams (Pragmatic SOLID - Module Boundaries)

**Philosophy**: Enforce structural discipline at module boundaries without falling into dogmatic OOP abstraction traps.

- Single Responsibility Principle (SRP): Enforce file weight and complexity ceilings (<1,000 LOC; <10,000 tokens). Every module owns one cohesive concept.
- Interface Segregation Principle (ISP): Consume tiny, focused interfaces (1-3 methods max, e.g. io.Reader, io.Closer) at package boundaries. Return concrete types, accept interfaces.
- Pragmatic Substitution (LSP): Prefer concrete struct composition over deep inheritance trees or multi-layer polymorphic hierarchies.
- Explicit Seams: Separate pure computational logic from I/O and persistence layers.
- *Anti-Pattern to Avoid*: God files (>1,000 LOC), kitchen-sink interfaces with a dozen methods, deep inheritance trees, abstract base classes that only have one concrete child.

## 4. Layer 4: Data Flow & Concurrency (Functional Programming - FP)

**Philosophy**: Structure data flow through pure transformations, append-only immutable logs, and strict concurrency safety.

- Pure Transforms: Write stateless transform functions at data, JSON, and I/O boundaries.
- Immutable Ledgers: Treat historical state and logs as append-only value objects.
- Zero Shared Mutable State: Eliminate shared mutable variables across concurrent goroutines, workers, or asynchronous tasks.
- Concurrency Discipline: Always bound concurrency, handle cancellation contexts, and ensure goroutine leak prevention.
- *Anti-Pattern to Avoid*: Unsynchronized global variables, in-place mutation of shared data structures across concurrent workers, unbuffered goroutines without context cancellation.

## 5. Layer 5: Epistemic Grounding (Digital Data Sovereignty)

**Philosophy**: Logs record claims; Code is truth. Every architectural claim and instruction must be grounded in live filesystem reality.

- Grounding Verification: All architectural claims, paths, and instructions must reference live files (file_refs_exist = 1).
- Never Hide Stale State: Flag, surface, or prune stale claims and drifting documentation rather than masking errors.
- Closed-Loop Attribution: Complete the verification loop for every task before claiming completion.
- Data Sovereignty: Respect private boundaries; never leak confidential intelligence, secrets, or keys into public repositories.
- *Anti-Pattern to Avoid*: Hallucinated file paths, phantom dependencies, unverified claims, blind copy-pasting from outdated documentation without local verification.

---

## Active 9-Verb Memory & Knowledge Protocol

AI coding agents MUST actively use Sila's 9-verb continuity toolchain during development:

1. **Instant Pre-Search (`sila recall "<symptom or error>"`)**:
   Before guessing, hallucinating, or initiating expensive external web searches, agents MUST run `sila recall "<symptom or error>"` to check if the problem was previously solved in this or any registered project.
2. **Upstream Specification Grounding (`sila doc search <query>`)**:
   Before writing integrations, scaffolding architectures, or changing external dependencies, agents MUST query `sila doc search <query>` to ground decisions in authoritative local docs.
3. **Durable Invariants Capture (`sila memo "<directive>" --files=...`)**:
   When an architectural boundary, policy invariant, or non-negotiable rule is decided, agents MUST persist it via `sila memo "<directive>" --files=...`.
4. **Empirical Bug-Fix Recording (`sila lesson "<symptom>" --cause="..." --fix="..." --files=...`)**:
   When resolving a non-trivial bug, race condition, or edge-case failure, agents MUST record a structured lesson via `sila lesson "<symptom>" --cause="..." --fix="..." --files=...`.
5. **Architectural Pattern Recording (`sila case "<intent>" --approach="..." --invariant="..." --files=...`)**:
   When introducing a structural design pattern, module seam, or system abstraction, agents MUST record a case via `sila case "<intent>" --approach="..." --invariant="..." --files=...`.
6. **Cryptographic Grounding Verification (`sila facts verify`)**:
   Before declaring a task or goal complete, agents MUST run `sila facts verify` to ensure all referenced files and dependencies exist on disk.
7. **Context Priming (`sila prime`)**:
   At the start of every session or task, agents run `sila prime` to load active context, ongoing goals, and architectural constraints.
8. **Deterministic Quality Gate (`sila gate`)**:
   Before committing, agents MUST run `sila gate` (or `sila gate --pre-commit`) to execute automated deterministic sentries, race detection, and style checks.
9. **Structured Session Handoff (`sila handoff "<summary>" --open="..."`)**:
   At the end of every work session, agents MUST persist progress, unblocked next steps, and open questions via `sila handoff "<summary>" --open="..."`.

---

## Stack-Specific Invariants (Go)

- 1. Code Review & Hygiene: maintain pure Go standard library with zero third-party dependencies, Pure-Go CGO-Free Portability with CGO_ENABLED=0, and hard file ceilings under 1,000 LOC.
- 2. Goroutine & Context Safety: require context.Context on all background routines; declare defer cancel() immediately after timeout context creation; ensure goroutine leak prevention.
- 3. Goal Scaffolding: break findings into discrete single-outcome packages under goals/<slug>/ with goal.md, facts.md, plan.md, and explicit dependencies.
- 4. Test-Driven Implementation: claim goal locks before editing; write table-driven tests verifying both acceptance assertions and negative rejection paths.
- 5. Quality Gates & Race Detector Discipline: verify changes locally with linter scripts and full race detection ('go test -race ./...') with zero tolerance for data races.
- 6. Attributed Commits: use semantic prefixes containing the goal slug: feat(scope): desc [goal: <slug>]; run sila goals audit to certify receipts.
- 7. Git Branch Discipline: develop on dedicated feature branches; never push unprompted to remotes; re-verify race tests on merged main.
- 8. Security & Trust Boundaries: validate all CLI flags and stdin readers; enforce UID ownership; use mode 0600 on sensitive temporary storage.
- 9. Release Lifecycle & Value Changelog: implement a 24-hour cached daily version check, offer direct upgrades, provide a human-readable benefits-focused changelog, and credit contributors.
<!-- sila:end -->
