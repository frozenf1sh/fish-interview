---
id: eng.kafka.consumer.commit-modes
kind: engineering
title: Kafka Commit API：自动、同步与异步提交
summary: Offset commit 保存的是 Consumer Group 的恢复位置；自动、同步和异步提交的等待与失败处理不同，都不能替代业务副作用的成功确认。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, offset, commit, reliability]
links: [eng.kafka.consumer.offset-commit, eng.kafka.consumer.processing-idempotence, eng.kafka.rebalance, eng.kafka.group.heartbeat-poll]
---

## 提交的对象

Consumer 的 `position` 是客户端下一次准备读取的位置；`committed offset` 是已经写入 Group 位点存储、可用于故障恢复的位置。处理 Offset `10` 后，通常提交 `11`，表示下一次从 `11` 开始，而不是给 Offset `10` 打一个业务 ACK。

```text
poll → position 前进 → 业务处理 → commit(next offset)
                         │
                         └── 失败前未提交：重启后可能重放
```

Kafka 不知道数据库、RPC 或文件写入是否完成。因此提交时机必须围绕业务副作用设计，而不能只围绕“已经调用过 `poll`”设计。

## 三种常见方式

| 方式 | 调用线程的等待 | 失败处理 | 适合用来理解什么 |
| --- | --- | --- | --- |
| 自动提交 | 客户端按周期后台提交 | 业务代码不直接掌控提交时机 | 简化消费循环，但容易和处理完成脱节 |
| `commitSync` | 等待提交结果 | 调用方能明确看到失败，通常会阻塞当前流程 | 把提交结果纳入当前批次的故障路径 |
| `commitAsync` | 发起后不必等待 | 通过回调观察结果；不能把回调当业务成功 | 减少等待，但要处理异步提交顺序和失败告警 |

自动提交尤其要警惕：记录已经被 `poll` 返回，并不等于业务处理成功；如果此时进程崩溃，恢复位点可能已经越过尚未完成的副作用，形成跳过。手动提交也不是“无重复”保证：处理成功后、提交完成前崩溃，仍会重放。

## 推荐的边界

对需要可靠处理的批次，先处理记录，再提交该批次每个 Partition 的下一位置；提交失败时不要把业务处理误判为失败，也不要盲目重做不可幂等的副作用。通常需要用事件 ID 或业务 ID 做去重、唯一约束或幂等 Upsert，见 [[eng.kafka.consumer.processing-idempotence]]。

Consumer 失去 Partition 归属后，旧成员的提交可能失败或不再有意义；提交与交接要结合 [[eng.kafka.rebalance]] 和 [[eng.kafka.group.heartbeat-poll]] 一起设计。
