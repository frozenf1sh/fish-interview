---
id: algo.dp.space-optimization
kind: concept
title: 线性 DP：空间优化规律
summary: 当新状态只读取固定数量、且已经确定的旧状态时，用滚动变量或滚动数组替代完整 DP 表。
parents: [algo.dp.linear]
tags: [dp, space-optimization, rolling]
links: [algo.dp.linear, algo.dp.lcs, algo.dp.path]
trace: space-rolling
---

## 例题：爬楼梯

[LeetCode 70 · 爬楼梯 ↗](https://leetcode.cn/problems/climbing-stairs/)

完整表里计算 `dp[i]` 只读取 `dp[i-1]` 与 `dp[i-2]`；更早的值不会再被访问。因此保留两个变量即可。

## 先画依赖，再删存储

$$
当前状态 = 已完成依赖的组合；只有最后一次读取之后，旧状态才可以被覆盖
$$

动画中蓝色变量是本次读取的旧状态，橙色变量是新写出的 `current`。只有 `current` 已经计算完，两个滚动变量才能整体前移。

判断能否压缩时，按下面顺序检查：

1. 写出每个新状态读取的旧状态集合。
2. 对每个旧状态标出**最后一次读取**发生在哪一轮、哪一列。
3. 选择不会在最后一次读取之前覆盖它的循环方向。
4. 若覆盖前还要读取旧的左上角或上方，先用临时变量保存。

爬楼梯的依赖窗口宽度为 2，因此只保留两个变量。LCS 的一维写法还要保存旧左上角；网格最小路径和从左到右更新时，数组中的 `dp[c]` 是上方、`dp[c-1]` 是左方。

## 分段实现

### 1. 保存仍会被读取的窗口

```go
previousTwo, previousOne := 1, 1 // 对应 dp[0] 与 dp[1]
```

### 2. 先写新状态，再按依赖失效顺序覆盖

```go
for i := 2; i <= n; i++ {
	current := previousOne + previousTwo // 新状态仍同时读取两个旧状态
	previousTwo, previousOne = previousOne, current // 按依赖顺序覆盖变量
}
return previousOne
```

## 二维压缩的方向判定

对于 `dp[r][c]` 压成一行，先把一维数组中的含义写出来：更新前 `dp[c]` 通常代表上一行，更新后代表当前行。

- 依赖上方和左方：`dp[c]` 要保留为上方，`dp[c-1]` 已更新为左方，列从左到右。
- 依赖上方和右方：列从右到左，避免右方提前被覆盖。
- 依赖左上角：覆盖 `dp[c]` 前先保存到 `diagonal`。

压缩只是删除不再需要的历史，不能改变原来的依赖图；方向无法保证时保留二维表。
