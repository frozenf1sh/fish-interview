---
id: eng.kafka.rebalance
kind: engineering
title: Kafka Rebalance：重新交接 Partition 归属
summary: Consumer Group 成员、订阅或 Partition 发生变化时，协调器会重新计算分配；交接期间可能暂停消费并重放未提交处理。
parents: [eng.kafka.group]
tags: [kafka, rebalance, coordination]
links: [eng.kafka.consumer-group, eng.kafka.group.coordinator-assignment, eng.kafka.group.assignment-strategies, eng.kafka.group.rebalance-lifecycle, eng.kafka.group.heartbeat-poll, eng.kafka.consumer.offset-commit]
---

## 什么时候发生

Consumer 加入、退出、崩溃、订阅变化以及 Topic Partition 数变化，都可能触发 Rebalance。协调器需要确认新成员集合，再把 Partition 重新分配给组内成员。

```text
变更前：C1 → P0、P1；C2 → P2、P3
加入 C3：重新分配 P0、P1、P2、P3
```

不同协议的暂停范围不同，但应用都应假定 Partition 可能暂时不可消费，且交接前没有妥善提交的处理可能被重新执行。

不要把所有版本的 Rebalance 都理解成完全相同的全组停顿：经典协议与增量/合作式协议在交接范围和时序上存在差异。通用生命周期见 [[eng.kafka.group.rebalance-lifecycle]]，具体 Assignment 策略的取舍见 [[eng.kafka.group.assignment-strategies]]。

## 为什么会频繁发生

处理线程长时间阻塞、没有按预期继续 `poll`、网络抖动、滚动发布和实例反复重启，都会让组认为成员不再健康。`max.poll.interval`、心跳和会话超时应与实际处理模型匹配，而不是只靠调大超时掩盖阻塞。

## 排查顺序

先看成员变更和协调器事件，再对照处理耗时、poll 间隔和提交失败，最后检查部署与网络。若重复消费只发生在分区交接窗口，应联查 [[eng.kafka.consumer.offset-commit]]。
