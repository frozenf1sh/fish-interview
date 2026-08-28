---
id: eng.kafka.producer.send
kind: engineering
title: Kafka Producer send：从业务对象到待发送 Record
summary: producer.send() 通常先完成序列化和分区决策，再把 Record 放入本地缓冲；它返回 Future 不等于 Broker 已确认。
parents: [eng.kafka.producer]
tags: [kafka, producer, send, serialization]
links: [eng.kafka.record, eng.kafka.key-partitioner, eng.kafka.producer.accumulator, eng.kafka.producer.sender]
---

## 调用发生了什么

```text
send(topic, key, value)
  → Key/Value Serializer
  → Partitioner 选择 Partition
  → 放入 RecordAccumulator
  → Sender Thread 发送
```

Kafka 客户端不理解业务对象。序列化失败、Topic 元数据暂时不可用或本地缓冲区不足，都可能在真正收到 Broker 响应前影响调用结果。

![Kafka Producer 从 send 到 Broker 的链路](/static/kafka-producer-path.drawio.svg)

## Future 的含义

`send()` 返回的 Future 表示这次异步发送最终可以得到结果。消息可能仍在 Producer 内存、等待 Batch、等待 Metadata 或等待 Broker 的 ACK；只有回调或 Future 完成并报告成功，才有 Producer 侧的确认。

Producer 应长期复用，让 Metadata、连接和 Batch 发挥作用，而不是每条消息都创建和关闭一个客户端。
