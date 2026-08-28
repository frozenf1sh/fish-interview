---
id: eng.kafka.consumer.commit-modes
kind: engineering
title: Kafka Commit API：三种提交方式怎么选
summary: 自动提交、同步提交和异步提交的区别在于谁决定时机、是否等待结果以及如何处理失败；它们都不等于业务 ACK。
parents: [eng.kafka.consumer]
tags: [kafka, consumer, offset, commit, reliability]
links: [eng.kafka.consumer.offset-commit, eng.kafka.consumer.processing-idempotence, eng.kafka.rebalance, eng.kafka.group.heartbeat-poll]
---

## 先回答“提交什么”

Consumer 处理 Offset `10` 后，通常提交下一次读取位置 `11`。它是在保存“这个 Group 下次从哪里恢复”，不是在告诉 Kafka“数据库已经成功写入”。

```text
poll → 得到 10 → 业务处理 → commit(11)
                            │
                            └── 只提交位点，不提交外部数据库事务
```

## 三种方式放在一起

| 方式 | 谁决定时机 | 是否等待结果 | 新手要记住的风险 |
| --- | --- | --- | --- |
| 自动提交 | 客户端按周期推进 | 应用通常不直接等待 | `poll` 已返回不等于业务处理已完成 |
| `commitSync` | 应用显式调用 | 等待提交结果 | 会阻塞当前流程，但失败更容易被看见 |
| `commitAsync` | 应用显式发起 | 不必立刻等待 | 结果在回调中返回，失败与提交顺序要处理 |

自动提交适合简单消费循环，但可能让恢复位置越过尚未完成的业务副作用。手动提交也不能消除所有重复：处理成功、提交完成前崩溃，重启仍可能再次处理。

## 可靠处理的基本顺序

对需要可靠处理的批次，通常先完成业务处理，再提交每个 Partition 的下一位置；提交失败时不要把业务结果误判成“肯定没执行”，也不要盲目重复执行不可幂等的副作用。事件 ID、唯一键或幂等 Upsert 等防线见 [[eng.kafka.consumer.processing-idempotence]]。

Consumer 失去 Partition 归属后，旧成员继续提交旧分区也可能失败；提交要和 [[eng.kafka.group.rebalance-lifecycle]]、[[eng.kafka.group.heartbeat-poll]] 一起设计。
