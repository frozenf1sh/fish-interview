---
id: eng.kafka.consumer.offset-commit
kind: engineering
title: Kafka Offset Commit：提交的是下一次读取位置
summary: Consumer Commit 记录某个 Group 在某个 Partition 的恢复位置；它不是对业务数据库写入的自动确认。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, offset, commit]
links: [eng.kafka.offset-retention, eng.kafka.consumer.commit-modes, eng.kafka.consumer.processing-idempotence, eng.kafka.consumer-group, eng.kafka.rebalance]
---

## 位置与状态

假设 Consumer 拉取并处理了 Offset `10`，通常提交的是下一次从 `11` 开始读取的位置。进程重启或分区重新分配时，客户端根据 Group 的 committed offset 恢复。

![Kafka Consumer 处理、提交与崩溃重放](/static/kafka-consumer-failure.drawio.svg)

```text
fetch 10 → process 10 → commit 11
```

要区分两个位置：Consumer 当前已经拉到的 position 可以领先于已持久化的 committed position。前者表示客户端准备继续读哪里，后者才是故障恢复时 Kafka 能提供的起点。

## 提交时机决定故障语义

- 先处理、后提交：崩溃窗口会重复处理，但不容易跳过未处理消息。
- 先提交、后处理：处理前崩溃可能造成消息被跳过。
- 自动提交：减少样板代码，但提交时机可能与业务副作用脱节。

Kafka 只知道位点提交到哪里，不知道数据库写入、RPC 或文件操作是否成功；这些边界见 [[eng.kafka.consumer.processing-idempotence]]。

自动、同步与异步提交的等待和失败语义，见 [[eng.kafka.consumer.commit-modes]]；分区交接时提交旧归属，还要结合 [[eng.kafka.group.rebalance-lifecycle]]。

面试追问“提交 Offset 是提交当前消息还是下一条消息”时，应回答：通常提交下一次要读取的位置；具体 API 以客户端的 offset 约定为准。
