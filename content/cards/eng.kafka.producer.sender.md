---
id: eng.kafka.producer.sender
kind: engineering
title: Kafka Producer Sender Thread：把 Batch 变成网络请求
summary: 业务线程负责准备 Record，后台 Sender Thread 负责从 Accumulator 取可发送 Batch、按 Broker 合并请求并处理响应。
parents: [eng.kafka.producer]
tags: [kafka, producer, sender, io]
links: [eng.kafka.producer.accumulator, eng.kafka.leader-routing, eng.kafka.replication.write-ack]
---

## 两类线程的分工

```text
业务线程：serialize → partition → append accumulator
Sender：  ready batch → metadata → ProduceRequest → response/retry
```

Sender 可以把同一 Broker 上多个 Partition 的 Batch 放进请求，避免每条消息一次网络往返。Broker 收到请求后，由目标 Partition Leader 追加日志并给新增 Record 分配 Offset。

## 什么时候算完成

Sender 收到 Broker 响应后，才会完成对应 Future 或回调。响应可能表示成功、可刷新 Metadata 的错误、可重试故障或不可重试失败；因此发送路径必须和 ACK、重试以及幂等性一起阅读。
