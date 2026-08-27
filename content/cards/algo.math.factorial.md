---
id: algo.math.factorial
kind: concept
title: 阶乘与模阶乘
summary: n! 是 1 到 n 的乘积；在取模题中预处理阶乘和逆阶乘，把重复组合计算降为 O(1)。
parents: [algo.math]
tags: [math, factorial, modulo]
links: [algo.math.combination]
---

## 基础计算

```go
func factorial(n int) int64 {
	ans := int64(1)
	for x := 2; x <= n; x++ { ans *= int64(x) }
	return ans
}
```

整数阶乘增长很快：`20!` 已接近 `int64` 上限。笔试中若答案要求取模 `mod`，每次乘法立刻 `% mod`。

## 模阶乘与逆阶乘

当 `mod` 是质数，费马小定理给出 `x^(mod-2) mod mod` 是 `x` 的乘法逆元。预处理后：

```go
fact[0], invFact[n] = 1, powMod(fact[n], mod-2, mod)
for i := 1; i <= n; i++ { fact[i] = fact[i-1] * int64(i) % mod }
for i := n; i >= 1; i-- { invFact[i-1] = invFact[i] * int64(i) % mod }
```

注意先完成 `fact`，再计算 `invFact[n]`；组合数见 [[algo.math.combination]]。
