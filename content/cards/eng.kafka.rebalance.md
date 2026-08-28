---
id: eng.kafka.rebalance
kind: engineering
title: Kafka Rebalance：重新决定谁读哪个 Partition
summary: Rebalance 是 Consumer Group 因成员或订阅变化而重新分配 Partition 的过程；它会影响消费暂停、位点提交和未提交处理的重放。
parents: [eng.kafka.group]
tags: [kafka, rebalance, coordination]
links: [eng.kafka.consumer-group, eng.kafka.group.coordinator-assignment, eng.kafka.group.assignment-strategies, eng.kafka.group.rebalance-lifecycle, eng.kafka.group.heartbeat-poll, eng.kafka.consumer.offset-commit]
---

## 先看一个成员变化

```text
原来：C1 → P0、P1；C2 → P2、P3
C3 加入
结果：重新决定 P0、P1、P2、P3 分给谁
```

**Rebalance** 的直译是“重新平衡”，在 Kafka 中就是 Consumer Group 重新确认 Partition owner。触发原因包括 Consumer 加入、退出、崩溃、订阅变化和 Topic 增加 Partition。

## 它为什么会造成重复

旧 Consumer 失去 P0 后，新 Consumer 会从 Group 已保存的 committed offset 继续。如果旧 Consumer 在处理完成后还没提交书签就被交接，同一条记录可能再次处理；这不是 Rebalance 自己复制了消息，而是恢复位置没有越过它。

不同协议的暂停范围不同：经典协议常见更明显的全组同步，增量/合作式协议可以缩小交接范围。因此不要笼统说“每次 Rebalance 都会全组停顿”，但应用必须接受分区暂时不可用和位点重放。

## 为什么线上会频繁发生

处理线程长时间阻塞、不再按时 `poll()`、心跳或网络异常、滚动发布和实例反复重启，都可能让 Group 认为成员不再健康。`max.poll.interval` 是两次 `poll`/处理循环之间的时间约束；它不是把慢业务自动变快的开关。

加入、撤销、重新分配和恢复的顺序见 [[eng.kafka.group.rebalance-lifecycle]]；排查提交窗口时看 [[eng.kafka.consumer.offset-commit]]。
