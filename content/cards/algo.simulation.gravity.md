---
id: algo.simulation.gravity
kind: algorithm-pattern
title: 网格模拟：带障碍的重力下落
summary: 对每个无障碍区间从落点方向扫描，用 write 指针记录下一块可落到的位置；固定块会重置 write。
parents: [algo.patterns.simulation]
tags: [simulation, grid, two-pointers]
links: [algo.dfs.grid]
trace: row-gravity
exam_signals:
  - company: netease
    year: 2025
    role: game-client
    confidence: medium
    source: https://www.nowcoder.com/feed/main/detail/23479f514ab9476082e96109a5ce14d9
---

## 例题

[LeetCode 1861 · 旋转后的盒子 ↗](https://leetcode.cn/problems/rotating-the-box/)。先让每行的石头向右下落，再旋转矩阵；这里先把每行下落的核心写对。

把 `#` 看作可下落块、`.` 看作空位、`*` 看作固定障碍。`#.*..` 向右下落后是 `.#*..`。

## 识别方式

题目让物体沿固定方向移动，途中有障碍，且最终状态只取决于同一行或同一列的局部相对顺序。不要逐格模拟“每一秒”的移动；直接把每个无障碍区间压实到落点一端。

## 写指针为什么正确

从右向左扫描时，`write` 始终指向当前无障碍区间最右的空位：

1. 遇到 `.`，不占位置。
2. 遇到 `#`，把它写到 `write`，`write--`。
3. 遇到 `*`，左右不能互相穿过，下一段的 `write` 重置为 `i-1`。

因为扫描方向与下落方向相反，每个 `#` 写入时，右侧已经是最终状态，不会覆盖还未读取的块。

## 分段实现

### 1. 一行内把块压到右侧

```go
func settleRight(row []byte) {
	write := len(row) - 1 // 当前无障碍区间的最右可落点
	for i := len(row) - 1; i >= 0; i-- {
		if row[i] == '*' { // 固定块切断区间，左侧重新找落点
			write = i - 1
			continue
		}
		if row[i] != '#' {
			continue // 空位保持给 write 使用
		}
		row[i], row[write] = '.', '#' // 交换不会丢掉未扫描位置
		write--
	}
}
```

### 2. 扩展到矩阵

如果重力向右，逐行调用 `settleRight`。向下则逐列从底向上扫描；向左、向上只需把扫描方向和 `write` 的初值一起翻转。旋转题通常先按原行完成下落，再按坐标变换写入新矩阵。

## 边界

- `*` 自己不移动，且它两侧的 `#` 永远不能互换位置。
- 同一无障碍区间内，所有 `#` 的相对顺序无关；只关心它们数量。
- 二维题先确认“重力发生在旋转前还是后”；方向取决于这一步，不能只看最终矩阵朝向。

[[algo.dfs.grid]] 处理的是连通性；本题不走图，核心是方向、障碍与写指针不变量。
