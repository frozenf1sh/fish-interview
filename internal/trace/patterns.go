package trace

// LinearDPClimbStairs traces the recurrence dp[i] = dp[i-1] + dp[i-2].
// The example answers: how many ways are there to reach step 5 with moves of 1 or 2?
func LinearDPClimbStairs() Trace {
	const steps = 5
	dp := make([]int, steps+1)
	dp[0], dp[1] = 1, 1
	result := Trace{
		Kind:  "dp-table",
		Title: "线性 DP：从小状态推到目标状态",
		Pseudocode: []string{
			"dp[0], dp[1] = 1, 1 // 空路径与走一步各有一种方式",
			"for i := 2; i <= n; i++ {",
			"    dp[i] = dp[i-1] + dp[i-2] // 最后一步来自 1 或 2",
			"}",
			"return dp[n]",
		},
	}
	result.Frames = append(result.Frames, dpFrame(dp, -1, 0, "先定义状态：dp[i] 表示到达第 i 级台阶的方案数。"))
	result.Frames = append(result.Frames, dpFrame(dp, 1, 0, "边界 dp[0]=1、dp[1]=1：后续转移从这里开始。"))
	for i := 2; i <= steps; i++ {
		result.Frames = append(result.Frames, dpFrame(dp, i, 1, "计算 dp["+itoa(i)+"]：最后一步只可能走 1 级或 2 级。"))
		dp[i] = dp[i-1] + dp[i-2]
		result.Frames = append(result.Frames, dpFrame(dp, i, 2, "dp["+itoa(i)+"] = "+itoa(dp[i-1])+" + "+itoa(dp[i-2])+" = "+itoa(dp[i])+"。"))
	}
	result.Frames = append(result.Frames, dpFrame(dp, steps, 4, "答案是 dp[5] = "+itoa(dp[steps])+"。"))
	return result
}

type dpTableState struct {
	Cells   []dpCell `json:"cells"`
	Current int      `json:"current"`
}

type dpCell struct {
	Index int    `json:"index"`
	Value int    `json:"value"`
	State string `json:"state"`
}

func dpFrame(values []int, current, line int, narration string) Frame {
	state := dpTableState{Current: current, Cells: make([]dpCell, len(values))}
	for i, value := range values {
		cellState := "pending"
		if i < current || (i == 0 || i == 1) {
			cellState = "ready"
		}
		if i == current {
			cellState = "current"
		}
		state.Cells[i] = dpCell{Index: i, Value: value, State: cellState}
	}
	return Frame{
		ActiveLine: line,
		Narration:  narration,
		Variables: map[string]string{
			"i":  itoa(current),
			"dp": "到达每一级的方案数",
		},
		State: state,
	}
}

// BinaryAnswerPartition traces minimizing the largest group sum for two groups.
func BinaryAnswerPartition() Trace {
	nums := []int{7, 2, 5, 10, 8}
	const groups = 2
	const minimum, maximum = 10, 32
	low, high := minimum, maximum
	result := Trace{
		Kind:  "binary-answer",
		Title: "二分答案：最小化最大分组和",
		Pseudocode: []string{
			"lo, hi := max(nums), sum(nums)",
			"for lo < hi {",
			"    mid := lo + (hi-lo)/2",
			"    if groupsNeeded(nums, mid) <= k {",
			"        hi = mid // 这个上限可行，继续缩小",
			"    } else { lo = mid + 1 } // 上限太小",
			"}",
			"return lo",
		},
	}
	result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, -1, high, 0, false, 0, "答案范围是 [10, 32]：至少容纳最大元素，至多放进一个组。"))
	for low < high {
		mid := low + (high-low)/2
		needed := groupsNeeded(nums, mid)
		feasible := needed <= groups
		result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, mid, high, needed, feasible, 2, "尝试最大组和 "+itoa(mid)+"，用贪心扫描计算需要多少组。"))
		result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, mid, high, needed, feasible, 3, "上限 "+itoa(mid)+" 需要 "+itoa(needed)+" 组；目标是最多 "+itoa(groups)+" 组。"))
		if feasible {
			high = mid
			result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, mid, high, needed, true, 4, "可行：答案可以更小，把右边界收为 "+itoa(high)+"。"))
		} else {
			low = mid + 1
			result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, mid, high, needed, false, 5, "不可行：答案必须更大，把左边界提到 "+itoa(low)+"。"))
		}
	}
	result.Frames = append(result.Frames, binaryFrame(nums, minimum, maximum, low, low, high, groupsNeeded(nums, low), true, 7, "最小可行最大组和是 "+itoa(low)+"。"))
	return result
}

type binaryAnswerState struct {
	Numbers  []int `json:"numbers"`
	Minimum  int   `json:"minimum"`
	Maximum  int   `json:"maximum"`
	Low      int   `json:"low"`
	Mid      int   `json:"mid"`
	High     int   `json:"high"`
	Groups   int   `json:"groups"`
	Feasible bool  `json:"feasible"`
}

func binaryFrame(numbers []int, minimum, maximum, low, mid, high, groups int, feasible bool, line int, narration string) Frame {
	return Frame{
		ActiveLine: line,
		Narration:  narration,
		Variables: map[string]string{
			"lo":     itoa(low),
			"mid":    itoa(mid),
			"hi":     itoa(high),
			"groups": itoa(groups),
		},
		State: binaryAnswerState{Numbers: append([]int(nil), numbers...), Minimum: minimum, Maximum: maximum, Low: low, Mid: mid, High: high, Groups: groups, Feasible: feasible},
	}
}

func groupsNeeded(numbers []int, limit int) int {
	groups, sum := 1, 0
	for _, number := range numbers {
		if sum+number > limit {
			groups++
			sum = number
			continue
		}
		sum += number
	}
	return groups
}
