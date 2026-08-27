---
id: algo.dp.interval
kind: algorithm-pattern
title: 区间 DP：按长度枚举
summary: 用左右端点表示一个连续区间，并按区间长度从短到长计算，保证子区间先于父区间完成。
parents: [algo.patterns.dp]
tags: [dp, interval]
links: [algo.dp.modeling]
trace: interval-dp
---

## 例题

[LeetCode 1000 · 合并石头的最低成本 ↗](https://leetcode.cn/problems/minimum-cost-to-merge-stones/)

先学习每次合并两堆的简化版本：石子 `[3,5,1,2]` 只能合并相邻两堆，求合并成一堆的最小代价。

## 状态转移

$$
区间 [l,r] 的最小代价 = 枚举切分点后，两段代价之和的最小值 + 当前整段重量
$$

> **适用条件**：一次决策把连续区间拆成更短区间，或答案由区间两端共同决定。

## 以合并石子为例

令 `dp[l][r]` 为合并闭区间 `[l,r]` 的最小代价。

最后一次合并会把 `[l,r]` 切在某个 `k`：左边 `[l,k]` 和右边 `[k+1,r]` 都已合并完，再支付整段重量。子区间长度小于当前长度，所以外层必须枚举 `length=2..n`；蓝色格子是当前最优切分读取的两个子区间。

## 分段实现

### 1. 预处理区间和并建立空表

```go
prefix := make([]int, n+1) // prefix 让任意区间和在 O(1) 得到
for i, value := range stones {
	prefix[i+1] = prefix[i] + value
}
dp := make([][]int, n) // dp[l][r]：把闭区间 [l,r] 合成一堆的最小代价
for i := range dp { dp[i] = make([]int, n) }
```

### 2. 由短到长枚举区间

```go
for length := 2; length <= n; length++ { // 先算短区间，供长区间引用
	for l := 0; l+length <= n; l++ {
		r, total := l+length-1, prefix[l+length]-prefix[l]
		dp[l][r] = math.MaxInt
		for k := l; k < r; k++ { // k 是最后一次合并前的切分点
			dp[l][r] = min(dp[l][r], dp[l][k]+dp[k+1][r]+total)
		}
	}
}
return dp[0][n-1] // 完整区间的答案
```

## 容易写错的位置

- 右端点应为 `r=l+length-1`。
- 切分点 `k` 的范围是 `[l,r-1]`，两侧区间都不能为空。
