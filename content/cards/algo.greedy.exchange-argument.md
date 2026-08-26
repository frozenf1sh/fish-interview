---
id: algo.greedy.exchange-argument
kind: concept
title: 贪心：交换论证
summary: 用局部选择替换任意最优解中的一个选择，且不使结果变差，从而证明贪心选择安全。
parents: [algo.greedy]
tags: [greedy, proof]
links: [algo.greedy.interval-scheduling]
---

## 识别信号

你已经有一个看似合理的局部选择，但还不能说明“局部最优一定导向全局最优”时，优先尝试交换论证。

## 建模与正确性直觉

取任意最优解。若它没有采用贪心选择，证明可以把其中冲突的选择替换为贪心选择，且可行性与目标值不变或更好。反复交换后，得到一个包含贪心选择的最优解。

## 常见误区

“当前看起来最好”不是证明。必须指出替换对象、替换后为什么仍可行，以及目标值为何不变差。

## 变体与关联

[[algo.greedy.interval-scheduling]] 是按结束时间选择区间的经典例子。

