---
id: eng.kafka.partition
kind: engineering
title: Kafka Partition：顺序、并行与 Offset 的边界
summary: Partition 是一条可追加的有序日志，也是 Kafka 扩展吞吐、分配消费者和维护副本的最小共同单位。
parents: [eng.kafka.model]
tags: [kafka, partition, offset, ordering]
links: [eng.kafka.topic, eng.kafka.broker, eng.kafka.offset-retention, eng.kafka.consumer-group, eng.kafka.ordering]
---

## 核心机制

每个 Partition 都是一条独立的、持续追加的日志。Broker 为其中的 Record 分配递增 Offset；Offset 只在该 Partition 内有意义，因此完整定位必须写成 `Topic + Partition + Offset`。

Partition 同时承担三种边界：

| 边界 | 含义 |
| --- | --- |
| 顺序 | Consumer 读取同一 Partition 时看到写入顺序 |
| 并行 | 不同 Partition 可以由不同 Consumer 并行处理 |
| 高可用 | 一个 Partition 的 Replica 分布在多个 Broker 上 |

## 为什么没有 Topic 全局顺序

如果 `orders` 有三个 Partition，`P0` 中的 `A → C` 和 `P1` 中的 `B → D` 各自有序，但 Kafka 不定义 `A、B、C、D` 的全局交错顺序。想让同一订单有序，应使用稳定的 `orderId` 作为 Key，让相关事件进入同一 Partition。

## 常见误区

- 增加 Consumer 不能突破单个 Partition 的串行边界。
- Replica 增加的是容错能力，不会让一个 Partition 自动变成多路并行。
- 扩大 Partition 数可能改变 Key 映射，跨扩容时的严格顺序需要单独设计。

Partition 的物理承载者是 [[eng.kafka.broker]]，消费归属由 [[eng.kafka.consumer-group]] 管理。
