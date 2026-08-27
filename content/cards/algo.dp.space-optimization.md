---
id: algo.dp.space-optimization
kind: concept
title: 线性 DP：空间优化规律
summary: 当新状态只读取固定数量、且已经确定的旧状态时，用滚动变量或滚动数组替代完整 DP 表。
parents: [algo.dp.linear]
tags: [dp, space-optimization, rolling]
links: [algo.dp.linear, algo.dp.lcs, algo.dp.path]
---

## 例题：爬楼梯

[LeetCode 70 · 爬楼梯 ↗](https://leetcode.cn/problems/climbing-stairs/)

完整表里计算 `dp[i]` 只读取 `dp[i-1]` 与 `dp[i-2]`；更早的值不会再被访问。因此保留两个变量即可。

## 更新规律

$$
下一项 = 前一项 + 前两项；更新后，旧的前一项成为新的前两项
$$

### 1. 先保存两个旧状态

```go
previousTwo, previousOne := 1, 1 // 对应 dp[0] 与 dp[1]
```

### 2. 计算新状态后整体向前滚动

```go
for i := 2; i <= n; i++ {
	current := previousOne + previousTwo // 新状态仍同时读取两个旧状态
	previousTwo, previousOne = previousOne, current // 按依赖顺序覆盖变量
}
return previousOne
```

## 先检查再压缩

1. 当前状态只依赖已经完成的更早状态。
2. 覆盖一个变量后，它不会再被本轮或后续转移读取。
3. 二维表常压缩为一维：循环方向必须保证被读取的格子还没有被当前轮覆盖。

LCS 与网格 DP 的一维压缩需要额外保存左上角旧值；没有先写清依赖，不要先压缩。
