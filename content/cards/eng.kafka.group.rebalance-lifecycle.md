---
id: eng.kafka.group.rebalance-lifecycle
kind: engineering
title: Kafka Rebalance 生命周期：从加入到重新消费
summary: Rebalance 可以按“成员加入—计算分工—交接分区—从书签恢复”的顺序理解；掌握这条链，就能解释暂停、提交失败和重复处理。
parents: [eng.kafka.group]
tags: [kafka, consumer-group, rebalance, assignment, offset]
links: [eng.kafka.rebalance, eng.kafka.group.coordinator-assignment, eng.kafka.group.assignment-strategies, eng.kafka.consumer.commit-modes, eng.kafka.group.heartbeat-poll]
---

## 先看一轮交接

C1 原来负责 P0。C2 加入后，P0 可能交给 C2：

```text
C1 报到/发现成员变化
  → Coordinator 确认当前成员
  → Assignment 计算新的分工
  → C1 撤销 P0
  → C2 接收 P0
  → C2 从已提交位置继续读取
```

**Coordinator** 是负责 Group 状态的协调窗口；**Assignment** 是分工结果；**Revoke** 是旧成员放弃分区；**Assign** 是新成员接收分区。不同协议的具体请求时序可能不同，但这几个动作能帮助理解应用生命周期。

## 每一步要防什么

| 阶段 | 新手应该问的问题 | 出错后常见表现 |
| --- | --- | --- |
| 成员加入 | 谁还活着、谁订阅了什么？ | 成员反复加入，Group 不稳定 |
| 计算分工 | P0 的新 owner 是谁？ | 分配不均、热点或大量迁移 |
| 撤销旧分区 | 旧 Consumer 是否停止处理并保存状态？ | 旧成员继续写副作用或提交失败 |
| 接收新分区 | 新 Consumer 是否初始化完成？ | 缓存未恢复、从错误位置读取 |
| 恢复消费 | 书签保存到哪里？ | 未提交记录被重放，错误提交造成跳过 |

经典协议可能让全组更明显地参与同步，增量/合作式协议可能只交接受影响的分区，所以不能把某一种暂停表现当成所有版本的规则。

## Rebalance 与重复消费的关系

```text
处理 P0@10
  ├── 先提交 P0@11，再交接 → 新 owner 从 11 继续
  └── 交接/崩溃先发生 → 新 owner 可能重放 P0@10
```

Rebalance 只是改变归属，不等于丢消息或重复消息。重复通常来自“业务已完成但位点未提交”或“提交结果不确定”，因此要把 [[eng.kafka.consumer.commit-modes]]、[[eng.kafka.consumer.processing-idempotence]] 和 [[eng.kafka.group.heartbeat-poll]] 一起看。
