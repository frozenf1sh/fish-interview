---
id: algo.bfs.multi-source
kind: concept
title: BFS：多源扩散与最短距离
summary: 多个起点同时入队，第一次到达某格的层数就是它到最近源的最短距离。
parents: [algo.patterns.bfs]
tags: [bfs, grid, multi-source]
links: [algo.bfs.shortest-path]
---

## 例题

[LeetCode 542 · 01 矩阵 ↗](https://leetcode.cn/problems/01-matrix/)。求每个 1 到最近 0 的距离。

## 为什么从所有 0 出发

逐个 1 做 BFS 会重复扫描。把全部 0 作为第 0 层同时入队，扩散到一个 1 的第一条路径必然来自最近 0。

```go
queue := make([][2]int, 0)
for r := range mat {
	for c := range mat[r] {
		if mat[r][c] == 0 { queue = append(queue, [2]int{r, c}) } else { mat[r][c] = -1 }
	}
}
for head := 0; head < len(queue); head++ {
	r, c := queue[head][0], queue[head][1]
	for _, d := range dirs {
		nr, nc := r+d[0], c+d[1]
		if nr < 0 || nr >= len(mat) || nc < 0 || nc >= len(mat[0]) || mat[nr][nc] != -1 { continue }
		mat[nr][nc] = mat[r][c] + 1
		queue = append(queue, [2]int{nr, nc})
	}
}
```

火焰扩散、腐烂橘子、最近出口都可先找“所有初始源”，再判断是否同层扩散。
