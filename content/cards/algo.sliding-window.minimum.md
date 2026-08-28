---
id: algo.sliding-window.minimum
kind: algorithm-pattern
title: 滑动窗口：满足条件的最少数
summary: 右端扩张直到窗口覆盖目标，再逐步移除左端字符并记录仍然合法的最短窗口。
parents: [algo.patterns.sliding-window]
tags: [sliding-window, minimum, hashmap]
links: [algo.string.window]
trace: sliding-window-minimum
---

## 例题

[LeetCode 76 · 最小覆盖子串 ↗](https://leetcode.cn/problems/minimum-window-substring/)。在 `ADOBECODEBANC` 中寻找覆盖 `ABC` 的最短连续区间。

```go
for right := range s {
	add(s[right])
	for coversTarget() {
		best = min(best, window(left, right))
		remove(s[left])
		left++
	}
}
```

覆盖成立后，扩张没有意义，答案只能从连续收缩中变短；动画把“记录候选”和“移除左端”拆成可观察步骤。
