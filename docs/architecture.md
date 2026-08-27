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
parents: [algo.patterns.greedy]
tags: [greedy, interval]
links: [algo.greedy.exchange-argument]
trace: interval-scheduling # required for algorithm-pattern cards
exam_signals:
  - company: meituan
    year: 2027
    role: backend
    confidence: medium
    source: https://example.invalid/source
```

Card IDs are stable public keys: file paths may change, IDs may not. `links` must resolve to card IDs. A taxonomy node can reference one card and still have child nodes: non-leaf concepts such as “贪心” and “Kafka” are first-class explanatory cards, not empty folders. The validator rejects duplicate IDs, broken references, tree cycles, and incomplete exam signals.

## Algorithm knowledge and pattern cards

The algorithms root is intentionally split into two branches:

```text
算法
├── 知识       # 原理、建模、证明语言和边界
│   ├── 贪心原理
│   └── 动态规划
└── 题型       # 可直接套用并验证的输入/输出模式
    ├── 贪心 / 区间调度
    ├── 动态规划 / 线性 DP
    └── 二分答案
```

A **knowledge card** explains one reusable mechanism. Start with its conclusion, derive it through one minimal example, then give the boundary where it stops applying. A non-leaf node may own such a card: for example, the “贪心原理” node explains safe local choices while its leaf cards hold specific templates.

An **algorithm-pattern card** is a runnable retrieval unit. It must contain a `trace`, a short recognition condition, one progressively derived example, annotated Go code, and only the nearest confusions. The player is rendered before the Markdown body, so the card must not repeat generic text such as “how to watch the animation.” The content validator rejects an `algorithm-pattern` without a trace.

Write code comments next to the state definition, base case, transition, and boundary update. Comments should explain the role of the line, not restate its syntax.

## Reusable animation model

Algorithms produce a `trace.Trace` in Go. A trace contains a renderer kind, pseudocode lines, and replayable frames. A frame contains the active line, narration, variable values, and a renderer-specific state. Browser code can therefore play, step backward, or jump without re-running the algorithm.

Current renderer kinds are `intervals`, `dp-table`, `dp-grid`, `stock-state`, `bitmask-state`, and `binary-red-blue`. `dp-grid` is shared by LCS, interval DP, and path DP; add a renderer only when several traces need a visual state that existing primitives cannot express. For a new pattern: (1) add its trace generator under `internal/trace`, (2) test its final frame and key state transition, (3) expose it from `internal/web/server.go`, and (4) reference it from the card frontmatter.

For binary search, use a named color invariant instead of a mix of closed-interval templates. The standard minimum-feasible template keeps `(red, blue]`: `red` is a known invalid sentinel, `blue` is known valid, and the answer is `blue`. If a variant reverses the objective, declare the colors and return endpoint before writing its loop.

## Engineering mechanism cards

Engineering trees use the same contract, but their leaves are concepts rather than exercises. A parent card explains a subsystem boundary and the mechanism that joins its children; each leaf answers one question such as “consumer group,” “partition,” or “rebalance.” Prefer one request flow, state transition, or failure table over a list of disconnected facts. Use `links` for cross-tree references so a card can point directly to the exact prerequisite or related failure mode.

## Extension rules

| Change | Required work |
| --- | --- |
| New knowledge card | Markdown card, taxonomy entry, one mechanism/example/boundary, validator pass |
| New exam signal | Source URL, year, role, confidence, validator pass |
| New algorithm pattern | Pattern card, embedded Go trace generator, trace tests, annotated Go template |
| New engineering concept | Leaf card, taxonomy entry, parent mechanism and `links` checked |
| New visual primitive | Renderer contract, a representative trace, browser accessibility check |

## Persistence and release model

Versioned knowledge lives in Git. Personal scratch notes belong in `notes/private/`, which is intentionally ignored. The app loads content at startup; a later static-content compiler can be added without changing card contracts.

Each milestone must be independently runnable and must pass `go test ./...`, `go vet ./...`, and `git diff --check` before commit.
