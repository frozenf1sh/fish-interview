---
id: eng.kafka.producer.retry-idempotence
kind: engineering
title: Kafka Producer 重试与幂等：处理 ACK 丢失
summary: Broker 可能已经写入但 ACK 在网络中丢失，重试会产生重复；幂等 Producer 用 Producer ID 与序列信息识别同一发送的重复批次。
parents: [eng.kafka.producer]
tags: [kafka, producer, retry, idempotence]
links: [eng.kafka.producer.sender, eng.kafka.replication.write-ack, eng.kafka.delivery-semantics, eng.kafka.transactions-eos]
---

## 重复从哪里来

```text
Producer → Broker：写入 m1
Broker：已追加 m1
ACK：在网络中丢失
Producer：认为超时，重试 m1
```

如果客户端和 Broker 没有重复识别机制，日志中可能出现两份 m1。单纯把 `retries` 调大只能提高暂时故障下的成功机会，不能自动消除重复。

## 幂等性的边界

启用幂等 Producer 后，Broker 可以结合 Producer ID 和 Partition 上的序列信息识别同一逻辑发送的重试，从而避免重复追加。并发请求数与重试还会影响顺序，因此不能把“开启幂等”理解成所有业务步骤的 Exactly-Once。

它主要覆盖 Producer 到 Kafka Log 的一段；Consumer 读出后写数据库或再写另一个 Topic，仍需看事务和业务幂等。
