---
feature: max-mode
status: designed
updated: 2026-09-12
branch:
commits:
---

# Max Mode

## Report

## [S1] Problem

A hard prompt gets one attempt. There is no best-of-N sampling and no judge to
pick the best candidate. MiMoCode enables best-of-N through
`experimental.maxMode`, with a judge that selects the winner.

## [S2] Design

A command runs several candidates and a judge picks one.

- `/max <n>` sets the candidate count for the next turn. `/max off` clears it.
  `/max` prints the current setting.
- Text mode: the plugin runs `n` calls to `ctx.generate.text` in parallel, then
  one judge call with the same `generate.text`, and posts the winning text as
  the answer. This is the first and fully supported mode.
- Edit mode is phase two. It would run `n` child sessions with a subagent, take
  their diffs, and judge them. Applying a winner needs a merge step, so it stays
  out until the text mode works end to end.
- Cost: the command states the candidate count and that cost multiplies. A hard
  cap on `n` prevents accidents.

## [S3] Out of Scope

- Routing candidates across different models.
- Edit mode and diff merging.
- Automatic selection of when to use max mode.
- A TUI control.

## Tasks

- [ ] T1: the /max command and per-session settings - acceptance: set, print,
      clear, and reject an out-of-range count in a test (covers: S2)
- [ ] T2: parallel candidate generation - acceptance: a fake-context test
      confirms n calls run and all results are collected (covers: S2; depends:
      T1)
- [ ] T3: the judge call and winner selection - acceptance: a test supplies a
      judge answer and confirms the winner is posted with `ctx.session.prompt`
      (covers: S2; depends: T2)
- [ ] T4: README and NOTICE - acceptance: both files exist and name MiMoCode's
      max mode as the inspiration (covers: S2; depends: T3)
