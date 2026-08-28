---
id: algo.greedy.weighted-intervals
kind: algorithm-pattern
title: 区间：带权调度的前驱 DP
summary: 带权区间不能沿用结束最早贪心；按结束时间排序后，用二分找每段的最近兼容前驱，在选与不选之间做 DP。
parents: [algo.patterns.greedy]
tags: [interval, dp, binary-search]
links: [algo.greedy.interval-scheduling, algo.sequence.lis]
trace: weighted-intervals
---

## 例题

[LeetCode 1235 · 规划兼职工作 ↗](https://leetcode.cn/problems/maximum-profit-in-job-scheduling/)。每份工作有开始、结束和收益，选择不重叠工作使收益最大。

## 为什么无权贪心失效

无权题选择结束早的区间，是为了给后续留最多数量的机会。带权题中一段晚结束的工作可能收益极高；例如 `A=[1,3],5`、`B=[2,5],100`、`C=[4,6],5`，结束最早选 A 只能得到 10，而 B 得到 100。

## 状态与前驱

按结束时间排序。`dp[i]` 表示前 `i` 份工作可获得的最大收益。对第 `i` 份工作：

$$
dp[i] = max(跳过第 i 份工作的 dp[i-1], 选择它的收益 + 最近兼容前驱的 dp)
$$

最近兼容前驱满足 `end <= start[i]`，在排序后的结束时间数组里二分查找。

## 实现

```go
sort.Slice(jobs, func(i, j int) bool { return jobs[i].end < jobs[j].end })
ends := make([]int, len(jobs))
dp := make([]int, len(jobs)+1) // dp[i]：前 i 份已排序工作
for i, job := range jobs {
	ends[i] = job.end
	prev := sort.Search(i, func(k int) bool { return ends[k] > job.start })
	take := job.profit + dp[prev] // prev 份工作都在 job 前结束
	dp[i+1] = max(dp[i], take)    // 跳过当前工作或选择它
}
return dp[len(jobs)]
```

它属于区间题的核心分支，但本质是 DP；树上把它放在贪心旁边是为了先提醒“无权规则不能直接迁移”。
