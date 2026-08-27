package trace

import "sort"

type sequenceState struct {
	Numbers []int  `json:"numbers"`
	Current int    `json:"current"`
	Tails   []int  `json:"tails"`
	Action  string `json:"action"`
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
	result.Frames = append(result.Frames, lisFrame(numbers, -1, nil, "长度为 0 的序列还没有结尾；tails 的下标加 1 就是长度。", 0))
	tails := []int{}
	for index, value := range numbers {
		position := sort.SearchInts(tails, value)
		if position == len(tails) {
			tails = append(tails, value)
			result.Frames = append(result.Frames, lisFrame(numbers, index, tails, "x="+itoa(value)+" 大于所有已有结尾，新增一个更长的长度。", 3))
			continue
		}
		old := tails[position]
		tails[position] = value
		result.Frames = append(result.Frames, lisFrame(numbers, index, tails, "x="+itoa(value)+" 替换长度 "+itoa(position+1)+" 的结尾 "+itoa(old)+"；长度不变，但更容易接后续数字。", 4))
	}
	result.Frames = append(result.Frames, lisFrame(numbers, -1, tails, "tails 长度为 "+itoa(len(tails))+"，即最长严格递增子序列长度。", 6))
	return result
}

func lisFrame(numbers []int, current int, tails []int, action string, line int) Frame {
	return Frame{
		ActiveLine: line,
		Narration:  action,
		Variables: map[string]string{
			"LIS 长度": itoa(len(tails)),
		},
		State: sequenceState{Numbers: append([]int(nil), numbers...), Current: current, Tails: append([]int(nil), tails...), Action: action},
	}
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
	cells := []string{"#", ".", "*", ".", "."}
	result.Frames = append(result.Frames, gravityFrame(cells, -1, 4, 0, "从右端开始，write 指向当前无障碍区间最右的可落位置。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 4, 4, 1, "空格不占位置，继续向左扫描。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 2, 1, 2, "遇到固定块 *：它切断两侧，左侧区间的新 write 变为它左边一格。"))
	result.Frames = append(result.Frames, gravityFrame(cells, 0, 1, 3, "遇到可下落块 #：移动到 write=1，随后 write 左移。"))
	cells = []string{".", "#", "*", ".", "."}
	result.Frames = append(result.Frames, gravityFrame(cells, -1, 0, 4, "这一行稳定为 . # * . .；逐行使用同一规则即可处理矩阵。"))
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
