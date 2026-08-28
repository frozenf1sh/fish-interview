---
id: algo.sliding-window.exact
kind: algorithm-pattern
title: 滑动窗口：恰好满足
summary: 右端扩张后，条件超限就逐步收缩；当窗口恰好满足条件时记录答案或完成一次计数。
parents: [algo.patterns.sliding-window]
tags: [sliding-window, exact, hashmap]
links: [algo.string.window]
trace: sliding-window-exact
---

## 例题

[LeetCode 992 · K 个不同整数的子数组 ↗](https://leetcode.cn/problems/subarrays-with-k-different-integers/)。动画用字符串 `abcad` 展示“恰好 3 种不同字符”的连续区间。

```go
for right := range s {
	count[s[right]]++
	for distinct > k {
		count[s[left]]--
		if count[s[left]] == 0 { distinct-- }
		left++
	}
	if distinct == k { record(left, right) }
}
```

`distinct == k` 是记录条件，不是允许继续扩张的上限；每次收缩都只移除一个左端字符。
