---
id: eng.kafka.producer.accumulator
kind: engineering
title: Kafka RecordAccumulator：Producer 的待发送收件箱
summary: RecordAccumulator 是 Producer 内存中的待发送缓冲区；记录按 Partition 聚成 Batch，满足大小或等待条件后交给 Sender。
parents: [eng.kafka.producer]
tags: [kafka, producer, batch, linger]
links: [eng.kafka.producer.send, eng.kafka.producer.sender, eng.kafka.performance.batch-compression]
---

## 先用“按目的地分拣包裹”理解

Producer 同时发送订单、支付和库存事件。它不会把所有事件混成一袋，而是先按目标 Partition 分开：

```text
RecordAccumulator（待发送收件箱）
├── P0：Batch A = m1、m2、m3
├── P1：Batch B = m4、m5
└── P2：Batch C = m6
```

**RecordAccumulator** 就是这个收件箱；**Batch** 是同一个 Partition 的一小批待发送 Record。因为不同 Partition 的 Leader 可能在不同 Broker，按 Partition 分组后才能正确路由、压缩和保持 Partition 内顺序。

## Batch 什么时候发

- `batch.size`：一批最多希望积累到多大，单位是字节，不是固定消息条数。
- `linger.ms`：流量较低时，为了等更多记录进来，最多允许等待的时间窗口。

```text
消息少：m1 → 等待一小段时间 → m1、m2、m3 → 一起发送
消息多：m1、m2、… 很快填满 Batch → 尽快发送
```

`linger.ms=10ms` 不等于每条消息必然睡眠 10ms；Batch 已满或 Sender 有机会时可能更早发送。等待越多，通常越容易摊薄网络开销和提高压缩效率，但延迟和内存占用也可能增加。

## 缓冲区也会满

Batch 还没成功发送就会占用 Producer 缓冲区。Broker 变慢、网络重试或流量突然增大时，缓冲区可能耗尽，业务线程会感受到阻塞、超时或发送失败。这是背压：下游处理不过来，压力反过来限制上游，而不是无限堆内存。
