---
id: eng.kafka.ordering
kind: engineering
title: Kafka 消息顺序：先确定业务要保护哪种顺序
summary: Kafka 只在一个 Partition 内保证追加顺序；如果同一订单的事件要按先后处理，应使用稳定 Key 并评估重试、扩容和并发的影响。
parents: [eng.kafka.semantics]
tags: [kafka, ordering, partition, key]
links: [eng.kafka.partition, eng.kafka.key-partitioner, eng.kafka.partition-expansion, eng.kafka.producer.retry-idempotence]
---

## 先问“谁的顺序”

业务通常不是要求所有订单排成一条长队，而是要求同一个订单的状态不能倒退：

```text
order-123：Created → Paid → Shipped
order-456：Created → Cancelled
```

Kafka 的顺序单位是 **Partition**。Partition 是 Topic 下的一册追加日志；同一册内有先后，不同册之间没有天然的全局先后。

## Kafka 能保证什么

如果事件被放进：

```text
P0：order-123 Created → Paid → Shipped
P1：order-456 Created → Cancelled
```

Consumer 读取 P0 时能看到 P0 内的追加顺序，但 Kafka 不定义 P0 和 P1 之间谁先被观察。面试回答“Kafka 保证消息顺序”是不完整的，应该补上“Partition 内”。

## 如何保护同一实体

Producer 用稳定的业务 **Key**（例如 `orderId`）让同一订单通常映射到同一 Partition：

```text
orderId = 123 → Partitioner → P0
```

还要注意三件事：扩容可能改变 Key 映射；并发发送和失败重试可能带来观察顺序风险；热点 Key 可能让一个 Partition 变慢。具体路由见 [[eng.kafka.key-partitioner]]，扩容见 [[eng.kafka.partition-expansion]]。
