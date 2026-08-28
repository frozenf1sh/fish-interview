---
id: eng.kafka.failure-diagnosis
kind: engineering
title: Kafka 故障排查：沿一条消息的旅行逐段定位
summary: Kafka 故障要沿“是否写入—写到哪里—是否复制—谁来读—是否处理—是否提交”逐段确认，不能只看一个 Lag 或错误码。
parents: [eng.kafka.operations]
tags: [kafka, troubleshooting, reliability]
links: [eng.kafka.producer.send, eng.kafka.leader-routing, eng.kafka.replication.isr-election, eng.kafka.lag-scaling, eng.kafka.rebalance]
---

## 先画出检查顺序

```text
Producer 是否得到成功结果
  → Record 进入哪个 Partition
  → Partition Leader / Replica 是否健康
  → Consumer Group 是否拥有该 Partition
  → Consumer 是否正常 poll
  → 业务处理是否完成
  → Offset 是否提交
```

每一步都要先解释它在问什么：**Leader** 是当前主要读写的副本，**Replica** 是数据副本，**Consumer Group** 是共享分工和书签的读者集合，**Offset Commit** 是保存恢复位置。

## 按症状进入链路

- Producer 超时：看 Metadata 是否过期、Leader 是否切换、ACK 是否返回、重试是否造成重复。
- Consumer Lag 增长：看是否分到 Partition、`poll` 是否按时、业务处理是否慢、提交是否失败。
- 重复处理：看崩溃发生在业务处理前后、是否处于 Rebalance 交接窗口、业务是否幂等。
- 数据“消失”：先确认是 Group 书签前进、Retention 清理、Compaction 清理旧 Key，还是根本没有写入成功。

排查的目标不是找到一个万能参数，而是确定消息在哪个边界停住，再只调整对应组件。
