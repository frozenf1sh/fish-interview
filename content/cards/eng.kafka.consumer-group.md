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

## 核心概念

Topic 被拆为多个 Partition。Consumer Group 将分区分配给组内消费者，以实现并行消费；一个分区不能被同组两个消费者同时处理。

## 故障定位

先区分“没有消息”、“没有分到分区”、“消费了但未提交”、“提交了但业务未成功”。Lag 只能说明积压，不足以单独说明根因。

## 关联

成员变化或订阅变化会触发 [[eng.kafka.rebalance]]。

