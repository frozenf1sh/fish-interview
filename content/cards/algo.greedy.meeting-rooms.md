---
id: algo.greedy.meeting-rooms
kind: algorithm-pattern
title: 区间：最少会议室的双排序指针
summary: 将所有开始与结束时间分别排序；每个开始若早于最早结束，就增加房间，否则复用已结束的房间。
parents: [algo.patterns.greedy]
tags: [greedy, interval, two-pointers]
links: [algo.greedy.interval-start-merge, algo.greedy.interval-scheduling]
trace: meeting-rooms
---

## 例题

[LeetCode 253 · 会议室 II ↗](https://leetcode.cn/problems/meeting-rooms-ii/)。`[[0,30],[5,10],[15,20]]` 至少需要 `2` 个会议室。

动画把会议画在同一条时间轴上，同时保留 starts、ends 两条排序轨道；扫描开始时间时，只移动对应指针和房间数。

## 变量的含义

把会议拆成 `starts` 与 `ends`。扫描每个开始时间 `start` 时，`ends[end]` 永远是尚未释放会议中最早的结束时间：

- `start < ends[end]`：新会议开始得更早，无法复用，需要一间新房。
- `start >= ends[end]`：至少有一间房已释放，`end++`，当前会议复用它。

这不是维护具体哪个会议室，而是统计同时进行会议的峰值。

## 实现

```go
starts, ends := make([]int, len(intervals)), make([]int, len(intervals))
for i, interval := range intervals {
	starts[i], ends[i] = interval[0], interval[1] // 拆开后分别排序
}
sort.Ints(starts)
sort.Ints(ends)

rooms, end := 0, 0
for _, start := range starts {
	if start < ends[end] {
		rooms++ // 当前最早结束的会议尚未结束，只能增加并发容量
	} else {
		end++ // 复用一间已结束的房；rooms 保持峰值
	}
}
return rooms
```

若需要返回具体房间编号，改用按结束时间排序的小根堆；房间数版本只需要两个有序数组。
