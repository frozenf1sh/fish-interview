---
id: algo.tree.path-sum
kind: concept
title: 树题：根到叶路径与路径和
summary: 根到叶路径题把当前路径和或 path 作为递归状态；叶子节点才是一次完整路径的结束。
parents: [algo.patterns.tree]
tags: [tree, path, dfs]
links: [algo.dfs.path-enumeration]
trace: flow-tree-path-sum
---

## 例题

[LeetCode 112 · 路径总和 ↗](https://leetcode.cn/problems/path-sum/)。是否存在一条根到叶路径，其节点和等于 target。

动画固定带数值的二叉树和 target=22；进入节点时路径与剩余和更新，走到叶子时单独判断并回退到下一个分支。

```go
func hasPathSum(node *TreeNode, remain int) bool {
	if node == nil { return false }
	remain -= node.Val
	if node.Left == nil && node.Right == nil { return remain == 0 } // 必须是叶子
	return hasPathSum(node.Left, remain) || hasPathSum(node.Right, remain)
}
```

要输出所有路径时维护 `path` 并在叶子复制；要统计任意起点到任意终点路径和时，用前缀和 map，进入节点加计数、退出节点减计数，防止兄弟子树互相污染。
