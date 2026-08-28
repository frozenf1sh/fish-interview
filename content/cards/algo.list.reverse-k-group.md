---
id: algo.list.reverse-k-group
kind: algorithm-pattern
title: 链表：k 个一组翻转
summary: 先定位完整分组并保存后继，再逐条反转组内指针，最后用 dummy 把临时链接回原链。
parents: [algo.patterns.linked-list]
tags: [linked-list, dummy, pointer]
links: [algo.list.dummy-rewire, algo.list.fast-slow, algo.list.merge-sort]
trace: list-k-group
---

## 例题

[LeetCode 25 · K 个一组翻转链表 ↗](https://leetcode.cn/problems/reverse-nodes-in-k-group/)。把 `1→2→3→4→5` 按 `k=3` 翻转为 `3→2→1→4→5`；不足一组的 `4→5` 保持原样。

动画固定展示 `dummy→1→2→3→4→5`：先逐步定位第 1 组并断开，再把每条 `Next` 指针反向接到临时链，最后覆盖回主链。每个完整组都重复同样过程。

## 操作顺序

1. 从 `groupPrev` 出发走 `k` 步，先保存 `groupNext`。
2. 令 `prev=groupNext`，逐次执行 `next=cur.Next`、`cur.Next=prev`，直到当前组翻完。
3. 用 `groupPrev.Next=kth` 接回新头，再把 `groupPrev` 移到旧头；剩余不足 `k` 个时直接返回。

```go
func reverseKGroup(dummy *ListNode, k int) *ListNode {
	groupPrev := dummy
	for {
		kth := getKth(groupPrev, k)
		if kth == nil { break }
		groupNext := kth.Next
		prev, cur := groupNext, groupPrev.Next
		for cur != groupNext {
			next := cur.Next
			cur.Next = prev
			prev, cur = cur, next
		}
		oldHead := groupPrev.Next
		groupPrev.Next = kth
		groupPrev = oldHead
	}
	return dummy.Next
}
```
