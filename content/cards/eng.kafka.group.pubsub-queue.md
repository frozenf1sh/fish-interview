---
id: eng.kafka.group.pubsub-queue
kind: engineering
title: Kafka Consumer Group：组间发布订阅，组内负载均衡
summary: Kafka 没有切换 Queue 或 Pub/Sub 的模式开关；同一个 Topic 配合不同 Consumer Group 自然形成组间独立消费、组内分工消费。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, pubsub, queue]
links: [eng.kafka.consumer-group, eng.kafka.partition, eng.kafka.event-streaming]
---

## 先看“一份日志，几枚书签”

```text
orders / P0：A、B、C、D
├── inventory Group：读到 C
├── recommendation Group：读到 A
└── risk Group：读到 D
```

**Partition** 是一份物理日志，Topic 是这些 Partition 的逻辑集合；**Consumer Group** 是一组共享分工和读取进度的读者。不同 Group 不会互相抢走记录，因为它们都从同一个 Partition 日志读取，只是各自保存 P0 的 Offset。

![Kafka 不复制消息，只为不同 Consumer Group 保存独立书签](/static/kafka-group-offset.drawio.svg)

## 两种语义从哪里来

- **组间像发布订阅**：inventory、recommendation、risk 都能从自己的位置读到 A、B、C、D。
- **组内像队列**：inventory Group 有 C1、C2 时，P0 只能在它们之间选择一个 owner；两人分摊不同 Partition，避免同一记录被组内重复处理。

这里的“队列”只是使用效果，不是 Topic 被读取后删除。日志仍由 Kafka 保存，保存多久由 Retention 或 Compaction 决定，见 [[eng.kafka.offset-retention]]。
