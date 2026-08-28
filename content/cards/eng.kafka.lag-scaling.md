---
id: eng.kafka.lag-scaling
kind: engineering
title: Kafka Consumer Lag：积压指标不是根因
summary: Lag 表示日志位置与 Consumer Group 已提交位置之间的差距，但不能单独说明是生产过快、未分配、处理慢还是提交失败。
parents: [eng.kafka.operations]
tags: [kafka, lag, scaling, troubleshooting]
links: [eng.kafka.consumer-group, eng.kafka.consumer.offset-commit, eng.kafka.performance.partition-parallelism, eng.kafka.failure-diagnosis]
---

## 先拆 Lag 的位置

排查时至少区分：

1. Producer 是否持续写入，写入是否集中到少数 Partition。
2. Consumer Group 是否拿到预期 Partition。
3. Consumer 是否成功 Poll，以及处理是否耗时过长。
4. 业务处理是否完成，Offset 是否按预期提交。

只增加 Consumer 只有在存在空闲 Partition 且处理能力可扩展时才有效。若热点来自 Key 倾斜，增加组成员不会把一个 Partition 拆开。
