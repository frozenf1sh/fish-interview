---
id: eng.kafka.replication.isr-election
kind: engineering
title: Kafka ISR 与 Leader Election：只让跟得上的副本接任
summary: ISR 是与 Leader 保持同步边界的 Replica 集合；Leader 故障后，选举优先从可接受的同步副本中选择，避免盲目提升严重落后的副本。
parents: [eng.kafka.replication]
tags: [kafka, isr, election, availability]
links: [eng.kafka.replication.leader-follower, eng.kafka.replication.write-ack, eng.kafka.leader-routing]
---

## ISR 表达什么

```text
Leader offset=100
Follower A offset=100  → ISR
Follower B offset=99   → 可能仍在 ISR
Follower C offset=20   → 可能被移出 ISR
```

ISR 不是“所有曾经配置过的 Replica”，而是当前满足同步条件的集合，Leader 也在其中。Follower 长时间落后或不可用时，可能退出 ISR；恢复并追上后才可能重新加入。

## 为什么不能随便选

如果旧 Leader 已写到 `G`，某个副本只追到 `C`，直接提升它可能让 `D、E、F、G` 暂时不可见甚至丢失。Leader 选举、Producer Metadata 刷新和重试因此是同一条故障链上的三个环节。
