package trace

import "sort"

type sequenceState struct {
	Numbers    []int    `json:"numbers"`
	Current    int      `json:"current"`
	Tails      []int    `json:"tails"`
	TailStates []string `json:"tailStates"`
	Action     string   `json:"action"`
}

// LISTrace shows how tails[length-1] keeps the smallest tail for each length.
func LISTrace() Trace {
	numbers := []int{10, 9, 2, 5, 3, 7, 101, 18}
	result := Trace{
		Kind:  "sequence-tails",
		Title: "LIS：用最小结尾的 tails 做二分",
		Pseudocode: []string{
			"tails := []int{}",
			"for i, x := range nums {",
			"    pos := sort.SearchInts(tails, x)",
			"    if pos == len(tails) { tails = append(tails, x) }",
			"    else { tails[pos] = x }",
			"}",
			"return len(tails)",
		},
	}
	result.Frames = append(result.Frames, lisFrame(numbers, -1, []int{}, []string{}, "长度为 0 的序列还没有结尾；tails 的下标加 1 就是长度。", 0))
	tails := []int{}
	for index, value := range numbers {
		position := sort.SearchInts(tails, value)
		beforeStates := stableTailStates(len(tails))
		if position < len(beforeStates) {
			beforeStates[position] = "dependency"
		}
		result.Frames = append(result.Frames, lisFrame(numbers, index, tails, beforeStates, "读入 x="+itoa(value)+"，二分定位第一个 >= x 的 tails 位置。蓝色结尾参与本次判断。", 2))
		if position == len(tails) {
			tails = append(tails, value)
			states := stableTailStates(len(tails))
			states[len(states)-1] = "current"
			result.Frames = append(result.Frames, lisFrame(numbers, index, tails, states, "x="+itoa(value)+" 大于所有已有结尾，新增长度 "+itoa(len(tails))+" 的最小结尾。", 3))
			continue
		}
		old := tails[position]
		tails[position] = value
		states := stableTailStates(len(tails))
		states[position] = "current"
		result.Frames = append(result.Frames, lisFrame(numbers, index, tails, states, "把长度 "+itoa(position+1)+" 的结尾 "+itoa(old)+" 写成 "+itoa(value)+"。长度不变，但后续更容易接数字。", 4))
	}
	result.Frames = append(result.Frames, lisFrame(numbers, -1, tails, stableTailStates(len(tails)), "tails 长度为 "+itoa(len(tails))+"，即最长严格递增子序列长度。", 6))
	return result
}

func lisFrame(numbers []int, current int, tails []int, tailStates []string, action string, line int) Frame {
	return Frame{
		ActiveLine: line,
		Narration:  action,
		Variables: map[string]string{
			"LIS 长度": itoa(len(tails)),
		},
		State: sequenceState{Numbers: append(make([]int, 0, len(numbers)), numbers...), Current: current, Tails: append(make([]int, 0, len(tails)), tails...), TailStates: append(make([]string, 0, len(tailStates)), tailStates...), Action: action},
	}
}

func stableTailStates(length int) []string {
	states := make([]string, length)
	for index := range states {
		states[index] = "ready"
	}
	return states
}

type gravityState struct {
	Cells  []string `json:"cells"`
	Cursor int      `json:"cursor"`
	Write  int      `json:"write"`
}

// RowGravityTrace demonstrates rightward gravity inside obstacle-separated segments.
func RowGravityTrace() Trace {
	result := Trace{
		Kind:  "row-gravity",
		Title: "带障碍的重力模拟：写指针从右向左",
		Pseudocode: []string{
			"write := len(row) - 1",
			"for i := len(row)-1; i >= 0; i-- {",
			"    if row[i] == '*' { write = i - 1 }",
			"    if row[i] == '#' { move '#' to write; write-- }",
			"}",
		},
	}
	cells := []string{"#", ".", ".", "*", "#", "."}
	result.Frames = append(result.Frames, gravityFrame(cells, -1, 5, 0, "从右端开始，write=5 指向右侧区间的最右可落位置。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 5, 5, 1, "读取 row[5]='.'：空位不写入，write 保持 5。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 4, 5, 3, "读取 row[4]='#'：先高亮当前块和蓝色落点 write=5。"))
	cells = []string{"#", ".", ".", "*", ".", "#"}
	result.Frames = append(result.Frames, gravityFrame(cells, 5, 4, 3, "执行交换：# 落到 5，原位置变为空位，write 左移到 4。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 3, 2, 2, "读取 row[3]='*'：障碍切断区间，左侧的 write 重置为 2。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 2, 2, 1, "读取 row[2]='.'：空位不写入，继续扫描。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 1, 2, 1, "读取 row[1]='.'：空位不写入，write 仍为 2。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 0, 2, 3, "读取 row[0]='#'：当前块将落到左侧区间的 write=2。"))
	cells = []string{".", ".", "#", "*", ".", "#"}
	result.Frames = append(result.Frames, gravityFrame(cells, 2, 1, 3, "执行交换并让 write 左移；右侧已稳定，不会再被改动。"))
	result.Frames = append(result.Frames, gravityFrame(cells, -1, 1, 4, "这一行稳定为 . . # * . #；逐行使用同一规则即可处理矩阵。"))
	return result
}

func gravityFrame(cells []string, cursor, write, line int, narration string) Frame {
	return Frame{
		ActiveLine: line,
		Narration:  narration,
		Variables: map[string]string{
			"write": itoa(write),
		},
		State: gravityState{Cells: append([]string(nil), cells...), Cursor: cursor, Write: write},
	}
}
