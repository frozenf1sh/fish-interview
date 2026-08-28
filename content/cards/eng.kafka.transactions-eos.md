---
id: eng.kafka.transactions-eos
kind: engineering
title: Kafka Transaction 与 Exactly-Once：先画清事务边界
summary: Kafka Transaction 可以把特定的 Kafka 输出和 Consumer 位点放进一个原子边界，但不会自动包住任意数据库、RPC 或业务副作用。
parents: [eng.kafka.semantics]
tags: [kafka, transaction, eos, idempotence]
links: [eng.kafka.delivery-semantics, eng.kafka.producer.retry-idempotence, eng.kafka.consumer.processing-idempotence]
---

## 先看一个“读 A 写 B”流程

```text
读取 Topic A
  → 计算或转换
  → 写入 Topic B
  → 提交 A 的 Offset
```

如果写 B 成功后进程崩溃、A 的位点没提交，重启会再读 A、再写 B。**Transaction（事务）** 的目标是把一组相关操作放进“要么一起对外可见、要么一起失败”的边界；**EOS** 是 Exactly-Once Semantics 的缩写，意思是恰好一次语义。

## Kafka 事务能覆盖什么

在 Kafka 支持的读写链路里，可以协调输出记录和输入位点，使下游按事务读取时不会把未提交的中间结果当成正常结果。它解决的是 Kafka 内部“输出记录 + 位点”的一致关系。

```text
Kafka A 的读取位置 + Kafka B 的输出
              ↓
       Kafka 事务边界
```

## 它不能自动覆盖什么

如果流程中还有“写数据库”“调用支付 RPC”“发邮件”，这些外部系统不自动加入 Kafka 事务。仍需业务幂等、Outbox（先把待发事件写入本地事务表）或其他跨系统一致性方案。

所以回答 Exactly-Once 时必须补充“对什么系统、从哪里到哪里”。Producer 幂等性只是发送阶段去重，消费业务的完整边界见 [[eng.kafka.delivery-semantics]]。
