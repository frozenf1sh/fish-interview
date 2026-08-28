---
id: eng.kafka.group.assignment-strategies
kind: engineering
title: Kafka Partition Assignment：组内如何分配分区
summary: Consumer Group 的 Assignment 决定成员各自读取哪些 Partition；策略在均衡性、迁移量和交接暂停范围之间做权衡，并且受客户端与协议版本影响。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, assignment, rebalance]
links: [eng.kafka.consumer-group, eng.kafka.group.coordinator-assignment, eng.kafka.rebalance, eng.kafka.group.rebalance-lifecycle, eng.kafka.performance.partition-parallelism]
---

## Assignment 解决什么问题

Producer 的 Partitioner 决定一条 Record 写入哪个 Partition；Consumer Group 的 Assignment 决定哪个成员读取哪个 Partition。这是两个方向相反、互不替代的决策：前者发生在写入路由，后者发生在组内消费协调。

```text
Topic 的 P0 P1 P2 P3
          │ 订阅 + 成员集合
          ▼
Group Assignment
          │
          ├── C1：P0、P2
          └── C2：P1、P3
```

经典 Consumer 协议中，协调器收集成员后，由组内的分配器计算结果，再通过同步阶段让成员应用各自的 Partition；较新的 Consumer Rebalance Protocol 可以把更多分配工作放到服务端，并支持更细粒度的增量交接。面试时应先说清协议边界，再谈具体策略名称。

## 常见策略的取舍

| 思路 | 优点 | 代价或风险 |
| --- | --- | --- |
| Range | 实现直观，按 Partition 范围分配 | 多 Topic 订阅时可能出现成员间偏斜 |
| Round-robin | 更倾向于均匀摊开分区 | 成员变化时可能搬动较多归属 |
| Sticky | 尽量保持旧归属，同时追求均衡 | 策略实现与版本行为需要核对 |
| Cooperative / 增量交接 | 一次只交接必要的分区，降低全组暂停 | 应用必须正确处理增量 revoke/assign 生命周期 |

不存在脱离场景的“最好策略”。Partition 数、Topic 订阅、成员数量、处理耗时以及是否允许短暂停顿，都会改变选择。还要注意：策略名称、默认值和服务端协议会随 Kafka 客户端/集群版本变化，不要把某个版本的默认行为当成 Kafka 永久规则。

## 选择策略前先问三个问题

1. 一个 Group 订阅几个 Topic，Partition 数是否能均匀覆盖成员？
2. Rebalance 时能否接受所有成员先释放再重新获取，还是需要增量交接？
3. 处理是否有本地缓存、连接或状态，分区迁移时怎样安全清理和恢复？

Assignment 只负责归属，不负责业务处理成功，也不负责保存业务进度；位点与交接顺序分别见 [[eng.kafka.consumer.commit-modes]] 和 [[eng.kafka.group.rebalance-lifecycle]]。
