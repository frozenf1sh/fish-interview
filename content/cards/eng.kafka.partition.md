---
id: eng.kafka.partition
kind: engineering
title: Kafka Partition：顺序、Offset 与复制的最小边界
summary: Partition 是追加日志、局部顺序、消费者归属和副本选主的共同边界。
parents: [eng.kafka]
tags: [kafka, partition, offset, replication]
links: [eng.kafka, eng.kafka.consumer-group]
---

## 核心机制

一个 Partition 是仅追加的有序日志。Broker 为其维护 Leader 和副本；客户端通常与 Leader 交互。offset 按 Partition 递增，因此 `offset=42` 必须连同 Topic 和 Partition 一起解释。

Producer 的确认级别、副本同步状态和最小同步副本数共同决定故障时可接受的持久性边界。Consumer 的并行度同样在这里受限：同组一个 Partition 一次只能交给一个消费者。

## 排查线索

出现热点时，先比较各 Partition 的写入速率、消息大小、key 分布和 consumer lag。若只有少数分区积压，盲目增加消费者通常无效；需要定位键倾斜或改变分区策略。

## 常见误区

副本数提高的是可用性与容错，不会让单 Partition 的读取顺序变成可并行处理。要并行处理同一键，必须先定义业务是否允许重排。

## 变体与关联

[[eng.kafka.consumer-group]] 直接以 Partition 为分配单位；Topic 只是这些分区的逻辑集合。

