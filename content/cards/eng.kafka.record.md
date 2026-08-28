---
id: eng.kafka.record
kind: engineering
title: Kafka Record：Key、Value 与定位信息
summary: Kafka 处理的是带有 Topic、Partition、Offset 和字节载荷的 Record；业务对象要先由客户端序列化。
parents: [eng.kafka.model]
tags: [kafka, record, serialization]
links: [eng.kafka.key-partitioner, eng.kafka.producer.send, eng.kafka.partition]
---

## Record 由什么组成

业务代码看到的可能是一个订单对象，但 Kafka 客户端最终处理的是类似下面的记录：

```text
Record = topic + partition? + key? + value + headers + timestamp
```

`value` 是业务载荷，`key` 常用于决定分区和表达同一业务实体；`partition` 可以由调用方显式指定，也可以交给 Partitioner；`offset` 则是在 Leader 追加时确定，不由 Producer 预先生成。

## Kafka 不理解业务对象

JSON、Protobuf、Avro 或自定义结构都要经 Serializer 转成字节，Consumer 再用对应 Deserializer 还原。序列化协议、Schema 演进和兼容性是业务契约，不应误认为 Kafka 自动提供。

完整定位一条 Record 时，要把 `Topic + Partition + Offset` 放在一起；单独一个 Offset 没有全局含义。
