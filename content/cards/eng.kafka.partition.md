---
id: eng.kafka.partition
kind: engineering
title: Kafka Partition：一册有序、可并行的日志
summary: Partition 是 Topic 下独立追加的日志；它同时决定记录的局部顺序、Consumer 的分配边界和 Kafka 能扩展的并行度。
parents: [eng.kafka.model]
tags: [kafka, partition, offset, ordering]
links: [eng.kafka.topic, eng.kafka.record, eng.kafka.broker, eng.kafka.offset-retention, eng.kafka.consumer-group, eng.kafka.ordering]
---

## 先把它想成“一册账本”

一个 Topic 可以有多册账本。每册只能在末尾追加新记录，记录一旦写入，就有一个只在本册内有意义的编号：**Offset**。

```text
orders / Partition 0
Offset： 0    1    2    3
记录：   创建  支付  拣货  发货
```

一条记录的完整地址要写成 `Topic + Partition + Offset`。单独说“Offset 3”不够，因为另一个 Partition 也可能有 Offset 3。

## Partition 同时提供三种边界

- **顺序边界**：同一册按追加顺序读取；Kafka 不承诺不同册之间谁先谁后。
- **并行边界**：不同册可以被不同 Consumer 并行处理；一个 Consumer Group 的并行度不能超过可用 Partition 数。
- **副本边界**：一册可以在多个 Broker 上保存副本，提升故障容忍度，但副本不是额外的消费并行通道。

## 为什么 Topic 没有全局顺序

如果 `A、B、C、D` 被分到两册：

```text
P0：A → C
P1：B → D
```

Kafka 只知道 P0 内 `A` 在 `C` 前、P1 内 `B` 在 `D` 前，不定义 `A、B、C、D` 的全局排列。如果订单的 `Created → Paid → Shipped` 必须有序，就用稳定的 `orderId` 作为 Key，让这些事件通常进入同一 Partition，见 [[eng.kafka.ordering]]。

## 常见误解

- 增加 Consumer，不能把同一 Partition 的串行处理变成多路并行。
- 增加 Replica 是为了容错，不是为了让同一条消息被组内多个 Consumer 同时处理。
- 增加 Partition 可能改变 Key 的映射；扩容前后的严格顺序要单独设计。

Partition 由 [[eng.kafka.broker]] 承载，具体哪一个 Consumer 负责读取由 [[eng.kafka.consumer-group]] 决定。
