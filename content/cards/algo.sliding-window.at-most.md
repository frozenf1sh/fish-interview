---
id: algo.sliding-window.at-most
kind: algorithm-pattern
title: 滑动窗口：满足的最多数
summary: 维护一个始终不超过上限的合法窗口，窗口恢复合法后用当前长度更新最长答案。
parents: [algo.patterns.sliding-window]
tags: [sliding-window, at-most, hashmap]
links: [algo.string.window]
trace: sliding-window-at-most
---

## 例题

[LeetCode 159 · 至多包含两个不同字符的最长子串 ↗](https://leetcode.cn/problems/longest-substring-with-at-most-two-distinct-characters/)。动画使用 `eceba` 和 `k=2`。

```go
for right := range s {
	count[s[right]]++
	for distinct > k {
		count[s[left]]--
		if count[s[left]] == 0 { distinct-- }
		left++
	}
	best = max(best, right-left+1)
}
```

“最多”题的关键是不合法时必须收缩到合法；之后的整个窗口都可以作为当前右端点的候选。
