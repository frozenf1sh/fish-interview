---
id: algo.dp.lcs
kind: algorithm-pattern
title: 线性 DP：最长公共子序列 LCS
summary: 用两个前缀长度定位子问题；字符相等时同时前进，不等时保留删去一侧字符后的更优答案。
parents: [algo.patterns.dp]
tags: [dp, sequence, lcs]
links: [algo.dp.modeling, algo.dp.linear]
trace: lcs-dp
---

## 状态转移

$$
dp_{i,j} =
\begin{cases}
dp_{i-1,j-1}+1, & a_{i-1}=b_{j-1} \\
\max(dp_{i-1,j},dp_{i,j-1}), & a_{i-1}\neq b_{j-1}
\end{cases}
$$

> **适用条件**：两个序列按前缀推进，目标允许跳过元素且要求保留原相对顺序。

## 从两个前缀开始

令 `dp[i][j]` 为 `a[:i]` 与 `b[:j]` 的 LCS 长度。例子 `a="abcde"`、`b="ace"`：当比较到相同字符 `c`，它可以接在两个更短前缀的公共子序列后面；字符不同则只能丢弃其中一侧末尾字符。

蓝色格子标记本次读取的旧状态：字符相等时取左上角；不等时比较左侧和上方。

## 分段实现

### 1. 分配带空前缀的表

```go
dp := make([][]int, len(a)+1) // dp[i][j]：a 前 i 个与 b 前 j 个的 LCS 长度
for i := range dp {
	dp[i] = make([]int, len(b)+1) // 第 0 行和第 0 列为 0：空序列没有公共字符
}
```

### 2. 按前缀长度填写转移

```go
for i := 1; i <= len(a); i++ {
	for j := 1; j <= len(b); j++ {
		if a[i-1] == b[j-1] { // 两个前缀的末尾字符可以同时保留
			dp[i][j] = dp[i-1][j-1] + 1
		} else {
			dp[i][j] = max(dp[i-1][j], dp[i][j-1]) // 删除一侧末尾字符，保留较优结果
		}
	}
}
return dp[len(a)][len(b)] // 完整前缀对应目标答案
```

## 相邻题型

- 需要连续子串时，状态不能继承左或上，转移会改变。
- 需要输出具体序列时，从 `dp[m][n]` 按转移来源回溯。

[[algo.dp.modeling]] 解释了状态句子、边界和计算顺序如何一起确定。
