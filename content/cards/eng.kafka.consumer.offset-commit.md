---
id: eng.kafka.consumer.offset-commit
kind: engineering
title: Kafka Offset Commit：提交下一次从哪里读
summary: Consumer Group 提交的是某个 Partition 的恢复位置，通常是下一条要读取的 Offset；它不是数据库写成功的自动确认。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, offset, commit]
links: [eng.kafka.offset-retention, eng.kafka.consumer.commit-modes, eng.kafka.consumer.processing-idempotence, eng.kafka.consumer-group, eng.kafka.rebalance]
---

## 先用书签理解

Partition 像一本不断追加的书，Consumer Group 像一个读者：

```text
读取：    Offset 10 的 Record
处理成功：业务副作用完成
提交：    Offset 11
含义：    下次从 11 继续，而不是重复从 10 开始
```

**Offset** 是 Partition 内的记录位置；**Commit** 是把 Group 的恢复位置保存下来。提交 `11` 通常表示 `10` 已经处理完，下一次读取从 `11` 开始。

![Kafka Consumer 处理、提交与崩溃重放](/static/kafka-consumer-failure.drawio.svg)

## 两个位置不要混

- **Position**：当前 Consumer 客户端准备继续读取的位置，可能已经随着 `poll()` 前进。
- **Committed Offset**：已经保存到 Kafka Group 位点存储的位置，Consumer 重启或 Partition 交接时会用它恢复。

Position 领先 Committed Offset 很正常：客户端可以先把数据拉到本地，处理完成后再提交。

## 提交时机决定故障结果

```text
先处理 → 再提交：崩溃时可能重复，但不容易跳过尚未处理的记录
先提交 → 再处理：处理前崩溃可能直接跳过记录
```

Kafka 只知道位点有没有提交，不知道数据库、RPC 或文件操作是否成功。自动、同步和异步提交的差异见 [[eng.kafka.consumer.commit-modes]]；业务副作用的去重见 [[eng.kafka.consumer.processing-idempotence]]。
