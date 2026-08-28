---
id: eng.kafka.key-partitioner
kind: engineering
title: Kafka Key 与 Partitioner：把相关事件放到同一分区
summary: Producer Client 的 Partitioner 根据显式 Partition、Key 或无 Key 策略决定 Record 的目标 Partition。
parents: [eng.kafka.routing]
tags: [kafka, key, partitioner, ordering]
links: [eng.kafka.partition, eng.kafka.record, eng.kafka.ordering, eng.kafka.partition-expansion]
---

## 一条订单事件如何路由

```text
orderId = 123
    → Partitioner 的映射策略
    → Partition 2
```

`OrderCreated`、`OrderPaid`、`OrderShipped` 使用同一个稳定 Key 时，通常会进入同一 Partition，于是可以利用 Partition 内顺序表达同一订单的事件顺序。Key 的选择应对应业务真正需要保持顺序的实体，例如 `orderId`、`userId` 或 `deviceId`。

## 两个边界

第一，`hash(key) % partitionCount` 只是概念模型，实际映射由客户端版本和 Partitioner 实现决定。第二，Kafka 默认 Key 映射不等于经典一致性哈希；增加 Partition 可能改变映射，跨扩容时不能自动承诺全局顺序。

没有 Key 时，Producer 仍会按客户端策略选择 Partition，但业务不能因此推导出跨 Partition 顺序。
