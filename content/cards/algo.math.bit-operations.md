---
id: algo.math.bit-operations
kind: concept
title: 通用位运算：掩码、低位与异或
summary: 用按位与测试/清位，用或置位，用异或翻转或消去相同数，再用 lowbit 处理最低位。
parents: [algo.math]
tags: [math, bit-operation, xor]
links: [algo.dp.bit-operations]
---

## 四个基础动作

```go
has := x&(1<<i) != 0 // 测试第 i 位
x |= 1 << i          // 置位
x &^= 1 << i         // 清位（Go 的 AND NOT）
x ^= 1 << i          // 翻转第 i 位
```

`x & -x` 得到最低位的 1；`x & (x-1)` 删除最低位的 1。它们适合统计 1 的个数、Fenwick 树和枚举集合元素。

## 异或的消去性质

`a^a=0`、`a^0=a`，且交换顺序不影响结果。因此“只有一个数出现一次、其余出现两次”可以线性异或：

```go
ans := 0
for _, x := range nums { ans ^= x }
```

位运算用于集合状态时，完整的子集枚举与 `full` 边界见 [[algo.dp.bit-operations]]。
