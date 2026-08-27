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

## 识别信号

目标是“选出尽可能多的互不重叠区间”，每次选择都会压缩后续可用时间，但没有区间权重。

## 建模与正确性直觉

按结束时间升序扫描。结束最早的可行区间给后续留下的空间最大；对任何最优方案中第一个结束更晚的区间，都可交换成它而不减少可选数量。

## Go 模板

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

## 关键动画说明

动画直接内嵌在本卡片下方。每一行区间都已按结束时间排序：绿色表示已经被选择；橙色表示当前候选；红色表示与 `end` 冲突而被跳过。观察 `B=[2,5)` 与 `D=[5,7)`：它们不是“长度更长所以不好”，而是开始时间早于当前边界，无法与既有选择共存。

## 常见误区

- 按开始时间排序不能保证留下最大剩余空间。
- “最少会议室”不是本题，应按开始时间扫描并维护结束时间最小堆。

## 变体与关联

[[algo.greedy.exchange-argument]]；继续观察本页下方的执行帧，把交换论证与实际筛选过程对应起来。
