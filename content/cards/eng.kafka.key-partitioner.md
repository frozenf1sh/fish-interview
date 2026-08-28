---
id: eng.kafka.key-partitioner
kind: engineering
title: Kafka Key 与 Partitioner：先决定放进哪一册
summary: Producer 用 Partitioner（分区器）把 Record 选择到某个 Partition；稳定的业务 Key 可以让同一实体的事件通常进入同一册，从而获得实体内顺序。
parents: [eng.kafka.routing]
tags: [kafka, key, partitioner, ordering]
links: [eng.kafka.partition, eng.kafka.record, eng.kafka.ordering, eng.kafka.partition-expansion]
---

## 先看订单为什么要有序

订单 `123` 先创建、再支付、最后发货。业务真正关心的是这个订单自己的顺序，不是所有订单混在一起的全局顺序。

```text
key = orderId = 123
OrderCreated → OrderPaid → OrderShipped
                         ↓
              通常进入同一个 Partition
```

这里的 **Key** 是 Record 上可选的关联字段；**Partitioner** 是 Producer Client 中负责“选择 Partition”的代码。Producer 根据 `key` 做确定性映射，就能让相同 Key 通常落到同一个 Partition。

## 选择过程

概念上可以这样看：

```text
Record
  → 有显式 Partition？直接使用
  → 否则有 Key？按 Partitioner 的 Key 规则映射
  → 仍没有 Key？按客户端的无 Key 策略选择
```

`hash(key) % partitionCount` 只是帮助理解的简化公式。实际结果由客户端版本和 Partitioner 实现决定；无 Key 也不意味着获得 Topic 全局顺序。

## Key 不是“一致性哈希”

**一致性哈希** 是一种把 Key 放到哈希环、尽量减少节点变化时迁移的算法；Kafka Producer 默认的 Key 映射不应直接叫它。一旦 Topic 从 3 个 Partition 扩成 6 个，旧 Key 的映射可能改变：

```text
扩容前：order-123 → P1
扩容后：order-123 → P4
```

P1 和 P4 各自仍有序，但跨两个 Partition 没有全局顺序。扩容前必须确认业务是否允许这种迁移，见 [[eng.kafka.partition-expansion]]。

## 读卡时的判断题

先问“谁需要有序”：订单用 `orderId`，用户事件用 `userId`，设备事件用 `deviceId`。Key 选择的是业务的顺序边界，不是越细越好；热点 Key 也可能让某个 Partition 变成瓶颈。
