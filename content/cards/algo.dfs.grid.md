---
id: algo.dfs.grid
kind: concept
title: DFS：网格连通块
summary: 把网格坐标当作图节点；DFS 负责访问一个连通块，外层扫描负责计数或选择起点。
parents: [algo.patterns.dfs]
tags: [dfs, grid, flood-fill]
links: [algo.bfs.multi-source]
trace: flow-dfs-grid
---

## 例题

[LeetCode 200 · 岛屿数量 ↗](https://leetcode.cn/problems/number-of-islands/)。每次遇到未访问陆地，就用 DFS 淹没整个岛，答案加一。

动画固定网格，把递归调用栈画在网格旁；进入陆地、检查相邻方向、标记和返回分别可见，整块连通区域会逐格被淹没。

```go
var dfs func(int, int)
dfs = func(r, c int) {
	if r < 0 || r >= m || c < 0 || c >= n || grid[r][c] != '1' { return }
	grid[r][c] = '0' // 进入节点时标记，防止绕环重复访问
	for _, d := range dirs { dfs(r+d[0], c+d[1]) }
}
count := 0
for r := 0; r < m; r++ {
	for c := 0; c < n; c++ {
		if grid[r][c] == '1' { count++; dfs(r, c) }
	}
}
```

网格很大、递归深度可能接近节点数时改显式栈。需要最短距离时用 BFS，DFS 只保证连通性，不保证先到最短。
