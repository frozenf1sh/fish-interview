---
id: eng.kafka.producer.retry-idempotence
kind: engineering
title: Kafka Producer 重试与幂等：ACK 丢了怎么办
summary: Broker 可能已经写入但响应在网络中丢失；Producer 重试会造成重复，幂等 Producer 用身份和序列信息过滤同一发送的重复追加。
parents: [eng.kafka.producer]
tags: [kafka, producer, retry, idempotence]
links: [eng.kafka.producer.sender, eng.kafka.replication.write-ack, eng.kafka.delivery-semantics, eng.kafka.transactions-eos]
---

## 先看一个无法判断的超时

```text
Producer → Broker：写入 m1
Broker：已经追加 m1
Broker → Producer：ACK 在网络中丢失
Producer：只看到 timeout，不知道 m1 到底写没写入
Producer → Broker：重试 m1
```

**Retry（重试）** 是失败或超时后再次发送；**ACK** 是 Broker 返回的确认。Producer 看到超时，无法区分“请求没到”与“写入成功但响应没到”，所以普通重试可能让日志出现两份 `m1`。

## Idempotent Producer 做了哪一层防护

**Idempotence（幂等）** 的意思是：同一个逻辑操作重复到达，结果仍像只做了一次。Kafka Producer 会带上 Producer 的身份（Producer ID）和 Partition 内的序列信息：

```text
首次：PID=7，sequence=12，写入 m1
重试：PID=7，sequence=12，再次到达
Broker：识别为同一发送，避免重复追加
```

它主要解决 Producer → Kafka Log 这一段的重试重复，不是所有业务链路的去重器。数据库写入、调用支付接口、消费 A 再生产 B，仍可能在进程崩溃后重复执行。

## 顺序为什么也会受影响

为了提高吞吐，Producer 可能让多个请求同时在路上；如果前一个请求失败、后一个请求先成功，重试的安排可能影响观察顺序。幂等与客户端的并发/重试约束需要一起考虑，不能只背一个参数就承诺所有情况下有序。

因此要把“发送不重复”与“业务只执行一次”分开回答，后者见 [[eng.kafka.delivery-semantics]] 和 [[eng.kafka.transactions-eos]]。
