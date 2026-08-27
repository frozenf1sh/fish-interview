---
id: algo.tree.lca
kind: concept
title: 树题：最近公共祖先 LCA
summary: 在普通二叉树中，递归返回“子树是否找到目标”；左右各找到一个时当前节点就是 LCA。
parents: [algo.patterns.tree]
tags: [tree, lca, dfs]
links: [algo.dfs.tree]
---

## 例题

[LeetCode 236 · 二叉树的最近公共祖先 ↗](https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/)。

## 返回值的含义

函数返回当前子树中找到的 `p`、`q` 或它们的 LCA；空表示两者都不在这棵子树。

```go
func lca(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q { return root }
	left := lca(root.Left, p, q)
	right := lca(root.Right, p, q)
	if left == nil { return right }
	if right == nil { return left }
	return root // p、q 分居左右子树
}
```

BST 的 LCA 可以利用值域从根迭代下降：两个值都小于根往左，都大于根往右，否则当前根即答案。
