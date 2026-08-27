---
id: algo.dp.lcs-space-optimization
kind: concept
title: LCS：一维数组空间优化
summary: 逐行更新 LCS 表时，保存左上角旧值即可把 O(mn) 空间压到 O(n)。
parents: [algo.dp.lcs]
tags: [dp, lcs, space-optimization]
links: [algo.dp.lcs, algo.dp.space-optimization]
trace: flow-lcs-space
---

## 例题

[LeetCode 1143 · 最长公共子序列 ↗](https://leetcode.cn/problems/longest-common-subsequence/)

仍以 `a="abcde"`、`b="ace"` 为例。若只要求 LCS 长度，不需要保留完整二维表来回溯答案。

## 只保留当前行

计算 `dp[i][j]` 时，二维转移读取三个值：上方 `dp[i-1][j]`、左方 `dp[i][j-1]`、左上角 `dp[i-1][j-1]`。

- 一维数组里的 `dp[j]` 在更新前是上方，更新后是当前格。
- `dp[j-1]` 已更新，是左方。
- 用 `diagonal` 在覆盖 `dp[j]` 前保存它，作为左上角。

因此内层循环必须从左到右；反向循环会使 `dp[j-1]` 不再是当前行的左方。

## 分段实现

### 1. 选择较短序列作为列

```go
if len(a) < len(b) { a, b = b, a } // 让一维数组长度为较短序列的长度
dp := make([]int, len(b)+1)
```

### 2. 在覆盖前保存左上角

```go
for i := 1; i <= len(a); i++ {
	diagonal := 0 // 本行开始时，对应 dp[i-1][0]
	for j := 1; j <= len(b); j++ {
		up := dp[j] // 覆盖前是 dp[i-1][j]
		if a[i-1] == b[j-1] {
			dp[j] = diagonal + 1
		} else {
			dp[j] = max(dp[j], dp[j-1]) // dp[j] 是上方，dp[j-1] 是左方
		}
		diagonal = up // 给下一列保留旧的左上角
	}
}
return dp[len(b)]
```

## 边界

空间从 `O(mn)` 降到 `O(min(m,n))`，时间仍是 `O(mn)`。一维数组只保存长度；若题目要求恢复具体公共子序列，通常仍要保留完整表或增加额外回溯结构。
