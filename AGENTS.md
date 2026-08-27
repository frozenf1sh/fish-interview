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

## Deployment

- k3s target: SSH `tc`, kubeconfig `~/.kube/config-tc.yaml`, namespace `fish-interview`, Traefik + cert-manager DNS-01.
- Fast path after a clean commit: `bash scripts/deploy-tc.sh`; for deliberate local/uncommitted testing use `IMAGE_TAG=local-<unique> bash scripts/deploy-tc.sh`.
- The script builds `linux/amd64`, imports the image into tc k3s, applies `deploy/`, and waits for the rollout. It does not change DNS.
- Cloudflare DNS: add a proxied `A` record `code` → `43.155.223.199` in [frozenf1sh.top DNS](https://dash.cloudflare.com/e4df346763044db08dfaa53e6aafcef8/frozenf1sh.top/dns/records). The Ingress requests `code-frozenf1sh-top-tls` from `letsencrypt-prod-dns01`.
- Verify: `kubectl --kubeconfig ~/.kube/config-tc.yaml -n fish-interview get pods,ingress,certificate`; then `curl -I https://code.frozenf1sh.top/`.
- Roll back app: `kubectl --kubeconfig ~/.kube/config-tc.yaml -n fish-interview rollout undo deployment/fish-interview`; DNS rollback means remove only the `code` record.
- A successful local build, image import, or Kubernetes rollout is not by itself public acceptance; verify DNS, certificate readiness, HTTPS redirect, and the public URL.

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
3. Keep algorithm concepts under `algo.knowledge`; add runnable templates only under `algo.patterns`. Every card beneath `algo.patterns` needs a `trace` and a matching Go trace test; use a state-level visualization when data structure state changes, and a flow trace only for template control flow.
4. Start every algorithm pattern with `## 例题`; link its matching LeetCode problem when one exists, then give a minimal progressive example and precise Go comments. Do not add standalone instructions for watching its embedded animation.
5. Make engineering nodes mechanism-first: a parent can explain a subsystem, while leaves explain one concrete concept and link to related cards by ID.
6. Use card IDs for `links` and record `source`, `year`, `role`, and `confidence` for every exam signal.
7. Run the content validator.
8. Visible card code and trace pseudocode must use literal Go operators (`<`, `>`, `&&`), never HTML entities. Render dynamic pseudocode through DOM text nodes, and scan all user-visible source before committing.
