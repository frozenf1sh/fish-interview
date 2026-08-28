---
id: eng.kafka.consumer.processing-idempotence
kind: engineering
title: Kafka 消费处理与业务幂等：重复是正常故障路径
summary: Consumer 处理成功但 Offset 尚未提交时崩溃，重启会再次读取同一 Record；业务副作用必须具备幂等或去重能力。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, idempotence, at-least-once]
links: [eng.kafka.consumer.offset-commit, eng.kafka.delivery-semantics, eng.kafka.producer.retry-idempotence, eng.kafka.transactions-eos]
---

## 重复处理路径

```text
读取 A → 数据库写入成功 → 进程崩溃 → Offset 未提交
重启 → 再次读取 A
```

常见防线包括事件 ID 或业务 ID 唯一键、去重表、幂等 Upsert 和状态机。选择哪一种取决于副作用能否天然重复、是否需要保序以及数据库事务边界。

## 不要把“只提交一次”当成万能答案

提交位点和外部数据库写入若不在同一个原子边界，仍可能出现重复或跳过。Kafka Transaction/EOS 可以覆盖一部分 Kafka 内部的读写链路，但不能自动包住任意外部系统。
