---
id: eng.kafka.rebalance
kind: engineering
title: Kafka Rebalance：重新分配分区的协调过程
summary: 消费者组成员、订阅或分区变化后，协调器会暂停并重新分配分区；频繁 Rebalance 会放大消费抖动。
parents: [eng.kafka]
tags: [kafka, rebalance, availability]
links: [eng.kafka.consumer-group]
---

## 核心概念

Rebalance 的目标是让分区归属与当前组成员一致。它是正常协调机制，但在处理时间过长、心跳异常或成员频繁启停时会成为可见的可用性问题。

## 排查顺序

先看组成员变更和协调器日志，再比较处理耗时与心跳/会话超时配置，最后检查部署滚动、网络抖动与订阅变更。

## 常见误区

不要只提高超时时间掩盖问题；若单条消息处理长期阻塞，根因仍是消费链路的处理模型。

## 关联

[[eng.kafka.consumer-group]]。

