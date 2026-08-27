---
id: algo.dp.path
kind: algorithm-pattern
title: 寻路 DP：网格最小路径和
summary: 当移动方向构成 DAG 时，按拓扑顺序累计到达每个格子的最优代价；网格右下移动时依赖上方和左方。
parents: [algo.patterns.dp]
tags: [dp, grid, path]
links: [algo.dp.modeling]
trace: path-dp
---

> **适用条件**：移动只能沿无环方向进行，且每个位置的最优答案由已完成的前驱位置决定。

## 从左上走到右下

网格 `[[1,3,1],[1,5,1],[4,2,1]]` 只能向右或向下，求最小路径和。令 `dp[r][c]` 为走到格子 `(r,c)` 的最小代价。

到 `(r,c)` 的最后一步来自上方或左方，因此取两者较小值再加当前格子值。首行只能从左进入，首列只能从上进入；动画按行填表，绿色格子是已经完成的前驱。

## Go 实现

```go
m, n := len(grid), len(grid[0])
dp := make([][]int, m) // dp[r][c]：到达 (r,c) 的最小路径和
for r := range dp { dp[r] = make([]int, n) }
dp[0][0] = grid[0][0] // 起点只包含自身代价
for c := 1; c < n; c++ { dp[0][c] = dp[0][c-1] + grid[0][c] } // 首行只能从左来
for r := 1; r < m; r++ {
	dp[r][0] = dp[r-1][0] + grid[r][0] // 首列只能从上来
	for c := 1; c < n; c++ {
		dp[r][c] = min(dp[r-1][c], dp[r][c-1]) + grid[r][c] // 在两个前驱中选更小代价
	}
}
return dp[m-1][n-1] // 右下角是完整路径的最优值
```

## 边界

允许四向移动时图可能有环，普通 DP 顺序失效；边权非负时改用 Dijkstra。
