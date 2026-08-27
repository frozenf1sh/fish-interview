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

Rebalance 的目标是让分区归属与当前组成员一致。成员加入/退出、订阅变化或分区变化时，协调器会触发新一轮分配。不同协议对暂停范围不同，但应用都应假定：分区在交接期间可能暂时不可消费，且未妥善提交的处理会被重新执行。

## 排查顺序

先看组成员变更和协调器日志，再比较处理耗时与 `max.poll.interval`、心跳和会话超时的关系，最后检查部署滚动、网络抖动与订阅变更。处理线程长时间阻塞却没有继续 poll，是频繁被踢出组的常见机制性原因。

## 常见误区

不要只提高超时时间掩盖问题；若单条消息处理长期阻塞，根因仍是消费链路的处理模型。

## 关联

[[eng.kafka.consumer-group]]。
