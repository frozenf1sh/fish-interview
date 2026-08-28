---
id: eng.kafka.producer.accumulator
kind: engineering
title: Kafka RecordAccumulator：按 Partition 聚合 Batch
summary: Producer 在内存中按目标 Partition 缓冲 Record，Batch 达到大小或等待窗口后交给 Sender 发送。
parents: [eng.kafka.producer]
tags: [kafka, producer, batch, linger]
links: [eng.kafka.producer.send, eng.kafka.producer.sender, eng.kafka.performance.batch-compression]
---

## 为什么按 Partition 缓冲

不同 Partition 可能由不同 Broker Leader 承载。Producer 需要把同一 Partition 的 Record 聚成一个 Batch，再按 Leader Broker 组织 ProduceRequest：

```text
P0 → Batch A ─┐
P2 → Batch B ─┴→ Broker 1
P1 → Batch C ───→ Broker 2
```

这样可以减少请求数、摊薄协议开销，并让压缩在一批相近数据上工作。

## `batch.size` 与 `linger.ms`

`batch.size` 描述一个 Partition Batch 的目标大小，`linger.ms` 描述为形成更好 Batch 所允许的等待窗口。流量高时 Batch 可能很快填满，低流量时等待可能增加延迟；它不是“每条消息固定睡眠这么久”。

Batch 未成功发送前仍占用 Producer 缓冲，缓冲耗尽时背压和超时会反过来影响业务线程。
