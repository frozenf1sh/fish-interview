---
id: eng.kafka.group.heartbeat-poll
kind: engineering
title: Kafka Consumer Heartbeat 与 poll：活着不等于处理得完
summary: 心跳、会话存活和 poll 处理间隔共同影响 Consumer 是否留在 Group；长时间业务阻塞可能导致失去分区并触发 Rebalance。
parents: [eng.kafka.group]
tags: [kafka, consumer, heartbeat, poll, rebalance]
links: [eng.kafka.consumer.pull-poll, eng.kafka.rebalance, eng.kafka.consumer.offset-commit]
---

## 三个时间问题

- 心跳用于让 Group 知道成员仍在运行。
- 会话超时用于判断成员是否失联。
- `max.poll.interval` 一类约束用于限制两次处理循环之间允许间隔。

如果一次业务处理耗时十分钟，而 Consumer 长时间不再进入 poll，组可能认为它不能继续承担分区。调大参数只能改变判定窗口，不能解决单条处理过重、线程模型错误或外部依赖阻塞。

## 排查顺序

对齐 poll 间隔、单批处理耗时、心跳与会话事件，再看是否发生分配交接和提交失败。不要只看“Consumer 进程还在”就判断它仍拥有 Partition。
