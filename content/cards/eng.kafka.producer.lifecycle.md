---
id: eng.kafka.producer.lifecycle
kind: engineering
title: Kafka Producer 生命周期：复用、Flush 与 Close
summary: KafkaProducer 是带连接、Metadata、缓冲和后台线程的重量级客户端，应长期复用并在进程退出前完成 Flush/Close。
parents: [eng.kafka.producer]
tags: [kafka, producer, flush, lifecycle]
links: [eng.kafka.producer.send, eng.kafka.producer.accumulator, eng.kafka.delivery-semantics]
---

## `flush()` 做什么

`flush()` 会推动此前积累的 Record 尽快发送，并等待这些发送完成。它适合明确的阶段边界，不适合每发一条消息就调用；后者会破坏 Batch 和异步 I/O 的优势。

## `close()` 做什么

进程直接退出时，Accumulator 中的 Record 可能还没有发出。正常关闭 Producer 要给 Sender 处理待发送数据、完成资源释放和失败回调的机会。

推荐的生命周期是：应用启动时创建，多个业务请求复用，优雅关闭时 Flush/Close。每条消息都 `new → send → close` 会反复丢弃连接、Metadata 和 Batch 的收益。
