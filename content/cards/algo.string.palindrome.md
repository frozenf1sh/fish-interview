---
id: algo.string.palindrome
kind: concept
title: 区间 DP：最长回文子串
summary: 用 dp[l][r] 表示闭区间是否为回文，按区间长度从短到长填表，让内部区间先于外层区间完成。
parents: [algo.dp.interval]
tags: [dp, interval, palindrome, string]
links: [algo.dp.interval]
trace: palindrome-interval-dp
---

## 例题

[LeetCode 5 · 最长回文子串 ↗](https://leetcode.cn/problems/longest-palindromic-substring/)。

```go
dp := make([][]bool, n)
for i := range dp {
	dp[i] = make([]bool, n)
	dp[i][i] = true
}
for length := 2; length <= n; length++ {
	for left := 0; left+length <= n; left++ {
		right := left + length - 1
		dp[left][right] = s[left] == s[right] && (length <= 2 || dp[left+1][right-1])
	}
}
```

外层按 `length` 递增，当前格只读取左下方的内部区间；时间复杂度 `O(n²)`，空间复杂度 `O(n²)`。这里的动画与区间 DP 一致：灰色是未使用的下三角，蓝色是依赖，橙色是当前写入格，绿色是已完成区间。
