---
id: algo.dp.bitmask
kind: algorithm-pattern
title: 状态压缩 DP：访问集合与最后位置
summary: 用二进制掩码表示已经处理的元素集合，再配合最后位置记录下一步受什么限制。
parents: [algo.patterns.dp.state]
tags: [dp, bitmask, subset]
links: [algo.dp.modeling, algo.dp.bit-operations]
trace: bitmask-dp
---

## 例题

[LeetCode 847 · 访问所有节点的最短路径 ↗](https://leetcode.cn/problems/shortest-path-visiting-all-nodes/)

这类题都要记住“已经选了哪些元素”和“最后停在哪里”。动画用 4 个城市 A、B、C、D：从 A 出发，每次访问一个未访问城市，目标是覆盖全集合并取得最小代价。

动画前半段固定显示 `mask` 的 4 位和城市状态，后半段固定显示 `dp[mask][last]` 的已写状态；每一次都按“读旧状态 → 检查 next → 算候选 → 写新状态”展开。

## 状态转移

$$
新集合、最后位置的最小代价 = 所有可到达旧状态接上 next 后的较小值
$$

> **适用条件**：元素数量较小，决策只与“已经选了哪些元素”和少量附加信息有关。常见规模是 `n≤20` 左右。

## 从旅行商问题建立状态

从城市 `0` 出发，访问所有城市且每个城市只访问一次。令 `dp[mask][last]` 为已访问集合为 `mask`、最后停在 `last` 的最短路程。

`mask` 的第 `i` 位为 1 表示城市 `i` 已访问。若去往未访问城市 `next`，新集合为 `mask | (1<<next)`；路径代价增加 `cost[last][next]`。动画展示一个集合如何逐位扩张。

位测试、置位、枚举子集与低位提取见 [[algo.dp.bit-operations]]；它们决定 `mask` 的写法，但不改变“集合 + 最后位置”的状态含义。

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
