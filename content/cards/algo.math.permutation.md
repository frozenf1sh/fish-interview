---
id: algo.math.permutation
kind: concept
title: 排列数：顺序不同就是不同结果
summary: 从 n 个不同元素中按顺序取 k 个，数量为 A(n,k)=n×(n-1)×…×(n-k+1)。
parents: [algo.math]
tags: [math, permutation, counting]
links: [algo.math.combination]
---

## 先问顺序是否影响结果

选队长再选副队长，与选副队长再选队长不同，因此是排列。第一个位置有 `n` 种，第二个位置剩 `n-1` 种，连续相乘得到：

$$
A(n,k) = n! / (n-k)!
$$

```go
func perm(n, k int64) int64 {
	if k < 0 || k > n { return 0 }
	ans := int64(1)
	for x := int64(0); x < k; x++ { ans *= n - x }
	return ans
}
```

若元素有重复，不能直接套公式；需要除掉每类重复元素的阶乘，或用计数型 DFS。忽略顺序时转到 [[algo.math.combination]]。
