---
feature: skip-permissions
status: designed
updated: 2026-09-12
branch:
commits:
---

# Skip Permissions

## Report

## [S1] Problem

OpenCode V2 permissions are configuration only. Turning them off at runtime
means editing config and restarting. MiMoCode offers a runtime toggle that
auto-allows actions while keeping explicit deny rules in force.

## [S2] Design

A command toggles auto-allow for one session.

- `/skip-permissions on` enables auto-allow for the current session.
  `/skip-permissions off` disables it. `/skip-permissions` prints the state.
- Auto-allow is applied with `ctx.permission.rules({ sessionID, permissions })`.
  Because session rules are evaluated after the agent rules and the last match
  wins, the plugin re-applies the configured deny rules after the allow rule so
  that denies still block.
- The plugin reads the configured deny rules from the resolved OpenCode config.
  A spike confirms the read path, since there is no rule-read API.
- Enabling prints a clear warning in the command result. The state is
  session-scoped and resets when the session ends. It never turns on by itself.

## [S3] Out of Scope

- The `--dangerously-skip-permissions` CLI flag and any global setting.
- Multi-session or persistent toggles.
- A TUI switch.
- Replacing the permission system.

## Tasks

- [ ] T0: spike how to read the configured deny rules - acceptance: the spec
      records the read path, or states that denies cannot be recovered and the
      command refuses to enable (covers: S2)
- [ ] T1: apply allow-all plus re-applied denies for one session - acceptance: a
      fake-context test asserts the rule order and that a deny still wins
      (covers: S2; depends: T0)
- [ ] T2: the /skip-permissions command with on, off, and status - acceptance: a
      test round-trips the state and shows the warning on enable (covers: S2;
      depends: T1)
- [ ] T3: README with the safety note and a NOTICE - acceptance: the README
      warns about the risk and the NOTICE names MiMoCode (covers: S2; depends: T2)
