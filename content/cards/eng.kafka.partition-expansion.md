---
id: eng.kafka.partition-expansion
kind: engineering
title: Kafka 扩大 Partition：吞吐增加也可能改变顺序
summary: 增加 Partition 会提高未来数据的并行机会，但可能改变 Key 映射、增加副本与协调成本，并破坏跨扩容时刻的实体顺序。
parents: [eng.kafka.operations]
tags: [kafka, partition, scaling, ordering]
links: [eng.kafka.key-partitioner, eng.kafka.ordering, eng.kafka.consumer-group, eng.kafka.performance.partition-parallelism]
---

## 先看一个订单跨扩容

扩容前，`order-123` 的事件都进入 P1；Topic 扩容后，分区器可能把新事件映射到 P4：

```text
扩容前：P1 = Created、Paid
扩容后：P4 = Shipped、Delivered
```

每个 Partition 内仍然有序，但 Kafka 不定义 P1 与 P4 之间的全局顺序。**Key 映射** 是 Producer 根据 Key 选择 Partition 的规则；Partition 数变化可能改变这个映射。

## 扩容带来什么收益和成本

收益是未来写入、读取和处理可以有更多并行通道；成本是更多 Replica、Metadata、连接、Group Assignment 和 Rebalance。已有数据不会自动均匀搬到新 Partition，新旧 Partition 还可能形成冷热不均。

扩容前先确认业务是否允许同一实体跨分区、是否可以停止写入完成切换，或是否需要版本化 Key/额外序列约束。不能只因为“机器还有资源”就扩大分区。
