---
id: eng.kafka.leader-routing
kind: engineering
title: Kafka Partition Leader 路由：地图过期后重新找路
summary: Producer 先选 Partition，再根据 Metadata 找到该 Partition 的 Leader；Leader 变化后，客户端刷新路由并按错误类型决定是否重试。
parents: [eng.kafka.routing]
tags: [kafka, leader, metadata, retry]
links: [eng.kafka.bootstrap-metadata, eng.kafka.replication.leader-follower, eng.kafka.producer.retry-idempotence]
---

## 先看一次写入怎么找到服务器

先假设 Producer 选中了 P2：

```text
Record
  → Partitioner 选择 P2
  → Metadata 说 P2 Leader 在 Broker 3
  → Producer 把 ProduceRequest 发给 Broker 3
```

**Leader** 是某个 Partition 当前负责主要读写的副本；**Metadata** 是 Producer 保存的路由表。Producer 不需要先把消息交给一台固定 Router，再由 Router 转发。

## Leader 挂掉时发生什么

```text
旧地图：P2 Leader = Broker 3
Broker 3 故障
集群选出新 Leader = Broker 1
Producer 刷新 Metadata
Producer → Broker 1 重试
```

这里的“选出新 Leader”由 Kafka 集群完成；Producer 的任务是发现旧地图失效并重新定位。刷新可能由错误响应、连接异常或定期更新触发。

## 不是所有错误都能重试

客户端需要区分：

- **可刷新路由**：说明当前 Leader 信息过期，应更新 Metadata。
- **暂时性故障**：例如短暂网络问题，可能适合重试。
- **不可重试失败**：例如记录过大、权限不足或配置错误，继续重试没有意义。

重试虽然能提高暂时故障下的成功率，也可能带来重复写入和顺序风险。要把这条路由故障链和 [[eng.kafka.producer.retry-idempotence]] 放在一起看。
