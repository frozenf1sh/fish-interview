---
id: eng.kafka.consumer.pull-poll
kind: engineering
title: Kafka Consumer Pull 与 poll：消费者自己来取
summary: Consumer 通过 poll 主动向已分配的 Partition 拉取一批 Record；主动拉取让客户端控制节奏，也把拉取、处理和提交分成不同边界。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, poll, fetch]
links: [eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.consumer.offset-commit, eng.kafka.group.heartbeat-poll]
---

## 先看一次消费循环

**Consumer** 是读取 Kafka 的客户端。它不是等 Broker 把每条消息推到业务函数，而是自己定期调用 `poll()`：

```text
Consumer
  → 向自己拥有的 Partition 请求数据
  → poll() 返回一批 Record
  → 业务处理这批 Record
  → 再次 poll()
```

**Pull** 是“消费者主动拉取”；**Fetch** 是客户端向 Broker 请求数据；`poll()` 是很多客户端把拉取结果交给应用、同时推进消费循环的入口。

## 为什么要主动拉取

Consumer 可以根据自己的处理能力控制每次拿多少、何时再拿，并且容易做 Batch（批量处理）和背压。Broker 不会因为某个 Consumer 读过就删除日志；Consumer 只能读取当前分配给它的 Partition，其他 Group 仍可从各自位置读取。

## `poll()` 不只是“拿数组”

在常见 Consumer 客户端中，`poll()` 还会参与 Group 协调、接收分区分配变化和维持消费循环。业务处理如果长时间不返回，可能导致心跳/处理间隔异常、失去 Partition 或触发 Rebalance；这与 [[eng.kafka.group.heartbeat-poll]] 一起看。
