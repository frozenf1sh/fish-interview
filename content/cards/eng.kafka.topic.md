---
id: eng.kafka.topic
kind: engineering
title: Kafka Topic：逻辑事件流与分区集合
summary: Topic 是一类事件的逻辑命名空间；真正存储、排序和扩缩容的基本单元是它包含的 Partition。
parents: [eng.kafka]
tags: [kafka, topic, partition]
links: [eng.kafka, eng.kafka.partition]
---

## 核心机制

Producer 向 Topic 写入记录，但最终每条记录都会落到一个 Partition。分区键相同的记录通常应进入同一分区，以维持该键上的顺序；没有稳定键时，顺序语义必须明确降级为“无全局顺序”。

Topic 的保留策略按时间或大小清理日志，与消费者是否已经读取没有直接的逐条删除关系。消费者的 offset 是自己的读取进度，不是 Topic 的删除开关。

## 常见误区

增加 Partition 通常能增加并行度，却可能改变键到分区的映射。若业务依赖某个键的严格历史顺序，需要先评估扩分区后的迁移影响。

## 变体与关联

顺序和并行的真实边界在 [[eng.kafka.partition]]；逻辑命名与存储策略在 Topic 层。

