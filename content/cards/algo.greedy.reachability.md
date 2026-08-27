---
id: algo.greedy.reachability
kind: concept
title: 贪心：维护最远可达边界
summary: 扫描数组时维护当前最远可达位置；只要当前位置未越界，就用它更新边界，适合跳跃与分层扩展题。
parents: [algo.patterns.greedy]
tags: [greedy, array, reachability]
links: [algo.greedy.interval-scheduling]
trace: flow-greedy-reachability
---

## 例题

[LeetCode 55 · 跳跃游戏 ↗](https://leetcode.cn/problems/jump-game/)。`nums[i]` 是从 i 最远可跳的距离，问能否到达末尾。

## 先定义边界

`farthest` 表示已扫描位置能够到达的最远下标。扫描到 `i` 前若 `i>farthest`，说明根本到不了 i，后续也无需看；否则更新 `farthest=max(farthest,i+nums[i])`。

```go
farthest := 0
for i, jump := range nums {
	if i > farthest { return false }              // 当前点不可达
	farthest = max(farthest, i+jump)              // 用可达点扩展边界
	if farthest >= len(nums)-1 { return true }    // 已覆盖终点
}
return true
```

求最少跳数时再维护当前层末尾 `end` 与下一层末尾 `farthest`：扫描到 `end` 才把步数加一。不要在每个位置都加步数。
