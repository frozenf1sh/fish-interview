---
id: eng.kafka.replication.write-ack
kind: engineering
title: "Kafka ACK 与 min.insync.replicas：写入确认的可靠性边界"
summary: "`acks` 决定 Producer 等待到什么程度，`min.insync.replicas` 限制 ISR 至少要有多少成员；两者共同表达可接受的写入风险。"
parents: [eng.kafka.replication]
tags: [kafka, acks, isr, durability]
links: [eng.kafka.producer.sender, eng.kafka.producer.retry-idempotence, eng.kafka.replication.isr-election, eng.kafka.delivery-semantics]
---

## 三种确认意图

| 设置 | Producer 等待的边界 | 主要风险 |
| --- | --- | --- |
| `acks=0` | 不等待 Broker 响应 | 发送失败难以及时感知 |
| `acks=1` | Leader 追加后响应 | Leader 尚未复制就故障可能丢失 |
| `acks=all` | ISR 满足确认条件后响应 | 可用性和延迟成本更高 |

`acks=all` 关注的是当前 ISR，不是所有历史 Replica。`min.insync.replicas=2` 可以要求 ISR 至少有两个成员，否则宁愿拒绝高可靠写入，也不把只剩 Leader 的状态伪装成可靠成功。

这只定义 Producer 侧写入确认；Consumer 是否处理完成，还要看 Offset Commit 和业务副作用。
