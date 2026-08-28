---
id: eng.kafka.performance.log-page-cache
kind: engineering
title: Kafka 顺序日志与 Page Cache：为什么磁盘也能跑得快
summary: Partition 主要追加而不是随机修改，Kafka 再利用操作系统的 Page Cache 缓存热点文件页，从而减少随机 I/O；性能来自整条链路而非单一技巧。
parents: [eng.kafka.performance]
tags: [kafka, performance, storage, page-cache]
links: [eng.kafka.partition, eng.kafka.performance.batch-compression, eng.kafka.offset-retention]
---

## 先比较两种写法

数据库频繁改历史位置，可能不断寻找不同磁盘位置；Kafka 的 Partition 更像不断在账本末尾写：

```text
append A → append B → append C → append D
```

**顺序写** 是连续追加，通常比到处寻找位置的随机写更容易获得高吞吐。Kafka 还把日志切成 **Segment（日志段）**，便于管理和按保留策略清理。

## Page Cache 是什么

**Page Cache** 是操作系统为文件内容保留的内存缓存。Kafka 写入日志文件后，最近访问的数据页可能仍在内存里，Consumer 读取热点数据时就不必每次都等待物理磁盘：

```text
日志文件 ↔ 操作系统 Page Cache ↔ Consumer 读取
```

这不等于 Kafka 只把数据放内存：数据仍由日志文件保存，缓存也可能被操作系统回收。

## 不要只背“顺序写很快”

真实吞吐还受 Batch、压缩、网络、复制、Partition 并行和 Consumer 处理速度影响。把这些因素串起来，才是 Kafka 的性能模型，见 [[eng.kafka.performance.batch-compression]]。
