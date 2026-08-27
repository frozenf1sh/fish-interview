---
id: algo.dp.bitmask
kind: algorithm-pattern
title: 状态压缩 DP：访问集合与最后位置
summary: 用二进制掩码表示已经处理的元素集合，再配合最后位置记录下一步受什么限制。
parents: [algo.patterns.dp.state]
tags: [dp, bitmask, subset]
links: [algo.dp.modeling]
trace: bitmask-dp
---

## 例题

[LeetCode 943 · 最短超级串 ↗](https://leetcode.cn/problems/find-the-shortest-superstring/)

这类题都要记住“已经选了哪些元素”和“最后停在哪里”；动画用 4 个城市展示同一套 `集合 + 最后位置` 状态。

## 状态转移

$$
新集合、最后位置的最小代价 = 所有可到达旧状态接上 next 后的较小值
$$

> **适用条件**：元素数量较小，决策只与“已经选了哪些元素”和少量附加信息有关。常见规模是 `n≤20` 左右。

## 从旅行商问题建立状态

从城市 `0` 出发，访问所有城市且每个城市只访问一次。令 `dp[mask][last]` 为已访问集合为 `mask`、最后停在 `last` 的最短路程。

`mask` 的第 `i` 位为 1 表示城市 `i` 已访问。若去往未访问城市 `next`，新集合为 `mask | (1<<next)`；路径代价增加 `cost[last][next]`。动画展示一个集合如何逐位扩张。

## 分段实现

### 1. 为每个集合和最后位置初始化代价

```go
const inf = int(1e9)
dp := make([][]int, 1<<n) // 第一维 mask，第二维 last
for mask := range dp {
	dp[mask] = make([]int, n)
	for last := range dp[mask] { dp[mask][last] = inf } // 未到达的状态保持无穷大
}
dp[1][0] = 0 // 只有城市 0 被访问，且最后位置为 0
```

### 2. 从已到达状态扩展一个未访问城市

```go
for mask := 1; mask < 1<<n; mask++ {
	for last := 0; last < n; last++ {
		if dp[mask][last] == inf { continue } // 只从已经可达的状态扩展
		for next := 0; next < n; next++ {
			if mask&(1<<next) != 0 { continue } // 集合中的城市不能重复访问
			nextMask := mask | (1 << next)
			dp[nextMask][next] = min(dp[nextMask][next], dp[mask][last]+cost[last][next])
		}
	}
}
```

## 复杂度边界

状态数约为 `n·2^n`，每个状态再枚举一个 `next`，时间为 `O(n^2·2^n)`。`n` 大时应寻找贪心、图算法或分治结构。
