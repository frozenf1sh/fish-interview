# Architecture

## Product boundary

Fish Interview is a local, searchable reference system for interview preparation. Its first job is to answer “I forgot this—where is the shortest reliable explanation?” It deliberately excludes accounts, cloud sync, gamification, and progress analytics.

Company information is secondary evidence. A company can be used to filter or aggregate cards, but cannot own a knowledge tree.

## Layers

```text
content cards + taxonomy + exam signals
              │
              ▼
      content parser and validator
              │
              ├── knowledge graph / search index
              └── page view models
                         │
                         ▼
              Go HTTP server and templates
                         │
                         ▼
          browser-native tree, cards, SVG player
```

`internal/content` owns all file parsing and content integrity. `internal/trace` owns algorithm state transitions. `internal/web` only renders view models and never parses raw content itself. The left navigation is a browser-rendered, vertically oriented node-and-edge canvas derived from taxonomy data; it is not a nested directory list. Its first two levels are visible by default; deeper branches are expanded through their child-count badge.

Card navigation uses same-origin HTML partials (`/partials/cards/{id}`): selecting a tree node replaces only the right content panel. The canvas, its scroll position, selected tree, and expansion state stay in place. Algorithm traces are embedded in their algorithm card rather than presented as separate study pages.

## Content contracts

Every card is Markdown with YAML frontmatter:

```yaml
id: algo.greedy.interval-scheduling
kind: algorithm-pattern # concept | algorithm-pattern | engineering | ai-coding
title: 区间调度：按结束时间选择
summary: 选择最多互不重叠区间时，优先保留结束最早的可行区间。
parents: [algo.greedy.interval]
tags: [greedy, interval]
links: [algo.greedy.exchange-argument]
trace: interval-scheduling # optional, connects an algorithm leaf to a Go Trace
exam_signals:
  - company: meituan
    year: 2027
    role: backend
    confidence: medium
    source: https://example.invalid/source
```

Card IDs are stable public keys: file paths may change, IDs may not. `links` must resolve to card IDs. A taxonomy node can reference one card and still have child nodes: non-leaf concepts such as “贪心” and “Kafka” are first-class explanatory cards, not empty folders. The validator rejects duplicate IDs, broken references, tree cycles, and incomplete exam signals.

## Algorithm card structure

Algorithm-pattern cards should use these headings in order:

1. 识别信号
2. 建模与正确性直觉
3. Go 模板
4. 常见误区
5. 变体与关联

The order matters: sorting rules are consequences of a model, not the model itself.

## Reusable animation model

Algorithms produce a `trace.Trace` in Go. A trace contains a renderer kind, pseudocode lines, and replayable frames. A frame contains the active line, narration, variable values, and a renderer-specific state. Browser code can therefore play, step backward, or jump without re-running the algorithm.

The initial renderer kinds are `intervals`, `array`, and `dp-table`. Add a renderer only when several traces need a visual state that existing primitives cannot express.

## Extension rules

| Change | Required work |
| --- | --- |
| New knowledge card | Markdown card, taxonomy entry, validator pass |
| New exam signal | Source URL, year, role, confidence, validator pass |
| New algorithm | Card plus Go trace generator and trace tests |
| New visual primitive | Renderer contract, a representative trace, browser accessibility check |

## Persistence and release model

Versioned knowledge lives in Git. Personal scratch notes belong in `notes/private/`, which is intentionally ignored. The app loads content at startup; a later static-content compiler can be added without changing card contracts.

Each milestone must be independently runnable and must pass `go test ./...`, `go vet ./...`, and `git diff --check` before commit.
