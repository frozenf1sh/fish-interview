---
id: algo.binary-search.answer
kind: algorithm-pattern
title: 二分答案：把最优值变成可判定问题
summary: 对答案范围二分，并构造单调 check(mid) 判断目标值是否可行。
parents: [algo.binary-search]
tags: [binary-search, monotonicity, greedy]
links: [algo.greedy.interval-scheduling]
exam_signals:
  - company: netease
    year: 2027
    role: backend
    confidence: low
    source: https://www.nowcoder.com/
---

## 识别信号

“最大化最小值”或“最小化最大值”尤其常见。关键不是答案是否有序，而是 `check(x)` 的真假是否只发生一次方向变化。

## 建模与正确性直觉

先写出候选答案的上下界，再明确 `x` 变大后可行性如何变化。只在能用贪心、前缀和或其他较低复杂度方法完成 `check` 时使用。

## Go 模板

```go
lo, hi := lowerBound, upperBound
for lo < hi {
	mid := lo + (hi-lo)/2
	if feasible(mid) {
		hi = mid
	} else {
		lo = mid + 1
	}
}
return lo
```

## 常见误区

没有先证明单调性；或把“最小化最大值”的真假方向写反。

## 变体与关联

本质是参数搜索。`check` 往往本身是一个贪心过程。

