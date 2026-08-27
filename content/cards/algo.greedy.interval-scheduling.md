---
id: algo.greedy.interval-scheduling
kind: algorithm-pattern
title: 区间调度：按结束时间选择
summary: 选择最多互不重叠区间时，优先保留结束最早的可行区间。
parents: [algo.greedy]
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

> **先做判断**：题目是“选最多个不重叠区间”，没有权重；每次选择只会影响后续可用的时间边界。满足这两个条件，再考虑本模板。

## 模板：结束最早优先

按结束时间升序扫描，维护已选择集合的最后结束时间 `end`。当前区间从 `start >= end` 开始，才可加入答案。

**为什么是结束时间**：结束越早，未来能接上的区间只会更多。任取一个最优方案，若它的第一个区间结束得更晚，可直接换成当前结束最早的可行区间，数量不会变少；这就是交换论证。

## Go 实现

```go
sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })
count, end := 0, math.MinInt
for _, in := range intervals {
	if in[0] >= end {
		count++
		end = in[1]
	}
}
```

> **不变量**：循环开始前，`count` 是已选择区间数，`end` 是最后一个已选区间的结束时间；所有已选区间两两不重叠。

## 动画应该看什么

动画直接内嵌在本卡片顶部。绿色表示已选择，橙色是当前候选，红色表示冲突。重点看 `B=[2,5)` 与 `D=[5,7)`：它们不是“长度更长所以不好”，而是开始时间早于当前 `end`，无法接在已选区间之后。

## 一眼排除的相似题

- **按开始时间排序**不能保证给未来留下最大空间。
- **最少会议室**不是本题：它要求按开始时间扫描并维护结束时间最小堆。
- **带权区间选择**也不是本题：目标变为总收益最大，通常应使用 DP。

## 下一步

先复习 [[algo.greedy.exchange-argument]] 的证明语言，再把这个模板和“会议室 / 带权区间”做对比。
