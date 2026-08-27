---
id: algo.greedy.lexicographic
kind: concept
title: 贪心：字典序最小与单调栈
summary: 当每个元素只能保留一次且希望字典序最小，遇到更小元素时弹出仍可在后面补回的栈顶。
parents: [algo.patterns.greedy]
tags: [greedy, monotonic-stack, string]
links: []
trace: flow-greedy-lexicographic
---

## 例题

[LeetCode 316 · 去除重复字母 ↗](https://leetcode.cn/problems/remove-duplicate-letters/)。保留每个字符一次，结果字典序最小。

## 弹栈的两个条件

当前字符 `ch` 比栈顶小，且栈顶在后面还会出现，才可以弹出栈顶。第二个条件保证不会丢失这个字符。

```go
last := make([]int, 26)
for i, ch := range s { last[ch-'a'] = i }
inStack := make([]bool, 26)
stack := make([]byte, 0, len(s))
for i := range s {
	ch := s[i]
	if inStack[ch-'a'] { continue }
	for len(stack) > 0 && ch < stack[len(stack)-1] && last[stack[len(stack)-1]-'a'] > i {
		inStack[stack[len(stack)-1]-'a'] = false
		stack = stack[:len(stack)-1]
	}
	stack, inStack[ch-'a'] = append(stack, ch), true
}
return string(stack)
```

若栈顶不会再出现，就算当前字符更小也不能弹；这是这类题的贪心安全条件。
