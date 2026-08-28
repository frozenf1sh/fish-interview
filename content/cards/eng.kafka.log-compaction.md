---
id: eng.kafka.log-compaction
kind: engineering
title: Kafka Log Compaction：保留每个 Key 的较新状态
summary: Log Compaction 按 Key 清理被新记录覆盖的旧记录，适合状态快照或配置流；它与按时间/大小清理完整历史的 Retention 不是一回事。
parents: [eng.kafka.operations]
tags: [kafka, compaction, retention, key]
links: [eng.kafka.offset-retention, eng.kafka.key-partitioner, eng.kafka.event-streaming]
---

## 先看同一个 Key 的多次更新

```text
key=user-123：name=小明
key=user-123：name=明明
key=user-123：name=小明明
```

**Log Compaction（日志压缩）** 会在后台寻找同一个 Key 的旧版本，让日志最终更关注每个 Key 的较新状态。它适合配置、用户资料、实体快照这类“恢复最新值”场景。

## 它和普通 Retention 不一样

- **Retention**：按时间、日志大小等条件删除旧日志，目标是控制保留空间和时间。
- **Compaction**：按 Key 清理被覆盖的旧值，目标是保留最新状态。

Compaction 是后台过程，不是写入后立刻只剩最后一条；Offset 也不会重新编号。带有删除语义的 Tombstone（墓碑记录）还需要经过保留窗口，才能让下游知道某个 Key 被删除。

如果业务需要完整的 `Created → Paid → Shipped` 历史，就不能只依赖 Compaction；事件流与状态视图的区别见 [[eng.kafka.event-streaming]]。
