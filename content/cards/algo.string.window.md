---
id: algo.string.window
kind: concept
title: 字符串：滑动窗口与频次计数
summary: 窗口右端扩张以覆盖条件，左端收缩以恢复条件；用频次表维护窗口状态，避免每次重新扫描子串。
parents: [algo.patterns.string]
tags: [string, sliding-window, hashmap]
links: [algo.string.golang]
---

## 例题

[LeetCode 3 · 无重复字符的最长子串 ↗](https://leetcode.cn/problems/longest-substring-without-repeating-characters/)。

## 通用骨架

```go
count := map[byte]int{}
left, best := 0, 0
for right := 0; right < len(s); right++ {
	count[s[right]]++                       // 扩张，纳入一个新字符
	for count[s[right]] > 1 {               // 条件被破坏就收缩
		count[s[left]]--
		left++
	}
	best = max(best, right-left+1)           // 此时窗口重新合法
}
```

“最多 k 个”“至少覆盖 target”“恰好满足”只改变窗口合法条件和何时记录答案。若只含小写字母，map 换 `[26]int`；Unicode 字符要以 rune 维护下标。
