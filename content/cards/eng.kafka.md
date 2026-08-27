---
id: eng.kafka
kind: engineering
title: Kafka：从日志、分区到消费者组的协作模型
summary: Kafka 用追加日志保存事件，以 Partition 建立顺序与并行边界，再通过 Consumer Group 协调消费归属和位点。
parents: [engineering]
tags: [kafka, messaging, distributed-systems]
links: [eng.kafka.topic, eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.rebalance]
---

## 核心机制

Kafka 不是“一个消息队列”这么简单。Topic 是逻辑主题；每个 Topic 被拆成多个有序的 [[eng.kafka.partition]]；Producer 选择一个分区追加记录；Consumer Group 将分区分配给不同消费者。顺序只在单个 Partition 内成立，并行度也受 Partition 数量上限约束。

每条记录带有递增的 offset。offset 是该 Partition 上的位置，而不是全局消息 ID。消费者处理成功后提交位点，表示“下一次从哪里继续”；提交策略与业务副作用如何配合，决定了至少一次、至多一次还是更强的语义边界。

## 排查主线

先定位问题落在何条链路：写入是否到达正确 Topic/Partition、记录是否可见、消费者是否拿到分区、业务是否实际处理、位点是否按预期提交。不要用一项 Lag 指标替代整条链路的判断。

## 关键边界

- 同一 Consumer Group 内，一个分区同一时刻只属于一个消费者。
- 多消费者不能提高单分区的消费并行度；需要更多 Partition 或拆分业务键。
- 副本解决可用性，不会让同一 Partition 内的业务顺序变成全局顺序。

## 变体与关联

按最小对象继续：[[eng.kafka.topic]]、[[eng.kafka.partition]]、[[eng.kafka.consumer-group]] 与 [[eng.kafka.rebalance]]。

