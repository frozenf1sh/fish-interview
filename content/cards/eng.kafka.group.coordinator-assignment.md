---
id: eng.kafka.group.coordinator-assignment
kind: engineering
title: Kafka Group Coordinator 与 Partition Assignment
summary: Group Coordinator 负责维护组状态，客户端参与加入、同步和分区分配，使每个 Partition 在组内有唯一 owner。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, coordinator, assignment]
links: [eng.kafka.consumer-group, eng.kafka.rebalance, eng.kafka.group.heartbeat-poll]
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
