---
id: eng.kafka.ordering
kind: engineering
title: Kafka 消息顺序：用 Partition 内有序换取水平扩展
summary: Kafka 的顺序保证以 Topic-Partition 为边界；需要实体内有序时，应使用稳定 Key，而不是宣称 Topic 全局有序。
parents: [eng.kafka.semantics]
tags: [kafka, ordering, partition, key]
links: [eng.kafka.partition, eng.kafka.key-partitioner, eng.kafka.partition-expansion, eng.kafka.producer.retry-idempotence]
---

## 可承诺的顺序

同一 Partition 中，Consumer 看到的是日志追加顺序。订单系统可以让同一 `orderId` 的 Created、Paid、Shipped 进入同一 Partition，从而维持该订单内的顺序。

不同 Partition 之间没有天然全局顺序。Producer 并发请求、失败重试和扩分区还可能影响跨时刻的观察顺序，因此“Kafka 保证消息顺序”不是完整回答。

## 面试回答

Kafka 只保证 Partition 内有序；业务需要某个实体有序时，使用稳定业务 Key 让相关消息进入同一 Partition，并控制 Producer 并发与重试造成的顺序风险。
