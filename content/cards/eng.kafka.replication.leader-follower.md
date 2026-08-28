---
id: eng.kafka.replication.leader-follower
kind: engineering
title: Kafka Replica、Leader 与 Follower：给一册日志留副本
summary: Replica 是同一个 Partition 的副本；Leader 负责主要读写，Follower 通过拉取日志追赶 Leader，副本让 Broker 故障时有机会恢复服务。
parents: [eng.kafka.replication]
tags: [kafka, replica, leader, follower]
links: [eng.kafka.partition, eng.kafka.broker, eng.kafka.replication.isr-election, eng.kafka.replication.write-ack]
---

## 先看为什么要复制

如果 P0 只有 Broker 1 一份，Broker 1 挂掉，P0 的数据和服务入口都可能一起消失。Kafka 可以把 P0 放在多个 Broker：

```text
P0 的三份 Replica
├── Broker 1：Leader，主要接收读写
├── Broker 2：Follower，保存复制日志
└── Broker 3：Follower，保存复制日志
```

**Replica** 就是副本；**Leader** 是当前承担主要请求的副本；**Follower** 是跟随 Leader 保存日志的副本。Broker 是服务器，Replica 是服务器上存放的一份 Partition 数据，不能混为一谈。

## Follower 如何追上

Follower 通常主动向 Leader 发 **Fetch** 请求，意思是“把我还没有的日志给我”：

```text
Follower：我目前到 Offset 100
        → Fetch：请给我 101 之后的记录
Leader：返回新增记录
```

所以不要想成 Leader 为每条消息单独推送给每个 Follower。Consumer 读取也采用类似的主动拉取思路。

## 复制解决什么、不解决什么

副本提高的是故障容忍度：Leader 故障后，符合条件的 Follower 可能接任。它不会把一个 Partition 变成多条可并行消费的队列，也不会改变同一 Partition 内的顺序。哪些副本可以接任，见 [[eng.kafka.replication.isr-election]]。
