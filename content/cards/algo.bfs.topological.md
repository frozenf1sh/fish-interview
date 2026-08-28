---
id: algo.bfs.topological
kind: concept
title: BFS：拓扑排序 Kahn 算法
summary: 有向图中反复取入度为 0 的节点；处理完所有节点说明无环，否则剩余节点构成环。
parents: [algo.patterns.bfs]
tags: [bfs, graph, topological-sort]
links: [algo.dfs.path-enumeration]
trace: flow-bfs-topological
---

## 例题

[LeetCode 207 · 课程表 ↗](https://leetcode.cn/problems/course-schedule/)。先修关系能否完成全部课程，本质是判断有向图是否有环。

动画固定课程依赖图和每个节点的入度；每次取入度为 0 的节点、删除它的出边、更新邻居入度各占一步。

## 模板

```go
graph := make([][]int, n)
indegree := make([]int, n)
for _, e := range edges { graph[e[0]] = append(graph[e[0]], e[1]); indegree[e[1]]++ }
queue := make([]int, 0)
for v, deg := range indegree { if deg == 0 { queue = append(queue, v) } }
count := 0
for head := 0; head < len(queue); head++ {
	v := queue[head]; count++
	for _, next := range graph[v] {
		indegree[next]--
		if indegree[next] == 0 { queue = append(queue, next) }
	}
}
return count == n
```

边方向要统一：`a → b` 表示完成 a 后才能做 b。若题目需要按层学习/最少轮数，在每轮开始记录当前队列长度。
