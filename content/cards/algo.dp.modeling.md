---
id: algo.dp.modeling
kind: concept
title: DP 建模：状态、转移与计算顺序
summary: 动态规划先定义“子问题答案是什么”，再确定从哪些更小状态转移，以及它们必须先被计算的顺序。
parents: [algo.dp]
tags: [dp, modeling]
links: [algo.dp.linear]
exam_signals:
  - company: meituan
    year: 2027
    role: backend
    confidence: medium
    source: https://www.nowcoder.com/
---

## 识别信号

问题存在重复子问题，且一个阶段的最优解可以由更小阶段的答案组合得到；题目通常要求最优值、方案数或可达性。

## 建模与正确性直觉

依次回答四件事：`dp[i]` 表示什么、最后一步发生什么、最小状态如何初始化、依赖关系决定何种遍历顺序。状态含义不清时，不要急着写转移。

## 常见误区

把下标解释写错、漏掉不可达初值、在原地压缩时使用了错误遍历方向。

## 变体与关联

[[algo.dp.linear]] 是最常见的第一层；区间、树、DAG 与状态压缩 DP 都沿用相同建模顺序。

