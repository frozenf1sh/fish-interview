---
id: eng.kafka.delivery-semantics
kind: engineering
title: Kafka 投递语义：一条消息可能被处理几次
summary: At-Most-Once、At-Least-Once 和 Exactly-Once 描述的是一次处理链路的故障取舍，不能简单等同于某一个 ACK 或 Producer 配置。
parents: [eng.kafka.semantics]
tags: [kafka, delivery-semantics, reliability]
links: [eng.kafka.replication.write-ack, eng.kafka.producer.retry-idempotence, eng.kafka.consumer.offset-commit, eng.kafka.consumer.processing-idempotence, eng.kafka.transactions-eos]
---

## 先从一次崩溃开始

Consumer 读到支付事件后写数据库，进程可能在“数据库已经写好、Offset 还没提交”之间崩溃。重启时它会再次读到这条事件。

这里的 **投递语义** 讨论的是：在网络、重试和崩溃下，一条逻辑消息最多可能被处理几次，以及结果能否被确认。它不是在问 Kafka 是否物理保存了多份消息。

## 三种常见取舍

| 名称 | 直白解释 | 常见做法与风险 |
| --- | --- | --- |
| At-Most-Once（至多一次） | 尽量不重复，但可以跳过 | 先提交位点再处理；处理中崩溃可能丢业务结果 |
| At-Least-Once（至少一次） | 尽量不跳过，但允许重放 | 先处理再提交；提交前崩溃会重复 |
| Exactly-Once（恰好一次） | 在明确边界内只让结果生效一次 | 需要事务或幂等设计，不能覆盖任意外部系统 |

**Offset Commit** 是保存 Group 书签；**ACK** 是 Producer 等待 Broker 写入结果；业务数据库是否成功是第三个边界。把 `acks=all` 或 Producer 幂等性直接说成整条链路 Exactly-Once，都是把边界混在一起。

工程上常见的落点是 At-Least-Once 加业务幂等，再在 Kafka 内部读写链路确有需要时使用事务，见 [[eng.kafka.consumer.processing-idempotence]] 和 [[eng.kafka.transactions-eos]]。
