---
id: eng.kafka.consumer.pull-poll
kind: engineering
title: Kafka Consumer Pull 与 poll：消费者主动拉取数据
summary: Kafka Consumer 通过 poll 从分配到的 Partition 拉取数据，消费节奏由客户端和处理能力共同控制，而不是 Broker 主动逐条 Push。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, poll, fetch]
links: [eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.consumer.offset-commit, eng.kafka.group.heartbeat-poll]
---

## 基本路径

```text
Consumer
  → 根据分配结果向 Partition Leader Fetch
  → poll() 返回一批 Record
  → 业务处理
  → 再次 poll()
```

Pull 模型让 Consumer 可以按自己的处理速度读取，也便于按 Batch 处理和控制背压。它不会绕过 Group 的分区归属：Consumer 只能读取当前分配给自己的 Partition。

## `poll` 不只是取消息

在常见客户端模型中，poll 还参与组协调、获取分配变化和推进消费循环。业务处理长期阻塞、不再调用 poll，可能同时触发处理超时、提交失败和 Rebalance；需要和 [[eng.kafka.group.heartbeat-poll]] 一起排查。
