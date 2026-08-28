---
id: eng.kafka.performance.batch-compression
kind: engineering
title: Kafka Batch 与 Compression：把固定开销摊薄
summary: Producer、Broker 和 Consumer 都倾向批量处理，压缩也通常以 Batch 为单位，从而减少请求、协议和网络开销。
parents: [eng.kafka.performance]
tags: [kafka, performance, batch, compression]
links: [eng.kafka.producer.accumulator, eng.kafka.producer.sender, eng.kafka.performance.partition-parallelism]
---

## 从单条请求到批量请求

单条发送会重复支付系统调用、请求头、网络往返和 Broker 请求处理成本。Producer 把同一 Partition 的多条 Record 聚成 Batch，Sender 再按 Broker 合并请求；压缩一批相近数据通常也比逐条压缩更有效。

Batch 越大不一定越好：它会增加等待、内存、压缩和失败重试的成本。`batch.size`、`linger.ms` 和压缩策略应结合消息大小、流量和延迟目标观察。
