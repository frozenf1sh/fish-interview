---
id: eng.kafka.broker
kind: engineering
title: Kafka Broker：承载 Partition 的服务器节点
summary: Broker 是 Kafka Server 实例，负责接收请求、存储本机副本、参与复制并响应客户端；它不是 Partition 本身。
parents: [eng.kafka.model]
tags: [kafka, broker, cluster]
links: [eng.kafka.partition, eng.kafka.replication.leader-follower, eng.kafka.bootstrap-metadata]
---

## 先分清两种对象

Broker 是物理或进程层面的 Kafka Server；Partition 是 Topic 的逻辑数据分片。一个 Broker 可以同时承载多个 Partition 的 Leader 和 Follower，一个 Partition 的不同 Replica 也可以分布在多个 Broker 上。

```text
Cluster
├── Broker 1：P0 Leader、P1 Follower
├── Broker 2：P1 Leader、P0 Follower
└── Broker 3：P0 Follower、P1 Follower
```

Producer 写入时最终要找到目标 Partition 的 Leader Broker；Consumer 拉取时也通常向该 Partition 的 Leader 请求。Broker 的数量影响承载和容错，Partition 的数量影响分片和并行度。

## 面试边界

不要回答“消息发给某个 Broker 后 Broker 再帮我路由到 Partition”作为默认模型。Producer 通过 Metadata 感知 Partition Leader，并直接向目标 Broker 发请求；Leader 负责追加和副本协作。
