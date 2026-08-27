---
id: algo.math.combination
kind: concept
title: 组合数：只选集合，不排顺序
summary: 从 n 个不同元素选 k 个但不区分顺序，C(n,k)=n!/(k!(n-k)!)，并满足杨辉递推。
parents: [algo.math]
tags: [math, combination, counting]
links: [algo.math.factorial, algo.math.permutation]
---

## 两条常用路线

直接计算适合单次、小范围：利用对称性先令 `k=min(k,n-k)`，边乘边除。模质数、多次查询则预处理阶乘与逆阶乘。

```go
func comb(n, k int64) int64 {
	if k < 0 || k > n { return 0 }
	if k > n-k { k = n-k }
	ans := int64(1)
	for i := int64(1); i <= k; i++ { ans = ans * (n-k+i) / i }
	return ans
}
```

无法做除法或需要整张表时，用杨辉递推：`C(n,k)=C(n-1,k-1)+C(n-1,k)`。一维压缩必须从右到左更新，否则当前行的值会污染左上角。

## 容易混淆

“从 5 人中选 2 人组成小组”是 `C(5,2)`；“选第 1、2 名”是 `A(5,2)`。区分点只有一个：交换两个位置后，答案是否算同一种。
