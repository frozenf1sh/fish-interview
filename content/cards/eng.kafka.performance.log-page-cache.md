---
id: eng.kafka.performance.log-page-cache
kind: engineering
title: Kafka 顺序日志与 Page Cache：高吞吐的存储基础
summary: Kafka 把写入组织成追加日志，结合操作系统 Page Cache 降低随机 I/O 和重复数据搬运成本。
parents: [eng.kafka.performance]
tags: [kafka, performance, storage, page-cache]
links: [eng.kafka.partition, eng.kafka.performance.batch-compression, eng.kafka.offset-retention]
---

## 为什么追加适合吞吐

Partition 主要执行 `append → append → append`，而不是在历史位置上频繁随机修改。顺序写和 Segment 化日志更容易利用磁盘及操作系统缓存；读取热点数据时，Page Cache 可能直接提供数据，未必每次都触发物理读盘。

这不是“Kafka 只靠内存”。数据仍由日志文件和保留策略管理，Page Cache 只是让常用数据访问更高效。

## 不能脱离全链路谈快

最终吞吐还取决于 Partition 并行、Batch、压缩、网络、复制和 Consumer 处理能力。单独背“顺序写”或“零拷贝”无法解释 Kafka 的整体性能模型。
