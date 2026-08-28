---
id: eng.kafka.group.assignment-strategies
kind: engineering
title: Kafka Partition Assignment：组内怎样分配分区
summary: Assignment 是 Consumer Group 的分工表；不同分配策略在均衡、搬迁数量和暂停范围之间取舍，并且具体行为受客户端与协议版本影响。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, assignment, rebalance]
links: [eng.kafka.consumer-group, eng.kafka.group.coordinator-assignment, eng.kafka.rebalance, eng.kafka.group.rebalance-lifecycle, eng.kafka.performance.partition-parallelism]
---

## 先区分两个“分配器”

Producer 的 Partitioner 决定“这条 Record 写进哪个 Partition”；Consumer Group 的 Assignment 决定“哪个 Consumer 读取哪个 Partition”。前者给数据找目的地，后者给成员分工。

```text
P0、P1、P2、P3 + C1、C2
              ↓
Assignment：C1 → P0、P2；C2 → P1、P3
```

**Assignment** 就是分配结果；**Assignor** 是计算这个结果的策略或组件。经典 Consumer 协议中，协调流程会收集成员、计算并同步分配；较新的 Consumer Rebalance Protocol 可以由服务端承担更多分配工作。具体 RPC 不同，但“成员获得 Partition owner”这条主线不变。

## 常见思路用人话解释

| 思路 | 怎么分 | 主要取舍 |
| --- | --- | --- |
| Range | 像切蛋糕一样按连续范围分 | 直观，但多 Topic 时可能偏斜 |
| Round-robin | 按顺序轮流发给成员 | 通常更平均，但成员变化可能搬动较多 |
| Sticky | 尽量保留旧分工，再做均衡 | 减少搬迁，但策略行为要看客户端 |
| Cooperative / 增量 | 只交接受影响的 Partition | 减少全组暂停，但生命周期处理更复杂 |

不要先背某个默认值。要先看 Partition 数、Topic 订阅、成员数量、状态缓存以及业务能否接受短暂停顿。策略名称和默认行为可能随 Kafka 客户端/集群版本变化，不能脱离版本作永久结论。

Assignment 只决定归属，不保存业务进度，也不证明处理成功；位点和交接分别见 [[eng.kafka.consumer.commit-modes]] 与 [[eng.kafka.group.rebalance-lifecycle]]。
