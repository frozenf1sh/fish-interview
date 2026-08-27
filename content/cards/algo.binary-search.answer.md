---
id: algo.binary-search.answer
kind: algorithm-pattern
title: 二分答案：把最优值变成可判定问题
summary: 对答案范围二分，并构造单调 check(mid) 判断目标值是否可行。
parents: [algo.binary-search]
tags: [binary-search, monotonicity, greedy]
links: [algo.greedy.interval-scheduling]
trace: binary-answer
exam_signals:
  - company: netease
    year: 2027
    role: backend
    confidence: low
    source: https://www.nowcoder.com/
---

> **先识别单调性**：给定一个最大组和 `x`，若能在最多 `k` 组内完成划分，则任何更大的 `x` 也可行。

## 用“分成两组”建立 check

数组 `[7,2,5,10,8]` 分成两组，目标是最小化其中较大的组和。答案至少为最大元素 `10`，至多为总和 `32`。

对候选值 `x` 从左到右累加；下一项会超过 `x` 时就新开一组。这个贪心过程给出达到上限 `x` 所需的最少组数。组数不超过 `k`，说明 `x` 可行。

## Go 实现

```go
lo, hi := max(nums), sum(nums) // 答案的最小与最大边界
for lo < hi {
	mid := lo + (hi-lo)/2       // 尝试允许的最大组和
	if groupsNeeded(nums, mid) <= k {
		hi = mid // mid 可行，继续寻找更小的可行值
	} else {
		lo = mid + 1 // mid 太小，必须扩大上限
	}
}
return lo // 第一个可行答案
```

> **关键方向**：本题把“可行”放在右侧，`mid` 可行时收缩 `hi`。写二分前先把这句话写出来，能避免更新方向颠倒。

## 常见误区

- 只看到“二分”就写代码，没有先定义 `check(x)`。
- 把“最小化最大值”的更新方向写反。
- `check` 本身过慢；通常需要贪心、前缀和或双指针把它降到 `O(n)`。

## 下一步

本质是参数搜索。`check` 往往本身是一个贪心过程。
