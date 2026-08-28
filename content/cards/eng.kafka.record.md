---
id: eng.kafka.record
kind: engineering
title: Kafka Record：一条事件到底带什么
summary: Record 是 Producer 要写入 Kafka 的一条记录；业务对象会先被序列化成字节，Partition 和 Offset 则分别在路由与追加时确定。
parents: [eng.kafka.model]
tags: [kafka, record, serialization]
links: [eng.kafka.key-partitioner, eng.kafka.producer.send, eng.kafka.partition]
---

## 先看一张“寄件单”

业务代码可能准备的是一个订单对象：

```text
orderId = 123
event   = OrderCreated
amount  = 100
```

Producer 发送给 Kafka 的 **Record** 可以理解成带地址和内容的一张寄件单：

```text
Record = topic + partition? + key? + value + headers + timestamp
```

- `topic`：写入哪类事件。
- `partition`：可由调用方指定，也可让客户端计算；问号表示它可能还没确定。
- `key`：可选的业务关联字段，例如 `orderId`。
- `value`：真正的业务内容。
- `headers`、`timestamp`：附加元数据。

## Kafka 为什么看不懂 Order 对象

Kafka Broker 主要保存字节，不理解 Java 对象、Go 结构体或业务字段。Producer 用 **Serializer（序列化器）** 把对象变成字节，Consumer 再用 **Deserializer（反序列化器）** 还原：

```text
OrderCreated 对象 → Serializer → bytes → Kafka
Kafka 中的 bytes → Deserializer → OrderCreated 对象
```

JSON、Protobuf、Avro 都是不同的编码方式。编码格式、字段兼容和版本升级属于业务契约，Kafka 不会自动替你解决。

## 谁决定哪个字段什么时候出现

Producer 的 Partitioner 可以利用 `key` 选择 Partition；Broker 的 Leader 把 Record 追加成功后才分配 Offset。也就是说，Key 影响写入分片，Offset 是日志位置，不是 Producer 自己提前生成的业务 ID，见 [[eng.kafka.key-partitioner]]。
