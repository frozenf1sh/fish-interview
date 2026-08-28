---
id: eng.kafka.consumer.offset-commit
kind: engineering
title: Kafka Offset Commit：提交的是下一次读取位置
summary: Consumer Commit 记录某个 Group 在某个 Partition 的恢复位置；它不是对业务数据库写入的自动确认。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, offset, commit]
links: [eng.kafka.offset-retention, eng.kafka.consumer.processing-idempotence, eng.kafka.consumer-group, eng.kafka.rebalance]
---

## 位置与状态

假设 Consumer 拉取并处理了 Offset `10`，通常提交的是下一次从 `11` 开始读取的位置。进程重启或分区重新分配时，客户端根据 Group 的 committed offset 恢复。

![Kafka Consumer 处理、提交与崩溃重放](/static/kafka-consumer-failure.drawio.svg)

```text
fetch 10 → process 10 → commit 11
```

## 提交时机决定故障语义

- 先处理、后提交：崩溃窗口会重复处理，但不容易跳过未处理消息。
- 先提交、后处理：处理前崩溃可能造成消息被跳过。
- 自动提交：减少样板代码，但提交时机可能与业务副作用脱节。

Kafka 只知道位点提交到哪里，不知道数据库写入、RPC 或文件操作是否成功；这些边界见 [[eng.kafka.consumer.processing-idempotence]]。
