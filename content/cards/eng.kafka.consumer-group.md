---
id: eng.kafka.consumer-group
kind: engineering
title: Kafka Consumer Group：谁来读哪个 Partition
summary: Consumer Group 是一组协同读取 Kafka 的 Consumer；组内分摊 Partition，组间各自保存读取进度，因此同时具备负载均衡和发布订阅效果。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, offset]
links: [eng.kafka.partition, eng.kafka.group.pubsub-queue, eng.kafka.group.coordinator-assignment, eng.kafka.consumer.offset-commit, eng.kafka.rebalance]
---

## 先看两组读者

Topic `orders` 有 P0、P1、P2：

```text
inventory Group：       C1 → P0、P1；C2 → P2
recommendation Group：  C3 → P0；C4 → P1；C5 → P2
```

**Consumer Group** 是一组使用同一个 `group.id`、一起分工的 Consumer。一个 Partition 在同一个 Group 内同一时刻只交给一个 Consumer；但不同 Group 可以同时读取同一个 Partition。

## 为什么既像队列又像发布订阅

- Group 内：成员瓜分 Partition，每条记录通常只由其中一个成员处理，像“多个工人分工”，也叫竞争消费者。
- Group 间：每个 Group 有自己的读取位置，都能看完整份日志，像发布订阅。

消息没有被复制成每组一份，Group 只是各自保存自己的“书签”。位点和重放见 [[eng.kafka.consumer.offset-commit]] 与 [[eng.kafka.group.pubsub-queue]]。

## 并行度上限

一个 Group 有 4 个 Partition、8 个 Consumer，最多 4 个 Consumer 同时拥有分区，另外 4 个会空闲。增加 Consumer 只有在还有未分配 Partition、并且处理能力确实能扩展时才有用；某个热点 Key 让 P0 很忙时，增加成员不会把 P0 自动拆开。

成员变化后，谁拥有哪个 Partition 会重新计算，这叫 **Rebalance（重新分配）**，见 [[eng.kafka.rebalance]]。
