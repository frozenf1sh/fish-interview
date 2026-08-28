---
id: algo.dfs.path-enumeration
kind: concept
title: DFS：路径枚举与回溯边界
summary: 当题目要求所有路径或存在性时，维护当前 path；进入节点加入，离开节点恢复，避免兄弟分支共享状态。
parents: [algo.patterns.dfs]
tags: [dfs, graph, path]
links: [algo.backtracking.choose-skip]
trace: flow-dfs-path
---

## 例题

[LeetCode 797 · 所有可能的路径 ↗](https://leetcode.cn/problems/all-paths-from-source-to-target/)。图是 DAG，枚举 0 到 n-1 的所有路径。

动画使用固定 DAG，从 0 沿一条边一步步进入；当前 path、递归栈和回退后的兄弟分支同时保留，能看到路径如何被枚举与撤销。

```go
path := []int{0}
var dfs func(int)
dfs = func(v int) {
	if v == n-1 {
		ans = append(ans, append([]int(nil), path...)) // 必须复制当前切片
		return
	}
	for _, next := range graph[v] {
		path = append(path, next) // 选择
		dfs(next)
		path = path[:len(path)-1] // 撤销，恢复给下一个兄弟
	}
}
dfs(0)
```

一般图要增加 `onPath` 防环；只问能否到达时找到答案就立刻返回 bool，不需要收集全部路径。
