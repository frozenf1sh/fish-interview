---
id: eng.kafka.group.rebalance-lifecycle
kind: engineering
title: Kafka Rebalance 生命周期：加入、同步与交接
summary: Rebalance 是 Consumer Group 重新确定成员与 Partition 归属的生命周期；理解加入、撤销、分配和恢复顺序，才能解释暂停、提交失败与重复消费。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, rebalance, assignment, offset]
links: [eng.kafka.rebalance, eng.kafka.group.coordinator-assignment, eng.kafka.group.assignment-strategies, eng.kafka.consumer.commit-modes, eng.kafka.group.heartbeat-poll]
---

## 一条通用生命周期

不同 Kafka 版本和 Rebalance Protocol 的具体 RPC 不完全相同，但可以用下面这条主线建立模型：

```text
加入或发现变更
      ↓
协调器确认 Group 成员与订阅
      ↓
计算新的 Assignment
      ↓
撤销不再拥有的 Partition
      ↓
成员应用新的 Assignment
      ↓
从已提交位置继续 Fetch / 处理 / Commit
```

经典协议通常体现为更明显的全组同步；增量或合作式协议可以只交接受影响的 Partition。不要据此断言“每次 Rebalance 都会全组停顿”，但应用始终要为 Partition 暂时不可用、归属变化和处理重放做好准备。

## 交接时各阶段关心什么

| 阶段 | 关键问题 | 失败表现 |
| --- | --- | --- |
| Join | 谁仍是 Group 成员，谁的订阅发生变化？ | 成员反复加入，Generation 不稳定 |
| Assignment | 每个 Partition 的新 owner 是谁？ | 分配不均、热点成员或大量迁移 |
| Revoke | 旧 owner 是否停止处理并保存必要状态？ | 旧成员继续处理或提交失效分区 |
| Assign | 新 owner 是否完成本地初始化？ | 读取位置错误、缓存未恢复 |
| Resume | 从哪个 committed offset 继续？ | 未提交处理被重放，或错误提交造成跳过 |

如果 Consumer 在失去归属后仍提交旧 Partition，提交可能被拒绝，也可能暴露更深的线程模型错误。处理线程、提交线程和分区生命周期必须有明确的所有权；不能让一个无界任务池在分区已经交接后继续写入副作用。

## 为什么会和重复消费一起出现

```text
处理 P0@10
      │
      ├── commit P0@11 成功 → 新 owner 从 11 继续
      └── Rebalance / Crash 先发生 → 新 owner 从 11 之前重放
```

Rebalance 本身不等于消息丢失或重复；它只是改变消费归属。重复通常来自“业务处理已完成，但位点尚未提交”或提交结果不确定。可靠处理要同时看 [[eng.kafka.consumer.commit-modes]]、[[eng.kafka.consumer.processing-idempotence]] 和 [[eng.kafka.group.heartbeat-poll]]。

## 排查顺序

先确认触发原因与 Group Generation，再对齐 `poll` 间隔、单批处理耗时、心跳/会话事件和提交结果，最后检查 Assignment 是否发生大规模迁移。只调大超时而不拆分处理或修正所有权，通常只能推迟下一次 Rebalance。
