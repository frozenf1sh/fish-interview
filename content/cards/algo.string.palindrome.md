---
id: algo.string.palindrome
kind: concept
title: 字符串：回文中心扩展
summary: 每个回文由一个中心向两侧对称扩展；奇数和偶数长度分别对应单点中心与双点中心。
parents: [algo.patterns.string]
tags: [string, palindrome, two-pointers]
links: []
---

## 例题

[LeetCode 5 · 最长回文子串 ↗](https://leetcode.cn/problems/longest-palindromic-substring/)。

```go
start, length := 0, 0
expand := func(left, right int) {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		if right-left+1 > length { start, length = left, right-left+1 }
		left--; right++
	}
}
for center := range s {
	expand(center, center)     // 奇数长度，例如 aba
	expand(center, center+1)   // 偶数长度，例如 abba
}
return s[start : start+length]
```

中心扩展为 `O(n²)` 时间、`O(1)` 额外空间，常适合笔试。需要统计回文、判断单次区间时可用同一中心思路；需要很多区间答案再考虑 DP 或 Manacher。
