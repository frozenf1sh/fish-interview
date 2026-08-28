---
id: eng.kafka.replication.isr-election
kind: engineering
title: Kafka ISR 与 Leader Election：只让跟得上的副本接任
summary: ISR 是当前跟得上 Leader 的副本集合；Leader 故障时优先从这个集合中选新 Leader，以降低已经确认数据丢失的风险。
parents: [eng.kafka.replication]
tags: [kafka, isr, election, availability]
links: [eng.kafka.replication.leader-follower, eng.kafka.replication.write-ack, eng.kafka.leader-routing]
---

## 先看三个副本的进度

假设 P0 的 Leader 已经写到 Offset 100：

```text
Broker 1：Leader，Offset 100
Broker 2：Follower，Offset 100  → 跟得上
Broker 3：Follower，Offset 20   → 落后很多
```

**ISR（In-Sync Replicas）** 的中文常译是“同步副本集合”。它表示当前满足同步条件、可以被视为跟得上的 Replica；Leader 自己也属于 ISR。落后太久的 Follower 可能暂时被移出，追上后再回来。

## 为什么不能随便选副本

如果 Leader 已经有 `A B C D E`，某个 Follower 只有 `A B C`，却让它直接当新 Leader，那么 `D E` 可能暂时不可见甚至丢失。ISR 的意义就是把“可以承担接任风险的副本”与“只是曾经配置过的副本”区分开。

## 故障链怎么串

```text
Leader 故障
  → 集群从可接受副本中选新 Leader
  → Producer 的旧 Metadata 失效
  → Producer 刷新路由并处理重试
```

ISR 不是“数据绝对不会丢”的单独保证；它要和 Producer 的 ACK、最小同步副本数量一起看，见 [[eng.kafka.replication.write-ack]]。
