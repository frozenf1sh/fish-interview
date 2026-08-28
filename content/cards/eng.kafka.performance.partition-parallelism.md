---
id: eng.kafka.performance.partition-parallelism
kind: engineering
title: Kafka Partition 并行：扩吞吐的基本单位
summary: Kafka 用多个 Partition 把写入、读取和复制工作分散到不同节点与 Consumer；但并行度受 Partition 数、热点 Key 和最慢分区限制。
parents: [eng.kafka.performance]
tags: [kafka, performance, partition, parallelism]
links: [eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.lag-scaling, eng.kafka.key-partitioner]
---

## 先看“多条传送带”

一个巨大的单日志像一条传送带；多个 Partition 像多条传送带：

```text
P0 → Consumer C1
P1 → Consumer C2
P2 → Consumer C3
```

Producer、Broker 和 Consumer 可以在不同 Partition 上并行工作，所以 Partition 是 Kafka 扩展吞吐的基本单位。

## 并行度不是无限增加

一个 Consumer Group 只有 4 个 Partition 时，放 8 个 Consumer，最多 4 个真正拥有分区；其余没有工作可分。更多 Partition 还意味着更多 Replica、连接、Metadata、Assignment 和 Rebalance 成本。

如果大量相同 Key 都落到 P0，P0 会成为热点：

```text
P0：██████████  很忙
P1：██
P2：█
```

此时增加 Consumer 不会拆开 P0。应先看 Key 分布、消息大小、写入速率和 Consumer Lag，再决定是否调整分区策略，见 [[eng.kafka.key-partitioner]]。
