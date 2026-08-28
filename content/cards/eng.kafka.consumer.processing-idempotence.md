---
id: eng.kafka.consumer.processing-idempotence
kind: engineering
title: Kafka 消费处理与业务幂等：重复是可预期的故障路径
summary: Consumer 处理成功但位点未提交时崩溃，重启会再次读取同一 Record；业务副作用要能安全重复或先去重。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, idempotence, at-least-once]
links: [eng.kafka.consumer.offset-commit, eng.kafka.delivery-semantics, eng.kafka.producer.retry-idempotence, eng.kafka.transactions-eos]
---

## 先看重复是怎么发生的

```text
读到订单支付事件
  → 数据库扣款状态写成功
  → Consumer 还没提交下一 Offset 就崩溃
  → 重启后从旧位点再次读到支付事件
```

**重复消费** 不一定表示 Kafka 把消息复制错了，通常是“业务结果已经发生，但 Group 的书签还没有推进”。Kafka 为了避免跳过未处理数据，允许这种恢复路径。

## 让业务重复一次也安全

常见方法包括：

- 用事件 ID 或业务 ID 建唯一约束，第二次写入变成已存在。
- 用幂等 Upsert，让同一个状态写两次结果相同。
- 用去重表记录已经处理过的事件。
- 用状态机拒绝不合法的重复状态迁移。

**幂等** 就是同一个业务操作执行一次或多次，最终状态一致；它不代表每次调用都没有副作用，而是副作用被设计成可重复确认。

## 不要把三件事混成一个词

Producer 幂等性主要处理发送重试；Consumer 提交主要保存恢复位置；业务幂等处理外部副作用。Kafka Transaction/EOS 可以覆盖一部分 Kafka 内部读写，但不会自动把任意数据库和 RPC 包进同一个原子操作，见 [[eng.kafka.transactions-eos]]。
