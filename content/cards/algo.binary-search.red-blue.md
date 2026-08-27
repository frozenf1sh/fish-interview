---
id: algo.binary-search.red-blue
kind: concept
title: 红蓝染色：用不变量固定二分边界
summary: 给单调判定函数染色，维护红色不可行点与蓝色可行点，就能从区间不变量直接推出循环条件和返回值。
parents: [algo.knowledge]
tags: [binary-search, monotonicity, invariant]
links: [algo.binary-search.answer]
---

> **统一语言**：先指定 `check(x)` 的含义，再指定红蓝颜色；边界、循环和返回值只从这个约定推导。

## 最小化可行答案

当 `check(x)` 随 `x` 增大从假变真时，左侧为红、右侧为蓝。取一个确定红的 `red` 和一个确定蓝的 `blue`，维护整数区间 `(red, blue]`。

- `red` 不属于候选区间，且始终红。
- `blue` 属于候选区间，且始终蓝。
- `mid` 为蓝时收 `blue=mid`；为红时收 `red=mid`。
- `blue-red=1` 时，唯一候选就是 `blue`。

这套约定避免同时混用“闭区间左右端点”“答案下标”和“可行区间”的含义。

## 更换目标时只换颜色

要找最大可行值，可以把“可行”染红、“不可行”染蓝，仍维护 `(red, blue]`，结束时返回 `red`。模板本身不变；先写出哪种颜色代表目标，再确定返回哪一端。

具体的分组检查见 [[algo.binary-search.answer]]。
