---
feature: distill
status: designed
updated: 2026-09-12
branch:
commits:
---

# Distill

## Report

## [S1] Problem

Repeated manual workflows stay manual. A user runs the same sequence of steps
again and again, and nothing turns it into a reusable skill, command, or agent.
MiMoCode's `/distill` finds repeated workflows in recent sessions and packages
the strong candidates.

## [S2] Design

A command proposes reusable artifacts from recent session history.

- `/distill` gathers recent sessions for the current project with
  `ctx.session.list` and reads their content with `ctx.session.context`.
- A model, through `ctx.generate.text`, finds repeated multi-step patterns and
  returns candidate artifacts: a skill, a command, or a subagent, each with a
  name, a purpose, and the steps it would encode.
- The command prints the candidates and their confidence. Nothing is written
  without user approval. A follow-up prompt lets the user pick which to write.
- Approved artifacts are written under `~/.config/opencode/` (`skills/<name>/SKILL.md`,
  `commands/<name>.md`, or `agents/<name>.md`) with the same frontmatter rules
  the other ports use.
- A candidate that repeats an existing skill or command is reported as a
  duplicate, not written.

## [S3] Out of Scope

- Automatic writes without approval.
- Memory updates from traces. That is the memory port's `/dream` direction.
- Ranking quality beyond a simple confidence field.
- Cross-project distillation.

## Tasks

- [ ] T1: history gathering with a size cap - acceptance: a fake-context test
      collects recent sessions and truncates at the cap (covers: S2)
- [ ] T2: pattern detection and candidate parsing - acceptance: a stub model
      result parses into skills, commands, and agents with confidence (covers:
      S2; depends: T1)
- [ ] T3: the approval prompt and the write step - acceptance: a test approves
      one candidate, writes the file, and skips a duplicate (covers: S2;
      depends: T2)
- [ ] T4: README and NOTICE - acceptance: both files exist and name MiMoCode's
      distill feature (covers: S2; depends: T3)
