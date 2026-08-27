---
id: eng.kafka.consumer-group
kind: engineering
title: Kafka Consumer Group：分区归属与位点
summary: 同一消费者组中，每个分区同一时刻只由一个消费者消费；消费进度通过已提交 offset 表达。
parents: [eng.kafka]
tags: [kafka, consumer-group, offset]
links: [eng.kafka.rebalance]
exam_signals:
  - company: mihoyo
    year: 2027
    role: backend
    confidence: low
    source: https://www.nowcoder.com/
---

## 核心机制

Topic 被拆为多个 Partition。Consumer Group 的协调器维护成员关系和分区归属：同组一个 Partition 同一时刻只能交给一个消费者，而不同分区可以并行处理。组内消费者数超过分区数时，多出的消费者会处于空闲状态。

消费者通过 poll 获取记录，业务处理完成后提交 offset。自动提交、同步提交和异步提交的取舍本质上是在吞吐、重复处理窗口和失败恢复确定性之间权衡。业务副作用不是 Kafka 自动事务的一部分：写数据库、发 RPC 等操作仍需幂等键、去重或外部事务设计。

## 故障定位

先区分“没有消息”、“没有分到分区”、“拉取到但未处理”、“处理了但未提交”、“提交了但业务副作用失败”。Lag 只能说明积压，不足以单独说明根因；还要看分区分配、处理耗时和提交位点。

## 关联

成员变化或订阅变化会触发 [[eng.kafka.rebalance]]。
