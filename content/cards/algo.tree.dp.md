---
id: algo.tree.dp
kind: concept
title: 树形 DP：每个节点返回多个状态
summary: 树形 DP 从叶子向根后序计算；状态描述“选当前节点/不选当前节点”等局部条件，父节点只组合孩子结果。
parents: [algo.patterns.tree]
tags: [tree, dp, postorder]
links: [algo.dfs.tree, algo.dp.modeling]
---

## 例题

[LeetCode 337 · 打家劫舍 III ↗](https://leetcode.cn/problems/house-robber-iii/)。相邻父子不能同时选，求最大金额。

## 让子树返回两种答案

`take` 表示选当前节点时的最大值，`skip` 表示不选当前节点时的最大值。孩子必须先算完，所以处理位置在后序。

```go
func dfs(node *TreeNode) (take, skip int) {
	if node == nil { return 0, 0 }
	leftTake, leftSkip := dfs(node.Left)
	rightTake, rightSkip := dfs(node.Right)
	take = node.Val + leftSkip + rightSkip // 选父节点，孩子必须不选
	skip = max(leftTake, leftSkip) + max(rightTake, rightSkip)
	return take, skip
}
take, skip := dfs(root)
return max(take, skip)
```

树形 DP 的关键不是“递归”，而是明确每种父节点状态对孩子状态的限制。换根 DP 还需要第二遍 DFS，把父方向贡献传给孩子。
