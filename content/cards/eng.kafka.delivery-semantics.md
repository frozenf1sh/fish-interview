---
id: eng.kafka.delivery-semantics
kind: engineering
title: Kafka 投递语义：At-Most-Once、At-Least-Once 与 Exactly-Once
summary: 投递语义不是单个参数的别名，而是 Producer 确认、Consumer 提交、业务处理和副作用边界共同形成的结果。
parents: [eng.kafka.semantics]
tags: [kafka, delivery-semantics, reliability]
links: [eng.kafka.replication.write-ack, eng.kafka.producer.retry-idempotence, eng.kafka.consumer.offset-commit, eng.kafka.consumer.processing-idempotence, eng.kafka.transactions-eos]
---

## 三种直觉模型

| 语义 | 常见故障取舍 |
| --- | --- |
| At-Most-Once | 先提交再处理，可能跳过，但减少重复 |
| At-Least-Once | 先处理再提交，尽量不跳过，但允许重放 |
| Exactly-Once | 需要在明确的事务边界内让结果只产生一次 |

网络超时会让 Producer 不知道 Broker 是否已经写入；Consumer 崩溃会让应用不知道业务副作用是否已完成。因而不能只凭 `acks` 或幂等 Producer 宣称整条业务链路 Exactly-Once。

实际系统通常选择 At-Least-Once 加业务幂等，再根据 Kafka 内部读写链路的需要引入事务。
