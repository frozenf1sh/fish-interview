---
id: algo.dp.bitmask
kind: algorithm-pattern
title: 状态压缩 DP：访问集合与最后位置
summary: 用二进制掩码表示已经处理的元素集合，再配合最后位置记录下一步受什么限制。
parents: [algo.patterns.dp]
tags: [dp, bitmask, subset]
links: [algo.dp.modeling]
trace: bitmask-dp
---

> **适用条件**：元素数量较小，决策只与“已经选了哪些元素”和少量附加信息有关。常见规模是 `n≤20` 左右。

## 从旅行商问题建立状态

从城市 `0` 出发，访问所有城市且每个城市只访问一次。令 `dp[mask][last]` 为已访问集合为 `mask`、最后停在 `last` 的最短路程。

`mask` 的第 `i` 位为 1 表示城市 `i` 已访问。若去往未访问城市 `next`，新集合为 `mask | (1<<next)`；路径代价增加 `cost[last][next]`。动画展示一个集合如何逐位扩张。

## Go 实现

```go
const inf = int(1e9)
dp := make([][]int, 1<<n) // 第一维 mask，第二维 last
for mask := range dp {
	dp[mask] = make([]int, n)
	for last := range dp[mask] { dp[mask][last] = inf } // 未到达的状态保持无穷大
}
dp[1][0] = 0 // 只有城市 0 被访问，且最后位置为 0
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
