---
id: eng.kafka.transactions-eos
kind: engineering
title: Kafka Transaction 与 Exactly-Once：明确事务覆盖范围
summary: Kafka 事务可以把特定的 Kafka 读写和 Offset 提交放进原子边界，但不会自动覆盖任意数据库或外部 RPC。
parents: [eng.kafka.semantics]
tags: [kafka, transaction, eos, idempotence]
links: [eng.kafka.delivery-semantics, eng.kafka.producer.retry-idempotence, eng.kafka.consumer.processing-idempotence]
---

## 先问“Exactly-Once 对什么”

典型的 Consume-Process-Produce 链路是：

```text
读 Topic A → 计算 → 写 Topic B → 提交 A 的 Offset
```

如果只靠普通发送和提交，任一步之后崩溃都可能重做。Kafka Transaction 可以在 Kafka 支持的事务边界内协调输出和位点，让下游看到一致结果；但“写 Kafka 后调用支付服务”仍然需要外部幂等、Outbox 或其他一致性方案。

事务消息还涉及下游的读取隔离：只有按照事务语义读取时，消费者才会把已提交事务的结果作为正常输出处理。它解决的是 Kafka 内部记录与位点的原子关系，不是让任意 Consumer 代码自动变成不可重复的函数。

所以 Idempotent Producer 是发送阶段能力，EOS 是更大的处理协议，两者不能画等号。
