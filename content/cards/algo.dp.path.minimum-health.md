---
id: algo.dp.path.minimum-health
kind: algorithm-pattern
title: 逆推寻路 DP：最低初始生命值
summary: 当路径要求全过程资源始终为正时，反推“进入当前格至少需要多少资源”，让未来约束成为当前状态的依赖。
parents: [algo.dp.path]
tags: [dp, grid, reverse, path]
links: [algo.dp.path, algo.dp.modeling]
trace: reverse-path-dp
---

## 例题

[LeetCode 174 · 地下城游戏 ↗](https://leetcode.cn/problems/dungeon-game/)

格子会增加或减少生命值；从左上走到右下，任意时刻生命值都必须大于 `0`。对于 `[[-2,-3,3],[-5,-10,1],[10,30,-5]]`，最低初始生命值为 `7`。

## 状态转移

$$
进入当前格所需最低生命值 = max(1, 右方与下方所需生命值的较小值 - 当前格效果)
$$

> **识别信号**：题目不是询问“走到这里已经积累多少”，而是询问“从这里开始至少要带多少，才能保证后续一直合法”。未来路径决定当前最低门槛，顺序需要从终点反推。

## 为什么正向累计不够

负数格并不表示路径必然更差：若后面有大量补给，它可以被抵消；问题在于每一步都不能让生命值降到 `0`。因此状态定义为 `need[r][c]`：进入 `(r,c)` 前至少拥有多少生命值，才能安全到终点。

终点之后只需生命值 `1`。对当前格，先选择右方、下方中要求较小的后继门槛，再减去当前格效果，并用 `1` 兜底。动画从右下角开始，蓝色格子是已经确定的后继门槛。

## 分段实现

### 1. 用边界哨兵统一终点

```go
m, n := len(dungeon), len(dungeon[0])
const inf = int(1e9)
need := make([][]int, m+1)
for r := range need {
	need[r] = make([]int, n+1)
	for c := range need[r] { need[r][c] = inf } // 不存在的后继不能被选中
}
need[m][n-1], need[m-1][n] = 1, 1 // 终点右侧与下侧都是“离开后仍需 1”
```

### 2. 从右下向左上反推

```go
for r := m - 1; r >= 0; r-- {
	for c := n - 1; c >= 0; c-- {
		nextNeed := min(need[r+1][c], need[r][c+1]) // 先选后续门槛更低的方向
		need[r][c] = max(1, nextNeed-dungeon[r][c]) // 当前格扣血时要预先多带
	}
}
return need[0][0]
```

## 边界

反向 DP 适合“从当前位置满足未来约束的最低资源、最小代价、可行性”。若状态只依赖过去累计值，例如普通最短路径和，按前驱正向推进更直接；见 [[algo.dp.path]]。
