---
id: eng.kafka.offset-retention
kind: engineering
title: Kafka Offset 与 Retention：书签会落后，但书页会被清理
summary: Offset 是 Partition 内的位置；Retention 会删除过旧日志并推进可读起点，所以 Kafka 支持的是保留范围内的 Replay，不是永久历史。
parents: [eng.kafka.model]
tags: [kafka, offset, retention, replay]
links: [eng.kafka.partition, eng.kafka.consumer.offset-commit, eng.kafka.log-compaction, eng.kafka.event-streaming]
---

## 先用“书签”理解 Offset

Partition 像一本不断往后写的书，Consumer Group 像一个读者，**committed offset** 是这个读者写下的书签：

```text
书页：     0    1    2    3    4    5
记录：     A    B    C    D    E    F
书签：                         ↑ 已提交到 3
```

Offset 只在当前 Partition 内编号。一条旧记录被清理后，后面的记录不会重新编号。

## Retention 会做什么

**Retention（保留策略）** 是 Kafka 按时间、日志大小等条件清理旧日志的规则。清理前 0、1、2 后，仍然存在的范围可能是：

```text
logStartOffset = 3
可读记录：3、4、5、6……
```

`3` 仍然叫 `3`，不会改名成 `0`。因此 Consumer 如果停得太久，书签可能指向已经被撕掉的页：

```text
已提交位置：1
当前最早位置：3
```

这时客户端要按 Offset Reset 策略处理，常见选择是从当前仍存在的最早位置（`earliest`）或当前末尾位置（`latest`）开始；具体行为要结合客户端版本和配置确认。

## Replay 的真正边界

把 Group 的读取位置移回去，就可以重新读取仍在日志中的数据，这叫 **Replay（重放）**。但 Replay 的前提是目标 Offset 还存在；它不是永久档案能力。

Retention 和 **Log Compaction（日志压缩）** 是两种不同目标：Retention 按时间/空间清理，Compaction 主要保留每个 Key 的较新状态，见 [[eng.kafka.log-compaction]]。位点的保存方式见 [[eng.kafka.consumer.offset-commit]]。
