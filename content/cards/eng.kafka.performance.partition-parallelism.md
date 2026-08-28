---
id: eng.kafka.performance.partition-parallelism
kind: engineering
title: Kafka Partition 并行：吞吐扩展的基本单位
summary: Kafka 通过多个 Partition 把写入、读取和副本工作分散到不同 Broker 与 Consumer，但并行度受分区数量和数据分布限制。
parents: [eng.kafka.performance]
tags: [kafka, performance, partition, parallelism]
links: [eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.lag-scaling, eng.kafka.key-partitioner]
---

## 扩展的两面

更多 Partition 可以让 Producer、Broker 和 Consumer 并行处理更多数据，但也会带来更多 Replica、连接、元数据、分配和 Rebalance 成本。Consumer Group 的有效并行度不会超过可消费的 Partition 数。

如果 Key 分布不均，某个 Partition 可能成为热点，整体吞吐会被最慢分区限制。此时先看消息大小、Key 分布、写入速率和 Lag，再决定是否改变分区策略。
