---
id: algo.string.kmp
kind: concept
title: 字符串：KMP 前缀函数
summary: next/pi[j] 记录模式串前缀与后缀的最长相等长度；失配时复用已匹配结构，主串指针不回退。
parents: [algo.patterns.string]
tags: [string, kmp, pattern-matching]
links: []
trace: flow-string-kmp
---

## 例题

[LeetCode 28 · 找出字符串中第一个匹配项的下标 ↗](https://leetcode.cn/problems/find-the-index-of-the-first-occurrence-in-a-string/)。

## 先构造 pi

```go
pi := make([]int, len(pattern))
for i, j := 1, 0; i < len(pattern); i++ {
	for j > 0 && pattern[i] != pattern[j] { j = pi[j-1] }
	if pattern[i] == pattern[j] { j++ }
	pi[i] = j
}
```

`j` 是当前已匹配前缀长度。失配后跳到 `pi[j-1]`，因为它是仍可能延续的最长前后缀。

## 在主串匹配

```go
for i, j := 0, 0; i < len(text); i++ {
	for j > 0 && text[i] != pattern[j] { j = pi[j-1] }
	if text[i] == pattern[j] { j++ }
	if j == len(pattern) { return i - len(pattern) + 1 }
}
return -1
```

pattern 为空按题意处理；KMP 的关键变量是“已匹配长度”，不是已经比较过的下标。
