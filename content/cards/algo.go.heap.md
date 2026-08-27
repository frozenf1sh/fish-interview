---
id: algo.go.heap
kind: concept
title: Go：container/heap 与优先队列
summary: container/heap 只维护堆序；自定义类型实现 Len、Less、Swap、Push、Pop 后，用堆顶完成最小/最大优先队列。
parents: [algo.go]
tags: [golang, heap, priority-queue]
links: [algo.tree.binary-properties]
---

## 最小堆模板

`Less(i,j)` 返回 true 表示 `i` 应更靠近堆顶。最小堆用 `<`，最大堆改为 `>`。

```go
import "container/heap"

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	x := old[len(old)-1] // heap 已把堆顶交换到切片末尾
	*h = old[:len(old)-1]
	return x
}
```

## 实际使用

```go
h := &IntHeap{3, 1, 4}
heap.Init(h)             // 把已有切片原地建堆
heap.Push(h, 2)
smallest := heap.Pop(h).(int)
top := (*h)[0]           // 只查看堆顶，不删除
_ = smallest; _ = top
```

优先队列通常存结构体，例如 `{value, priority}`，让 `Less` 比较 `priority`。需要修改任意元素时维护 `index` 字段并调用 `heap.Fix`；只做入队/出队时不要增加这层复杂度。
