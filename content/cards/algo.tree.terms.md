---
id: algo.tree.terms
kind: concept
title: 树的术语与结构
summary: 树是连通且无环的图；根、父子、深度、高度、子树都是围绕一条唯一简单路径定义的。
parents: [algo.tree]
tags: [tree, graph, terminology]
links: [algo.tree.degree-count, algo.tree.traversal]
---

## 先看一棵树

```text
        A                 A 是根，深度 0
      /   \
     B     C               B、C 是 A 的孩子
    / \     \
   D   E     F             D、E、F 是叶子
```

- **父/子**：一条根到节点路径中相邻的上一层/下一层。
- **深度**：根到该节点的边数；上图 `depth(E)=2`。
- **高度**：该节点到其最远叶子的边数；叶子高度为 0。
- **子树**：某节点及全部后代构成的树；`B` 的子树有 `B,D,E`。
- **祖先/后代**：在根到节点路径上的前驱/后继。

根树中任意两点路径唯一。写 DFS 时，`parent` 参数就是为了避免沿无向边走回父节点。

## 根与边的方向

输入常给无向边，算法里先任选根并建立父子方向。树形 DP、LCA、子树和都依赖这个方向；未建根时“子树”没有确定含义。
