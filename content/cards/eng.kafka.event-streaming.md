---
id: eng.kafka.event-streaming
kind: engineering
title: Kafka Event Streaming：保存发生过什么
summary: Kafka 把事件按时间追加到可读取的日志中，使多个系统能独立消费、重放并从事件重新计算状态。
parents: [eng.kafka.semantics]
tags: [kafka, event-streaming, replay, event-sourcing]
links: [eng.kafka.offset-retention, eng.kafka.group.pubsub-queue, eng.kafka.delivery-semantics]
---

## 事件和状态

账户状态可以由 `AccountCreated`、`Deposit`、`Withdraw` 等事件累积得到。Kafka 更擅长保存“发生过什么”，Consumer 再按自己的 Offset 把事件折叠成当前状态。

这使 Kafka 适合事件流、CDC、数据管道和流处理：新系统可以从保留范围内的历史 Offset 开始建立自己的视图。它不是永久档案库，Retention 和 Compaction 仍然决定哪些历史可被读取。
