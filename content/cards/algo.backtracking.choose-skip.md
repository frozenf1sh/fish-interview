---
id: algo.backtracking.choose-skip
kind: concept
title: 回溯：选或不选二叉决策树
summary: 每个位置只有“选/不选”两种独立决策时，用 index 沿数组推进，天然对应 2^n 个子集。
parents: [algo.patterns.backtracking]
tags: [backtracking, subset, recursion]
links: [algo.backtracking.enumeration]
trace: flow-backtracking-choose-skip
---

## 例题

[LeetCode 78 · 子集 ↗](https://leetcode.cn/problems/subsets/)。每个数只面对两种选择：放入当前集合，或跳过。

```go
path := make([]int, 0)
var dfs func(int)
dfs = func(index int) {
	if index == len(nums) {
		ans = append(ans, append([]int(nil), path...))
		return
	}
	dfs(index + 1)                         // 不选 nums[index]
	path = append(path, nums[index])
	dfs(index + 1)                         // 选 nums[index]
	path = path[:len(path)-1]
}
dfs(0)
```

它适合“每个元素独立二选一”。需要固定长度时，`len(path)==k` 提前收集；剩余元素不足以凑满时可以剪枝。
