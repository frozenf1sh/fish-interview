---
id: eng.kafka.offset-retention
kind: engineering
title: Kafka Offset 与 Retention：可重放但不是永久保存
summary: Offset 是 Partition 内的位置，Retention 会推进可读日志的起点；旧消息清理后不能再从原位置 Replay。
parents: [eng.kafka.model]
tags: [kafka, offset, retention, replay]
links: [eng.kafka.partition, eng.kafka.consumer.offset-commit, eng.kafka.log-compaction, eng.kafka.event-streaming]
---

## Offset 不会重新编号

假设原日志是 `0、1、2、3、4、5`。Retention 清理前面的 Segment 后，可读范围可能从 `3` 开始，但剩余记录仍然是 `3、4、5`，不会把 `3` 改成 `0`。可以把变化理解为 `logStartOffset` 向前推进。

## Replay 有边界

Consumer Group 的 committed offset 可能还停在 `1`，但如果 Kafka 已经清理到 `3`，再次启动时就不能读取 `1`。此时需要按客户端的 Offset Reset 策略从当前仍存在的最早位置或最新位置开始，具体语义要结合客户端配置确认。

Retention 按时间、大小等策略清理，不关心某个 Consumer 是否已经读取。Replay 的前提是目标数据仍在保留范围内；若使用 Compaction，还要区分“保留事件历史”和“保留每个 Key 的最新状态”。
