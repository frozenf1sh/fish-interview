---
id: algo.list.fast-slow
kind: concept
title: 链表：快慢指针
summary: fast 每次两步、slow 每次一步；它们的相对速度为 1，可用于中点、环检测和环入口。
parents: [algo.patterns.linked-list]
tags: [linked-list, two-pointers, cycle]
links: [algo.list.dummy-rewire]
---

## 例题

[LeetCode 141 · 环形链表 ↗](https://leetcode.cn/problems/linked-list-cycle/)。有环时 fast 最终会追上 slow；无环时 fast 会到 nil。

```go
slow, fast := head, head
for fast != nil && fast.Next != nil {
	slow = slow.Next
	fast = fast.Next.Next
	if slow == fast { return true }
}
return false
```

找中点时同样循环；`fast=head` 时偶数长度 slow 偏右，`fast=head.Next` 时偏左。找环入口时相遇后令一个指针回 head，两者每次一步，再次相遇即入口。
