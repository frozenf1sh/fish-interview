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
	Index      int    `json:"index"`
	Value      int    `json:"value"`
	State      string `json:"state"`
	Dependency bool   `json:"dependency"`
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
		state.Cells[i] = dpCell{Index: i, Value: value, State: cellState, Dependency: current >= 2 && (i == current-1 || i == current-2)}
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

// BinaryRedBluePartition traces the left-open, right-closed red-blue template.
func BinaryRedBluePartition() Trace {
	nums := []int{7, 2, 5, 10, 8}
	const groups = 2
	const firstAnswer, maximum = 10, 32
	red, blue := firstAnswer-1, maximum
	result := Trace{
		Kind:  "binary-red-blue",
		Title: "二分答案：红蓝染色",
		Pseudocode: []string{
			"red, blue := max(nums)-1, sum(nums)",
			"for red+1 < blue {",
			"    mid := red + (blue-red)/2",
			"    if groupsNeeded(nums, mid) <= k {",
			"        blue = mid // mid 为蓝色可行点",
			"    } else { red = mid } // mid 为红色不可行点",
			"}",
			"return blue",
		},
	}
	result.Frames = append(result.Frames, redBlueFrame(nums, red, maximum, red, -1, blue, 0, false, 0, "red=9 确定不可行，blue=32 确定可行；候选区间为 (red, blue]。"))
	for red+1 < blue {
		mid := red + (blue-red)/2
		needed := groupsNeeded(nums, mid)
		feasible := needed <= groups
		result.Frames = append(result.Frames, redBlueFrame(nums, firstAnswer-1, maximum, red, mid, blue, needed, feasible, 2, "尝试 "+itoa(mid)+"：用贪心扫描计算所需分组数。"))
		result.Frames = append(result.Frames, redBlueFrame(nums, firstAnswer-1, maximum, red, mid, blue, needed, feasible, 3, "最多允许 "+itoa(groups)+" 组；当前需要 "+itoa(needed)+" 组，因此 mid 被染色。"))
		if feasible {
			blue = mid
			result.Frames = append(result.Frames, redBlueFrame(nums, firstAnswer-1, maximum, red, mid, blue, needed, true, 4, "mid 为蓝色可行点：收缩右端点 blue="+itoa(blue)+"，右端点保持闭区间。"))
		} else {
			red = mid
			result.Frames = append(result.Frames, redBlueFrame(nums, firstAnswer-1, maximum, red, mid, blue, needed, false, 5, "mid 为红色不可行点：收缩左端点 red="+itoa(red)+"，左端点保持开区间。"))
		}
	}
	result.Frames = append(result.Frames, redBlueFrame(nums, firstAnswer-1, maximum, red, -1, blue, groupsNeeded(nums, blue), true, 7, "blue-red=1，区间 (red, blue] 只剩 blue="+itoa(blue)+"；它是第一个蓝点。"))
	return result
}

type redBlueState struct {
	Numbers  []int `json:"numbers"`
	Minimum  int   `json:"minimum"`
	Maximum  int   `json:"maximum"`
	Red      int   `json:"red"`
	Mid      int   `json:"mid"`
	Blue     int   `json:"blue"`
	Groups   int   `json:"groups"`
	Feasible bool  `json:"feasible"`
}

func redBlueFrame(numbers []int, minimum, maximum, red, mid, blue, groups int, feasible bool, line int, narration string) Frame {
	return Frame{
		ActiveLine: line,
		Narration:  narration,
		Variables: map[string]string{
			"red":    itoa(red),
			"mid":    itoa(mid),
			"blue":   itoa(blue),
			"groups": itoa(groups),
		},
		State: redBlueState{Numbers: append([]int(nil), numbers...), Minimum: minimum, Maximum: maximum, Red: red, Mid: mid, Blue: blue, Groups: groups, Feasible: feasible},
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
