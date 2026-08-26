---
id: algo.dp.linear
kind: algorithm-pattern
title: 线性 DP：从前缀状态推导当前位置
summary: 状态沿数组、字符串或时间轴单向推进，当前位置只依赖已完成的前缀状态。
parents: [algo.dp]
tags: [dp, linear]
links: [algo.dp.modeling]
---

## 识别信号

对象天然按下标或时间顺序推进；“截至 i”、“以 i 结尾”、“前 i 个元素”是常见状态语言。

## 建模与正确性直觉

先区分 `dp[i]` 是“前 i 个元素”的答案，还是“必须以 i 结束”的答案。二者常常有不同的转移与最终答案位置。

## Go 模板

```go
dp := make([]int, n+1)
dp[0] = base
for i := 1; i <= n; i++ {
	dp[i] = transition(dp, i)
}
```

## 常见误区

不要把 `dp[0]` 当成无意义的占位符；它通常承担空前缀或起点的边界定义。

## 变体与关联

[[algo.dp.modeling]]；后续可扩展到背包、区间、树形与 DAG DP。

