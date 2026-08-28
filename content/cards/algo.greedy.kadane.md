---
id: algo.greedy.kadane
kind: algorithm-pattern
title: 贪心：最大子数组 Kadane
summary: 扫描数组时维护必须以当前位置结尾的最大和；当前和为负贡献时直接从当前元素重开，同时维护全局最大值。
parents: [algo.patterns.greedy]
tags: [greedy, dp, array]
links: [algo.dp.linear, algo.greedy.reachability]
trace: kadane
---

## 例题

[LeetCode 53 · 最大子数组和 ↗](https://leetcode.cn/problems/maximum-subarray/)。`[-2,1,-3,4,-1,2,1,-5,4]` 的最大和为 `6`，子数组为 `[4,-1,2,1]`。

动画固定整条数组和当前连续区间；读入每个数字时只更新当前和、是否从该数字重新开始，以及全局最大区间。

## 状态句子

`current`：**必须以当前下标结尾**的最大子数组和。读入 `x` 时只有两种来源：

$$
current = max(x, 旧 current + x)
$$

如果旧 `current` 是负数，接上它只会拖累 `x`，所以从 `x` 重开。`best` 记录所有已经完成位置中的最大 `current`。

## 实现

```go
current, best := nums[0], nums[0] // 第一个元素也是一个合法子数组
for _, x := range nums[1:] {
	current = max(x, current+x) // 只比较重开与延续，不枚举起点
	best = max(best, current)   // 每轮都用新结尾状态更新全局答案
}
return best
```

全负数组不应初始化为 0；否则 `[-3]` 会错误地返回空子数组。若题目允许空子数组，再单独把答案下界设为 0。
