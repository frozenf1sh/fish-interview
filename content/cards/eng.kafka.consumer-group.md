---
id: eng.kafka.consumer-group
kind: engineering
title: Kafka Consumer Group：消费进度与组内分工
summary: 同一个 Consumer Group 共享一套分区归属和消费进度；一个 Partition 同一时刻只分给组内一个 Consumer。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, offset]
links: [eng.kafka.partition, eng.kafka.group.pubsub-queue, eng.kafka.group.coordinator-assignment, eng.kafka.consumer.offset-commit, eng.kafka.rebalance]
---

## 两个维度

Consumer Group 是一组协同消费的实例。对同一个 Group：

- 一个 Partition 同一时刻只由一个 Consumer 负责；多个 Partition 可以被多个 Consumer 并行处理。
- Group 自己保存每个 Partition 的 committed offset；另一个 Group 有另一套独立位点。

```text
Topic orders
P0 ── Group inventory 的 C1
P0 ── Group recommendation 的 C3
```

这里没有给 Group 复制一份消息。两个 Group 都从同一个 Partition Log 拉取，只是书签不同。Group 内部是负载分摊，Group 之间是独立订阅。

## 并行度上限

一个 Group 的有效消费并行度最多受订阅 Partition 数限制。四个 Partition 配八个 Consumer 时，最多四个 Consumer 有分区，其余处于空闲状态。若只有少数 Partition 热，增加 Consumer 也不会自动解决热点。

## 位点不是业务成功证明

Consumer 的本地读取位置、已经提交到 Kafka 的位点、业务数据库写入成功，可能处在不同阶段。提交策略和业务副作用必须在 [[eng.kafka.consumer.offset-commit]] 与 [[eng.kafka.consumer.processing-idempotence]] 中一起设计。

成员变化会进入 [[eng.kafka.rebalance]]；Group 与 Pub/Sub 的关系见 [[eng.kafka.group.pubsub-queue]]。
