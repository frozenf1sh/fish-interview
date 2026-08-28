---
id: eng.kafka.replication.leader-follower
kind: engineering
title: Kafka Replica、Leader 与 Follower：一个分区的多份副本
summary: Partition 的 Replica 分布在多个 Broker 上，Leader 承担主要读写，Follower 通过复制追赶日志并在故障时参与接任。
parents: [eng.kafka.replication]
tags: [kafka, replica, leader, follower]
links: [eng.kafka.partition, eng.kafka.broker, eng.kafka.replication.isr-election, eng.kafka.replication.write-ack]
---

## 物理结构

```text
P0 Replica
├── Broker 1：Leader
├── Broker 2：Follower
└── Broker 3：Follower
```

Producer 的写入和 Consumer 的读取通常围绕 Partition Leader 进行；Follower 保存自己的日志副本并主动从 Leader Fetch 数据。Replica Factor 描述副本数量，Broker 是副本实际所在的节点。

## 这解决什么问题

单个 Broker 故障时，只要有满足条件的副本，Partition 仍可能选出新的 Leader。复制提高故障容忍度，但不会改变 Partition 内顺序，也不会让一个 Partition 产生多路消费并行度。
