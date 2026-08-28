---
id: eng.kafka.log-compaction
kind: engineering
title: Kafka Log Compaction：保留 Key 的最新状态
summary: Log Compaction 按 Key 清理较旧记录，目标是保留每个 Key 的最终状态；它与按时间或大小删除历史的 Retention 不是同一策略。
parents: [eng.kafka.operations]
tags: [kafka, compaction, retention, key]
links: [eng.kafka.offset-retention, eng.kafka.key-partitioner, eng.kafka.event-streaming]
---

## 两种保留目标

普通 Retention 关心记录存在了多久或日志占了多大空间；Compaction 关心同一 Key 的历史更新，后台 Cleaner 可以删除被更新版本覆盖的旧记录。

Compaction 不会把 Offset 重新编号，也不应被理解成“立刻只剩最后一条”。清理是后台过程，记录顺序和删除时机还受配置及删除标记保留规则影响。

因此 Compacted Topic 更适合保存状态快照、配置或实体最新值；若业务需要完整事件历史，就不能只依赖 Compaction。
