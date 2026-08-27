---
id: algo.dp.modeling
kind: concept
title: DP 建模：状态、转移与计算顺序
summary: DP 建模先把“一个状态到底代表什么”说清楚，再由最后一步推导转移和计算顺序。
parents: [algo.dp]
tags: [dp, modeling]
links: [algo.dp.linear, algo.dp.lcs, algo.dp.interval, algo.dp.stock, algo.dp.bitmask, algo.dp.path]
exam_signals:
  - company: meituan
    year: 2027
    role: backend
    confidence: medium
    source: https://www.nowcoder.com/
---

> **先写状态句子**：`dp[i]` 用一句完整的话描述一个子问题的答案。转移和循环方向都从这句话推出。

## 用爬楼梯建立第一张 DP 表

每次可走 1 或 2 级，求到第 `n` 级的方案数。令 `dp[i]` 表示“到达第 i 级的方案数”。到达 `i` 的最后一步来自 `i-1` 或 `i-2`，于是：

$$
dp_{i} = dp_{i-1} + dp_{i-2}
$$

边界是 `dp[0]=1`、`dp[1]=1`。`dp[0]` 表示什么都不走这一种起始方案，它让 `dp[2]` 的转移自然成立。

## 每次建模按这个顺序

1. 状态：`dp[...]` 表示哪一个子问题的答案？
2. 最后一步：到达该状态前，发生了哪些互斥选择？
3. 转移：每个选择依赖哪些已经求出的状态？
4. 边界：最小规模问题的答案是多少？
5. 顺序：依赖谁，就先计算谁。

[[algo.dp.linear]] 从固定前序状态开始；LCS、区间、股票、状态压缩和网格寻路分别改变了状态坐标或状态含义，但都遵循这五步。
