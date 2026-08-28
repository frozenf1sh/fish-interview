---
id: algo.dfs.tree
kind: concept
title: DFS：树的递归骨架
summary: 递归函数先定义“当前节点要返回什么”，再决定在递归前、中、后处理当前节点。
parents: [algo.patterns.dfs]
tags: [dfs, tree, recursion]
links: [algo.tree.traversal]
trace: flow-dfs-tree
---

## 例题

[LeetCode 104 · 二叉树的最大深度 ↗](https://leetcode.cn/problems/maximum-depth-of-binary-tree/)。每棵子树返回自身高度，父节点取左右较大值加一。

动画使用固定二叉树展示真实递归调用栈：进入节点、走左/右子树、回到节点组合高度分别占一步，节点和边不会被抽象成一串文字。

## 递归三问

1. 参数代表哪棵子树？
2. 空节点返回什么？
3. 当前节点如何组合左右子树结果？

```go
func maxDepth(node *TreeNode) int {
	if node == nil { return 0 } // 空树高度
	left := maxDepth(node.Left)
	right := maxDepth(node.Right)
	return max(left, right) + 1 // 后序位置组合孩子答案
}
```

需要收集根到叶路径时，在递归前加入当前节点、返回前撤销；需要维护父节点时，把 `parent` 显式作为参数。
