---
id: algo.string.window
kind: concept
title: 滑动窗口：三类条件总览
summary: 滑动窗口统一由右端扩张、左端收缩和窗口不变量组成；具体题目要先判断是恰好满足、最多满足还是最少覆盖。
parents: [algo.patterns.sliding-window]
tags: [sliding-window, hashmap, interval]
links: [algo.sliding-window.exact, algo.sliding-window.at-most, algo.sliding-window.minimum]
---

## 三类窗口

- [[algo.sliding-window.exact]]：窗口必须恰好满足条件，常在收缩结束后记录。
- [[algo.sliding-window.at-most]]：窗口最多满足某个上限，合法时更新最长答案。
- [[algo.sliding-window.minimum]]：窗口覆盖目标后不断收缩，记录最短答案。

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

三类题只改变合法条件和记录时机；动画统一把字符串位置放在一条区间轴上，扩张和收缩各占一个步骤。若只含小写字母，map 可换成 `[26]int`；Unicode 字符要以 rune 维护下标。
