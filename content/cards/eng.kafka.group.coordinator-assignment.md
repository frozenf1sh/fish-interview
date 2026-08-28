---
id: eng.kafka.group.coordinator-assignment
kind: engineering
title: Kafka Group Coordinator 与 Assignment：谁负责协调分工
summary: Group Coordinator 是负责维护某个 Consumer Group 状态的 Kafka 端协调者；Assignment 是把 Partition 分给组内成员的结果，不是 Producer 的写入路由。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, coordinator, assignment]
links: [eng.kafka.consumer-group, eng.kafka.group.assignment-strategies, eng.kafka.rebalance, eng.kafka.group.rebalance-lifecycle, eng.kafka.group.heartbeat-poll]
---

## 先看一个新成员加入

C1、C2 已经在消费，C3 刚启动并使用同一个 `group.id`。Kafka 需要知道当前有哪些成员、它们订阅什么，然后决定 P0、P1、P2 分给谁。

**Group Coordinator** 可以理解为这个 Group 的“登记和协调窗口”：它维护成员状态、接收加入和心跳、推动新一轮分配。**Assignment** 是最终的分工表，例如 `C1 → P0、P2`、`C2 → P1`。

```text
Consumer 成员
  → 向 Coordinator 报到
  → 确认订阅和成员集合
  → 计算 Assignment
  → 每个成员拿到自己的 Partition
```

## 它不是中心消息 Router

Coordinator 不负责把每条消息逐条转发给 Consumer。Consumer 取得自己的分区归属后，仍然直接向对应 Partition 的 Leader Fetch；Producer 的写入路由则由 Partitioner 和 Metadata 决定，二者不要混淆。

成员加入、离开、失联或订阅变化，会触发 [[eng.kafka.rebalance]]。具体怎么做到尽量均衡、少搬迁，见 [[eng.kafka.group.assignment-strategies]]；完整交接时序见 [[eng.kafka.group.rebalance-lifecycle]]。
