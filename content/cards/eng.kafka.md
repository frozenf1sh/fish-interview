---
id: eng.kafka
kind: engineering
title: Kafka：从日志、分区到消费者组的协作模型
summary: Kafka 用追加日志保存事件，以 Partition 建立顺序与并行边界，再通过 Consumer Group 协调消费归属和位点。
parents: [engineering]
tags: [kafka, messaging, distributed-systems]
links: [eng.kafka.topic, eng.kafka.partition, eng.kafka.consumer-group, eng.kafka.rebalance]
---

> **核心结论**：Kafka 的顺序、并行度、offset 与副本选主，都以 **Partition** 为边界；不是以 Topic 或 Consumer 为边界。

## 一条消息的最短链路

`Producer → Topic → Partition → Broker Log → Consumer Group → 业务处理 → offset 提交`

Topic 是逻辑主题；每条记录最终落在一个有序的 [[eng.kafka.partition]]。Consumer Group 把分区分给组内成员：顺序只在单 Partition 内成立，并行度也受 Partition 数量限制。

## 必须区分的两个状态

| 状态 | 它表达什么 | 不能表达什么 |
| --- | --- | --- |
| Log offset | 记录在某个 Partition 中的位置 | 全局消息 ID |
| Consumer offset | 某个组下一次准备从何处读取 | 业务副作用已经可靠完成 |

消费后何时提交位点，决定故障时的重复处理窗口；写数据库、RPC 等外部副作用仍需要幂等键、去重或事务边界设计。

## 排查顺序

1. 写入是否进入预期 Topic / Partition？
2. Consumer Group 是否拿到了该分区？
3. 拉取后是否实际完成业务处理？
4. 位点提交与业务成功是否一致？

> **不要只看 Lag。**它只能说明积压，不能单独说明是生产、分配、处理还是提交出了问题。

## 三条边界

- 同一组内，一个 Partition 同一时刻只属于一个消费者。
- 多加消费者不能提高单分区并行度；热点应先查 key 分布和分区策略。
- 副本提升容错能力，不会把局部顺序变成全局顺序。

## 继续拆分

按最小对象继续：[[eng.kafka.topic]]、[[eng.kafka.partition]]、[[eng.kafka.consumer-group]] 与 [[eng.kafka.rebalance]]。
