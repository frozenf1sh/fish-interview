---
id: algo.list.dummy-rewire
kind: algorithm-pattern
title: 链表：dummy、断开与头插重连
summary: dummy 统一头部边界；重连时先保存 next，再断开、头插、恢复，避免丢失后续链表。
parents: [algo.patterns.linked-list]
tags: [linked-list, dummy, reverse]
links: [algo.list.fast-slow, algo.list.merge]
trace: linked-list-rewire
---

## 例题

[LeetCode 92 · 反转链表 II ↗](https://leetcode.cn/problems/reverse-linked-list-ii/)。将 `1→2→3→4→5` 的第 2 到 4 个节点反转为 `1→4→3→2→5`。

## 操作顺序

$$
保存 next → 从原链断开 next → next 头插到 pre 后
$$

> **识别信号**：题目改变一段链表的连接关系，且可能修改 head。先放 dummy；每次改 `Next` 前先保存仍要访问的节点。

## dummy 为什么必要

`left=1` 时，待翻转段前没有实际节点。让 `dummy.Next=head` 后，`pre` 永远存在，循环与返回值都不再需要分支。

动画中的主链、暂离节点和指针标签对应每一次重连。`cur` 固定在已翻转段的尾部；每轮从它后面抽一个 `next` 头插到 `pre` 后。

## 分段实现

### 1. 定位翻转段前驱

```go
dummy := &ListNode{Next: head}
pre := dummy
for step := 1; step < left; step++ {
	pre = pre.Next // pre 最终停在第 left 个节点之前
}
cur := pre.Next
```

### 2. 每轮抽取一个节点并头插

```go
for step := 0; step < right-left; step++ {
	next := cur.Next       // 先保存，否则断开后无法找到后续节点
	cur.Next = next.Next   // next 从原位置脱离，cur 保持为尾部
	next.Next = pre.Next   // next 指向当前翻转段头
	pre.Next = next        // next 插到 pre 后，成为新段头
}
return dummy.Next
```

## 同一套骨架

删除倒数第 k 个节点时让快指针先走 k 步，慢指针从 dummy 同步走；合并链表时让 `tail` 从 dummy 向后追加。两者都依赖 dummy 消除头节点特判。
