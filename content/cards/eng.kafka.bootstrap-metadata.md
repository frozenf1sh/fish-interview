---
id: eng.kafka.bootstrap-metadata
kind: engineering
title: Kafka Bootstrap Server 与 Metadata：Producer 的集群路由表
summary: bootstrap.servers 只是客户端进入集群并获取 Metadata 的入口，Producer 随后依据路由表直接访问目标 Partition Leader。
parents: [eng.kafka.routing]
tags: [kafka, bootstrap, metadata, producer]
links: [eng.kafka.broker, eng.kafka.leader-routing, eng.kafka.producer.send]
---

## Bootstrap 不是固定网关

Producer 启动时并不知道整个集群的 Broker、Topic Partition 或 Leader。它先尝试连接 `bootstrap.servers` 中的一个或多个地址，获取类似这样的 Metadata：

```text
orders / P0 → Broker 1
orders / P1 → Broker 3
orders / P2 → Broker 2
```

配置多个入口的目的主要是提高初始发现的成功率，不是把 P0 永久发给第一台、P1 永久发给第二台。拿到 Metadata 后，Producer 会根据目标 Partition 找 Leader。

## Metadata 会变化

Leader 切换、Topic 扩分区或 Broker 拓扑变化都会使本地 Metadata 过期。客户端会根据错误、定期刷新或下一次访问触发 Metadata 更新，再重新选择目标 Broker。Bootstrap 地址因此是引导信息，不是整个发送路径中的中心 Router。
