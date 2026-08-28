---
id: eng.kafka.failure-diagnosis
kind: engineering
title: Kafka 故障排查：沿消息生命周期定位
summary: Kafka 问题应沿写入、路由、复制、分配、拉取、处理和提交逐段定位，不能只看一个 Lag 或错误码。
parents: [eng.kafka.operations]
tags: [kafka, troubleshooting, reliability]
links: [eng.kafka.producer.send, eng.kafka.leader-routing, eng.kafka.replication.isr-election, eng.kafka.lag-scaling, eng.kafka.rebalance]
---

## 一条可复用的排查链

```text
是否写入
  → 写入了哪个 Partition
  → Leader 和 ISR 是否健康
  → Group 是否分到该 Partition
  → Consumer 是否成功 Poll
  → 业务是否处理完成
  → Offset 是否提交
```

Producer 超时要联查 Metadata、Leader 和 ACK；消费积压要联查分配、处理耗时和提交；频繁重复要看崩溃位置、Rebalance 和幂等边界。每一步都要区分“数据没有发生”和“数据发生但观察位置没有推进”。
