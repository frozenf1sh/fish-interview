---
id: eng.kafka.producer.send
kind: engineering
title: Kafka Producer send：一次发送不是一次网络请求
summary: producer.send() 先把业务对象编码、选择 Partition 并放入本地缓冲，后台线程再批量发送；返回 Future 不代表 Broker 已确认。
parents: [eng.kafka.producer]
tags: [kafka, producer, send, serialization]
links: [eng.kafka.record, eng.kafka.key-partitioner, eng.kafka.producer.accumulator, eng.kafka.producer.sender]
---

## 先看业务代码背后发生了什么

业务代码可能只有一行：

```text
producer.send(topic="orders", key="order-123", value=OrderCreated)
```

但 Kafka Producer（生产者客户端）通常会按这条路径处理：

```text
业务对象
  → Serializer：对象变成字节
  → Partitioner：决定放进哪个 Partition
  → Accumulator：放进本地待发送 Batch
  → Sender Thread：后台线程发网络请求
  → Broker：追加并返回结果
```

**Serializer** 是序列化器，负责把对象编码成 Kafka 能保存的字节；**Partitioner** 是分区器，负责选择 Partition；**Batch** 是为了批量发送而暂存的一组 Record；**Future** 是以后可以拿到发送结果的对象。

![Kafka Producer 从 send 到 Broker 的链路](/static/kafka-producer-path.drawio.svg)

## 为什么 `send()` 看起来很快

`send()` 通常先把 Record 放入 Producer 内存中的缓冲区，不必为每条消息立即发一次网络请求。后台 **Sender Thread**（发送线程）稍后把同一 Partition 的记录组成 Batch，再发给对应 Leader。

所以三件事不是一回事：

1. `send()` 方法被调用并返回。
2. Record 离开 Producer，真正发到网络。
3. Broker 追加成功并返回 ACK。

只有 Future 或 Callback 报告成功，才说明 Producer 侧拿到了成功响应；如果应用直接退出，仍在缓冲区的记录可能没有发出。缓冲细节见 [[eng.kafka.producer.accumulator]]。

## 使用边界

Producer 是带连接、路由表、缓冲区和后台线程的长期客户端。通常应用启动时创建、多个请求复用、退出时关闭，不要每条消息都创建一个 Producer。完整生命周期见 [[eng.kafka.producer.lifecycle]]。
