---
id: algo.dp.linear
kind: algorithm-pattern
title: 斐波那契 DP：只依赖前两个状态
summary: 当当前位置只依赖前面固定少数状态时，先写边界，再按依赖方向依次填表。
parents: [algo.patterns.dp]
tags: [dp, linear]
links: [algo.dp.modeling, algo.dp.space-optimization]
trace: linear-dp
---

## 例题

[LeetCode 70 · 爬楼梯 ↗](https://leetcode.cn/problems/climbing-stairs/)

每次可走 1 或 2 级，求到第 `n` 级的方案数。动画使用 `n=5`。

## 状态转移

$$
第 i 级的方案数 = 前一级方案数 + 前两级方案数；边界：第 0、1 级都为 1
$$

> **适用条件**：状态沿单一方向推进，且 `dp[i]` 只读取固定数量的前序状态。

## 用爬楼梯推导状态

最后一步只有两种来源：从 `i-1` 走 1 级，或从 `i-2` 走 2 级。因此 `dp[i]` 定义为“到达第 i 级的方案数”。

蓝色格子是本次转移读取的两个旧状态；橙色格子是刚写入的新状态。

## 分段实现

### 1. 定义状态并写边界

```go
dp := make([]int, n+1)  // dp[i]：到达第 i 级的方案数
dp[0], dp[1] = 1, 1    // 空路径和走一级各算一种
```

### 2. 依赖前两个状态推进

```go
for i := 2; i <= n; i++ {
	dp[i] = dp[i-1] + dp[i-2] // 最后一步来自 i-1 或 i-2
}
return dp[n] // 目标状态就是答案
```

> **状态检查**：写转移前先用一句完整的话读出 `dp[i]`；如果读不通，状态定义通常还不准确。

## 何时转向其他 DP

- 两个序列的前缀相互比较，进入 [[algo.dp.lcs]]。
- 决策由左右端点限定，进入 [[algo.dp.interval]]。
- 同一位置需要持有、空仓等不同条件，进入 [[algo.dp.stock]]。
- 只读取固定数量旧状态时，可用 [[algo.dp.space-optimization]] 压缩空间。

## 下一步

[[algo.dp.modeling]]；这一类的核心是让计算顺序服从依赖方向。
