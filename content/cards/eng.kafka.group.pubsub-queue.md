---
id: eng.kafka.group.pubsub-queue
kind: engineering
title: Kafka Consumer Group：发布订阅与组内负载均衡
summary: Kafka 没有 Queue、Pub/Sub 模式开关；Topic、Partition 和 Consumer Group 的组合自然产生两种消费语义。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, pubsub, queue]
links: [eng.kafka.consumer-group, eng.kafka.partition, eng.kafka.event-streaming]
---

## 两个维度

```text
Topic orders / P0
├── Group inventory：C1 读取
├── Group recommendation：C2 读取
└── Group risk：C3 读取
```

不同 Group 都可以从同一物理 Partition 读取，各自维护自己的 Offset，所以表现为发布订阅。同一 Group 内，多个 Consumer 瓜分 Partition，一个分区同一时刻只交给一个成员，所以表现为队列式负载均衡。

![Kafka 不复制消息，只为不同 Consumer Group 保存独立书签](/static/kafka-group-offset.drawio.svg)

关键不是复制消息，而是复制“读取视角”：消息只有一份日志，Group 各自拥有书签。
