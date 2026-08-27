---
id: algo.math.gcd
kind: concept
title: 最大公因数 GCD：欧几里得算法
summary: gcd(a,b) 等于 gcd(b, a mod b)，每次取余都会严格缩小第二个数。
parents: [algo.math]
tags: [math, gcd, number-theory]
links: [algo.math.lcm]
---

## 核心等式

若 `a = qb + r`，那么 `a` 与 `b` 的公因数集合，和 `b` 与 `r` 的公因数集合完全相同。因此不断把 `(a,b)` 换成 `(b,a%b)`，直到余数为 0。

```go
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b // 新的第二项严格变小
	}
	return a
}
```

`gcd(48,18)` 的过程是 `(48,18) → (18,12) → (12,6) → (6,0)`，答案为 `6`。时间复杂度为 `O(log min(a,b))`。

## 常见用途

- 约分：分子分母同时除以 GCD。
- 同余、分组、周期题：先把共同单位缩到最小。
- 求 LCM：先算 GCD，见 [[algo.math.lcm]]。

输入可能为负数时，先取绝对值；`gcd(0,x)=|x|`，两个数同时为 0 时需按题意单独定义。
