---
id: eng.kafka.replication.write-ack
kind: engineering
title: Kafka ACK 与 min.insync.replicas：什么时候算写成功
summary: ACK 是 Producer 对 Broker 写入结果的等待边界；min.insync.replicas 是允许成功写入所需的最小同步副本数，两者共同决定故障时的取舍。
parents: [eng.kafka.replication]
tags: [kafka, acks, isr, durability]
links: [eng.kafka.producer.sender, eng.kafka.producer.retry-idempotence, eng.kafka.replication.isr-election, eng.kafka.delivery-semantics]
---

## 先区分“写进 Leader”和“有副本确认”

Producer 发来一批记录后，Leader 可以先把它追加到自己的日志；Follower 再通过 Fetch 把数据复制过去。**ACK** 就是 Producer 在请求中表达的“我等到哪一步才算成功”。

| 设置 | 等待到哪里 | 故障时的直觉风险 |
| --- | --- | --- |
| `acks=0` | 不等 Broker 响应 | 失败难以及时发现 |
| `acks=1` | Leader 追加后响应 | Leader 尚未复制就故障，最新写入可能丢失 |
| `acks=all` | ISR 满足确认条件后响应 | 等待和不可写时间可能增加 |

![Kafka ISR、Follower Fetch 与 ACK 确认边界](/static/kafka-replication-ack.drawio.svg)

## `all` 不是“所有历史副本”

`acks=all` 关注的是当前 **ISR**，也就是当前跟得上的副本集合，不是所有曾经配置过、但已经落后的 Replica。`min.insync.replicas=2` 可以再加一道门槛：至少有两个同步副本时才允许这种写入成功。

```text
ISR = Leader + Follower 1 + Follower 2 → 可以写
ISR = Leader + Follower 1             → 仍可能写，取决于门槛
ISR = Leader                          → min=2 时拒绝写入
```

拒绝写入牺牲的是故障时的可用性，换来更强的副本确认边界。继续提高 `acks` 也不会自动保证 Consumer 已经处理业务；那是另一段 Offset Commit 和业务幂等问题。

## 面试时完整回答

不要把 `acks=1` 简化成“不落盘”、把 `acks=all` 简化成“绝不丢”。还要继续问：ISR 是否健康、Leader 是否切换、ACK 是否丢失、Producer 是否重试，以及下游是否已经提交处理结果。
