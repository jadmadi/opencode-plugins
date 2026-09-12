---
feature: workflow-runner
status: designed
updated: 2026-09-12
branch:
commits:
---

# Workflow Runner

## Report

## [S1] Problem

OpenCode V2 has no deterministic multi-agent workflows. Every run is
conversational, so a known phase sequence cannot run unattended. MiMoCode ships
JavaScript workflows (deep-research, fact-check, research-experiment, compose)
that encode fixed phases with bounded retries and automatic parallelization.

## [S2] Design

A plugin provides a small runner and two built-in workflows.

- A `workflow` tool and a `/workflow` command start a run. Input names the
  workflow and the task text.
- The runner executes ordered phases. A phase may fan out to parallel child
  sessions with a concurrency cap, then join.
- Each child session is created with `ctx.session.create`, given an agent with
  `ctx.session.switchAgent`, prompted with `ctx.session.prompt`, and awaited
  with `ctx.session.wait`. Child results are read with `ctx.session.context`.
- Phases write artifacts to a run directory under the data dir, so a run is
  resumable and auditable.
- Bounded retries: a phase retries on a failed result up to a small cap, then
  stops the run and reports.
- Built-ins to ship first: deep-research (brief, plan, research, reflect, write,
  review) and fact-check (plan, search, extract, group, crosscheck, report).
  Both read and write files and need no user interaction.

## [S3] Out of Scope

- A sandboxed JavaScript runtime for user-authored workflows. Phase two.
- The full compose pipeline and the research-experiment loop.
- A TUI progress view.
- Cost accounting.

## Tasks

- [ ] T1: runner core with ordered phases, parallel fan-out, a join, and
      `ctx.session.wait` - acceptance: a fake-context test runs a two-phase
      workflow with one fan-out and asserts the phase order and the joined
      result (covers: S2)
- [ ] T2: artifact directory and resume support - acceptance: a test writes an
      artifact per phase and a re-run resumes from the last completed phase
      (covers: S2; depends: T1)
- [ ] T3: bounded retries and a clear failure report - acceptance: a test forces
      a phase to fail twice and confirms the run stops with a report (covers:
      S2; depends: T1)
- [ ] T4: the deep-research workflow - acceptance: a fake-context run reaches
      the write phase and produces one report file (covers: S2; depends: T2)
- [ ] T5: the fact-check workflow - acceptance: a fake-context run reaches the
      report phase and produces a verdict file (covers: S2; depends: T2)
- [ ] T6: README, NOTICE, and tests for both built-ins - acceptance: files
      exist, tests pass (covers: S2; depends: T4, T5)
