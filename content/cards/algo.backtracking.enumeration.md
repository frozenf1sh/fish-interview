---
id: algo.backtracking.enumeration
title: 回溯：枚举下一个候选
kind: concept
summary: 当每一步要从一组剩余候选中选一个时，用 for 循环枚举；排列用 used，组合用 start 去重。
parents: [algo.patterns.backtracking]
tags: [backtracking, permutation, combination]
links: [algo.backtracking.choose-skip, algo.math.permutation, algo.math.combination]
trace: flow-backtracking-enumeration
---

## 排列、组合、子集的分界

- **排列**：位置有顺序，同一元素每条路径只能用一次，使用 `used[i]`。
- **组合**：只关心集合，递归传 `start`，下一层从 `i+1` 开始，避免 `[1,2]` 与 `[2,1]` 重复。
- **子集**：可用选或不选；也可用组合式枚举并在每层收集。

## 排列模板

[LeetCode 46 · 全排列 ↗](https://leetcode.cn/problems/permutations/)。

```go
used := make([]bool, len(nums))
var dfs func()
dfs = func() {
	if len(path) == len(nums) { ans = append(ans, append([]int(nil), path...)); return }
	for i, x := range nums {
		if used[i] { continue }
		used[i], path = true, append(path, x)
		dfs()
		path, used[i] = path[:len(path)-1], false
	}
}
```

有重复元素时，先排序；同一层跳过相同值，条件是 `i>0 && nums[i]==nums[i-1] && !used[i-1]`。组合题则在同层跳过重复值，并递归 `i+1`。
