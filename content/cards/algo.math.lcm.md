---
id: algo.math.lcm
kind: concept
title: 最小公倍数 LCM：先除后乘
summary: 对非零整数，lcm(a,b)=a/gcd(a,b)*b；先除后乘可避免中间溢出。
parents: [algo.math]
tags: [math, lcm, number-theory]
links: [algo.math.gcd]
---

## 公式从哪里来

把两个数分解质因数时，GCD 取每个质因数的较小指数，LCM 取较大指数，所以：

$$
最大公因数 × 最小公倍数 = |a × b|
$$

实现时不要直接计算 `a*b/gcd`，乘积可能先溢出。

```go
func lcm(a, b int) int {
	if a == 0 || b == 0 { return 0 }
	return a / gcd(a, b) * b // 先约掉公因子
}
```

多个数的 LCM 从左到右累计：`ans = lcm(ans, x)`。若题目给出上界，乘前检查 `ans/gcd(ans,x) > limit/x`。
