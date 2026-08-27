---
id: algo.greedy.interval-scheduling
kind: algorithm-pattern
title: 区间调度：按结束时间选择
summary: 选择最多互不重叠区间时，优先保留结束最早的可行区间。
parents: [algo.patterns.greedy]
tags: [greedy, interval, scheduling]
links: [algo.greedy, algo.greedy.exchange-argument]
trace: interval-scheduling
exam_signals:
  - company: meituan
    year: 2027
    role: backend
    confidence: medium
    source: https://www.nowcoder.com/
---

## 例题

[LeetCode 435 · 无重叠区间 ↗](https://leetcode.cn/problems/non-overlapping-intervals/)

给定若干区间，移除最少区间使剩余区间互不重叠；等价地，尽量保留更多互不重叠区间。

> **先做判断**：题目是“选最多个不重叠区间”，没有权重；每次选择只会影响后续可用的时间边界。满足这两个条件，再考虑本模板。

## 模板：结束最早优先

按结束时间升序扫描，维护已选择集合的最后结束时间 `end`。当前区间从 `start >= end` 开始，才可加入答案。

**为什么是结束时间**：结束越早，未来能接上的区间只会更多。任取一个最优方案，若它的第一个区间结束得更晚，可直接换成当前结束最早的可行区间，数量不会变少；这就是交换论证。

## Go 实现

```go
sort.Slice(intervals, func(i, j int) bool {
	return intervals[i][1] < intervals[j][1] // 结束早的区间排在前面
})
count, end := 0, math.MinInt // end：最后一个已选区间的结束时间
for _, in := range intervals {
	if in[0] >= end { // 当前区间能接在已选集合之后
		count++        // 记录一次有效选择
		end = in[1]    // 更新后续区间必须满足的边界
	}
}
```

> **不变量**：循环开始前，`count` 是已选择区间数，`end` 是最后一个已选区间的结束时间；所有已选区间两两不重叠。

## 一眼排除的相似题

- **按开始时间排序**不能保证给未来留下最大空间。
- **最少会议室**要求按开始时间扫描并维护结束时间最小堆。
- **带权区间选择**的目标是总收益最大，通常应使用 DP。

## 下一步

先复习 [[algo.greedy.exchange-argument]] 的证明语言，再把这个模板和“会议室 / 带权区间”做对比。
