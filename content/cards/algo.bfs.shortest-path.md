---
id: algo.bfs.shortest-path
kind: concept
title: BFS：无权图最短路
summary: 边权相同或都为 1 时，队列按距离层扩展；节点第一次出队或入队时确定最短距离。
parents: [algo.patterns.bfs]
tags: [bfs, graph, shortest-path]
links: [algo.bfs.multi-source]
---

## 例题

[LeetCode 1091 · 二进制矩阵中的最短路径 ↗](https://leetcode.cn/problems/shortest-path-in-binary-matrix/)。每一步代价相同，问从起点到终点最少走几步。

## 使用信号

状态之间每次移动代价一致；目标是最少步数、最少操作数，或最早到达时间。若边权不同，改用 Dijkstra 或 0-1 BFS。

```go
queue := []State{start}
dist := map[State]int{start: 0}
for head := 0; head < len(queue); head++ {
	cur := queue[head]
	if cur == target { return dist[cur] }
	for _, next := range neighbors(cur) {
		if _, seen := dist[next]; seen { continue } // 入队时标记，避免重复入队
		dist[next] = dist[cur] + 1
		queue = append(queue, next)
	}
}
return -1
```

队列的 `head` 下标避免 `queue=queue[1:]` 造成长期引用；状态若是网格坐标，可用二维 `dist` 或 `visited` 代替 map。
