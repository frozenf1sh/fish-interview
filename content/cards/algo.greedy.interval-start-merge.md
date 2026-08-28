---
id: algo.greedy.interval-start-merge
kind: algorithm-pattern
title: 区间：按开始时间排序与合并
summary: 需要合并、覆盖或判断交集时按开始时间升序扫描；当前合并段只维护最远右端，遇到断开才输出新段。
parents: [algo.patterns.greedy]
tags: [greedy, interval, sorting]
links: [algo.greedy.interval-scheduling, algo.greedy.interval-endpoints]
trace: interval-start-merge
---

## 例题

[LeetCode 56 · 合并区间 ↗](https://leetcode.cn/problems/merge-intervals/)。`[[1,3],[2,6],[8,10],[15,18]]` 合并为 `[[1,6],[8,10],[15,18]]`。

## 先确定排序键

目标是把重叠区间压成不重叠的段。按开始时间升序后，未来区间的开始不会回到当前段左侧；因此只要维护最后一段的右端 `last[1]`：

- `current.start <= last.end`：有交集，右端扩到二者较大值。
- `current.start > last.end`：已经断开，当前段固定，开启新段。

按结束时间排序适合“选择最多互不重叠区间”；这里要保留并扩张覆盖范围，排序键应换成开始时间。

## 分段实现

```go
sort.Slice(intervals, func(i, j int) bool {
	return intervals[i][0] < intervals[j][0] // 让后续区间只从右侧进入
})
merged := make([][]int, 0, len(intervals))
for _, current := range intervals {
	if len(merged) == 0 || current[0] > merged[len(merged)-1][1] {
		merged = append(merged, []int{current[0], current[1]}) // 与最后一段断开，复制新段
		continue
	}
	last := merged[len(merged)-1]
	last[1] = max(last[1], current[1]) // 只扩右端，左端已经最小
}
return merged
```

## 边界

`current.start == last.end` 是否合并取决于题目区间定义。闭区间通常合并；半开区间 `[a,b)` 中，`b` 与下段的 `b` 不重叠时应使用严格比较。
