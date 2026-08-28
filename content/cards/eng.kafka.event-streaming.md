---
id: eng.kafka.event-streaming
kind: engineering
title: Kafka Event Streaming：保存发生过什么
summary: Kafka 按顺序保存不断到来的事件，多个系统可以从各自位置读取并重放，再把事件累积成自己的状态视图。
parents: [eng.kafka.semantics]
tags: [kafka, event-streaming, replay, event-sourcing]
links: [eng.kafka.offset-retention, eng.kafka.group.pubsub-queue, eng.kafka.delivery-semantics]
---

## 先区分事件和当前状态

账户当前余额是 `120`，有两种保存方式：

```text
只保存状态：balance = 120
保存事件：AccountCreated → Deposit(+100) → Withdraw(-30) → Deposit(+50)
```

**事件** 描述“发生过什么”；**状态** 是把一系列事件按顺序计算后的结果。Kafka 更适合保存前者，Consumer 可以从某个 Offset 重新读取事件，构建报表、缓存或其他服务自己的状态视图。

## 为什么多个系统都能建立自己的视图

库存、推荐、数仓使用不同 Consumer Group，各自维护书签，所以可以从同一份事件日志独立读取。新系统也可以从保留范围内的较早 Offset 开始 Replay（重放）。

这使 Kafka 适合事件流、CDC（数据库变更捕获）、数据管道和流处理；但它不是永久档案库，历史能保留多久仍由 [[eng.kafka.offset-retention]] 和 [[eng.kafka.log-compaction]] 决定。
