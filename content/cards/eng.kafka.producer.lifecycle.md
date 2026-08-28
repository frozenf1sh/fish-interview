---
id: eng.kafka.producer.lifecycle
kind: engineering
title: Kafka Producer 生命周期：复用、Flush 与 Close
summary: Producer 不只是一个发送函数，而是持有连接、路由表、缓冲区和后台线程的客户端；应长期复用，并在退出时完成 Flush/Close。
parents: [eng.kafka.producer]
tags: [kafka, producer, flush, lifecycle]
links: [eng.kafka.producer.send, eng.kafka.producer.accumulator, eng.kafka.delivery-semantics]
---

## 先看为什么不能每条消息创建一次

Producer 第一次使用时要建立连接、获取 Metadata（集群路由表）、启动 Sender Thread（后台发送线程）并积累 Batch。合理生命周期是：

```text
应用启动 → 创建一个长期复用的 Producer
业务运行 → send、批量发送、处理回调
应用退出 → Flush / Close → 释放资源
```

如果每条消息都执行 `new → send → close`，刚建立的连接、路由表和批量机会马上被丢掉，吞吐和资源使用都会变差。

## `flush()` 是什么

`flush()` 可以理解为“把此前交给 Producer 的待发送记录尽快推完，并等待结果”。它适合批处理阶段结束、测试或明确的边界，不适合每发一条就调用：

```text
send(m1); flush(); send(m2); flush()
```

这样会把许多小消息强行拆成小批次，削弱 Kafka 的异步和 Batch 优势。

## `close()` 是什么

进程直接退出时，Accumulator 中可能还有未发出的 Record。正常 `close()` 会让客户端处理剩余发送、完成回调并释放连接；如果关闭时已经发生不可恢复失败，应用仍要记录并处理失败结果，不能把 Close 当成成功证明。
