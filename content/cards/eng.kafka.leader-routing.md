---
id: eng.kafka.leader-routing
kind: engineering
title: Kafka Partition Leader 路由：过期拓扑下的重试
summary: Producer 先选 Partition，再从 Metadata 找 Leader；遇到 NotLeader 或连接失败时刷新拓扑并按错误语义重试。
parents: [eng.kafka.routing]
tags: [kafka, leader, metadata, retry]
links: [eng.kafka.bootstrap-metadata, eng.kafka.replication.leader-follower, eng.kafka.producer.retry-idempotence]
---

## 请求路径

```text
Record
  → Partitioner 选 P2
  → Metadata：P2 Leader = Broker 3
  → ProduceRequest 发给 Broker 3
```

如果 Broker 3 下线，集群完成选举后 P2 可能改由 Broker 1 负责。Producer 原来的路由表仍可能指向 Broker 3，直到收到可重试错误或发现连接异常，再刷新 Metadata 并定位新 Leader。

## 这不是无条件重发

客户端要区分可刷新 Metadata 的错误、可重试的暂时故障和不可重试的业务或配置错误。重试还会带来超时、重复写入和顺序风险，必须和 [[eng.kafka.producer.retry-idempotence]] 一起理解。
