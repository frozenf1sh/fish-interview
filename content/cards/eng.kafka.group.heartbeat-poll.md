---
id: eng.kafka.group.heartbeat-poll
kind: engineering
title: Kafka Heartbeat 与 poll：进程活着不等于还能消费
summary: Heartbeat 告诉 Group 成员仍在线，poll/处理间隔决定它是否还能按时承担分区；业务阻塞可能让成员失去归属并触发 Rebalance。
parents: [eng.kafka.group]
tags: [kafka, consumer, heartbeat, poll, rebalance]
links: [eng.kafka.consumer.pull-poll, eng.kafka.rebalance, eng.kafka.consumer.offset-commit]
---

## 先区分三个时间概念

- **Heartbeat（心跳）**：Consumer 定期向 Group 表示“我还活着”。
- **Session timeout（会话超时）**：Group 多久收不到有效心跳后，认为成员失联。
- **poll 间隔约束**：两次消费循环之间允许间隔多久；过久可能说明成员处理不了当前任务。

心跳解决“进程是否还在线”，不等于业务处理已经完成；`poll()` 也不只是取数据，还可能推动协调和分区变化。

## 一个十分钟处理任务

```text
poll 拿到一批消息
  → 业务处理阻塞 10 分钟
  → 长时间没有下一次 poll
  → Group 认为成员不能继续及时工作
  → P0 可能被交给其他 Consumer
```

调大超时只能把判定推迟，不能解决单条任务过重、外部依赖卡住或线程模型错误。常见做法是让 Poll 线程按时循环，再把有限任务交给受控处理池；但同一 Partition 的顺序、提交边界和分区撤销仍必须由应用保证。

排查时对齐 poll 间隔、单批耗时、心跳/会话事件、分配交接和提交失败，不要只看 Consumer 进程还在就断言它仍拥有 Partition。
