---
id: algo.sequence.lis
kind: algorithm-pattern
title: 序列：最长递增子序列 LIS 的二分优化
summary: 用 tails[len-1] 保存每个长度能达到的最小结尾；新数替换第一个不小于它的结尾，或延长一个长度。
parents: [algo.patterns.sequence]
tags: [sequence, lis, binary-search, dp]
links: [algo.binary-search.answer, algo.dp.linear]
trace: lis
exam_signals:
  - company: alibaba
    year: 2025
    role: algorithm
    confidence: medium
    source: https://www.nowcoder.com/discuss/786631717116207104
---

## 例题

[LeetCode 300 · 最长递增子序列 ↗](https://leetcode.cn/problems/longest-increasing-subsequence/)。

`[10,9,2,5,3,7,101,18]` 的答案是 `4`，可取 `2,3,7,18`。

## 什么时候用

题目要求从原序列中保留相对顺序，并让值严格递增；`n` 足以让 `O(n²)` 的“枚举前一个位置”超时。若题目要输出具体方案，还需记录每个位置的前驱；这里只求长度。

## 先保存“同长度里最有前途的结尾”

令 `tails[k]` 表示：长度为 `k+1` 的严格递增子序列中，最小的结尾值。`tails` 自身严格递增，所以能二分。

读入 `x` 后，找第一个 `>= x` 的位置 `pos`：

- 找到：`tails[pos]=x`。长度没变，但结尾变小，后面更容易接数字。
- 找不到：`x` 比所有结尾大，追加它，LIS 长度加一。

例如处理到 `3` 时已有 `tails=[2,5]`。替换得到 `[2,3]`；它不代表原序列里存在 `2,3` 之后的所有展示路径，只表示长度 2 可以用更小结尾 `3` 达成。

## 分段实现

### 1. 建立每个长度的最小结尾

```go
tails := make([]int, 0, len(nums)) // tails[k]：长度 k+1 的最小结尾
```

### 2. 二分定位并替换或追加

```go
for _, x := range nums {
	pos := sort.SearchInts(tails, x) // 第一个 >= x 的位置；严格递增因此不能找 > x
	if pos == len(tails) {
		tails = append(tails, x) // x 接在最长长度后，答案增加 1
	} else {
		tails[pos] = x // 同一长度换成更小结尾，为后续数字留空间
	}
}
return len(tails)
```

`sort.SearchInts` 的边界是 `[0,len(tails))`，返回 `len(tails)` 表示“所有位置都小于 x”。这正好对应追加分支。

## 容易混淆的地方

- **严格递增**：找第一个 `>= x`；相等值必须替换，不能把长度增加。
- **非递减**：找第一个 `> x`，相等值才允许续在后面。
- **求方案**：额外维护 `parent[i]` 和每个长度对应的末尾下标；不要从 `tails` 直接恢复原序列。

[[algo.binary-search.answer]] 解释的是红蓝二分的区间约定；LIS 这里用的是标准库的下界搜索，目标是第一个满足比较条件的位置。
