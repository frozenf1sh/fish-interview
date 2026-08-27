---
id: algo.tree.traversal
kind: concept
title: 树遍历与序列性质
summary: 前中后序分别在访问左右子树前、中、后处理根；层序是 BFS。中序有序只对二叉搜索树成立。
parents: [algo.tree]
tags: [tree, traversal, bfs, dfs]
links: [algo.tree.terms]
---

## 三种 DFS 的根位置

```go
func dfs(node *TreeNode) {
	if node == nil { return }
	// 前序：这里处理 node
	dfs(node.Left)
	// 中序：这里处理 node
	dfs(node.Right)
	// 后序：这里处理 node
}
```

前序适合复制/序列化，后序适合“先得到子树结果再算父节点”，中序在 BST 中按升序给出键。仅凭前序和后序一般不能唯一还原二叉树；前序加中序且节点值不同可以。

## 层序遍历

队列当前长度就是本层节点数：

```go
queue := []*TreeNode{root}
for len(queue) > 0 {
	size := len(queue)
	for ; size > 0; size-- {
		node := queue[0]; queue = queue[1:]
		if node.Left != nil { queue = append(queue, node.Left) }
		if node.Right != nil { queue = append(queue, node.Right) }
	}
}
```

切片头删会保留底层数组；数据很大时用下标 `head` 代替反复截短。
