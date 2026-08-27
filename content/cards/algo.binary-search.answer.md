---
id: algo.binary-search.answer
kind: algorithm-pattern
title: 二分答案：红蓝染色法
summary: 约定红色为不可行、蓝色为可行，始终维护左开右闭区间 `(red, blue]`，最后返回第一个蓝点。
parents: [algo.binary-search]
tags: [binary-search, monotonicity, greedy]
links: [algo.binary-search.red-blue, algo.greedy.interval-scheduling]
trace: binary-red-blue
exam_signals:
  - company: netease
    year: 2027
    role: backend
    confidence: low
    source: https://www.nowcoder.com/
---

## 例题

[LeetCode 410 · 分割数组的最大值 ↗](https://leetcode.cn/problems/split-array-largest-sum/)

数组 `[7,2,5,10,8]` 分成两组，目标是最小化其中较大的组和，答案为 `18`。

> **固定约定**：红色 `check(x)=false`，蓝色 `check(x)=true`。维护 `(red, blue]`：左端点开区间，右端点闭区间。

## 用“分成两组”建立 check

答案至少为最大元素 `10`，至多为总和 `32`。令 `check(x)` 表示“最大组和不超过 `x` 时，是否能在两组内完成”。

对候选值 `x` 从左到右累加；下一项会超过 `x` 时就新开一组。这个贪心过程给出达到上限 `x` 所需的最少组数。组数不超过 `k`，说明 `x` 可行。

## 模板：左开右闭 `(red, blue]`

`red = 9` 是人为放在答案范围外的红色哨兵；`blue = 32` 是确定可行的蓝色端点。循环不变量始终是 `check(red)=false`、`check(blue)=true`。当 `blue-red=1` 时，区间内只剩第一个蓝点。

```go
red, blue := max(nums)-1, sum(nums) // red 必须不可行，blue 必须可行
for red+1 < blue {                  // 区间 (red, blue] 至少有两个候选
	mid := red + (blue-red)/2        // mid 严格落在两个端点之间
	if groupsNeeded(nums, mid) <= k { // check(mid)=true，mid 染蓝
		blue = mid // 保留左侧，寻找更小的蓝点
	} else {
		red = mid // mid 染红，排除它与左侧
	}
}
return blue // 右闭端点就是第一个蓝点
```

> **边界只记一句**：红点放左边并开区间，蓝点放右边并闭区间；答案取 `blue`。想找最后一个红点时，仍维护同一不变量，循环结束取 `red`。

## 常见误区

- 只看到“二分”就写代码，没有先定义 `check(x)`。
- 把 `red` 初始化为一个实际可行值，破坏左端点红色不变量。
- 循环写成 `red < blue`，最后只剩一个候选时无法收缩。
- `check` 本身过慢；通常需要贪心、前缀和或双指针把它降到 `O(n)`。

## 下一步

本质是参数搜索。`check` 往往本身是一个贪心过程。
