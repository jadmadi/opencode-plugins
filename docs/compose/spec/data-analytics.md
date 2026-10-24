# data-analytics spec

Status: delivered as a skills-only pack. Six analysis skills ship under
`skills/`, one directory per `da-` id. There is no plugin wrapper: one existed
through 0.2.1 and was removed, because the runtime command contract offers no
output channel for a status command and the skills already load by intent.

Map-level facts: six skills, no tools, no limits, no environment overrides.
Each skill sets `slash: true`, so it appears in the command catalog as
`/da-analyze-data-quality`, `/da-design-kpis`, `/da-metric-diagnostics`,
`/da-product-biz-analysis`, `/da-validate-data`, and `/da-visualize-data`. The
delivered copy lives at `docs/compose/spec/data-analytics.md` in
[jadmadi/opencode-data-analytics](https://github.com/jadmadi/opencode-data-analytics)
and carries the report and journey log.

## Skills

| Slash entry                | Skill folder               | Use it for                             |
| -------------------------- | -------------------------- | -------------------------------------- |
| `/da-analyze-data-quality` | `da-analyze-data-quality`  | Dataset trust, profiles, defect checks |
| `/da-design-kpis`          | `da-design-kpis`           | KPI frameworks, targets, measure plans |
| `/da-metric-diagnostics`   | `da-metric-diagnostics`    | Why a metric moved or disagrees        |
| `/da-product-biz-analysis` | `da-product-biz-analysis`  | Evidence-backed decisions              |
| `/da-validate-data`        | `da-validate-data`         | Correctness audits, normal or deep     |
| `/da-visualize-data`       | `da-visualize-data`        | Chart selection and visual standards   |

## Deferred

Dashboard, report, and data context stay deferred. They need app-building
runtime OpenCode does not have. Revisit when OpenCode surfaces cover durable
data apps.
