---
id: algo.tree.binary-properties
kind: concept
title: 二叉树常见性质
summary: 第 d 层最多 2^d 个节点，高度 h 的二叉树最多 2^(h+1)-1 个节点；完全二叉树适合数组下标。
parents: [algo.tree]
tags: [tree, binary-tree, complete-tree]
links: [algo.tree.degree-count, algo.tree.traversal]
---

## 层数、节点数、高度

约定根深度为 0：第 `d` 层最多 `2^d` 个节点；高度为 `h` 的二叉树最多有 `2^(h+1)-1` 个节点。达到上界时每层都满。

完全二叉树只允许最后一层不满，且节点从左到右连续。1 下标数组中：

```go
left, right := 2*i, 2*i+1
parent := i / 2
```

0 下标数组改为 `left=2*i+1`、`right=2*i+2`。堆正是完全二叉树的数组表示。

## 三种容易混淆的树

- **满二叉树**：每个节点要么 0 个孩子，要么 2 个孩子。
- **完美二叉树**：所有内部节点都有 2 个孩子，所有叶子深度相同。
- **完全二叉树**：按层从左到右填充；不要求每层都满。

遇到“数组存储、父子下标”优先想到完全二叉树；遇到“叶子数和内部节点数”优先确认是否满二叉树。
