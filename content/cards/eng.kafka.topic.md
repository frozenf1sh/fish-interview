---
id: eng.kafka.topic
kind: engineering
title: Kafka Topic：给一类事件命名
summary: Topic 是组织事件的逻辑名字，不是一条被 Consumer 取走就清空的队列；真正保存数据、提供顺序和并行能力的是它下面的 Partition。
parents: [eng.kafka.model]
tags: [kafka, topic, log]
links: [eng.kafka, eng.kafka.partition, eng.kafka.record, eng.kafka.offset-retention, eng.kafka.key-partitioner]
---

## 先用一个例子

订单系统把订单相关事件都写到名为 `orders` 的 Topic：

```text
Topic: orders
├── Partition 0：一册有顺序的日志
├── Partition 1：另一册有顺序的日志
└── Partition 2：另一册有顺序的日志
```

**Topic** 可以理解为“这类事件的名字”或“一个逻辑目录”。一条具体事件叫 **Record**，它最终会被追加到 Topic 下面的某一个 **Partition**，而不是直接存在于 Topic 这个抽象名字里。

## Topic 不是消费队列

在传统队列的直觉里，消费者取走消息后，消息可能从队列中消失。Kafka 不用这个规则：

```text
orders / Partition 0：A、B、C、D
       ├── inventory Group：读到 C
       ├── recommendation Group：读到 A
       └── risk Group：读到 D
```

这里的 **Consumer Group** 是一组共同消费的客户端。三个 Group 都可以从同一份日志读取，只是各自的阅读位置不同；一个 Group 读过 `A`，不会让另一个 Group 看不到 `A`。消息仍不仍存在，由 **Retention**（按时间/空间保留）或 **Compaction**（按 Key 清理旧值）决定，而不是由某个 Consumer 的 ACK 决定。

## 为什么 Topic 还要拆 Partition

如果所有订单都只有一册日志，写入和读取容易被一台机器、一个顺序边界卡住。拆成多册后，不同 Partition 可以分布到不同 Broker，并由不同 Consumer 并行处理。

代价是：Topic 不再天然拥有全局顺序；Partition 越多，副本、连接、分配和 Rebalance 的管理成本也越高。Partition 的细节见 [[eng.kafka.partition]]。

## 读卡时记住这句

Topic 回答“这是什么类型的事件”，Partition 回答“这条事件具体放在哪、和谁有序、能并行多少”。需要把同一订单的事件放在一起时，再看 [[eng.kafka.key-partitioner]]。
