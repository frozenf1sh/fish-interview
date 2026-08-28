---
id: algo.list.merge-sort
kind: algorithm-pattern
title: 链表：归并排序
summary: 用快慢指针断开链表，递归排序左右子链，再用 dummy 与 tail 按序合并；链表无需随机访问即可完成 O(n log n) 排序。
parents: [algo.patterns.linked-list]
tags: [linked-list, merge-sort, divide-and-conquer]
links: [algo.list.merge, algo.list.fast-slow]
trace: list-merge-sort
---

## 例题

[LeetCode 148 · 排序链表 ↗](https://leetcode.cn/problems/sort-list/)。把 `4→2→1→3` 排成 `1→2→3→4`。

```go
func sort(head *ListNode) *ListNode {
	if head == nil || head.Next == nil { return head }
	mid := splitWithFastSlow(head)
	left, right := sort(head), sort(mid)
	return merge(left, right)
}
```

拆分必须先断开 `slow.Next`，否则递归不会缩短；合并阶段沿用 [[algo.list.merge]] 的不变量：`dummy.Next` 是结果头，`tail` 始终指向结果尾，输入头只消费一个节点。
