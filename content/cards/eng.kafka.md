---
id: eng.kafka
kind: engineering
title: Kafka：先看一条消息如何旅行
summary: Kafka 可以先想成一组不会因“被读过”就消失的事件日志；本卡从订单事件出发，串起 Topic、Partition、Broker、Consumer Group 和 Offset。
parents: [engineering]
tags: [kafka, messaging, distributed-systems]
links: [eng.kafka.topic, eng.kafka.partition, eng.kafka.broker, eng.kafka.record, eng.kafka.offset-retention, eng.kafka.key-partitioner, eng.kafka.bootstrap-metadata, eng.kafka.producer.send, eng.kafka.replication.write-ack, eng.kafka.consumer.pull-poll, eng.kafka.consumer-group, eng.kafka.delivery-semantics]
---

## 先看一个订单事件

订单服务创建订单时，发出一条 `OrderCreated` 事件。库存、推荐和风控都想看到它，而且推荐系统以后还可能想把过去 7 天的订单重新计算一遍。

如果消息被第一个消费者取走就删除，其他系统就看不到它，也很难重跑。Kafka 的做法是：先把事件追加到日志里，各个消费者只记录自己的阅读位置。

```text
订单服务 → Kafka 日志：OrderCreated
                         ├── 库存系统：读到第 100 页
                         ├── 推荐系统：读到第 80 页
                         └── 风控系统：读到第 50 页
```

## 把名词翻译成人话

- **Topic**：给一类事件起的逻辑名字，例如 `orders`；它像“订单事件这本书的书名”，不是一个被取走消息的队列。
- **Record**：书中的一条记录，也就是一次具体事件。
- **Partition**：把一本很大的书分成几册；每册内部有顺序，也可以并行读写。
- **Offset**：某一册中的页码，用来表示一条记录的位置。
- **Broker**：保存这些书册的 Kafka Server 节点。
- **Producer**：写入 Record 的客户端，例如订单服务。
- **Consumer**：读取 Record 的客户端，例如库存服务。
- **Consumer Group**：一组协同工作的 Consumer；组内分摊 Partition，组与组之间各自保存阅读进度。
- **Replica**：同一个 Partition 的副本，用来应对 Broker 故障。

![Kafka 数据模型：Cluster、Broker、Topic、Partition 与 Replica](/static/kafka-model.drawio.svg)

## 一条消息经过哪些边界

```text
Producer
  → 选择 Topic 下的某个 Partition
  → Partition 的 Leader Broker 追加记录并分配 Offset
  → 其他 Replica 复制这条记录
  → Consumer Group 从自己的 Offset 读取
  → 业务处理后提交下一次读取位置
```

这里有三个容易混淆的事实：消息物理上只写入一个 Partition；不同 Group 可以分别读取同一条记录；消息何时还能读取由 Retention 或 Compaction 决定，而不是由某个 Consumer 的 ACK 决定。**Retention** 是按时间或空间清理旧日志，**Compaction** 是按 Key 清理被新值覆盖的旧记录，先记住名字，细节见 [[eng.kafka.offset-retention]] 和 [[eng.kafka.log-compaction]]。

## 先按这条路线学习

先看 [[eng.kafka.topic]]、[[eng.kafka.partition]]、[[eng.kafka.record]] 和 [[eng.kafka.broker]]，建立“数据是什么、放在哪里”的模型；再看 [[eng.kafka.key-partitioner]] 和 [[eng.kafka.bootstrap-metadata]]，理解 Producer 如何找到位置；最后沿 [[eng.kafka.producer.send]]、[[eng.kafka.replication.write-ack]]、[[eng.kafka.consumer.pull-poll]] 和 [[eng.kafka.consumer-group]] 走完整生命周期。

## 面试时不要说过头

- Kafka 首先是可追加、可按位置读取的分布式日志，也可以通过 Consumer Group 表现出队列式负载均衡。
- Kafka 默认只保证一个 Partition 内的顺序，不自动保证整个 Topic 的全局顺序。
- Producer 幂等性只解决发送重试的一部分重复问题，不等于跨数据库和外部 RPC 的 Exactly-Once。
