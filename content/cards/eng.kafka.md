---
id: eng.kafka
kind: engineering
title: Kafka：从追加日志到消息生命周期
summary: Kafka 先用 Partition 保存可重放的追加日志，再用 Producer、Replica、Consumer Group 和 Offset 把写入、复制与消费连接起来。
parents: [engineering]
tags: [kafka, messaging, distributed-systems]
links: [eng.kafka.topic, eng.kafka.partition, eng.kafka.broker, eng.kafka.key-partitioner, eng.kafka.bootstrap-metadata, eng.kafka.producer.send, eng.kafka.replication.write-ack, eng.kafka.consumer.pull-poll, eng.kafka.consumer-group, eng.kafka.delivery-semantics]
---

## 先记住一条主线

Kafka 可以先理解成一个**分布式、可持久化、可复制、按 Offset 读取的追加日志系统**。消息写入某个 Topic 的某个 Partition 后，不会因为一个 Consumer 读过就被拿走；Kafka 保存数据，Consumer Group 保存自己的读取位置。

```text
Producer
  → Partitioner 选择 Partition
  → Leader Broker 追加 Log 并分配 Offset
  → Follower Replica 复制
  → Consumer Group 按自己的 Offset 拉取
  → 处理后提交位点
```

这条链上的边界必须分开：Topic 是逻辑命名空间，Partition 是顺序和并行的基本单位，Broker 是服务器节点，Replica 是副本，Consumer Group 是消费进度的隔离边界。

![Kafka 数据模型：Cluster、Broker、Topic、Partition 与 Replica](/static/kafka-model.drawio.svg)

## 为什么不是“消息被取走”

同一个 Partition 可以被多个 Consumer Group 分别读取。库存组、推荐组和风控组各自维护 `P0` 的 committed offset；它们竞争的是组内的 Partition，不是 Topic 中的物理消息。

因此 Kafka 同时具备两种使用语义：不同 Group 之间像发布订阅，同一 Group 内像竞争消费者。Replay 也不是把消息恢复出来，而是把某个 Group 的读取位置移回仍在 Retention 范围内的 Offset。

## 面试时的边界回答

- Kafka 只保证 Partition 内有序，不保证 Topic 全局有序。
- 同一个稳定 Key 通常会被映射到同一个 Partition，但扩容可能改变映射，不能把默认分区逻辑称为经典一致性哈希。
- Producer 的幂等性主要处理发送重试造成的重复写入，不等于跨 Kafka、数据库和外部副作用的 Exactly-Once。
- `send()` 返回不等于 Broker 已确认；Producer 可能还在序列化、缓冲或等待 Sender Thread 发送。
- `acks=all`、ISR 和 `min.insync.replicas`共同定义写入确认边界；“写成功”仍需结合 Consumer 提交和业务幂等来讨论。

## 按什么顺序继续

先看 [[eng.kafka.topic]]、[[eng.kafka.partition]] 和 [[eng.kafka.broker]]，再沿 [[eng.kafka.key-partitioner]]、[[eng.kafka.bootstrap-metadata]] 进入 Producer 链路；之后连接 [[eng.kafka.replication.write-ack]]、[[eng.kafka.consumer.pull-poll]] 与 [[eng.kafka.consumer-group]]。
