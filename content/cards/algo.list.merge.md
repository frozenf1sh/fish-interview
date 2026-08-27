---
id: algo.list.merge
kind: concept
title: 链表：合并两个有序链表
summary: 用 dummy 固定结果头，用 tail 只向后移动；每步接入两个链表当前较小节点。
parents: [algo.patterns.linked-list]
tags: [linked-list, merge, dummy]
links: [algo.list.dummy-rewire]
---

## 例题

[LeetCode 21 · 合并两个有序链表 ↗](https://leetcode.cn/problems/merge-two-sorted-lists/)。

```go
dummy := &ListNode{}
tail := dummy
for a != nil && b != nil {
	if a.Val <= b.Val {
		tail.Next, a = a, a.Next
	} else {
		tail.Next, b = b, b.Next
	}
	tail = tail.Next // 结果尾部只前进，不回头
}
if a != nil { tail.Next = a } else { tail.Next = b }
return dummy.Next
```

dummy 消除了“第一个节点是谁”的特判。归并排序链表、合并 k 个有序链表都复用这个尾插骨架；k 路时把每条链表头放进优先队列。
