---
id: eng.kafka.group.coordinator-assignment
kind: engineering
title: Kafka Group Coordinator 与 Partition Assignment
summary: Group Coordinator 负责维护组状态，客户端参与加入、同步和分区分配，使每个 Partition 在组内有唯一 owner。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, coordinator, assignment]
links: [eng.kafka.consumer-group, eng.kafka.group.assignment-strategies, eng.kafka.rebalance, eng.kafka.group.rebalance-lifecycle, eng.kafka.group.heartbeat-poll]
---

## 协调过程

新 Consumer 使用同一个 `group.id` 加入后，组需要经历成员发现、订阅协商和 Assignment，把 Topic Partition 分发给成员。分配结果决定谁可以向哪个 Partition Fetch；它不是 Producer 写入路由。

```text
Group Coordinator
  → 收集成员
  → 生成 Assignment
  → 通知成员各自负责的 Partition
```

成员变化会生成新一轮分配。Consumer 只应处理当前拥有的分区，失去归属后继续提交旧分区位点可能失败，这也是交接逻辑必须被认真处理的原因。

## 和 Producer 路由的区别

| 问题 | Producer | Consumer Group |
| --- | --- | --- |
| 决定什么 | Record 写入哪个 Partition | 哪个 Consumer 读取哪个 Partition |
| 依据什么 | Partitioner、Key、显式 Partition | 成员、订阅和 Assignment 策略 |
| 变化后果 | Metadata 刷新和重试 | Rebalance、分区交接和位点恢复 |

Group Coordinator 负责组状态和协调，不是一个把所有消息转发给 Consumer 的中心 Router。

Assignment 策略的均衡与迁移取舍见 [[eng.kafka.group.assignment-strategies]]；加入、撤销、重新分配和恢复的时序见 [[eng.kafka.group.rebalance-lifecycle]]。
