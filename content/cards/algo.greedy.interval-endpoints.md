---
id: algo.greedy.interval-endpoints
kind: concept
title: 贪心：区间端点的选择规则
summary: 最大兼容数量按结束时间早的优先；最少箭数或最少点覆盖则维护所有已覆盖区间的交集右端点。
parents: [algo.patterns.greedy]
tags: [greedy, interval, sorting]
links: [algo.greedy.interval-scheduling]
---

## 两类题先分清目标

- **选最多不重叠区间**：[LeetCode 435](https://leetcode.cn/problems/non-overlapping-intervals/)，按结束时间升序，能接上就选。
- **用最少点命中所有区间**：[LeetCode 452](https://leetcode.cn/problems/minimum-number-of-arrows-to-burst-balloons/)，按开始时间排序并维护当前交集的最右端。

## 最少点覆盖模板

```go
sort.Slice(points, func(i, j int) bool { return points[i][0] < points[j][0] })
arrows, right := 1, points[0][1]
for _, p := range points[1:] {
	if p[0] > right {                    // 与当前交集没有重合
		arrows, right = arrows+1, p[1]
	} else {
		right = min(right, p[1])           // 点必须留在交集内
	}
}
```

端点是否允许相等决定比较符号：闭区间 `[a,b]` 用 `p[0] > right` 才新开；半开区间常用 `>=`。先写清区间语义。
