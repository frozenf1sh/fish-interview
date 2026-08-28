---
id: eng.kafka.partition-expansion
kind: engineering
title: Kafka 扩大 Partition：吞吐增加与顺序风险
summary: 增加 Partition 可以提高未来数据的并行度，但可能改变 Key 映射、增加协调成本，并破坏跨扩容时刻的实体顺序。
parents: [eng.kafka.operations]
tags: [kafka, partition, scaling, ordering]
links: [eng.kafka.key-partitioner, eng.kafka.ordering, eng.kafka.consumer-group, eng.kafka.performance.partition-parallelism]
---

## 为什么会影响顺序

假设扩容前 `key=A → P1`，扩容后同一个 Partitioner 把新消息映射到 `P4`：

```text
P1：A1、A2、A3
P4：A4、A5、A6
```

每个 Partition 内仍然有序，但 Kafka 不定义 P1 与 P4 之间的交错顺序。扩容前应评估业务是否允许迁移、是否需要版本化 Key 或建立其他序列约束。

扩容还会增加 Replica、Metadata、Group Assignment 和 Rebalance 的规模，不能只按“机器够不够”决定。
