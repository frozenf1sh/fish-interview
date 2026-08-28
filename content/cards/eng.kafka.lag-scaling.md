---
id: eng.kafka.lag-scaling
kind: engineering
title: Kafka Consumer Lag：书签落后了，但原因不止一个
summary: Consumer Lag 表示日志末尾与 Group 已提交位置之间的差距；它是积压信号，不直接说明生产过快、处理过慢、未分配还是提交失败。
parents: [eng.kafka.operations]
tags: [kafka, lag, scaling, troubleshooting]
links: [eng.kafka.consumer-group, eng.kafka.consumer.offset-commit, eng.kafka.performance.partition-parallelism, eng.kafka.failure-diagnosis]
---

## 先用书签理解 Lag

Partition 当前写到 Offset `100`，Consumer Group 已提交到 `80`，可以粗略说它还有约 20 个位置没追上。**Lag** 就是“日志最新位置”和“Group 书签”之间的落后量。

```text
日志：    80 81 82 … 100
书签：    ↑ 80
Lag：     还没追上的距离
```

## 看到 Lag 先问四件事

1. Producer 是否突然写快，或者所有数据都集中到一个 Partition？
2. Consumer Group 是否真的拿到了预期的 Partition？
3. Consumer 是否正常 `poll`，单批处理是否太慢或外部依赖阻塞？
4. 业务是否已经完成，Offset 是否因为提交失败没有推进？

只有存在空闲 Partition、且处理能水平扩展时，增加 Consumer 才有用。热点 Key、Rebalance 或业务数据库变慢时，盲目加实例可能没有效果，甚至增加协调成本。
