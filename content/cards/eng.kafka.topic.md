---
id: eng.kafka.topic
kind: engineering
title: Kafka Topic：事件流的逻辑命名空间
summary: Topic 只负责组织一类事件，真正承载存储、顺序、并行和复制的是它下面的 Partition。
parents: [eng.kafka.model]
tags: [kafka, topic, log]
links: [eng.kafka, eng.kafka.partition, eng.kafka.offset-retention, eng.kafka.key-partitioner]
---

## 核心机制

Topic 是 Producer 和 Consumer 共同使用的逻辑名称，例如 `orders`。它不是一条单独的物理队列，而是一组 Partition 的集合：每条 Record 最终只追加到其中一个 Partition。

```text
Topic: orders
├── Partition 0：有序追加日志
├── Partition 1：有序追加日志
└── Partition 2：有序追加日志
```

Topic 层适合表达“这类事件属于哪里”；Partition 层才表达“事件存在哪里、以什么顺序读取、能并行到什么程度”。

## 不要把 Topic 当成消费队列

一个 Consumer Group 读过某条 Record，并不会让其他 Group 看不到它。不同 Group 的读取位置独立保存，消息是否仍可读取由 Retention 或 Compaction 决定，而不是由某个 Consumer 的 ACK 决定。

## 设计边界

- 需要同一业务实体有序时，先确定 Key，再评估 Key 到 Partition 的映射。
- 需要更多吞吐时通常增加 Partition，但这会增加副本、连接、分配和 Rebalance 成本。
- 需要 Replay 时必须先确认目标 Offset 仍未被 [[eng.kafka.offset-retention]] 清理。

Topic 只是入口；顺序和并行的真实边界在 [[eng.kafka.partition]]。
