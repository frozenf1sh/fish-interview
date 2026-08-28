---
id: eng.kafka.producer.sender
kind: engineering
title: Kafka Producer Sender Thread：把 Batch 送到 Leader
summary: Sender Thread 从本地缓冲取出可发送 Batch，按 Broker 合并 ProduceRequest，接收 ACK 并推动重试、Metadata 刷新和 Future 完成。
parents: [eng.kafka.producer]
tags: [kafka, producer, sender, io]
links: [eng.kafka.producer.accumulator, eng.kafka.leader-routing, eng.kafka.replication.write-ack]
---

## 先看谁在做网络发送

调用 `send()` 的业务线程主要负责准备 Record；Producer 内部的 **Sender Thread**（发送线程）负责持续从 Accumulator 取货：

```text
业务线程：序列化 → 选 Partition → 放入缓冲
Sender：  取 Batch → 找 Leader → 发请求 → 等响应 → 成功/重试/失败
```

如果 P0 和 P2 的 Leader 都在 Broker 1，Sender 可以把这两个 Partition 的 Batch 放进同一个 **ProduceRequest**（写入请求）；P1 的 Batch 则发给另一个 Broker。这样比一条 Record 一个请求更省网络往返。

## Broker 返回什么

Partition Leader 收到请求后追加日志，并为 Record 分配 Offset；等到满足 Producer 的 ACK 条件后返回结果。Sender 再把结果交给对应 Future 或 Callback：

```text
成功 → Future 完成
Leader 信息过期 → 刷新 Metadata 后再尝试
暂时故障 → 按重试策略尝试
不可恢复错误 → 报告失败
```

**Metadata** 是 Producer 的集群路由表，**ACK** 是 Broker 对写入请求的确认。它们的解释分别见 [[eng.kafka.bootstrap-metadata]] 和 [[eng.kafka.replication.write-ack]]。

## 为什么 Sender 是性能关键

它把多个 Partition 的 Batch 组织成较少的网络请求，同时处理连接、路由、响应和重试。Kafka 的高吞吐不是 `send()` 这个方法本身神奇，而是业务线程、缓冲、批量和后台 I/O 共同形成的流水线。
