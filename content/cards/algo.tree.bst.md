---
id: algo.tree.bst
kind: concept
title: 树题：BST 范围与中序
summary: 二叉搜索树满足左子树小于根、右子树大于根；中序遍历有序，递归边界应携带允许值域。
parents: [algo.patterns.tree]
tags: [tree, bst, inorder]
links: [algo.tree.traversal]
trace: flow-tree-bst
---

## 例题

[LeetCode 98 · 验证二叉搜索树 ↗](https://leetcode.cn/problems/validate-binary-search-tree/)。只比较节点和直接孩子不够，因为祖先也会限制当前节点范围。

动画固定一棵 BST 候选树，在节点旁显示当前 `(low, high)` 开区间；递归向左或向右时只收紧对应端点。

## 范围递归

`low, high` 是当前节点允许落入的开区间。向左收紧上界为 root 值，向右收紧下界为 root 值。

```go
func valid(node *TreeNode, low, high int64) bool {
	if node == nil { return true }
	v := int64(node.Val)
	if v <= low || v >= high { return false }
	return valid(node.Left, low, v) && valid(node.Right, v, high)
}
return valid(root, math.MinInt64, math.MaxInt64)
```

中序序列严格递增也是等价检查，但要保存前一个访问值。题目允许重复值时先确认“等于”只能放哪一边，再相应调整开闭边界。
