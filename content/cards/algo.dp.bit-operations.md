---
id: algo.dp.bit-operations
kind: concept
title: 状态压缩：位运算技巧
summary: 把集合映射为 mask 后，用测试位、置位、低位提取和子集枚举把集合操作写成常数时间循环。
parents: [algo.dp]
tags: [dp, bitmask, bit-operation]
links: [algo.dp.bitmask]
---

## 最小约定

第 `i` 个元素对应第 `i` 位：`1<<i`。`mask` 的该位为 1 表示元素已经在集合中。先固定这一层语义，后面的表达式才不会混乱。

```go
has := mask&(1<<i) != 0      // 第 i 位是否为 1
mask |= 1 << i               // 把第 i 位设为 1
mask &^= 1 << i              // 把第 i 位清为 0
available := full ^ mask     // full 范围内尚未选择的元素
```

`full` 应写成 `(1<<n)-1`，不要直接对 `mask` 取反；高位并不属于题目集合。

## 枚举一个集合里的元素

`lowbit := x & -x` 取出最低的 1。它适合逐个处理已选或未选元素，且不需要扫描全部 `0..n-1`。

```go
for remain := mask; remain != 0; remain &= remain - 1 {
	bit := remain & -remain               // 当前最低位的 1
	i := bits.TrailingZeros(uint(bit))    // bit 对应的元素下标
	_ = i
}
```

Go 中需要 `import "math/bits"`。若 `n` 很小，直接枚举 `i := 0; i < n; i++` 更直观；不要为省常数牺牲状态含义。

## 枚举 mask 的所有子集

子集划分、集合拆分和“从已选集合取一部分”会用到下面的循环。它会依次给出所有非空子集，复杂度由所有 `mask` 汇总时为 `O(3^n)`。

```go
for sub := mask; sub > 0; sub = (sub - 1) & mask {
	other := mask ^ sub // sub 与 other 恰好划分 mask
	_ = other
}
```

若空子集也有意义，在循环后单独处理 `sub=0`；这样边界更清楚。

## 与状态定义配合

位运算只回答“集合如何读写”。状态压缩 DP 还需要回答：除了集合外，未来决策是否依赖最后位置、剩余次数或其他维度。只用 `dp[mask]` 丢掉必要信息时，位操作再正确也无法得到正确转移；见 [[algo.dp.bitmask]]。
