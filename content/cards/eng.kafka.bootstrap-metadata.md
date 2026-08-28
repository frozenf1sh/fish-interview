---
id: eng.kafka.bootstrap-metadata
kind: engineering
title: Kafka Bootstrap Server 与 Metadata：Producer 先找到集群地图
summary: bootstrap.servers 是 Producer 进入 Kafka 集群的起点；Metadata 是客户端维护的路由信息，告诉它每个 Partition 的 Leader 在哪。
parents: [eng.kafka.routing]
tags: [kafka, bootstrap, metadata, producer]
links: [eng.kafka.broker, eng.kafka.leader-routing, eng.kafka.producer.send]
---

## 先看 Producer 刚启动时什么都不知道

假设配置里写着：

```text
bootstrap.servers = broker-1:9092, broker-2:9092
```

**Bootstrap Server** 就是“引导地址”：Producer 先尝试连接其中一个可用 Broker，请它返回集群信息。它不是一个永远站在中间转发所有消息的网关。

## Metadata 是什么

**Metadata** 可以理解成 Producer 手里的集群地图，里面至少需要知道：有哪些 Broker、某个 Topic 有哪些 Partition、每个 Partition 当前由哪个 **Leader** 副本负责写入。

```text
orders / P0 → Broker 2（Leader）
orders / P1 → Broker 3（Leader）
orders / P2 → Broker 1（Leader）
```

Producer 先由 [[eng.kafka.key-partitioner]] 选 P2，再查这张地图，直接把请求发给 Broker 1。

## 为什么通常写多个引导地址

如果只写 `broker-1`，恰好它重启，Producer 可能连不上集群、拿不到第一张地图；写多个地址是为了提高初始发现成功率，不是规定 P0 永远去第一台、P1 永远去第二台。

## 地图会过期

Broker 故障、Partition Leader 切换、Topic 增加 Partition，都可能让本地 Metadata 变旧。Producer 收到“你找错 Leader”一类响应，或发现连接失败后，会刷新地图，再把请求发往新的 Leader。完整路径见 [[eng.kafka.leader-routing]]。
