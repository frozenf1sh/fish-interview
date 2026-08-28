---
id: eng.kafka.broker
kind: engineering
title: Kafka Broker：真正运行 Kafka 的服务器
summary: Broker 是一个 Kafka Server 实例，负责接收请求、保存本机的 Partition 副本、参与复制并响应客户端；它不是 Partition。
parents: [eng.kafka.model]
tags: [kafka, broker, cluster]
links: [eng.kafka.partition, eng.kafka.replication.leader-follower, eng.kafka.bootstrap-metadata]
---

## 先区分“机器”和“数据册”

**Broker** 是运行 Kafka 进程的服务器节点；**Partition** 是 Topic 的一册数据日志。一个 Broker 可以同时保存很多 Partition，也可以在同一个 Partition 上保存 Leader 或 Follower 副本。

```text
Kafka Cluster（集群）
├── Broker 1：P0 Leader、P1 Follower
├── Broker 2：P1 Leader、P0 Follower
└── Broker 3：P0 Follower、P1 Follower
```

**Cluster** 就是多个 Broker 共同组成的 Kafka 集群。**Leader** 是某个 Partition 当前负责主要读写的副本，**Follower** 是跟着 Leader 复制日志的副本；这两个词只描述副本角色，不是两种服务器。

## Broker 收到什么请求

- Producer 找到目标 Partition 的 Leader，发送写入请求。
- Leader 把 Record 追加到日志，并把结果复制给 Follower。
- Consumer 通常向 Partition Leader 拉取数据。
- 集群在 Leader 故障时重新选择副本接任，并告诉客户端新的位置。

因此，Producer 不是把消息随便丢给一台 Broker 再等它转发。Producer 会先取得集群信息，找到目标 Partition 的 Leader，见 [[eng.kafka.bootstrap-metadata]]。

## 记住两种数量

Broker 数量更多，通常意味着可承载的节点和故障余量更多；Partition 数量更多，意味着更细的分片和并行边界。两者不是同一个旋钮，复制关系见 [[eng.kafka.replication.leader-follower]]。
