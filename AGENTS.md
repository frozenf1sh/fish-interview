# Fish Interview Agent Guide

## Goal

Build an offline-first, Go-based interview knowledge explorer. It is a reference system, not a progress dashboard.

## Navigation

- `content/`: versioned knowledge cards, taxonomy, and exam-signal metadata.
- `internal/content/`: parsing, validation, links, and search indexing.
- `internal/web/`: HTTP handlers, templates, embedded static assets.
- `internal/trace/`: Go-generated algorithm execution traces.
- `docs/architecture.md`: boundaries, schemas, and extension rules.

## Constraints

- Keep Go as the only application language; use browser-native HTML/CSS and minimal vanilla JS only for interaction.
- Preserve stable content IDs. Do not use file paths as identifiers.
- A non-leaf taxonomy node may own a detailed mechanism card; do not model it as an empty folder.
- Keep the left navigation as an interactive node-and-edge tree canvas, with responsive fallback for small screens.
- Tree navigation must replace only the right content panel; do not reload or reset the tree canvas when opening a card.
- Company data is evidence-tagged metadata, never the primary information architecture.
- Do not add accounts, remote services, a database, progress tracking, or a frontend build chain without explicit approval.
- Prefer small, coherent commits. Do not combine content changes with unrelated refactors.

## Checks

Run before a code commit:

```sh
mise exec go@1.26 -- gofmt -w .
mise exec go@1.26 -- go test ./...
mise exec go@1.26 -- go vet ./...
git diff --check
```

## Adding content

1. Add a card under `content/cards/` with frontmatter matching the documented schema.
2. Add or verify its node in `content/taxonomy/tree.yaml`.
3. Keep algorithm concepts under `algo.knowledge`; add runnable templates only under `algo.patterns`. Every `algorithm-pattern` needs a `trace` and a matching Go trace test.
4. Start every algorithm pattern with `## 例题`; link its matching LeetCode problem when one exists, then give a minimal progressive example and precise Go comments. Do not add standalone instructions for watching its embedded animation.
5. Make engineering nodes mechanism-first: a parent can explain a subsystem, while leaves explain one concrete concept and link to related cards by ID.
6. Use card IDs for `links` and record `source`, `year`, `role`, and `confidence` for every exam signal.
7. Run the content validator.
