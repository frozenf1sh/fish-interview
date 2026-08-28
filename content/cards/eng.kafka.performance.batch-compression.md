---
id: eng.kafka.performance.batch-compression
kind: engineering
title: Kafka Batch 与 Compression：一次处理一批数据
summary: Kafka 把多条 Record 聚成 Batch 再发送和压缩，用更少的请求摊薄固定开销；Batch 过大又会增加等待、内存和失败重试成本。
parents: [eng.kafka.performance]
tags: [kafka, performance, batch, compression]
links: [eng.kafka.producer.accumulator, eng.kafka.producer.sender, eng.kafka.performance.partition-parallelism]
---

## 先看为什么不要每条都发请求

一次网络请求有固定成本：系统调用、请求头、网络往返和 Broker 处理。若每条 500 字节的消息都单独请求，1000 条消息就要承担 1000 次固定成本。

**Batch（批次）** 是把多条 Record 捆在一起；Producer 按 Partition 聚成 Batch，Sender 再按 Broker 发送：

```text
m1、m2、m3、…、m100 → 一个 Batch → 一次 ProduceRequest
```

**Compression（压缩）** 是把数据编码得更小。多条相似记录一起压缩，通常比每条单独压缩更划算，也能减少网络流量。

## Batch 不是越大越好

Batch 越大，通常越能摊薄开销；但低流量时要等更久，内存占用、单次失败重试量和尾延迟也会增加。`batch.size` 控制目标大小，`linger.ms` 控制为凑批次允许等待的窗口，应该结合消息大小、流量和延迟目标观察。

Batch 还需要依赖 [[eng.kafka.producer.accumulator]] 和 [[eng.kafka.producer.sender]] 才能真正发出。
