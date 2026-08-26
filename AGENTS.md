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
3. Use card IDs for `links` and record `source`, `year`, `role`, and `confidence` for every exam signal.
4. Run the content validator.

