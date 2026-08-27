---
id: algo.dp.stock
kind: algorithm-pattern
title: 状态 DP：股票买卖
summary: 同一天的最优收益取决于账户状态；先枚举状态含义，再只允许符合交易规则的状态转移。
parents: [algo.patterns.dp.state]
tags: [dp, state, stock]
links: [algo.dp.modeling]
trace: stock-dp
---

## 例题

[LeetCode 122 · 买卖股票的最佳时机 II ↗](https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-ii/)

价格 `[7,1,5,3,6,4]` 可多次买卖且任意时刻最多持有一股，最大收益为 `7`。

## 状态转移

$$
今天持有的最大收益 = 昨天持有，与昨天空仓后今天买入，取较大 \\
今天空仓的最大收益 = 昨天空仓，与昨天持有后今天卖出，取较大
$$

> **适用条件**：位置相同但“是否持有、是否冷冻、剩余次数”等条件不同，后续选择的价值也不同。

## 先区分持有与空仓

令：

- `hold`：当天结束后持有一股时的最大收益。
- `cash`：当天结束后空仓时的最大收益。

到价格 `p` 时，新的 `hold` 来自继续持有或今天买入；新的 `cash` 来自继续空仓或今天卖出。两个新状态都必须读取昨天的值，避免同一天先卖再买污染转移。

## 分段实现

### 1. 初始化第一天的两个状态

```go
hold, cash := -prices[0], 0 // 第一天买入后持仓；不操作时空仓收益为 0
```

### 2. 每天只读取前一天的状态

```go
for _, price := range prices[1:] {
	prevHold, prevCash := hold, cash // 保存昨天状态，保证同一天只发生一次状态转移
	hold = max(prevHold, prevCash-price) // 继续持有，或从空仓买入
	cash = max(prevCash, prevHold+price) // 继续空仓，或卖出昨天持有的股票
}
return cash // 最后持仓没有意义，答案必须处于空仓状态
```

## 规则变化时如何扩展

- 有冷冻期：增加 `cooldown` 状态。
- 最多 `k` 笔交易：状态加入交易次数维度。
- 有手续费：在买入或卖出的边上扣除手续费，保持一处扣除即可。

蓝色列是本次计算读取的前一天 `hold / cash`，橙色列是新状态。
