package trace

// The traces in this file are deliberately concrete replays. Each frame is a
// small observable transition: inspect an input, read a dependency, write a
// state, or return from a recursive call. A frame is not a summary slide.

func SpaceOptimizationTrace() Trace    { return redesignedSpaceOptimizationTrace() }
func BitmaskTrace() Trace              { return redesignedBitmaskTrace() }
func BinaryRedBluePartition() Trace    { return redesignedBinaryRedBlueTrace() }
func LISTrace() Trace                  { return redesignedLISTrace() }
func RowGravityTrace() Trace           { return redesignedGravityTrace() }
func StartSortedIntervalsTrace() Trace { return redesignedStartSortedIntervalsTrace() }
func MeetingRoomsTrace() Trace         { return redesignedMeetingRoomsTrace() }
func WeightedIntervalsTrace() Trace    { return redesignedWeightedIntervalsTrace() }
func KadaneTrace() Trace               { return redesignedKadaneTrace() }

func deepExampleFrame(line int, narration, caption string, variables map[string]string, lanes ...exampleLane) Frame {
	values := cloneStringMap(variables)
	if values == nil {
		values = make(map[string]string)
	}
	values["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: values, State: exampleState{Caption: caption, Lanes: lanes}}
}

func deepCells(rows []string, states map[string]string) []matrixCell {
	return matrixFromRows(rows, states)
}

func greedyRangeFrame(line int, narration, caption string, min, max int, variables map[string]string, elements ...any) Frame {
	tracks := make([]greedyRangeTrack, 0)
	markers := make([]greedyRangeMarker, 0)
	for _, element := range elements {
		switch value := element.(type) {
		case greedyRangeTrack:
			tracks = append(tracks, value)
		case greedyRangeMarker:
			markers = append(markers, value)
		}
	}
	values := cloneStringMap(variables)
	if values == nil {
		values = make(map[string]string)
	}
	values["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: values, State: greedyRangeState{Caption: caption, Min: min, Max: max, Tracks: tracks, Markers: markers}}
}

func makeGreedyRangeTrack(label string, segments ...greedyRangeSegment) greedyRangeTrack {
	return greedyRangeTrack{Label: label, Segments: append([]greedyRangeSegment{}, segments...)}
}

func makeGreedyRangeSegment(start, end int, label, state, kind string) greedyRangeSegment {
	return greedyRangeSegment{Start: start, End: end, Label: label, State: state, Kind: kind}
}

func makeGreedyRangeMarker(track, label string, position int, state string) greedyRangeMarker {
	return greedyRangeMarker{Track: track, Label: label, Position: position, State: state}
}

func greedyRangeItems(values []string, states map[int]string) []greedyRangeSegment {
	segments := make([]greedyRangeSegment, len(values))
	for index, value := range values {
		state := "ready"
		if configured, ok := states[index]; ok {
			state = configured
		}
		segments[index] = makeGreedyRangeSegment(index, index+1, value, state, "item")
	}
	return segments
}

func greedyRangeIntervals(values []Interval, states map[int]string) []greedyRangeSegment {
	segments := make([]greedyRangeSegment, len(values))
	for index, value := range values {
		state := "ready"
		if configured, ok := states[index]; ok {
			state = configured
		}
		segments[index] = makeGreedyRangeSegment(value.Start, value.End, value.Label, state, "range")
	}
	return segments
}

func greedyStackItems(stack string, states map[int]string) []greedyRangeSegment {
	values := make([]string, 0, len(stack))
	for _, character := range stack {
		values = append(values, string(character))
	}
	return greedyRangeItems(values, states)
}

func redesignedGreedyReachabilityTrace() Trace {
	code := []string{"farthest := 0", "for i, jump := range nums {", "    if i > farthest { return false }", "    farthest = max(farthest, i+jump)", "}", "return true"}
	values := []string{"2", "3", "1", "1", "4"}
	positions := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("位置 / 跳数", greedyRangeItems(values, states)...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 nums=[2,3,1,1,4]。横轴是下标；绿色带状区表示已经证明可达的范围。", "跳跃游戏：最远可达边界", 0, 5, map[string]string{"farthest": "0", "i": "-"}, positions(map[int]string{0: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 1, "0", "current", "range")), makeGreedyRangeTrack("下一候选"), makeGreedyRangeMarker("位置 / 跳数", "起点", 0, "current")),
		greedyRangeFrame(1, "读取下标 0 和 jump=2；0 没有越过当前可达边界。", "检查当前位置", 0, 5, map[string]string{"i": "0", "jump": "2", "farthest": "0"}, positions(map[int]string{0: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 1, "0", "dependency", "range")), makeGreedyRangeTrack("下一候选"), makeGreedyRangeMarker("位置 / 跳数", "i=0", 0, "current")),
		greedyRangeFrame(3, "候选终点是 0+2=2；橙色范围从当前位置铺到下标 2，写入新的边界。", "扩展到下标 2", 0, 5, map[string]string{"i": "0", "jump": "2", "candidate": "2", "farthest": "2"}, positions(map[int]string{0: "dependency"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 3, "0..2", "current", "range")), makeGreedyRangeTrack("下一候选", makeGreedyRangeSegment(0, 3, "0→2", "current", "range")), makeGreedyRangeMarker("位置 / 跳数", "i=0", 0, "current")),
		greedyRangeFrame(1, "读取下标 1 和 jump=3；它落在绿色边界内，因此可以继续贡献范围。", "检查第二个位置", 0, 5, map[string]string{"i": "1", "jump": "3", "farthest": "2"}, positions(map[int]string{0: "ready", 1: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 3, "0..2", "dependency", "range")), makeGreedyRangeTrack("下一候选", makeGreedyRangeSegment(1, 5, "1→4", "current", "range")), makeGreedyRangeMarker("位置 / 跳数", "i=1", 1, "current")),
		greedyRangeFrame(3, "候选边界 1+3=4 越过旧边界 2；把绿色范围扩展到终点下标 4。", "覆盖终点", 0, 5, map[string]string{"i": "1", "candidate": "4", "farthest": "4"}, positions(map[int]string{0: "ready", 1: "dependency"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 5, "0..4", "current", "range")), makeGreedyRangeTrack("下一候选", makeGreedyRangeSegment(1, 5, "1→4", "current", "range")), makeGreedyRangeMarker("位置 / 跳数", "i=1", 1, "current")),
		greedyRangeFrame(1, "下标 2 已被覆盖；它的候选只能到 3，不会缩短已经证明的最远边界。", "边界内继续扫描", 0, 5, map[string]string{"i": "2", "jump": "1", "candidate": "3", "farthest": "4"}, positions(map[int]string{0: "ready", 1: "ready", 2: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 5, "0..4", "ready", "range")), makeGreedyRangeTrack("下一候选", makeGreedyRangeSegment(2, 4, "2→3", "dependency", "range")), makeGreedyRangeMarker("位置 / 跳数", "i=2", 2, "current")),
		greedyRangeFrame(1, "下标 3 仍在绿色范围内；候选终点与边界相同，不产生新的覆盖。", "继续扫描", 0, 5, map[string]string{"i": "3", "jump": "1", "candidate": "4", "farthest": "4"}, positions(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 5, "0..4", "ready", "range")), makeGreedyRangeTrack("下一候选", makeGreedyRangeSegment(3, 5, "3→4", "dependency", "range")), makeGreedyRangeMarker("位置 / 跳数", "i=3", 3, "current")),
		greedyRangeFrame(5, "绿色边界覆盖最后一个下标 4；扫描完成，返回 true。", "可达边界贪心：答案", 0, 5, map[string]string{"farthest": "4", "answer": "true"}, positions(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "current"}), makeGreedyRangeTrack("可达边界", makeGreedyRangeSegment(0, 5, "0..4", "ready", "range")), makeGreedyRangeTrack("下一候选"), makeGreedyRangeMarker("位置 / 跳数", "终点", 4, "current")),
	}
	return concreteTrace("greedy-range", "可达边界贪心：跳跃游戏", code, frames...)
}

func redesignedGreedyLexicographicTrace() Trace {
	code := []string{"last := lastIndex(s)", "for i, ch := range s {", "    if inStack[ch] { continue }", "    for top > ch && last[top] > i { pop() }", "    push(ch)", "}"}
	input := []string{"c", "b", "a", "c", "d", "c", "b", "c"}
	inputTrack := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("输入位置", greedyRangeItems(input, states)...)
	}
	stackTrack := func(stack string, states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("单调栈", greedyStackItems(stack, states)...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 s=cbacdcbc。横轴保留原字符串位置；栈在另一条固定轨道上增长和回退。", "去重字典序：准备", 0, 8, map[string]string{"stack": "[]", "last": "a→2 b→6 c→7 d→4"}, inputTrack(nil), stackTrack("", nil), makeGreedyRangeTrack("可补回")),
		greedyRangeFrame(1, "i=0 读 c；栈为空，没有需要比较的栈顶。", "读取 c", 0, 8, map[string]string{"i": "0", "ch": "c", "stack": "[]"}, inputTrack(map[int]string{0: "current"}), stackTrack("", nil), makeGreedyRangeTrack("可补回"), makeGreedyRangeMarker("输入位置", "i=0", 0, "current")),
		greedyRangeFrame(4, "c 没有更大的栈顶可弹出，先把它放入栈的第一个位置。", "压入 c", 0, 8, map[string]string{"i": "0", "ch": "c", "stack": "[c]"}, inputTrack(map[int]string{0: "dependency"}), stackTrack("c", map[int]string{0: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(1, 8, "c 后面仍会出现", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=0", 0, "current")),
		greedyRangeFrame(3, "读 b：栈顶 c>b，且 c 的后续范围仍在输入轴上，先执行 pop。", "执行 pop：c 可补回", 0, 8, map[string]string{"i": "1", "ch": "b", "top": "c", "stack": "[]"}, inputTrack(map[int]string{0: "dependency", 1: "current"}), stackTrack("c", map[int]string{0: "dependency"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(2, 8, "c 后面仍会出现", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=1", 1, "current")),
		greedyRangeFrame(4, "把 b 放入栈；此时栈轨道只保留 b，c 暂时移出但仍可由后面的 c 补回。", "压入 b", 0, 8, map[string]string{"i": "1", "ch": "b", "stack": "[b]"}, inputTrack(map[int]string{0: "dependency", 1: "current"}), stackTrack("b", map[int]string{0: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(2, 8, "c 可补回", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=1", 1, "current")),
		greedyRangeFrame(3, "读 a：栈顶 b>a，且 b 的后续范围仍未结束，继续弹出 b。", "执行 pop：b 可补回", 0, 8, map[string]string{"i": "2", "ch": "a", "top": "b", "stack": "[]"}, inputTrack(map[int]string{1: "dependency", 2: "current"}), stackTrack("b", map[int]string{0: "dependency"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(3, 8, "b 后面仍会出现", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=2", 2, "current")),
		greedyRangeFrame(4, "压入 a；从现在开始，a 固定在栈底，后面的字符不能覆盖它。", "压入 a", 0, 8, map[string]string{"i": "2", "ch": "a", "stack": "[a]"}, inputTrack(map[int]string{1: "rejected", 2: "current"}), stackTrack("a", map[int]string{0: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(3, 8, "b / c 可补回", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=2", 2, "current")),
		greedyRangeFrame(4, "i=3 读 c；栈顶 a<c，不触发弹栈，把 c 接到栈尾。", "压入 c", 0, 8, map[string]string{"i": "3", "ch": "c", "stack": "[a,c]"}, inputTrack(map[int]string{2: "dependency", 3: "current"}), stackTrack("ac", map[int]string{0: "ready", 1: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(5, 8, "c 可补回", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=3", 3, "current")),
		greedyRangeFrame(4, "i=4 读 d；d 大于栈顶 c，直接压入，栈轨道变成 a→c→d。", "压入 d", 0, 8, map[string]string{"i": "4", "ch": "d", "stack": "[a,c,d]"}, inputTrack(map[int]string{2: "dependency", 3: "ready", 4: "current"}), stackTrack("acd", map[int]string{0: "ready", 1: "ready", 2: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(5, 8, "c 可补回", "dependency", "range")), makeGreedyRangeMarker("输入位置", "i=4", 4, "current")),
		greedyRangeFrame(2, "i=5 再读 c；c 已经在栈中，重复字符直接跳过，主轨道和栈都保持不变。", "跳过重复 c", 0, 8, map[string]string{"i": "5", "ch": "c", "stack": "[a,c,d]"}, inputTrack(map[int]string{5: "rejected"}), stackTrack("acd", map[int]string{0: "ready", 1: "ready", 2: "ready"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(5, 8, "已在栈中", "rejected", "range")), makeGreedyRangeMarker("输入位置", "i=5", 5, "rejected")),
		greedyRangeFrame(3, "i=6 读 b：d>b，但 d 的最后位置已经过去；红色表示不能弹出 d。", "拒绝弹出 d", 0, 8, map[string]string{"i": "6", "ch": "b", "top": "d", "stack": "[a,c,d]"}, inputTrack(map[int]string{6: "current"}), stackTrack("acd", map[int]string{0: "ready", 1: "ready", 2: "rejected"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(7, 8, "d 不可补回", "rejected", "range")), makeGreedyRangeMarker("输入位置", "i=6", 6, "current")),
		greedyRangeFrame(4, "把 b 接到栈尾；最后的 c 再次跳过，栈保持 [a,c,d,b]。", "完成扫描", 0, 8, map[string]string{"i": "7", "stack": "[a,c,d,b]", "answer": "acdb"}, inputTrack(map[int]string{7: "rejected"}), stackTrack("acdb", map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"}), makeGreedyRangeTrack("可补回", makeGreedyRangeSegment(7, 8, "c 已在栈中", "rejected", "range")), makeGreedyRangeMarker("输入位置", "i=7", 7, "rejected")),
	}
	return concreteTrace("greedy-range", "字典序贪心：删除重复字母", code, frames...)
}

func redesignedGreedyEndpointsTrace() Trace {
	code := []string{"sort intervals by start", "right := first.end", "for _, in := range intervals[1:] {", "    if in.start > right { arrows++; right = in.end }", "    else { right = min(right, in.end) }", "}"}
	intervals := []Interval{{Label: "A", Start: 1, End: 6}, {Label: "B", Start: 2, End: 8}, {Label: "C", Start: 7, End: 12}, {Label: "D", Start: 10, End: 16}}
	all := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("气球区间", greedyRangeIntervals(intervals, states)...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "四段气球已按起点排序；先把 A 的右端 6 作为第一组共同交集。", "区间端点：最少箭数", 0, 16, map[string]string{"right": "6", "arrows": "1"}, all(map[int]string{0: "current"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(1, 6, "[1,6]", "current", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "A", 1, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready")),
		greedyRangeFrame(2, "读取 B=[2,8]；它的起点 2 没越过 right=6，两个区间仍有共同位置。", "检查 B 的重叠", 0, 16, map[string]string{"start": "2", "end": "8", "right": "6", "arrows": "1"}, all(map[int]string{0: "ready", 1: "current"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(1, 6, "right=6", "dependency", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "B", 2, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready")),
		greedyRangeFrame(4, "共同交集右端收紧为 min(6,8)=6；第一支箭继续放在 6。", "收紧第一组交集", 0, 16, map[string]string{"right": "6", "arrows": "1"}, all(map[int]string{0: "ready", 1: "ready"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(2, 6, "[2,6]", "current", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "箭=6", 6, "current"), makeGreedyRangeMarker("箭", "6", 6, "current")),
		greedyRangeFrame(2, "读取 C=[7,12]；7 越过 right=6，旧交集为空，必须新开一支箭。", "发现断开", 0, 16, map[string]string{"start": "7", "end": "12", "right": "6", "arrows": "1"}, all(map[int]string{0: "ready", 1: "ready", 2: "current"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(1, 6, "旧交集", "rejected", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "C", 7, "rejected"), makeGreedyRangeMarker("箭", "6", 6, "ready")),
		greedyRangeFrame(3, "新增第二支箭，新的共同交集从 C 的范围开始，右端先设为 12。", "开启第二组", 0, 16, map[string]string{"right": "12", "arrows": "2"}, all(map[int]string{0: "ready", 1: "ready", 2: "dependency"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(7, 12, "[7,12]", "current", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "新组", 7, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready"), makeGreedyRangeMarker("箭", "12", 12, "current")),
		greedyRangeFrame(2, "读取 D=[10,16]；10 没越过第二组 right=12，仍可由同一支箭命中。", "检查 D 的重叠", 0, 16, map[string]string{"start": "10", "end": "16", "right": "12", "arrows": "2"}, all(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(7, 12, "right=12", "dependency", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "D", 10, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready"), makeGreedyRangeMarker("箭", "12", 12, "current")),
		greedyRangeFrame(4, "第二组右端保持 min(12,16)=12；两支箭稳定在 6 和 12。", "稳定第二组", 0, 16, map[string]string{"right": "12", "arrows": "2"}, all(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready"}), makeGreedyRangeTrack("当前交集", makeGreedyRangeSegment(10, 12, "[10,12]", "current", "range")), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "箭=12", 12, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready"), makeGreedyRangeMarker("箭", "12", 12, "current")),
		greedyRangeFrame(5, "扫描结束，返回 arrows=2；横轴上的两条竖线就是最终选择。", "区间端点贪心：答案", 0, 16, map[string]string{"arrows": "2", "positions": "6,12"}, all(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready"}), makeGreedyRangeTrack("当前交集"), makeGreedyRangeTrack("箭"), makeGreedyRangeMarker("气球区间", "答案", 12, "current"), makeGreedyRangeMarker("箭", "6", 6, "ready"), makeGreedyRangeMarker("箭", "12", 12, "current")),
	}
	return concreteTrace("greedy-range", "区间端点贪心：最少箭数", code, frames...)
}

func redesignedReachabilityTrace() Trace {
	code := []string{"farthest := 0", "for i, jump := range nums {", "    if i > farthest { return false }", "    farthest = max(farthest, i+jump)", "}", "return true"}
	nums := []string{"2", "3", "1", "1", "4"}
	frames := []Frame{
		deepExampleFrame(0, "例题 nums=[2,3,1,1,4]。先只建立起点不变量：下标 0 可达。", "跳跃游戏：可达边界", map[string]string{"i": "-", "farthest": "0"}, tokenRow("nums", nums, nil), lane("边界", item("[0,0]", "current"))),
		deepExampleFrame(1, "进入循环读取 i=0、jump=2；先比较 i 与 farthest，0<=0。", "检查当前位置", map[string]string{"i": "0", "jump": "2", "farthest": "0"}, tokenRow("nums", nums, map[int]string{0: "current"}), lane("边界", item("0", "dependency"))),
		deepExampleFrame(3, "候选终点是 i+jump=0+2=2；橙色写入新的最远边界。", "计算候选边界", map[string]string{"i": "0", "jump": "2", "candidate": "2", "farthest": "2"}, tokenRow("nums", nums, map[int]string{0: "dependency"}), lane("比较", item("0+2=2", "current")), lane("边界", item("2", "current"))),
		deepExampleFrame(1, "读取 i=1、jump=3。i=1 仍在已知边界 2 内，可以继续扩展。", "读取第二个位置", map[string]string{"i": "1", "jump": "3", "farthest": "2"}, tokenRow("nums", nums, map[int]string{1: "current"}), lane("边界", item("2", "dependency"))),
		deepExampleFrame(3, "候选边界 1+3=4 大于旧边界 2，写 farthest=4。", "扩展到终点", map[string]string{"i": "1", "jump": "3", "candidate": "4", "farthest": "4"}, tokenRow("nums", nums, map[int]string{0: "ready", 1: "dependency"}), lane("比较", item("max(2,4)=4", "current")), lane("边界", item("4", "current"))),
		deepExampleFrame(1, "i=2 的 jump=1 仍在边界内；它只能给出候选 3，不会覆盖 4。", "边界覆盖中的位置", map[string]string{"i": "2", "jump": "1", "candidate": "3", "farthest": "4"}, tokenRow("nums", nums, map[int]string{2: "current"}), lane("比较", item("max(4,3)=4", "current")), lane("边界", item("4", "ready"))),
		deepExampleFrame(1, "i=3 的 jump=1 也已覆盖；候选 4 与边界相同。", "继续扫描", map[string]string{"i": "3", "jump": "1", "candidate": "4", "farthest": "4"}, tokenRow("nums", nums, map[int]string{3: "current"}), lane("边界", item("4", "ready"))),
		deepExampleFrame(5, "farthest=4 覆盖最后一个下标，扫描完成，返回 true。", "可达边界贪心：答案", map[string]string{"farthest": "4", "answer": "true"}, tokenRow("nums", nums, map[int]string{4: "current"}), lane("答案", item("true", "current"))),
	}
	return concreteTrace("example-state", "可达边界贪心：跳跃游戏", code, frames...)
}

func redesignedLexicographicTrace() Trace {
	code := []string{"last := lastIndex(s)", "for i, ch := range s {", "    if inStack[ch] { continue }", "    for top > ch && last[top] > i { pop() }", "    push(ch)", "}"}
	input := []string{"c", "b", "a", "c", "d", "c", "b", "c"}
	frames := []Frame{
		deepExampleFrame(0, "例题 s=cbacdcbc。先记录每个字符最后一次出现的位置，后续才能判断弹出是否可补回。", "去重字典序：准备", map[string]string{"stack": "[]"}, tokenRow("输入", input, nil), lane("last", item("a→2 b→6 c→7 d→4", "ready")), lane("stack", item("[]", "current"))),
		deepExampleFrame(1, "i=0 读 c；栈为空，没有可比较的栈顶。", "读取 c", map[string]string{"i": "0", "ch": "c", "stack": "[]"}, tokenRow("输入", input, map[int]string{0: "current"}), lane("stack", item("[]", "dependency"))),
		deepExampleFrame(4, "c 没有更大的栈顶可弹出，先写入栈。", "压入 c", map[string]string{"i": "0", "ch": "c", "stack": "[c]"}, tokenRow("输入", input, map[int]string{0: "dependency"}), lane("stack", item("c", "current"))),
		deepExampleFrame(1, "i=1 读 b；比较栈顶 c>b，且 last[c]=7>1，c 后面还能补回。", "判断是否弹出 c", map[string]string{"i": "1", "ch": "b", "top": "c", "last[top]": "7", "stack": "[c]"}, tokenRow("输入", input, map[int]string{1: "current"}), lane("栈顶", item("c>b", "current")), lane("未来", item("c 仍会出现", "dependency"))),
		deepExampleFrame(3, "执行 pop：c 暂时离开栈，但它不是丢失，因为后面还有 c。", "弹出可补回字符", map[string]string{"i": "1", "ch": "b", "stack": "[]"}, tokenRow("输入", input, map[int]string{0: "rejected", 1: "current"}), lane("栈", item("c→[]", "current"))),
		deepExampleFrame(4, "把 b 压入；当前最小前缀是 b。", "压入 b", map[string]string{"i": "1", "ch": "b", "stack": "[b]"}, tokenRow("输入", input, map[int]string{1: "dependency"}), lane("栈", item("b", "current"))),
		deepExampleFrame(1, "i=2 读 a：栈顶 b>a，且 b 的最后位置 6>2，继续弹出。", "为 a 让位", map[string]string{"i": "2", "ch": "a", "top": "b", "stack": "[b]"}, tokenRow("输入", input, map[int]string{2: "current"}), lane("栈顶", item("b>a，可补回", "dependency"))),
		deepExampleFrame(3, "弹出 b 后压入 a；栈从 [b] 变成 [a]。", "写入 a", map[string]string{"i": "2", "ch": "a", "stack": "[a]"}, tokenRow("输入", input, map[int]string{1: "rejected", 2: "current"}), lane("栈", item("a", "current"))),
		deepExampleFrame(4, "i=3 读 c，栈顶 a<c，不能破坏递增机会，直接压入。", "压入 c", map[string]string{"i": "3", "ch": "c", "stack": "[a,c]"}, tokenRow("输入", input, map[int]string{3: "current"}), lane("栈", item("a", "ready"), item("c", "current"))),
		deepExampleFrame(4, "i=4 读 d，d 大于栈顶 c，压入后得到 [a,c,d]。", "压入 d", map[string]string{"i": "4", "ch": "d", "stack": "[a,c,d]"}, tokenRow("输入", input, map[int]string{4: "current"}), lane("栈", item("a", "ready"), item("c", "ready"), item("d", "current"))),
		deepExampleFrame(2, "i=5 再读 c；c 已经在栈中，重复字符直接跳过，不改变栈。", "跳过重复 c", map[string]string{"i": "5", "ch": "c", "stack": "[a,c,d]"}, tokenRow("输入", input, map[int]string{5: "rejected"}), lane("栈", item("a", "ready"), item("c", "ready"), item("d", "ready"))),
		deepExampleFrame(1, "i=6 读 b：d>b 且 last[d]=4 已经过去，d 不能补回，所以不能弹 d。", "不能弹出 d", map[string]string{"i": "6", "ch": "b", "top": "d", "last[top]": "4", "stack": "[a,c,d]"}, tokenRow("输入", input, map[int]string{6: "current"}), lane("栈顶", item("d>b，但不可补回", "rejected"))),
		deepExampleFrame(3, "把 b 压到 d 后面，栈成为 [a,c,d,b]；最后的 c 再次跳过。", "完成扫描", map[string]string{"i": "7", "stack": "[a,c,d,b]", "answer": "acdb"}, tokenRow("输入", input, map[int]string{7: "rejected"}), lane("答案", item("acdb", "current"))),
	}
	return concreteTrace("example-state", "字典序贪心：删除重复字母", code, frames...)
}

func redesignedEndpointsTrace() Trace {
	code := []string{"sort intervals by start", "right := first.end", "for _, in := range intervals[1:] {", "    if in.start > right { arrows++; right = in.end }", "    else { right = min(right, in.end) }", "}"}
	intervals := []string{"[1,6]", "[2,8]", "[7,12]", "[10,16]"}
	frames := []Frame{
		deepExampleFrame(0, "例题四段气球区间已按起点排序；先把第一组共同交集的右端设为 6。", "区间端点：最少箭数", map[string]string{"right": "6", "arrows": "1"}, tokenRow("区间", intervals, map[int]string{0: "current"}), lane("交集", item("[1,6]", "current")), lane("箭", item("6", "ready"))),
		deepExampleFrame(2, "读取 [2,8]，先比较 start=2 与 right=6；2<=6，仍有共同命中位置。", "检查 [2,8]", map[string]string{"start": "2", "end": "8", "right": "6", "arrows": "1"}, tokenRow("区间", intervals, map[int]string{1: "current"}), lane("交集", item("right=6", "dependency"))),
		deepExampleFrame(4, "交集右端写成 min(6,8)=6；第一支箭继续放在 6。", "收紧第一组交集", map[string]string{"right": "6", "arrows": "1"}, tokenRow("区间", intervals, map[int]string{0: "ready", 1: "ready"}), lane("交集", item("[2,6]", "current")), lane("箭", item("6", "current"))),
		deepExampleFrame(2, "读取 [7,12]；7>6，新的区间越过共同交集，旧组必须结算。", "发现断开", map[string]string{"start": "7", "end": "12", "right": "6", "arrows": "1"}, tokenRow("区间", intervals, map[int]string{2: "current"}), lane("判断", item("7>6", "rejected"))),
		deepExampleFrame(3, "新增第二支箭并把新组 right 写为 12。", "开启第二组", map[string]string{"right": "12", "arrows": "2"}, tokenRow("区间", intervals, map[int]string{2: "dependency"}), lane("箭", item("6", "ready"), item("12", "current"))),
		deepExampleFrame(2, "读取 [10,16]；10<=12，仍与第二组相交。", "检查最后区间", map[string]string{"start": "10", "end": "16", "right": "12", "arrows": "2"}, tokenRow("区间", intervals, map[int]string{3: "current"}), lane("交集", item("right=12", "dependency"))),
		deepExampleFrame(4, "第二组右端保持 min(12,16)=12；两支箭分别位于 6 和 12。", "稳定第二组", map[string]string{"right": "12", "arrows": "2"}, tokenRow("区间", intervals, map[int]string{3: "ready"}), lane("箭", item("6", "ready"), item("12", "current"))),
		deepExampleFrame(5, "扫描结束，返回 arrows=2。", "区间端点贪心：答案", map[string]string{"arrows": "2", "positions": "6,12"}, lane("答案", item("2 支箭", "current"))),
	}
	return concreteTrace("example-state", "区间端点贪心：最少箭数", code, frames...)
}

func redesignedBFSShortestTrace() Trace {
	code := []string{"queue := []Node{A}; dist[A] = 0", "for head < len(queue) {", "    cur := queue[head]; head++", "    for _, next := range graph[cur] {", "        if unseen(next) { dist[next] = dist[cur]+1; enqueue(next) }", "    }", "}"}
	values := map[string]string{"dist[A]": "0", "queue": "[A]"}
	frames := []Frame{
		nodeFrameDetail(0, "例题图 A→B、A→C、B→D、C→D。起点 A 入队并写 dist[A]=0。", "BFS：无权图最短路", map[string]string{"head": "0", "queue": "[A]"}, graphNodes(map[string]string{"A": "current"}), graphLinks, nil, []string{"dist[A]=0"}, nil, values),
		nodeFrameDetail(2, "head 指向 A，出队 A；现在只处理 A 的出边。", "出队当前节点", map[string]string{"head": "1", "cur": "A", "queue": "[]"}, graphNodes(map[string]string{"A": "current", "B": "dependency", "C": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}, {From: "A", To: "C"}}, []string{"A"}, nil, values),
		nodeFrameDetail(3, "读取边 A→B，B 尚未出现；候选距离是 dist[A]+1=1。", "检查邻居 B", map[string]string{"cur": "A", "next": "B", "candidate": "1"}, graphNodes(map[string]string{"A": "dependency", "B": "current", "C": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}}, []string{"A"}, nil, values),
		nodeFrameDetail(4, "写 dist[B]=1，再把 B 追加到队尾；入队时就标记，避免重复入队。", "写入 B 的距离", map[string]string{"queue": "[B]", "dist[B]": "1"}, graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}}, []string{"B"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "queue": "[B]"}),
		nodeFrameDetail(3, "读取另一条边 A→C；C 未访问，候选距离同样是 1。", "检查邻居 C", map[string]string{"cur": "A", "next": "C", "candidate": "1"}, graphNodes(map[string]string{"A": "dependency", "B": "ready", "C": "current"}), graphLinks, []nodeLink{{From: "A", To: "C"}}, []string{"A"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "queue": "[B]"}),
		nodeFrameDetail(4, "写 dist[C]=1，队列按发现顺序变成 [B,C]。", "写入 C 的距离", map[string]string{"queue": "[B,C]", "dist[C]": "1"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current"}), graphLinks, []nodeLink{{From: "A", To: "C"}}, []string{"C"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "dist[C]": "1", "queue": "[B,C]"}),
		nodeFrameDetail(2, "出队 B；B 的邻居 D 还未访问，准备沿 B→D 扩展。", "处理距离 1 的 B", map[string]string{"cur": "B", "queue": "[C]"}, graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency", "D": "dependency"}), graphLinks, []nodeLink{{From: "B", To: "D"}}, []string{"B"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "dist[C]": "1", "queue": "[C]"}),
		nodeFrameDetail(4, "写 dist[D]=2 并入队。之后即使 C 也连到 D，首次写入仍是最短距离。", "首次到达终点", map[string]string{"dist[D]": "2", "queue": "[C,D]"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "dependency", "D": "current"}), graphLinks, []nodeLink{{From: "B", To: "D"}}, []string{"D"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "dist[C]": "1", "dist[D]": "2", "queue": "[C,D]"}),
		nodeFrameDetail(4, "C 再看到 D 时发现 D 已有距离 2，红色表示冲突候选被拒绝，原值不覆盖。", "拒绝重复到达", map[string]string{"cur": "C", "next": "D", "existing": "2", "queue": "[D]"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current", "D": "rejected"}), graphLinks, []nodeLink{{From: "C", To: "D"}}, []string{"C"}, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "dist[C]": "1", "dist[D]": "2"}),
		nodeFrameDetail(6, "队列为空，dist[D]=2 是从 A 到 D 的最短步数。", "BFS：最终答案", map[string]string{"answer": "2", "dist[D]": "2"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "ready", "D": "current"}), graphLinks, nil, nil, nil, map[string]string{"dist[A]": "0", "dist[B]": "1", "dist[C]": "1", "dist[D]": "2"}),
	}
	return concreteTrace("node-link-state", "BFS：无权图最短路", code, frames...)
}

func redesignedBFSMultiSourceTrace() Trace {
	code := []string{"enqueue every source with dist=0", "for head < len(queue) {", "    cur := queue[head]; head++", "    for _, next := range neighbors(cur) {", "        if dist[next] == -1 { dist[next] = dist[cur]+1; enqueue(next) }", "    }", "}"}
	rows := []string{"0..", "...", "..0"}
	frames := []Frame{
		matrixFrameDetail(0, "例题 0 在 (0,0) 与 (2,2)。先把两个源同时写为距离 0 并入队。", "多源 BFS：同时建立第 0 层", map[string]string{"queue": "[(0,0),(2,2)]", "layer": "0"}, 3, 3, deepCells([]string{"0..", "...", "..0"}, map[string]string{"0:0": "current", "2:2": "current"})),
		matrixFrameDetail(2, "出队 (0,0)，只看右邻居与下邻居；它们是下一层的候选。", "扩散源 (0,0)", map[string]string{"cur": "(0,0)", "queue": "[(2,2)]"}, 3, 3, deepCells(rows, map[string]string{"0:0": "current", "0:1": "dependency", "1:0": "dependency", "2:2": "ready"})),
		matrixFrameDetail(4, "(0,1) 尚未写入，读取源距离 0，候选值为 1。", "计算 (0,1)", map[string]string{"cur": "(0,0)", "next": "(0,1)", "candidate": "1"}, 3, 3, deepCells([]string{"01.", "...", "..0"}, map[string]string{"0:1": "current", "0:0": "dependency"})),
		matrixFrameDetail(4, "写 distance(0,1)=1 并入队；一次只写一个新格，便于观察扩散前沿。", "写入第一格", map[string]string{"next": "(0,1)", "distance": "1", "queue": "[(2,2),(0,1)]"}, 3, 3, deepCells([]string{"01.", "...", "..0"}, map[string]string{"0:1": "current", "0:0": "ready"})),
		matrixFrameDetail(4, "同一源继续处理 (1,0)，写入距离 1。", "写入第二格", map[string]string{"next": "(1,0)", "distance": "1", "queue": "[(2,2),(0,1),(1,0)]"}, 3, 3, deepCells([]string{"01.", "1..", "..0"}, map[string]string{"1:0": "current", "0:0": "ready"})),
		matrixFrameDetail(2, "再出队右下源 (2,2)，上方与左方是它的扩散候选。", "扩散源 (2,2)", map[string]string{"cur": "(2,2)", "layer": "0"}, 3, 3, deepCells([]string{"01.", "1..", "..0"}, map[string]string{"2:2": "current", "1:2": "dependency", "2:1": "dependency"})),
		matrixFrameDetail(4, "写 distance(1,2)=1；两个源的前沿从两侧同时推进。", "写入右侧前沿", map[string]string{"next": "(1,2)", "distance": "1"}, 3, 3, deepCells([]string{"01.", "1.1", "..0"}, map[string]string{"1:2": "current", "2:2": "dependency"})),
		matrixFrameDetail(4, "写 distance(2,1)=1；未访问的格子只会被第一次到达。", "写入下侧前沿", map[string]string{"next": "(2,1)", "distance": "1"}, 3, 3, deepCells([]string{"01.", "1.1", ".10"}, map[string]string{"2:1": "current", "2:2": "dependency"})),
		matrixFrameDetail(4, "四个距离 1 的格子成为新的队列层，中心格 (1,1) 等待它们扩散。", "进入距离 1 层", map[string]string{"queue": "[(0,1),(1,0),(1,2),(2,1)]", "layer": "1"}, 3, 3, deepCells([]string{"01.", "1.1", ".10"}, map[string]string{"0:1": "ready", "1:0": "ready", "1:2": "ready", "2:1": "ready"})),
		matrixFrameDetail(4, "从任一距离 1 邻居到达中心，候选距离都是 2；只写入一次。", "写入中心", map[string]string{"next": "(1,1)", "candidate": "2"}, 3, 3, deepCells([]string{"012", "121", "210"}, map[string]string{"1:1": "current", "0:1": "dependency", "1:0": "dependency", "1:2": "dependency", "2:1": "dependency"})),
		matrixFrameDetail(6, "所有格子的第一次写入都来自最近源，最终距离矩阵为 012/121/210。", "多源 BFS：最终距离", map[string]string{"answer": "中心距离 2", "layer": "2"}, 3, 3, deepCells([]string{"012", "121", "210"}, map[string]string{"1:1": "current"})),
	}
	return concreteTrace("matrix-state", "多源 BFS：最近源距离", code, frames...)
}

func redesignedBFSTopologicalTrace() Trace {
	code := []string{"enqueue every indegree-0 node", "for head < len(queue) {", "    v := queue[head]; head++", "    for _, to := range graph[v] {", "        indegree[to]--", "        if indegree[to] == 0 { enqueue(to) }", "    }", "}"}
	values := map[string]string{"indegree": "A:0 B:1 C:1 D:2", "queue": "[A]", "order": "[]"}
	frames := []Frame{
		nodeFrameDetail(0, "例题依赖 A→B、A→C、B→D、C→D。只有 A 入度为 0，先进入队列。", "Kahn：初始化", map[string]string{"queue": "[A]", "order": "[]"}, graphNodes(map[string]string{"A": "current"}), graphLinks, nil, nil, nil, values),
		nodeFrameDetail(2, "取出 A；A 加入拓扑序，接下来逐条删除 A 的出边。", "出队 A", map[string]string{"v": "A", "queue": "[]", "order": "[A]"}, graphNodes(map[string]string{"A": "current", "B": "dependency", "C": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}, {From: "A", To: "C"}}, []string{"A"}, nil, values),
		nodeFrameDetail(4, "删除 A→B：B 的入度从 1 减到 0，B 现在变成可执行节点。", "释放 B", map[string]string{"edge": "A→B", "indegree[B]": "0", "queue": "[B]"}, graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}}, []string{"A"}, nil, map[string]string{"indegree": "A:0 B:0 C:1 D:2", "queue": "[B]", "order": "[A]"}),
		nodeFrameDetail(4, "删除 A→C：C 的入度也从 1 减到 0，追加到 B 后面。", "释放 C", map[string]string{"edge": "A→C", "indegree[C]": "0", "queue": "[B,C]"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current"}), graphLinks, []nodeLink{{From: "A", To: "C"}}, []string{"A"}, nil, map[string]string{"indegree": "A:0 B:0 C:0 D:2", "queue": "[B,C]", "order": "[A]"}),
		nodeFrameDetail(2, "出队 B，拓扑序变为 [A,B]；B→D 只删除 D 的一条前置边。", "出队 B", map[string]string{"v": "B", "queue": "[C]", "order": "[A,B]", "indegree[D]": "2"}, graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency", "D": "dependency"}), graphLinks, []nodeLink{{From: "B", To: "D"}}, []string{"A", "B"}, nil, map[string]string{"indegree": "B:0 C:0 D:2", "queue": "[C]", "order": "[A,B]"}),
		nodeFrameDetail(4, "删除 B→D：D 的入度 2→1，仍不能入队，因为 C 还未处理。", "保留 D 的前置约束", map[string]string{"edge": "B→D", "indegree[D]": "1", "queue": "[C]"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "dependency", "D": "dependency"}), graphLinks, []nodeLink{{From: "B", To: "D"}}, []string{"A", "B"}, nil, map[string]string{"indegree": "C:0 D:1", "queue": "[C]", "order": "[A,B]"}),
		nodeFrameDetail(2, "出队 C，删除 C→D 后 D 的入度 1→0，所有前置课程都已完成。", "释放 D", map[string]string{"edge": "C→D", "indegree[D]": "0", "queue": "[D]", "order": "[A,B,C]"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current", "D": "current"}), graphLinks, []nodeLink{{From: "C", To: "D"}}, []string{"A", "B", "C"}, nil, map[string]string{"indegree": "D:0", "queue": "[D]", "order": "[A,B,C]"}),
		nodeFrameDetail(7, "D 出队后所有 4 个节点都进入拓扑序；处理数量等于节点总数，图无环。", "Kahn：完成", map[string]string{"queue": "[]", "order": "[A,B,C,D]", "answer": "true"}, graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "ready", "D": "current"}), graphLinks, nil, nil, nil, map[string]string{"indegree": "全部 0", "order": "[A,B,C,D]"}),
	}
	return concreteTrace("node-link-state", "BFS：Kahn 拓扑排序", code, frames...)
}

func redesignedDFSTreeTrace() Trace {
	code := []string{"if node == nil { return 0 }", "left := dfs(node.Left)", "right := dfs(node.Right)", "return left + right + node.Val"}
	frames := []Frame{
		nodeFrameDetail(0, "例题对树求和。进入根 3，函数还没有左、右返回值，不能提前计算。", "DFS 树：进入根", map[string]string{"node": "3", "left": "?", "right": "?"}, treeNodes(map[string]string{"3": "current"}), treeLinks, nil, []string{"dfs(3)"}, nil, nil),
		nodeFrameDetail(1, "执行 left := dfs(3.Left)，沿 3→5 进入左子树；橙色只表示当前调用，蓝色表示等待返回。", "递归左子树", map[string]string{"node": "3", "child": "5", "return": "等待"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, []nodeLink{{From: "3", To: "5"}}, []string{"dfs(3)", "dfs(5)"}, nil, nil),
		nodeFrameDetail(1, "在 dfs(5) 中同样先递归左孩子 6，调用栈继续增长。", "深入到 6", map[string]string{"node": "5", "child": "6", "return": "等待"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, []nodeLink{{From: "5", To: "6"}}, []string{"dfs(3)", "dfs(5)", "dfs(6)"}, nil, nil),
		nodeFrameDetail(0, "dfs(6) 的左孩子为空；空节点返回 0，但树上的 6 仍保持当前调用。", "空左孩子返回", map[string]string{"node": "6", "child": "nil", "return": "0"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)", "dfs(6)"}, nil, map[string]string{"nil.left": "0"}),
		nodeFrameDetail(0, "dfs(6) 的右孩子同样为空，第二个依赖也返回 0。", "空右孩子返回", map[string]string{"node": "6", "left": "0", "right": "0"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)", "dfs(6)"}, nil, map[string]string{"left(6)": "0", "right(6)": "0"}),
		nodeFrameDetail(3, "两个孩子都返回后，后序组合 sum(6)=0+0+6=6；结果挂在节点 6 上。", "组合节点 6", map[string]string{"node": "6", "left": "0", "right": "0", "return": "6"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)", "dfs(6)"}, nil, map[string]string{"sum(6)": "6"}),
		nodeFrameDetail(1, "返回到 dfs(5)，读取 left=6；5 还需要处理右孩子 2。", "返回到 5", map[string]string{"node": "5", "left": "6", "right": "?"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "dependency", "2": "dependency"}), treeLinks, []nodeLink{{From: "5", To: "2"}}, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"sum(6)": "6"}),
		nodeFrameDetail(1, "沿 5→2 进入右孩子；调用栈显示递归不是一次跳到答案，而是逐层等待。", "深入到 2", map[string]string{"node": "2", "return": "等待"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "2": "current"}), treeLinks, []nodeLink{{From: "5", To: "2"}}, []string{"dfs(3)", "dfs(5)", "dfs(2)"}, nil, map[string]string{"sum(6)": "6"}),
		nodeFrameDetail(3, "2 的两个空孩子各返回 0，组合后 sum(2)=2。", "组合节点 2", map[string]string{"node": "2", "left": "0", "right": "0", "return": "2"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "2": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)", "dfs(2)"}, nil, map[string]string{"sum(6)": "6", "sum(2)": "2"}),
		nodeFrameDetail(3, "5 已拿到左右结果 6、2，后序组合 sum(5)=6+2+5=13，并向根返回。", "组合节点 5", map[string]string{"node": "5", "left": "6", "right": "2", "return": "13"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "ready", "2": "ready"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"sum(6)": "6", "sum(2)": "2", "sum(5)": "13"}),
		nodeFrameDetail(1, "根 3 读取 left=13，继续沿 3→1 递归右子树。", "根等待右子树", map[string]string{"node": "3", "left": "13", "right": "?"}, treeNodes(map[string]string{"3": "dependency", "1": "current", "5": "ready"}), treeLinks, []nodeLink{{From: "3", To: "1"}}, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"sum(5)": "13"}),
		nodeFrameDetail(1, "节点 1 先处理左孩子 0；叶子返回 sum(0)=0。", "处理节点 0", map[string]string{"node": "1", "child": "0", "return": "0"}, treeNodes(map[string]string{"3": "dependency", "1": "dependency", "0": "current"}), treeLinks, []nodeLink{{From: "1", To: "0"}}, []string{"dfs(3)", "dfs(1)", "dfs(0)"}, nil, map[string]string{"sum(5)": "13"}),
		nodeFrameDetail(1, "节点 1 再处理右孩子 8；返回 sum(8)=8。", "处理节点 8", map[string]string{"node": "1", "left": "0", "child": "8", "return": "8"}, treeNodes(map[string]string{"3": "dependency", "1": "dependency", "0": "ready", "8": "current"}), treeLinks, []nodeLink{{From: "1", To: "8"}}, []string{"dfs(3)", "dfs(1)", "dfs(8)"}, nil, map[string]string{"sum(5)": "13", "sum(0)": "0"}),
		nodeFrameDetail(3, "1 的两个孩子结果为 0、8，组合 sum(1)=0+8+1=9。", "组合节点 1", map[string]string{"node": "1", "left": "0", "right": "8", "return": "9"}, treeNodes(map[string]string{"3": "dependency", "1": "current", "0": "ready", "8": "ready"}), treeLinks, nil, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"sum(1)": "9"}),
		nodeFrameDetail(3, "根拿到 left=13、right=9，最后才组合 sum(3)=13+9+3=25。", "DFS 树：最终返回", map[string]string{"node": "3", "left": "13", "right": "9", "answer": "25"}, treeNodes(map[string]string{"3": "current", "5": "ready", "1": "ready", "6": "ready", "2": "ready", "0": "ready", "8": "ready"}), treeLinks, nil, []string{"dfs(3)"}, nil, map[string]string{"sum(3)": "25"}),
	}
	return concreteTrace("node-link-state", "DFS：树的后序组合", code, frames...)
}

func redesignedDFSGridTrace() Trace {
	code := []string{"if out of bounds or grid[r][c] != '1' { return }", "grid[r][c] = '0'", "for _, dir := range dirs {", "    dfs(r+dir[0], c+dir[1])", "}"}
	frames := []Frame{
		matrixFrameDetail(0, "例题 grid=110/010/011。外层第一次发现 (0,0)=1，岛屿计数加 1，准备启动 DFS。", "DFS 网格：启动一个连通块", map[string]string{"start": "(0,0)", "islands": "1"}, 3, 3, deepCells([]string{"110", "010", "011"}, map[string]string{"0:0": "current"})),
		matrixFrameDetail(0, "检查 (0,0) 在边界内且是陆地，允许进入递归。", "检查当前格", map[string]string{"cell": "(0,0)", "value": "1"}, 3, 3, deepCells([]string{"110", "010", "011"}, map[string]string{"0:0": "current"})),
		matrixFrameDetail(1, "进入 (0,0) 立刻写成 0；这是写访问标记，防止从邻居绕回。", "标记 (0,0)", map[string]string{"cell": "(0,0)", "grid[0][0]": "0"}, 3, 3, deepCells([]string{"010", "010", "011"}, map[string]string{"0:0": "current"})),
		matrixFrameDetail(2, "枚举第一个方向：向右检查 (0,1)，它仍是未访问陆地，作为下一次递归入口。", "枚举右邻居", map[string]string{"from": "(0,0)", "next": "(0,1)", "dir": "right"}, 3, 3, deepCells([]string{"010", "010", "011"}, map[string]string{"0:0": "dependency", "0:1": "current"})),
		matrixFrameDetail(1, "写 (0,1)=0，然后继续从它的四个方向扩展。", "标记 (0,1)", map[string]string{"cell": "(0,1)", "grid[0][1]": "0"}, 3, 3, deepCells([]string{"000", "010", "011"}, map[string]string{"0:1": "current"})),
		matrixFrameDetail(0, "从 (0,1) 向上越界，返回；越界调用不改变网格。", "边界调用返回", map[string]string{"cell": "(-1,1)", "reason": "out of bounds"}, 3, 3, deepCells([]string{"000", "010", "011"}, map[string]string{"0:1": "dependency"})),
		matrixFrameDetail(2, "从 (0,1) 向下进入 (1,1)，它是同一连通块的下一格。", "深入 (1,1)", map[string]string{"from": "(0,1)", "next": "(1,1)"}, 3, 3, deepCells([]string{"000", "010", "011"}, map[string]string{"0:1": "dependency", "1:1": "current"})),
		matrixFrameDetail(1, "标记 (1,1)=0；之后从 (1,1) 回看上方时会命中 0 并返回。", "标记 (1,1)", map[string]string{"cell": "(1,1)", "grid[1][1]": "0"}, 3, 3, deepCells([]string{"000", "000", "011"}, map[string]string{"1:1": "current"})),
		matrixFrameDetail(2, "继续向下进入 (2,1)，把这条分支接到同一个 DFS 调用栈。", "深入 (2,1)", map[string]string{"from": "(1,1)", "next": "(2,1)"}, 3, 3, deepCells([]string{"000", "000", "011"}, map[string]string{"1:1": "dependency", "2:1": "current"})),
		matrixFrameDetail(1, "标记 (2,1)=0；向右的 (2,2) 仍是陆地。", "标记 (2,1)", map[string]string{"cell": "(2,1)", "grid[2][1]": "0"}, 3, 3, deepCells([]string{"000", "000", "001"}, map[string]string{"2:1": "current"})),
		matrixFrameDetail(2, "从 (2,1) 向右进入 (2,2)，这是本岛最后一个未访问格。", "深入 (2,2)", map[string]string{"from": "(2,1)", "next": "(2,2)"}, 3, 3, deepCells([]string{"000", "000", "001"}, map[string]string{"2:1": "dependency", "2:2": "current"})),
		matrixFrameDetail(1, "标记 (2,2)=0；它的四个邻居随后都只会触发返回。", "标记 (2,2)", map[string]string{"cell": "(2,2)", "grid[2][2]": "0"}, 3, 3, deepCells([]string{"000", "000", "000"}, map[string]string{"2:2": "current"})),
		matrixFrameDetail(0, "递归逐层返回，所有相邻陆地都被标记；外层扫描不会再次计数这座岛。", "连通块完成", map[string]string{"islands": "1", "remaining": "0"}, 3, 3, deepCells([]string{"000", "000", "000"}, map[string]string{"0:0": "ready", "0:1": "ready", "1:1": "ready", "2:1": "ready", "2:2": "ready"})),
	}
	return concreteTrace("matrix-state", "DFS：网格连通块", code, frames...)
}

func redesignedDFSPathTrace() Trace {
	code := []string{"if node == target { collect copy(path); return }", "for _, next := range graph[node] {", "    path = append(path, next)", "    dfs(next)", "    path = path[:len(path)-1]", "}"}
	frames := []Frame{
		nodeFrameDetail(1, "例题枚举 A 到 D 的所有路径。先把起点 A 放入 path，再进入 dfs(A)。", "路径 DFS：初始化", map[string]string{"node": "A", "path": "[A]", "answers": "[]"}, graphNodes(map[string]string{"A": "current"}), graphLinks, nil, []string{"dfs(A)"}, []string{"A"}, nil),
		nodeFrameDetail(1, "枚举 A 的第一个邻居 B；边 A→B 是当前候选，尚未写入 path。", "选择 B 前", map[string]string{"node": "A", "next": "B", "path": "[A]"}, graphNodes(map[string]string{"A": "current", "B": "dependency"}), graphLinks, []nodeLink{{From: "A", To: "B"}}, []string{"dfs(A)"}, []string{"A"}, nil),
		nodeFrameDetail(2, "执行 append，把 B 写入 path=[A,B]。", "选择 B", map[string]string{"path": "[A,B]", "next": "B"}, graphNodes(map[string]string{"A": "dependency", "B": "current"}), graphLinks, []nodeLink{{From: "A", To: "B"}}, []string{"dfs(A)", "dfs(B)"}, []string{"A", "B"}, nil),
		nodeFrameDetail(1, "在 B 的邻居中选择 D，path 先扩展成 [A,B,D]。", "深入 D", map[string]string{"path": "[A,B,D]", "next": "D"}, graphNodes(map[string]string{"A": "dependency", "B": "dependency", "D": "current"}), graphLinks, []nodeLink{{From: "B", To: "D"}}, []string{"dfs(A)", "dfs(B)", "dfs(D)"}, []string{"A", "B", "D"}, nil),
		nodeFrameDetail(0, "命中 target=D，复制 path 到答案；答案必须是独立切片。", "收集第一条路径", map[string]string{"path": "[A,B,D]", "answers": "[[A,B,D]]"}, graphNodes(map[string]string{"A": "dependency", "B": "dependency", "D": "current"}), graphLinks, nil, []string{"dfs(A)", "dfs(B)", "dfs(D)"}, []string{"A", "B", "D"}, map[string]string{"answer": "[A,B,D]"}),
		nodeFrameDetail(4, "从 D 返回 B，执行 pop；path 恢复为 [A,B]，给 B 的下一个兄弟分支留出现场。", "撤销 D", map[string]string{"path": "[A,B]", "answers": "[[A,B,D]]"}, graphNodes(map[string]string{"A": "dependency", "B": "current", "D": "ready"}), graphLinks, nil, []string{"dfs(A)", "dfs(B)"}, []string{"A", "B"}, nil),
		nodeFrameDetail(4, "B 的邻居已经探索完，继续 pop B；回到 A 时 path 只剩 [A]。", "撤销 B", map[string]string{"path": "[A]", "answers": "[[A,B,D]]"}, graphNodes(map[string]string{"A": "current", "B": "ready", "D": "ready"}), graphLinks, nil, []string{"dfs(A)"}, []string{"A"}, nil),
		nodeFrameDetail(1, "A 的第二个邻居 C 成为候选；重新 append 得到 [A,C]，上一条路径不会污染这一分支。", "选择 C", map[string]string{"path": "[A,C]", "next": "C"}, graphNodes(map[string]string{"A": "dependency", "C": "current"}), graphLinks, []nodeLink{{From: "A", To: "C"}}, []string{"dfs(A)", "dfs(C)"}, []string{"A", "C"}, nil),
		nodeFrameDetail(2, "C 只有一个下一步 D，写入 path=[A,C,D]。", "深入第二个 D", map[string]string{"path": "[A,C,D]"}, graphNodes(map[string]string{"A": "dependency", "C": "dependency", "D": "current"}), graphLinks, []nodeLink{{From: "C", To: "D"}}, []string{"dfs(A)", "dfs(C)", "dfs(D)"}, []string{"A", "C", "D"}, nil),
		nodeFrameDetail(0, "再次命中 D，复制第二条路径；此时答案有两条互不共享的记录。", "收集第二条路径", map[string]string{"path": "[A,C,D]", "answers": "[[A,B,D],[A,C,D]]"}, graphNodes(map[string]string{"A": "dependency", "C": "dependency", "D": "current"}), graphLinks, nil, []string{"dfs(A)", "dfs(C)", "dfs(D)"}, []string{"A", "C", "D"}, map[string]string{"answer": "[A,B,D],[A,C,D]"}),
		nodeFrameDetail(4, "回溯 D、C 后 path 恢复为 [A]，所有邻居都处理完，最终返回全部路径。", "路径 DFS：完成", map[string]string{"path": "[A]", "answers": "2 条"}, graphNodes(map[string]string{"A": "current", "B": "ready", "C": "ready", "D": "ready"}), graphLinks, nil, []string{"dfs(A)"}, []string{"A"}, map[string]string{"answer": "[[A,B,D],[A,C,D]"}),
	}
	return concreteTrace("node-link-state", "DFS：路径枚举与回溯", code, frames...)
}

func redesignedChooseSkipTrace() Trace {
	code := []string{"if index == len(nums) { collect copy(path); return }", "dfs(index + 1) // 不选", "path = append(path, nums[index])", "dfs(index + 1) // 选", "path = path[:len(path)-1]"}
	nums := []string{"1", "2"}
	frames := []Frame{
		deepExampleFrame(0, "例题 nums=[1,2]。每层只决定一个位置；path 是当前递归分支的可变现场。", "回溯：选或不选", map[string]string{"index": "0", "path": "[]", "answers": "[]"}, tokenRow("nums", nums, map[int]string{0: "current"}), lane("决策", item("index 0", "current")), lane("path", item("[]", "ready"))),
		deepExampleFrame(1, "index=0 先走“不选 1”边；path 不变，递归进入 index=1。", "不选 1", map[string]string{"index": "1", "choice": "skip 1", "path": "[]"}, tokenRow("nums", nums, map[int]string{0: "rejected"}), lane("path", item("[]", "current"))),
		deepExampleFrame(1, "index=1 再走“不选 2”边；所有位置已决定，准备收集空集。", "不选 2", map[string]string{"index": "2", "choice": "skip 2", "path": "[]"}, tokenRow("nums", nums, map[int]string{1: "rejected"}), lane("path", item("[]", "current"))),
		deepExampleFrame(0, "到达叶子 index=2，复制 path=[]，答案新增空集。", "收集空集", map[string]string{"index": "2", "answers": "[[]]"}, lane("答案", item("[]", "current")), lane("path", item("[]", "dependency"))),
		deepExampleFrame(2, "返回 index=1，改走“选择 2”边，先写 path=[2]。", "选择 2", map[string]string{"index": "1", "choice": "take 2", "path": "[2]"}, tokenRow("nums", nums, map[int]string{1: "current"}), lane("path", item("2", "current"))),
		deepExampleFrame(1, "选择 2 后进入叶子；因为 index 已等于 2，收集 [2]。", "收集 [2]", map[string]string{"index": "2", "answers": "[[],[2]"}, lane("答案", item("[]", "ready"), item("[2]", "current")), lane("path", item("[2]", "dependency"))),
		deepExampleFrame(4, "从 index=1 返回时 pop 2，path 恢复为空；这是尝试另一条边的前提。", "撤销 2", map[string]string{"index": "1", "path": "[]"}, lane("path", item("[]", "current")), lane("答案", item("[] [2]", "ready"))),
		deepExampleFrame(2, "回到 index=0，执行“选择 1”，写 path=[1]，然后递归决定 2。", "选择 1", map[string]string{"index": "0", "choice": "take 1", "path": "[1]"}, tokenRow("nums", nums, map[int]string{0: "current"}), lane("path", item("1", "current"))),
		deepExampleFrame(1, "在 path=[1] 下不选 2，叶子收集 [1]。", "不选 2", map[string]string{"index": "2", "path": "[1]", "answers": "[[],[2],[1]]"}, tokenRow("nums", nums, map[int]string{1: "rejected"}), lane("答案", item("[1]", "current")), lane("path", item("1", "dependency"))),
		deepExampleFrame(2, "回到 index=1，选择 2，把 path 扩成 [1,2]。", "选择 2", map[string]string{"index": "1", "path": "[1,2]"}, tokenRow("nums", nums, map[int]string{1: "current"}), lane("path", item("1", "ready"), item("2", "current"))),
		deepExampleFrame(0, "叶子收集 [1,2]；四个叶子正好对应两层二叉决策树。", "收集 [1,2]", map[string]string{"index": "2", "answers": "[[],[2],[1],[1,2]]"}, lane("答案", item("[]", "ready"), item("[2]", "ready"), item("[1]", "ready"), item("[1,2]", "current"))),
		deepExampleFrame(4, "最后撤销 2、1，path 回到 []；回溯结束，现场没有泄漏到调用者。", "回溯完成", map[string]string{"index": "0", "path": "[]", "answerCount": "4"}, lane("path", item("[]", "current")), lane("答案", item("[] [2] [1] [1,2]", "ready"))),
	}
	return concreteTrace("example-state", "回溯：选或不选子集", code, frames...)
}

func redesignedEnumerationTrace() Trace {
	code := []string{"if len(path) == len(nums) { collect copy(path); return }", "for i, value := range nums {", "    if used[i] { continue }", "    used[i] = true; path = append(path, value)", "    dfs(); path = path[:len(path)-1]; used[i] = false", "}"}
	values := []string{"1", "2", "3"}
	frames := []Frame{
		deepExampleFrame(0, "例题 nums=[1,2,3] 求排列。path 为空，used=[F,F,F]。", "排列：建立搜索树", map[string]string{"path": "[]", "answers": "0"}, tokenRow("候选", values, nil), lane("used", item("F", "ready"), item("F", "ready"), item("F", "ready"))),
		deepExampleFrame(1, "第 0 层循环先读 i=0、value=1；used[0]=false，可以选择。", "枚举候选 1", map[string]string{"i": "0", "value": "1", "path": "[]"}, tokenRow("候选", values, map[int]string{0: "current"}), lane("used", item("F", "current"), item("F", "ready"), item("F", "ready"))),
		deepExampleFrame(3, "写 used[0]=true 并 append 1；进入下一层 path=[1]。", "选择 1", map[string]string{"path": "[1]", "used": "T,F,F"}, tokenRow("候选", values, map[int]string{0: "current"}), lane("used", item("T", "current"), item("F", "ready"), item("F", "ready")), lane("path", item("1", "current"))),
		deepExampleFrame(2, "下一层再次看到 i=0，但 used[0]=true；红色表示跳过，不能重复使用 1。", "跳过已用 1", map[string]string{"i": "0", "path": "[1]"}, tokenRow("候选", values, map[int]string{0: "rejected", 1: "dependency", 2: "dependency"}), lane("path", item("1", "ready"))),
		deepExampleFrame(3, "选择 2，写 path=[1,2]、used=[T,T,F]。", "选择 2", map[string]string{"path": "[1,2]", "used": "T,T,F"}, tokenRow("候选", values, map[int]string{1: "current"}), lane("used", item("T", "ready"), item("T", "current"), item("F", "ready")), lane("path", item("1", "ready"), item("2", "current"))),
		deepExampleFrame(3, "最后一个未使用候选是 3，写入后 path=[1,2,3] 达到长度 3。", "补齐排列", map[string]string{"path": "[1,2,3]", "used": "T,T,T"}, tokenRow("候选", values, map[int]string{2: "current"}), lane("path", item("1", "ready"), item("2", "ready"), item("3", "current"))),
		deepExampleFrame(0, "到达递归终点，复制 [1,2,3] 到答案。", "收集排列", map[string]string{"path": "[1,2,3]", "answers": "[[1,2,3]]"}, lane("答案", item("[1,2,3]", "current")), lane("path", item("[1,2,3]", "dependency"))),
		deepExampleFrame(4, "返回一层，pop 3 并恢复 used[2]=false；path 回到 [1,2]。", "撤销 3", map[string]string{"path": "[1,2]", "used": "T,T,F"}, lane("path", item("[1,2]", "current")), lane("used", item("T", "ready"), item("T", "ready"), item("F", "current"))),
		deepExampleFrame(4, "继续回退 2，恢复 path=[1]、used=[T,F,F]；现在可以尝试 3。", "回退到兄弟分支", map[string]string{"path": "[1]", "used": "T,F,F"}, lane("path", item("[1]", "current")), lane("used", item("T", "ready"), item("F", "current"), item("F", "ready"))),
		deepExampleFrame(3, "选择 3，深入 path=[1,3]；下一层选择 2 得到第二条排列。", "交换尾部顺序", map[string]string{"path": "[1,3,2]", "answers": "2"}, tokenRow("候选", values, map[int]string{1: "current", 2: "dependency"}), lane("path", item("1", "ready"), item("3", "ready"), item("2", "current"))),
		deepExampleFrame(4, "收集并回退后，第 0 层解除 1；循环才会尝试开头为 2 的分支。", "完成 1 开头分支", map[string]string{"path": "[]", "used": "F,F,F", "answers": "2"}, lane("answers", item("[1,2,3] [1,3,2]", "ready")), lane("used", item("F", "current"), item("F", "ready"), item("F", "ready"))),
		deepExampleFrame(3, "依次选择 2、1、3 与 2、3、1；每条叶子都经历选择、递归、撤销。", "尝试 2 开头", map[string]string{"path": "[2,1,3]", "answers": "4"}, tokenRow("候选", values, map[int]string{2: "current"}), lane("path", item("[2,1,3]", "current"))),
		deepExampleFrame(3, "再尝试 3 开头的两条排列，答案数量达到 6。", "完成全部排列", map[string]string{"path": "[3,2,1]", "answers": "6"}, tokenRow("候选", values, map[int]string{0: "ready", 1: "current", 2: "dependency"}), lane("答案", item("6 个排列", "current"))),
		deepExampleFrame(5, "最后一层恢复 used 与 path；所有分支结束，搜索树没有残留状态。", "回溯完成", map[string]string{"path": "[]", "used": "F,F,F", "answers": "6"}, lane("path", item("[]", "current")), lane("used", item("F,F,F", "ready"))),
	}
	return concreteTrace("example-state", "回溯：枚举下一个候选", code, frames...)
}

func cycleListNodes(states map[string]string) []nodeVisual {
	base := []nodeVisual{
		{ID: "1", Label: "1", X: 42, Y: 92},
		{ID: "2", Label: "2", X: 128, Y: 42},
		{ID: "3", Label: "3", X: 232, Y: 42},
		{ID: "4", Label: "4", X: 318, Y: 92},
	}
	return withNodes(base, states)
}

var cycleListLinks = []nodeLink{{From: "1", To: "2"}, {From: "2", To: "3"}, {From: "3", To: "4"}, {From: "4", To: "2"}}

func redesignedFastSlowTrace() Trace {
	code := []string{"slow, fast := head, head", "for fast != nil && fast.Next != nil {", "    slow = slow.Next", "    fast = fast.Next.Next", "    if slow == fast { return true }", "}", "return false"}
	frames := []Frame{
		cycleListFrame(0, "例题 1→2→3→4，4→2。把真实的 Next 边画出来，环不是把节点 2 复制到尾部。", "快慢指针：真实环结构", map[string]string{"slow": "1", "fast": "1"}, cycleListNodes(map[string]string{"1": "current"}), cycleListLinks, map[string]string{"slow": "1", "fast": "1"}, nil),
		cycleListFrame(1, "fast=1 且 fast.Next=2 都非空，循环条件成立；这一轮会同时移动两个指针。", "检查循环条件", map[string]string{"slow": "1", "fast": "1", "next": "2"}, cycleListNodes(map[string]string{"1": "current", "2": "dependency"}), cycleListLinks, map[string]string{"slow": "1", "fast": "1"}, nil),
		cycleListFrame(3, "一步同时移动：慢指针沿 1→2，快指针沿 1→2→3；两者落在 2 和 3。", "第 1 轮：同时移动", map[string]string{"slow": "2", "fast": "3", "path": "slow:1→2；fast:1→2→3"}, cycleListNodes(map[string]string{"2": "current", "3": "current"}), cycleListLinks, map[string]string{"slow": "2", "fast": "3"}, nil),
		cycleListFrame(4, "比较慢=2 与快=3，不相等；这次比较读取的两个节点保留蓝色。", "第 1 轮：比较", map[string]string{"slow": "2", "fast": "3", "equal": "false"}, cycleListNodes(map[string]string{"2": "dependency", "3": "rejected"}), cycleListLinks, map[string]string{"slow": "2", "fast": "3"}, nil),
		cycleListFrame(3, "一步同时移动：慢指针 2→3，快指针 3→4→2；快指针沿回边回到 2。", "第 2 轮：同时移动", map[string]string{"slow": "3", "fast": "2", "path": "slow:2→3；fast:3→4→2"}, cycleListNodes(map[string]string{"2": "current", "3": "current", "4": "dependency"}), cycleListLinks, map[string]string{"slow": "3", "fast": "2"}, nil),
		cycleListFrame(4, "比较慢=3 与快=2，仍不相等；访问过不代表相遇，必须是同一时刻同一节点。", "第 2 轮：比较", map[string]string{"slow": "3", "fast": "2", "equal": "false"}, cycleListNodes(map[string]string{"2": "rejected", "3": "dependency"}), cycleListLinks, map[string]string{"slow": "3", "fast": "2"}, nil),
		cycleListFrame(3, "一步同时移动：慢指针 3→4，快指针 2→3→4；两者现在都落在 4。", "第 3 轮：同时移动", map[string]string{"slow": "4", "fast": "4", "path": "slow:3→4；fast:2→3→4"}, cycleListNodes(map[string]string{"4": "current"}), cycleListLinks, map[string]string{"slow": "4", "fast": "4"}, nil),
		cycleListFrame(4, "读取 slow=4、fast=4，条件成立；相遇只在比较帧确认。", "第 3 轮：确认相遇", map[string]string{"slow": "4", "fast": "4", "equal": "true"}, cycleListNodes(map[string]string{"4": "current"}), cycleListLinks, map[string]string{"slow": "4", "fast": "4"}, nil),
		cycleListFrame(4, "slow==fast，返回 true；图中始终保留同一套节点和 Next 边，只移动带文字的快、慢指针。", "快慢指针：检测到环", map[string]string{"slow": "4", "fast": "4", "answer": "true"}, cycleListNodes(map[string]string{"4": "current"}), cycleListLinks, map[string]string{"slow": "4", "fast": "4"}, nil),
	}
	return concreteTrace("cycle-list-state", "链表：快慢指针判环", code, frames...)
}

func mergeItems(values []string, active string) []exampleItem {
	items := make([]exampleItem, len(values))
	for index, value := range values {
		state := "ready"
		if index == 0 && value == active {
			state = "current"
		}
		items[index] = item(value, state)
	}
	return items
}

func redesignedMergeListsTrace() Trace {
	code := []string{"dummy := new(ListNode); tail := dummy", "for a != nil && b != nil {", "    choose smaller head", "    tail.Next = chosen; tail = tail.Next", "}", "tail.Next = remaining"}
	caption := "链表：dummy 尾插合并"
	frames := []Frame{
		mergeListFrame(0, "例题 A=1→3→5，B=2→4→6。dummy 固定结果头，tail 先指向 dummy。", caption, map[string]string{"tail": "dummy"}, mergeItems([]string{"1", "3", "5"}, ""), mergeItems([]string{"2", "4", "6"}, ""), nil, "dummy", ""),
		mergeListFrame(2, "读取 A、B 两条链的当前头 1 和 2；它们都是候选，先比较值。", caption, map[string]string{"a": "1", "b": "2", "tail": "dummy"}, mergeItems([]string{"1", "3", "5"}, "1"), mergeItems([]string{"2", "4", "6"}, "2"), nil, "dummy", ""),
		mergeListFrame(3, "1 更小，选择 A 的 1，但还没改动输入头；当前写入对象是 tail.Next。", caption, map[string]string{"chosen": "A:1", "tail": "dummy"}, mergeItems([]string{"1", "3", "5"}, "1"), mergeItems([]string{"2", "4", "6"}, "2"), nil, "dummy", "A:1"),
		mergeListFrame(3, "执行 dummy.Next=1；结果链第一次出现节点，tail 仍需随后前进。", caption, map[string]string{"tail.Next": "1", "tail": "dummy"}, mergeItems([]string{"3", "5"}, ""), mergeItems([]string{"2", "4", "6"}, ""), mergeItems([]string{"1"}, "1"), "dummy", "A:1"),
		mergeListFrame(3, "tail 前进到 1，同时 A 的未合并头前进到 3；两个指针职责分离。", caption, map[string]string{"tail": "1", "a": "3", "b": "2"}, mergeItems([]string{"3", "5"}, "3"), mergeItems([]string{"2", "4", "6"}, "2"), mergeItems([]string{"1"}, "1"), "1", ""),
		mergeListFrame(2, "比较 3 与 2，选择 B 的 2。", caption, map[string]string{"a": "3", "b": "2", "tail": "1"}, mergeItems([]string{"3", "5"}, "3"), mergeItems([]string{"2", "4", "6"}, "2"), mergeItems([]string{"1"}, "1"), "1", ""),
		mergeListFrame(3, "把 2 接到 1 后，并让 tail 前进；B 的当前头随后变为 4。", caption, map[string]string{"tail.Next": "2", "tail": "2", "b": "4"}, mergeItems([]string{"3", "5"}, ""), mergeItems([]string{"4", "6"}, ""), mergeItems([]string{"1", "2"}, "2"), "2", "B:2"),
		mergeListFrame(2, "继续比较 3 与 4，选择 3；每一轮都只消费一个输入头。", caption, map[string]string{"a": "3", "b": "4", "tail": "2"}, mergeItems([]string{"3", "5"}, "3"), mergeItems([]string{"4", "6"}, "4"), mergeItems([]string{"1", "2"}, "2"), "2", ""),
		mergeListFrame(3, "接入 3 后结果为 1→2→3，A 的头变成 5。", caption, map[string]string{"tail": "3", "a": "5"}, mergeItems([]string{"5"}, ""), mergeItems([]string{"4", "6"}, ""), mergeItems([]string{"1", "2", "3"}, "3"), "3", "A:3"),
		mergeListFrame(2, "比较 5 与 4，选择 B 的 4；结果尾从 3 继续向后。", caption, map[string]string{"a": "5", "b": "4", "tail": "3"}, mergeItems([]string{"5"}, "5"), mergeItems([]string{"4", "6"}, "4"), mergeItems([]string{"1", "2", "3"}, "3"), "3", ""),
		mergeListFrame(3, "接入 4，结果为 1→2→3→4；B 的头前进到 6。", caption, map[string]string{"tail": "4", "b": "6"}, mergeItems([]string{"5"}, ""), mergeItems([]string{"6"}, ""), mergeItems([]string{"1", "2", "3", "4"}, "4"), "4", "B:4"),
		mergeListFrame(2, "比较 5 与 6，选择 A 的 5；随后 A 耗尽。", caption, map[string]string{"a": "5", "b": "6", "tail": "4"}, mergeItems([]string{"5"}, "5"), mergeItems([]string{"6"}, "6"), mergeItems([]string{"1", "2", "3", "4"}, "4"), "4", ""),
		mergeListFrame(5, "A 只剩 5，接入后 A 为空；循环条件不再成立。", caption, map[string]string{"tail": "5", "a": "nil", "b": "6"}, nil, mergeItems([]string{"6"}, ""), mergeItems([]string{"1", "2", "3", "4", "5"}, "5"), "5", "A:5"),
		mergeListFrame(5, "把 B 的剩余链整体接到 tail.Next；剩余链本身已经有序，无需逐个比较。", caption, map[string]string{"tail.Next": "6", "remaining": "B:6"}, nil, mergeItems([]string{"6"}, "6"), mergeItems([]string{"1", "2", "3", "4", "5", "6"}, "6"), "5", "B:6"),
		mergeListFrame(5, "返回 dummy.Next，得到完整有序链 1→2→3→4→5→6。", caption, map[string]string{"answer": "1→2→3→4→5→6"}, nil, nil, mergeItems([]string{"1", "2", "3", "4", "5", "6"}, ""), "6", ""),
	}
	return concreteTrace("linked-list-merge", "链表：合并两个有序链表", code, frames...)
}

func redesignedBSTTrace() Trace {
	code := []string{"if node == nil { return true }", "if node.Val <= low || node.Val >= high { return false }", "left := valid(node.Left, low, node.Val)", "right := valid(node.Right, node.Val, high)", "return left && right"}
	links := []nodeLink{{From: "5", To: "1"}, {From: "5", To: "4"}, {From: "4", To: "3"}, {From: "4", To: "6"}}
	nodes := func(states map[string]string) []nodeVisual {
		return withNodes([]nodeVisual{{ID: "5", Label: "5", X: 180, Y: 28}, {ID: "1", Label: "1", X: 90, Y: 95}, {ID: "4", Label: "4", X: 270, Y: 95}, {ID: "3", Label: "3", X: 230, Y: 160}, {ID: "6", Label: "6", X: 315, Y: 160}}, states)
	}
	frames := []Frame{
		nodeFrameDetail(0, "例题树根为 5；合法范围从所有祖先共同给出，初始是 (-∞,+∞)。", "BST：进入根", map[string]string{"node": "5", "range": "(-∞,+∞)"}, nodes(map[string]string{"5": "current"}), links, nil, []string{"valid(5, -∞, +∞)"}, nil, map[string]string{"5": "range=(-∞,+∞)"}),
		nodeFrameDetail(1, "5 落在范围内；递归左子树时上界收紧为 5。", "收紧左侧范围", map[string]string{"node": "5", "leftRange": "(-∞,5)"}, nodes(map[string]string{"5": "dependency", "1": "current"}), links, []nodeLink{{From: "5", To: "1"}}, []string{"valid(5)", "valid(1, -∞, 5)"}, nil, map[string]string{"5": "合法"}),
		nodeFrameDetail(1, "节点 1 在 (-∞,5) 内，左、右孩子都会继承新的祖先约束。", "检查节点 1", map[string]string{"node": "1", "range": "(-∞,5)", "valid": "true"}, nodes(map[string]string{"5": "dependency", "1": "current"}), links, nil, []string{"valid(5)", "valid(1)"}, nil, map[string]string{"1": "true"}),
		nodeFrameDetail(1, "回到根后向右进入 4；虽然 4<5，但它必须满足根传下来的下界 5。", "进入错误右子树", map[string]string{"node": "4", "range": "(5,+∞)"}, nodes(map[string]string{"5": "dependency", "4": "current"}), links, []nodeLink{{From: "5", To: "4"}}, []string{"valid(5)", "valid(4, 5, +∞)"}, nil, map[string]string{"1": "true"}),
		nodeFrameDetail(2, "检查 4<=low(5)，条件成立；红色表示立即违反祖先范围，而不是只比较父节点。", "拒绝节点 4", map[string]string{"node": "4", "low": "5", "value": "4", "valid": "false"}, nodes(map[string]string{"5": "dependency", "4": "rejected"}), links, []nodeLink{{From: "5", To: "4"}}, []string{"valid(5)", "valid(4)"}, nil, map[string]string{"1": "true", "4": "false"}),
		nodeFrameDetail(4, "右子树返回 false，短路返回；节点 3、6 不需要被错误地当作根的合法修复。", "向上返回 false", map[string]string{"node": "5", "left": "true", "right": "false", "answer": "false"}, nodes(map[string]string{"5": "current", "1": "ready", "4": "rejected"}), links, nil, []string{"valid(5)"}, nil, map[string]string{"1": "true", "4": "false"}),
	}
	return concreteTrace("node-link-state", "BST：祖先范围验证", code, frames...)
}

func redesignedLCATrace() Trace {
	code := []string{"if root == nil || root == p || root == q { return root }", "left := lca(root.Left, p, q)", "right := lca(root.Right, p, q)", "if left == nil { return right }", "if right == nil { return left }", "return root"}
	frames := []Frame{
		nodeFrameDetail(0, "例题在树中寻找 p=5、q=1 的最近公共祖先。根 3 不是目标，必须搜索两侧。", "LCA：进入根", map[string]string{"root": "3", "p": "5", "q": "1"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, nil, []string{"lca(3,5,1)"}, nil, nil),
		nodeFrameDetail(1, "先调用左子树 lca(5)；边 3→5 是当前递归方向。", "搜索左子树", map[string]string{"root": "5", "return": "等待"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, []nodeLink{{From: "3", To: "5"}}, []string{"lca(3)", "lca(5)"}, nil, nil),
		nodeFrameDetail(0, "5 命中目标 p，按边界条件直接返回 5，不再向下搜索它的孩子。", "命中 p", map[string]string{"root": "5", "return": "5"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, nil, []string{"lca(3)", "lca(5)"}, nil, map[string]string{"left": "5"}),
		nodeFrameDetail(1, "根 3 再搜索右子树 1；它命中 q 并直接返回。", "搜索并命中 q", map[string]string{"root": "1", "return": "1"}, treeNodes(map[string]string{"3": "dependency", "1": "current"}), treeLinks, []nodeLink{{From: "3", To: "1"}}, []string{"lca(3)", "lca(1)"}, nil, map[string]string{"left": "5"}),
		nodeFrameDetail(2, "根收到 left=5、right=1；两个返回值都非空，说明两个目标分别来自两侧。", "读取左右返回值", map[string]string{"left": "5", "right": "1"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, []nodeLink{{From: "3", To: "5"}, {From: "3", To: "1"}}, []string{"lca(3)"}, nil, map[string]string{"left": "5", "right": "1"}),
		nodeFrameDetail(5, "两个分支第一次在根 3 汇合，因此 return root=3；这就是最近公共祖先。", "LCA：返回 3", map[string]string{"answer": "3", "left": "5", "right": "1"}, treeNodes(map[string]string{"3": "current", "5": "ready", "1": "ready"}), treeLinks, nil, []string{"lca(3)"}, nil, map[string]string{"LCA": "3"}),
	}
	return concreteTrace("node-link-state", "树题：最近公共祖先", code, frames...)
}

func redesignedTreePathSumTrace() Trace {
	code := []string{"if node == nil { return false }", "remain -= node.Val", "if leaf { return remain == 0 }", "return dfs(node.Left, remain) || dfs(node.Right, remain)"}
	frames := []Frame{
		nodeFrameDetail(0, "例题 target=14。进入根 3，remain 从 14 开始；当前节点还未扣除。", "路径和：进入根", map[string]string{"node": "3", "remain": "14", "path": "[]"}, treeNodes(map[string]string{"3": "current"}), treeLinks, nil, []string{"dfs(3,14)"}, []string{"3"}, nil),
		nodeFrameDetail(1, "扣除当前值 3，remain=11；只有从根到叶的完整路径才有资格结算。", "扣除节点 3", map[string]string{"node": "3", "remain": "11", "path": "[3]"}, treeNodes(map[string]string{"3": "current"}), treeLinks, nil, []string{"dfs(3,11)"}, []string{"3"}, nil),
		nodeFrameDetail(3, "3 不是叶子，沿左边 3→5 继续；把 remain=11 传给孩子。", "选择左分支", map[string]string{"node": "5", "remain": "11", "path": "[3]"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, []nodeLink{{From: "3", To: "5"}}, []string{"dfs(3,11)", "dfs(5,11)"}, []string{"3", "5"}, nil),
		nodeFrameDetail(1, "进入 5 后扣除 5，remain 从 11 变成 6；路径现场变为 [3,5]。", "扣除节点 5", map[string]string{"node": "5", "remain": "6", "path": "[3,5]"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5,6)"}, []string{"3", "5"}, nil),
		nodeFrameDetail(3, "5 不是叶子，先尝试叶子 6；把剩余目标 6 传下去。", "进入叶子候选 6", map[string]string{"node": "6", "remain": "6", "path": "[3,5]"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, []nodeLink{{From: "5", To: "6"}}, []string{"dfs(3)", "dfs(5,6)", "dfs(6,6)"}, []string{"3", "5", "6"}, nil),
		nodeFrameDetail(1, "扣除 6 后 remain=0；6 是叶子，叶子条件成立，返回 true。", "叶子命中目标", map[string]string{"node": "6", "remain": "0", "leaf": "true", "return": "true"}, treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)", "dfs(6)"}, []string{"3", "5", "6"}, map[string]string{"path": "3→5→6=14"}),
		nodeFrameDetail(3, "左分支已经返回 true，父调用短路，不需要再展开 2；答案沿调用栈向上返回。", "路径和：短路返回", map[string]string{"path": "[3,5,6]", "remain": "0", "answer": "true"}, treeNodes(map[string]string{"3": "current", "5": "ready", "6": "ready"}), treeLinks, nil, []string{"dfs(3)"}, []string{"3", "5", "6"}, map[string]string{"answer": "true"}),
	}
	return concreteTrace("node-link-state", "树题：根到叶路径和", code, frames...)
}

func redesignedTreeDPTrace() Trace {
	code := []string{"leftTake, leftSkip := dfs(node.Left)", "rightTake, rightSkip := dfs(node.Right)", "take := node.Val + leftSkip + rightSkip", "skip := max(leftTake,leftSkip) + max(rightTake,rightSkip)", "return take, skip"}
	values := map[string]string{"take(3)": "?", "skip(3)": "?"}
	frames := []Frame{
		nodeFrameDetail(0, "例题打家劫舍 III。每个节点返回两个状态：take 表示选当前，skip 表示不选当前。", "树形 DP：进入根", map[string]string{"node": "3", "take": "?", "skip": "?"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, nil, []string{"dfs(3)"}, nil, values),
		nodeFrameDetail(0, "先递归左孩子 5；根 3 暂时只能等待，不可提前把 5 当作普通值相加。", "后序进入左孩子", map[string]string{"node": "5", "parent": "3"}, treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks, []nodeLink{{From: "3", To: "5"}}, []string{"dfs(3)", "dfs(5)"}, nil, values),
		nodeFrameDetail(2, "5 的左孩子 6 返回 (take=6, skip=0)；选 5 时必须使用孩子的 skip。", "读取 6 的双状态", map[string]string{"node": "5", "leftTake": "6", "leftSkip": "0"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "dependency"}), treeLinks, []nodeLink{{From: "5", To: "6"}}, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"6": "take=6, skip=0"}),
		nodeFrameDetail(2, "5 的右孩子 2 返回 (take=2, skip=0)；两个孩子状态都已准备好。", "读取 2 的双状态", map[string]string{"node": "5", "rightTake": "2", "rightSkip": "0"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "ready", "2": "dependency"}), treeLinks, []nodeLink{{From: "5", To: "2"}}, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"6": "take=6, skip=0", "2": "take=2, skip=0"}),
		nodeFrameDetail(3, "计算 take(5)=5+skip(6)+skip(2)=5；选 5 会禁止直接选两个孩子。", "计算 take(5)", map[string]string{"take(5)": "5"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "dependency", "2": "dependency"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"5": "take=5"}),
		nodeFrameDetail(4, "计算 skip(5)=max(6,0)+max(2,0)=8；不选 5 时孩子可以独立取较优状态。", "计算 skip(5)", map[string]string{"skip(5)": "8", "take(5)": "5"}, treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "ready", "2": "ready"}), treeLinks, nil, []string{"dfs(3)", "dfs(5)"}, nil, map[string]string{"5": "take=5, skip=8"}),
		nodeFrameDetail(1, "左子树返回 (5,8) 给根；根开始递归右孩子 1。", "返回左状态", map[string]string{"leftTake": "5", "leftSkip": "8", "node": "1"}, treeNodes(map[string]string{"3": "dependency", "5": "ready", "1": "current"}), treeLinks, []nodeLink{{From: "3", To: "1"}}, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"5": "take=5, skip=8"}),
		nodeFrameDetail(2, "1 的左孩子 0 返回 (0,0)，右孩子 8 返回 (8,0)；继续保存两种状态。", "读取右子树孩子状态", map[string]string{"left": "0,0", "right": "8,0"}, treeNodes(map[string]string{"3": "dependency", "1": "current", "0": "dependency", "8": "dependency"}), treeLinks, []nodeLink{{From: "1", To: "0"}, {From: "1", To: "8"}}, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"5": "5,8", "0": "0,0", "8": "8,0"}),
		nodeFrameDetail(3, "计算 take(1)=1+skip(0)+skip(8)=1。", "计算 take(1)", map[string]string{"take(1)": "1"}, treeNodes(map[string]string{"3": "dependency", "1": "current", "0": "dependency", "8": "dependency"}), treeLinks, nil, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"1": "take=1"}),
		nodeFrameDetail(4, "计算 skip(1)=max(0,0)+max(8,0)=8；右子树返回 (1,8)。", "计算 skip(1)", map[string]string{"take(1)": "1", "skip(1)": "8"}, treeNodes(map[string]string{"3": "dependency", "1": "current", "0": "ready", "8": "ready"}), treeLinks, nil, []string{"dfs(3)", "dfs(1)"}, nil, map[string]string{"1": "take=1, skip=8"}),
		nodeFrameDetail(2, "根读取左右双状态；它们是本次两种转移的全部依赖。", "读取根的依赖", map[string]string{"left": "5,8", "right": "1,8"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, nil, []string{"dfs(3)"}, nil, map[string]string{"5": "5,8", "1": "1,8"}),
		nodeFrameDetail(3, "选根 3：take(3)=3+skip(5)+skip(1)=19。", "计算 take(3)", map[string]string{"take(3)": "19"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, nil, []string{"dfs(3)"}, nil, map[string]string{"take(3)": "19"}),
		nodeFrameDetail(4, "不选根 3：skip(3)=max(5,8)+max(1,8)=16。", "计算 skip(3)", map[string]string{"take(3)": "19", "skip(3)": "16"}, treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks, nil, []string{"dfs(3)"}, nil, map[string]string{"take(3)": "19", "skip(3)": "16"}),
		nodeFrameDetail(4, "根向上返回 (19,16)，最终答案取 max=19；状态含义比单一节点颜色更重要。", "树形 DP：最终答案", map[string]string{"take(3)": "19", "skip(3)": "16", "answer": "19"}, treeNodes(map[string]string{"3": "current", "5": "ready", "1": "ready", "6": "ready", "2": "ready", "0": "ready", "8": "ready"}), treeLinks, nil, []string{"dfs(3)"}, nil, map[string]string{"root": "take=19, skip=16"}),
	}
	return concreteTrace("node-link-state", "树形 DP：选与不选", code, frames...)
}

func redesignedStringWindowTrace() Trace {
	code := []string{"for right := range s {", "    count[s[right]]++", "    for count[s[right]] > 1 { count[s[left]]--; left++ }", "    answer = max(answer, right-left+1)", "}"}
	input := []string{"a", "b", "c", "a", "b", "c", "b", "b"}
	frames := []Frame{
		deepExampleFrame(0, "例题 s=abcabcbb。left=0、answer=0，窗口还为空。", "滑动窗口：建立窗口", map[string]string{"left": "0", "right": "-1", "answer": "0"}, tokenRow("字符串", input, nil), lane("窗口", item("[]", "current")), lane("count", item("{}", "ready"))),
		deepExampleFrame(1, "right=0 纳入 a，频次 count[a] 从 0 写成 1。", "纳入 a", map[string]string{"left": "0", "right": "0", "count[a]": "1"}, tokenRow("字符串", input, map[int]string{0: "current"}), lane("窗口", item("a", "current")), lane("count", item("a:1", "current"))),
		deepExampleFrame(3, "窗口 [0,0] 合法，写 answer=max(0,1)=1。", "记录长度 1", map[string]string{"left": "0", "right": "0", "answer": "1"}, tokenRow("字符串", input, map[int]string{0: "ready"}), lane("窗口", item("a", "ready")), lane("answer", item("1", "current"))),
		deepExampleFrame(1, "right=1 纳入 b，b 的频次为 1，窗口 ab 仍合法。", "纳入 b", map[string]string{"left": "0", "right": "1", "count": "a:1,b:1"}, tokenRow("字符串", input, map[int]string{1: "current"}), lane("窗口", item("a b", "current"))),
		deepExampleFrame(3, "right=2 纳入 c，窗口 abc 合法，answer 写成 3。", "扩展到 abc", map[string]string{"left": "0", "right": "2", "answer": "3"}, tokenRow("字符串", input, map[int]string{2: "current"}), lane("窗口", item("a b c", "current")), lane("answer", item("3", "current"))),
		deepExampleFrame(1, "right=3 再纳入 a，count[a] 从 1 变为 2；重复使窗口暂时非法。", "重复字符进入窗口", map[string]string{"left": "0", "right": "3", "count[a]": "2"}, tokenRow("字符串", input, map[int]string{0: "dependency", 3: "current"}), lane("窗口", item("a b c a", "current")), lane("count", item("a:2", "current"))),
		deepExampleFrame(2, "while 条件成立，先移除 left=0 的旧 a，count[a] 写回 1。", "收缩左端一步", map[string]string{"left": "0", "right": "3", "count[a]": "1"}, tokenRow("字符串", input, map[int]string{0: "rejected", 3: "current"}), lane("窗口", item("a b c a", "dependency")), lane("写入", item("remove a", "current"))),
		deepExampleFrame(2, "窗口恢复合法，left 前进到 1；当前 bca 长度 3，answer 保持 3。", "窗口重新合法", map[string]string{"left": "1", "right": "3", "answer": "3"}, tokenRow("字符串", input, map[int]string{1: "ready", 2: "ready", 3: "current"}), lane("窗口", item("b c a", "current"))),
		deepExampleFrame(1, "right=4 纳入 b，b 重复；随后移除旧 b，left=2，窗口 cab。", "处理第二次重复", map[string]string{"left": "2", "right": "4", "count": "c:1,a:1,b:1"}, tokenRow("字符串", input, map[int]string{1: "rejected", 4: "current"}), lane("窗口", item("c a b", "current"))),
		deepExampleFrame(1, "right=5 纳入 c，重复 c 触发同样的收缩；窗口变为 abc，最大长度仍是 3。", "重复规则复用", map[string]string{"left": "3", "right": "5", "answer": "3"}, tokenRow("字符串", input, map[int]string{2: "rejected", 5: "current"}), lane("窗口", item("a b c", "current"))),
		deepExampleFrame(1, "right=6 纳入 b 后移除旧 b，left=5；right=7 再纳入 b，再移除旧 b，窗口最终只剩 b。", "连续重复的收缩", map[string]string{"left": "7", "right": "7", "answer": "3"}, tokenRow("字符串", input, map[int]string{7: "current"}), lane("窗口", item("b", "current")), lane("answer", item("3", "ready"))),
		deepExampleFrame(4, "扫描结束，answer=3；关键是每次先完成收缩，再读取合法窗口长度。", "滑动窗口：完成", map[string]string{"answer": "3", "window": "[a,b,c]"}, lane("答案", item("3", "current")), lane("不变量", item("窗口内每个字符频次≤1", "ready"))),
	}
	return concreteTrace("example-state", "字符串：无重复滑动窗口", code, frames...)
}

func redesignedStringGoTrace() Trace {
	code := []string{"if asciiOnly(s) { use s[i] }", "chars := []rune(s)", "var builder strings.Builder", "for _, part := range parts { builder.WriteString(part) }", "return builder.String()"}
	frames := []Frame{
		deepExampleFrame(0, "例题 s=Go中。先明确题目按字节还是按 Unicode 字符处理。", "Go 字符串：选择语义", map[string]string{"input": "Go中", "len(bytes)": "4", "len(chars)": "3"}, tokenRow("bytes", []string{"G", "o", "中·3bytes"}, map[int]string{2: "current"}), lane("问题", item("按字符索引", "current"))),
		deepExampleFrame(0, "输入包含非 ASCII 字符，不能把 byte 下标直接当作字符下标。", "识别 Unicode 输入", map[string]string{"asciiOnly": "false", "decision": "use []rune"}, tokenRow("bytes", []string{"G", "o", "中·3bytes"}, map[int]string{2: "rejected"}), lane("选择", item("[]rune", "current"))),
		deepExampleFrame(1, "执行 chars := []rune(s)，把 UTF-8 字节序列转换成三个字符元素 G、o、中。", "建立 rune 切片", map[string]string{"chars": "[G,o,中]"}, tokenRow("runes", []string{"G", "o", "中"}, map[int]string{2: "current"}), lane("byte 与 rune", item("4 bytes → 3 chars", "dependency"))),
		deepExampleFrame(2, "输出 parts=[Go,中]；创建 Builder，当前结果为空。", "创建 Builder", map[string]string{"parts": "[Go,中]", "builder": "\"\""}, tokenRow("parts", []string{"Go", "中"}, map[int]string{0: "dependency"}), lane("builder", item("\"\"", "current"))),
		deepExampleFrame(3, "循环读取第一段 Go，WriteString 只追加这一段。", "写入 Go", map[string]string{"part": "Go", "builder": "Go"}, tokenRow("parts", []string{"Go", "中"}, map[int]string{0: "current"}), lane("builder", item("Go", "current"))),
		deepExampleFrame(3, "读取第二段 中，继续追加而不是用 + 创建新的中间字符串。", "写入 中", map[string]string{"part": "中", "builder": "Go中"}, tokenRow("parts", []string{"Go", "中"}, map[int]string{1: "current"}), lane("builder", item("Go", "ready"), item("中", "current"))),
		deepExampleFrame(4, "循环结束，一次性取出 builder.String()，得到 Go中。", "构造完成", map[string]string{"result": "Go中"}, lane("builder", item("Go中", "current"))),
	}
	return concreteTrace("example-state", "Go 字符串：byte、rune 与 Builder", code, frames...)
}

func redesignedPalindromeTrace() Trace {
	code := []string{"for center := range s {", "    expand(center, center)", "    expand(center, center+1)", "}", "for l >= 0 && r < len(s) && s[l] == s[r] { l--; r++ }"}
	input := []string{"b", "a", "b", "a", "d"}
	frames := []Frame{
		deepExampleFrame(0, "例题 s=babad。最长回文的中心可能是字符，也可能是字符间隙；先从 center=0 开始。", "中心扩展：建立候选", map[string]string{"center": "0", "best": "b"}, tokenRow("字符串", input, map[int]string{0: "current"}), lane("最长", item("b", "ready"))),
		deepExampleFrame(4, "奇数中心 (0,0) 向外比较越界，停止；偶数中心 (0,1) 比较 b 与 a，失配。", "处理中心 0", map[string]string{"odd": "b", "even": "无", "best": "b"}, tokenRow("字符串", input, map[int]string{0: "dependency", 1: "rejected"})),
		deepExampleFrame(0, "移动到 center=1，单字符 a 是奇数回文候选。", "处理中心 1", map[string]string{"center": "1", "best": "b"}, tokenRow("字符串", input, map[int]string{1: "current"}), lane("候选", item("a", "current"))),
		deepExampleFrame(4, "偶数中心 (1,2) 比较 a 与 b，失配；奇数扩展先比较左右 b 与 b。", "中心 1 的第一次比较", map[string]string{"l": "0", "r": "2", "match": "b=b"}, tokenRow("字符串", input, map[int]string{0: "dependency", 1: "current", 2: "dependency"})),
		deepExampleFrame(4, "b=b 匹配，边界向外扩展到 l=-1、r=3；下一次越界，得到 bab。", "得到 bab", map[string]string{"candidate": "bab", "best": "bab"}, tokenRow("字符串", input, map[int]string{0: "ready", 1: "ready", 2: "ready"}), lane("最长", item("bab", "current"))),
		deepExampleFrame(0, "center=2 的奇数候选 b；先比较相邻 a 与 a，匹配。", "处理中心 2", map[string]string{"center": "2", "l": "1", "r": "3"}, tokenRow("字符串", input, map[int]string{1: "dependency", 2: "current", 3: "dependency"})),
		deepExampleFrame(4, "继续比较 b 与 d，失配；中心 2 得到 aba，长度与 bab 相同，不覆盖答案。", "中心 2 停止", map[string]string{"candidate": "aba", "best": "bab", "reason": "长度相同"}, tokenRow("字符串", input, map[int]string{0: "dependency", 4: "rejected"}), lane("最长", item("bab", "ready"), item("aba", "rejected"))),
		deepExampleFrame(4, "center=3 的奇数扩展先得到 a，向外比较 b 与 d 失配；偶数中心 (3,4) 也失配。", "处理中心 3", map[string]string{"center": "3", "best": "bab"}, tokenRow("字符串", input, map[int]string{2: "rejected", 3: "current", 4: "rejected"})),
		deepExampleFrame(0, "center=4 只得到 d；所有中心都完成，保留长度 3 的 bab。", "中心扩展：完成", map[string]string{"best": "bab", "length": "3"}, tokenRow("字符串", input, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready"}), lane("答案", item("bab", "current"))),
	}
	return concreteTrace("example-state", "字符串：回文中心扩展", code, frames...)
}

func redesignedKMPTrace() Trace {
	code := []string{
		"pi := make([]int, len(pattern))",
		"for i, j := 1, 0; i < len(pattern); i++ {",
		"    for j > 0 && pattern[i] != pattern[j] { j = pi[j-1] }",
		"    if pattern[i] == pattern[j] { j++ }",
		"    pi[i] = j",
		"}",
		"for i, j := range text {",
		"    for j > 0 && text[i] != pattern[j] { j = pi[j-1] }",
		"    if text[i] == pattern[j] { j++ }",
		"    if j == len(pattern) { report i-j+1 }",
		"}",
	}
	pattern := []string{"a", "b", "a", "b", "a", "c", "a"}
	pi := []string{"0", "0", "1", "2", "3", "0", "1"}
	patternRow := func(states map[int]string) exampleLane { return tokenRow("pattern", pattern, states) }
	piRow := func(computed int, states map[int]string) exampleLane {
		values := make([]string, len(pi))
		for index := range values {
			values[index] = "_"
			if index < computed {
				values[index] = pi[index]
			}
		}
		return tokenRow("pi", values, states)
	}
	textInput := []string{"z", "z", "a", "b", "a", "b", "a", "c", "a"}
	frames := []Frame{
		deepExampleFrame(0, "例题 pattern=ababaca、text=zzababaca。先创建与 pattern 等长的 pi 数组，所有位置还未计算。", "KMP：建立 pi 数组", map[string]string{"pattern": "ababaca", "pi": "[_,_,_,_,_,_,_]"}, patternRow(nil), piRow(0, nil)),
		deepExampleFrame(2, "i=1：比较 pattern[1]=b 与 pattern[0]=a，失配且 j=0，不能继续回退。", "pi[1] 比较失配", map[string]string{"i": "1", "j": "0", "compare": "b != a"}, patternRow(map[int]string{1: "rejected", 0: "dependency"}), piRow(0, map[int]string{1: "current"})),
		deepExampleFrame(4, "写 pi[1]=0；数组保留已经计算的前缀 [0,0]。", "写入 pi[1]", map[string]string{"i": "1", "pi[1]": "0", "pi": "[0,0,_,_,_,_,_]"}, patternRow(map[int]string{1: "current"}), piRow(2, map[int]string{0: "ready", 1: "current"})),
		deepExampleFrame(3, "i=2：pattern[2]=a 与 pattern[0]=a 相同，j 从 0 推进为 1。", "pi[2] 匹配推进", map[string]string{"i": "2", "j": "1", "compare": "a == a"}, patternRow(map[int]string{2: "current", 0: "dependency"}), piRow(2, map[int]string{1: "ready", 2: "current"})),
		deepExampleFrame(4, "写 pi[2]=1；长度为 2 的前缀 ab 已经可以作为后缀复用。", "写入 pi[2]", map[string]string{"i": "2", "pi[2]": "1", "pi": "[0,0,1,_,_,_,_]"}, patternRow(map[int]string{2: "current"}), piRow(3, map[int]string{0: "ready", 1: "ready", 2: "current"})),
		deepExampleFrame(3, "i=3：pattern[3]=b 与 pattern[j=1]=b 相同，j 推进为 2。", "pi[3] 匹配推进", map[string]string{"i": "3", "j": "2", "compare": "b == b"}, patternRow(map[int]string{3: "current", 1: "dependency"}), piRow(3, map[int]string{2: "dependency", 3: "current"})),
		deepExampleFrame(4, "写 pi[3]=2；当前前后缀都是 ab。", "写入 pi[3]", map[string]string{"i": "3", "pi[3]": "2", "pi": "[0,0,1,2,_,_,_]"}, patternRow(map[int]string{3: "current"}), piRow(4, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"})),
		deepExampleFrame(3, "i=4：pattern[4]=a 与 pattern[j=2]=a 相同，j 推进为 3。", "pi[4] 匹配推进", map[string]string{"i": "4", "j": "3", "compare": "a == a"}, patternRow(map[int]string{4: "current", 2: "dependency"}), piRow(4, map[int]string{3: "dependency", 4: "current"})),
		deepExampleFrame(4, "写 pi[4]=3；abab 的最长相等前后缀长度为 3。", "写入 pi[4]", map[string]string{"i": "4", "pi[4]": "3", "pi": "[0,0,1,2,3,_,_]"}, patternRow(map[int]string{4: "current"}), piRow(5, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "current"})),
		deepExampleFrame(2, "i=5：pattern[5]=c 与 pattern[j=3]=b 失配；先读取 pi[j-1]=pi[2]=1，把 j 从 3 回退到 1。", "pi[5] 第一次回退", map[string]string{"i": "5", "j": "3→1", "fallback": "pi[2]=1"}, patternRow(map[int]string{5: "rejected", 3: "dependency"}), piRow(5, map[int]string{2: "dependency", 5: "current"})),
		deepExampleFrame(2, "回退后再次比较 c 与 pattern[1]=b，仍失配；j=1 再读取 pi[0]=0。", "pi[5] 第二次回退", map[string]string{"i": "5", "j": "1→0", "fallback": "pi[0]=0"}, patternRow(map[int]string{5: "rejected", 1: "dependency"}), piRow(5, map[int]string{0: "dependency", 5: "current"})),
		deepExampleFrame(4, "j=0 后无法再回退，写 pi[5]=0；失配只改变模式前缀长度。", "写入 pi[5]", map[string]string{"i": "5", "pi[5]": "0", "pi": "[0,0,1,2,3,0,_]"}, patternRow(map[int]string{5: "current"}), piRow(6, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready", 5: "current"})),
		deepExampleFrame(3, "i=6：pattern[6]=a 与 pattern[0]=a 匹配，j 推进为 1。", "pi[6] 匹配推进", map[string]string{"i": "6", "j": "1", "compare": "a == a"}, patternRow(map[int]string{6: "current", 0: "dependency"}), piRow(6, map[int]string{5: "dependency", 6: "current"})),
		deepExampleFrame(4, "pi 计算完成，数组为 [0,0,1,2,3,0,1]；下面复用它扫描主串。", "KMP：pi 完成", map[string]string{"pi": "[0,0,1,2,3,0,1]"}, patternRow(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready", 5: "ready", 6: "ready"}), piRow(7, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready", 5: "ready", 6: "current"})),
		deepExampleFrame(7, "扫描 text=zzababaca：前两个 z 都与 pattern[0] 失配，j 保持 0；主串指针继续向右。", "扫描主串前缀", map[string]string{"textIndex": "0→1", "j": "0", "pi": "[0,0,1,2,3,0,1]"}, tokenRow("text", textInput, map[int]string{0: "rejected", 1: "rejected"}), patternRow(map[int]string{0: "dependency"}), piRow(7, nil)),
		deepExampleFrame(9, "读 text[2..6]=ababa，依次匹配 pattern[0..4]，j 从 0 推进到 5；pi 数组保持在画面上作为已知结构。", "主串匹配到前缀", map[string]string{"textIndex": "2→6", "j": "5", "pi": "[0,0,1,2,3,0,1]"}, tokenRow("text", textInput, map[int]string{2: "ready", 3: "ready", 4: "ready", 5: "ready", 6: "current"}), patternRow(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "dependency"}), piRow(7, map[int]string{2: "dependency", 3: "dependency", 4: "dependency"})),
		deepExampleFrame(8, "text[7]=c 与 pattern[5]=c 匹配，j=6；继续读取 text[8]。", "继续主串匹配", map[string]string{"textIndex": "7", "j": "6"}, tokenRow("text", textInput, map[int]string{7: "current"}), patternRow(map[int]string{5: "dependency"}), piRow(7, map[int]string{4: "dependency"})),
		deepExampleFrame(9, "text[8]=a 与 pattern[6]=a 匹配，j=7 达到模式长度；报告匹配起点 2。", "报告匹配", map[string]string{"textIndex": "8", "j": "7", "matchStart": "2", "answer": "true"}, tokenRow("text", textInput, map[int]string{2: "ready", 3: "ready", 4: "ready", 5: "ready", 6: "ready", 7: "ready", 8: "current"}), patternRow(map[int]string{6: "current"}), piRow(7, map[int]string{6: "dependency"})),
	}
	return concreteTrace("example-state", "字符串：KMP 与 pi 前缀函数", code, frames...)
}

func redesignedLCSSpaceTrace() Trace {
	code := []string{"dp := make([]int, len(b)+1)", "for i := 1; i <= len(a); i++ {", "    diagonal := 0", "    for j := 1; j <= len(b); j++ {", "        up := dp[j]", "        dp[j] = transition(diagonal, up, dp[j-1])", "        diagonal = up", "    }", "}"}
	frames := []Frame{
		deepExampleFrame(0, "例题 a=ab、b=ac。只保留一行 dp=[0,0,0]，但必须额外保存被覆盖前的左上角 diagonal。", "LCS 一维压缩：初始化", map[string]string{"dp": "[0,0,0]", "diagonal": "-"}, tokenRow("b", []string{"∅", "a", "c"}, nil), lane("dp", item("0", "ready"), item("0", "ready"), item("0", "ready"))),
		deepExampleFrame(2, "开始 i=1 对应 a[0]=a，diagonal 在新行起点设为 0。", "开始新的一行", map[string]string{"i": "1", "char": "a", "diagonal": "0"}, tokenRow("a", []string{"a", "b"}, map[int]string{0: "current"}), lane("diagonal", item("0", "current"))),
		deepExampleFrame(4, "j=1 读取覆盖前 dp[1]=0 作为 up，同时读取左方 dp[0]=0 与左上角 diagonal=0。", "读取第一格依赖", map[string]string{"i": "1", "j": "1", "up": "0", "left": "0", "diagonal": "0"}, tokenRow("dp", []string{"0", "0", "0"}, map[int]string{0: "dependency", 1: "dependency"}), lane("依赖", item("up=0", "dependency"), item("diag=0", "dependency"))),
		deepExampleFrame(5, "a[0]=b[0]=a，写 dp[1]=diagonal+1=1；旧 up 仍要保存在手里。", "写入 dp[1]", map[string]string{"i": "1", "j": "1", "dp[1]": "1", "oldUp": "0"}, tokenRow("dp", []string{"0", "1", "0"}, map[int]string{1: "current"}), lane("旧 up", item("0", "dependency"))),
		deepExampleFrame(6, "把旧 up=0 写入 diagonal；下一列的左上角正是这一格覆盖前的值。", "推进 diagonal", map[string]string{"i": "1", "j": "1", "diagonal": "0"}, tokenRow("dp", []string{"0", "1", "0"}, map[int]string{2: "dependency"}), lane("diagonal", item("0", "current"))),
		deepExampleFrame(4, "j=2 读取 up=旧 dp[2]=0，left=新 dp[1]=1，diagonal=0；a 与 c 不同。", "读取第二格依赖", map[string]string{"i": "1", "j": "2", "up": "0", "left": "1", "diagonal": "0"}, tokenRow("dp", []string{"0", "1", "0"}, map[int]string{1: "dependency", 2: "dependency"}), lane("依赖", item("max(up,left)=1", "dependency"))),
		deepExampleFrame(5, "失配取 max(up,left)=1，写 dp[2]=1；再把本格旧 up=0 交给下一行位置。", "写入 dp[2]", map[string]string{"i": "1", "j": "2", "dp[2]": "1", "diagonal": "0"}, tokenRow("dp", []string{"0", "1", "1"}, map[int]string{2: "current"})),
		deepExampleFrame(2, "开始 i=2 对应 b，新的 diagonal 再次归零；dp 保留上一行 [0,1,1]。", "开始第二行", map[string]string{"i": "2", "char": "b", "dp": "[0,1,1]", "diagonal": "0"}, tokenRow("a", []string{"a", "b"}, map[int]string{1: "current"}), lane("dp", item("[0,1,1]", "ready"))),
		deepExampleFrame(4, "处理 b 与 a：up=1、left=0、diagonal=0，失配后 dp[1]=1。", "第二行第一格", map[string]string{"i": "2", "j": "1", "up": "1", "left": "0", "diagonal": "0", "dp[1]": "1"}, tokenRow("dp", []string{"0", "1", "1"}, map[int]string{1: "current"})),
		deepExampleFrame(4, "处理 b 与 c：保存旧 up=1；字符失配，max(up=1,left=1)=1，dp[2] 保持 1。", "第二行第二格", map[string]string{"i": "2", "j": "2", "up": "1", "left": "1", "diagonal": "1", "dp[2]": "1"}, tokenRow("dp", []string{"0", "1", "1"}, map[int]string{2: "current", 1: "dependency"})),
		deepExampleFrame(8, "一维数组最终为 [0,1,1]，答案 dp[2]=1；空间压缩依赖的关键是保存旧 up。", "LCS 一维压缩：完成", map[string]string{"answer": "1", "dp": "[0,1,1]"}, lane("答案", item("dp[2]=1", "current")), lane("不变量", item("diagonal=左上角旧值", "ready"))),
	}
	return concreteTrace("example-state", "LCS：一维数组空间优化", code, frames...)
}

func redesignedSpaceOptimizationTrace() Trace {
	code := []string{"previousTwo, previousOne := 1, 1", "for i := 2; i <= n; i++ {", "    current := previousTwo + previousOne", "    previousTwo, previousOne = previousOne, current", "}", "return previousOne"}
	frames := []Frame{rollingFrame(1, 1, 1, 0, false, "ready", 0, "例题爬 5 级台阶。两个变量分别保存 dp[i-2] 与 dp[i-1]，当前还没有要写的状态。"),
		rollingFrame(2, 1, 1, 0, false, "read", 1, "i=2：先读取 previousTwo=1、previousOne=1；它们是本轮唯一依赖。"),
		rollingFrame(2, 1, 1, 2, true, "write", 2, "只计算 current=1+1=2，不急着覆盖旧依赖。"),
		rollingFrame(2, 1, 2, 2, true, "roll", 3, "计算完成后整体前移：previousTwo 接住旧 previousOne，previousOne 接住 current。"),
		rollingFrame(3, 1, 2, 0, false, "read", 1, "i=3：窗口已经是 dp[1]=1、dp[2]=2，先读后写。"),
		rollingFrame(3, 1, 2, 3, true, "write", 2, "写 current=1+2=3；旧 previousTwo/One 仍可复核。"),
		rollingFrame(3, 2, 3, 3, true, "roll", 3, "滚动后窗口变为 dp[2]=2、dp[3]=3。"),
		rollingFrame(4, 2, 3, 0, false, "read", 1, "i=4：读取 2 与 3，不能把 previousTwo 误更新成更早的值。"),
		rollingFrame(4, 2, 3, 5, true, "write", 2, "写 current=5；当前值与旧窗口分开显示。"),
		rollingFrame(4, 3, 5, 5, true, "roll", 3, "滚动后窗口变为 3、5。"),
		rollingFrame(5, 3, 5, 0, false, "read", 1, "i=5：最后一轮读取 dp[3]=3 与 dp[4]=5。"),
		rollingFrame(5, 3, 5, 8, true, "write", 2, "写 current=8；如果此时先覆盖旧变量，下一轮依赖就会丢失。"),
		rollingFrame(5, 5, 8, 8, true, "roll", 3, "最终窗口滚动为 previousTwo=5、previousOne=8。"),
		rollingFrame(5, 5, 8, 8, true, "ready", 5, "返回 previousOne=8；滚动变量只保留计算所需的两个历史状态。")}
	return concreteTrace("rolling-dependency", "空间优化：依赖读取与变量覆盖顺序", code, frames...)
}

func bitmaskDeepFrame(line int, narration string, mask, last, previousLast, cost, candidate int, candidates, states []string) Frame {
	lastName, candidateName := "—", "—"
	if last >= 0 {
		lastName = string(rune('A' + last))
	}
	if candidate >= 0 {
		candidateName = string(rune('A' + candidate))
	}
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"mask": binaryMask(mask, 4), "last": lastName, "candidate": candidateName, "cost": itoa(cost)}, State: bitmaskState{Names: []string{"A", "B", "C", "D"}, Mask: mask, Last: last, PreviousLast: previousLast, Cost: cost, Candidate: candidate, Candidates: append([]string(nil), candidates...), States: append([]string(nil), states...)}}
}

func redesignedBitmaskTrace() Trace {
	code := []string{"dp[1<<0][0] = 0", "for mask := 1; mask < 1<<n; mask++ {", "    for last := 0; last < n; last++ {", "        if dp[mask][last] == inf { continue }", "        for next := 0; next < n; next++ {", "            if mask&(1<<next) != 0 { continue }", "            nextMask := mask | (1 << next)", "            candidate := dp[mask][last] + cost[last][next]", "            dp[nextMask][next] = min(dp[nextMask][next], candidate)", "        }", "    }", "}"}
	frames := []Frame{
		bitmaskDeepFrame(0, "例题把 4 个城市 A、B、C、D 的访问顺序压进 DP；从 A 出发访问所有城市，边代价用 A→B=2、B→C=3、C→D=1 展示。mask 记录集合，last 记录当前位置。", 1, 0, -1, 0, -1, []string{"A→B=2", "A→C=9", "A→D=8"}, []string{"dp[0001][A]=0"}),
		bitmaskDeepFrame(1, "当前状态是 mask=0001、last=A；A 已访问，B、C、D 的位为 0，所以三条边都可以作为下一步。", 1, 0, -1, 0, -1, []string{"B 未访问", "C 未访问", "D 未访问"}, []string{"当前：0001/A=0"}),
		bitmaskDeepFrame(4, "检查 next=B：mask 的 B 位为 0，读取 A→B=2，得到新集合 0011。", 1, 0, 0, 0, 1, []string{"A→B +2", "新 mask=0011"}, []string{"读：0001/A=0"}),
		bitmaskDeepFrame(7, "候选代价是旧值 0 加边代价 2；把 dp[0011][B] 写成 2。", 3, 1, 0, 2, 1, []string{"0 + 2 = 2", "写 0011/B"}, []string{"0001/A → 0011/B = 2"}),
		bitmaskDeepFrame(4, "回到同一个 0001/A，再检查 next=C；这次得到不同的最后位置 C，不能和 last=B 的状态合并。", 1, 0, 0, 0, 2, []string{"A→C +9", "新 mask=0101"}, []string{"已有：0011/B=2", "候选：0101/C=9"}),
		bitmaskDeepFrame(7, "写入 dp[0101][C]=9；集合相似时，last 不同仍代表不同的下一条边。", 5, 2, 0, 9, 2, []string{"0 + 9 = 9", "写 0101/C"}, []string{"0011/B=2", "0101/C=9"}),
		bitmaskDeepFrame(1, "取出较短的 0011/B=2；现在集合是 A、B，last=B，下一步只能从 B 的出边继续。", 3, 1, 0, 2, -1, []string{"C 未访问", "D 未访问", "A 已访问"}, []string{"当前：0011/B=2"}),
		bitmaskDeepFrame(4, "检查 next=C：C 位仍为 0，读取 B→C=3；新集合变成 0111。", 3, 1, 0, 2, 2, []string{"B→C +3", "新 mask=0111"}, []string{"读：0011/B=2"}),
		bitmaskDeepFrame(7, "写入 dp[0111][C]=2+3=5；这一步同时扩展集合，并更新最后位置为 C。", 7, 2, 1, 5, 2, []string{"2 + 3 = 5", "写 0111/C"}, []string{"0011/B → 0111/C = 5"}),
		bitmaskDeepFrame(4, "从 0111/C 扫描 next：A、B 已访问，D 未访问；只有 D 能让 mask 变成 1111。", 7, 2, 1, 5, 3, []string{"A 已访问", "B 已访问", "C 已访问", "D 未访问"}, []string{"当前：0111/C=5"}),
		bitmaskDeepFrame(7, "读取 C→D=1，候选为 5+1=6；写入全集合状态 dp[1111][D]。", 15, 3, 2, 6, 3, []string{"5 + 1 = 6", "写 1111/D"}, []string{"0111/C → 1111/D = 6"}),
		bitmaskDeepFrame(2, "mask=1111 表示四个城市都访问过；此时仍按 last 分开保存，最后在 dp[1111][last] 中取最小值。", 15, 3, 2, 6, -1, []string{"last=A", "last=B", "last=C", "last=D"}, []string{"全集合：等待比较"}),
		bitmaskDeepFrame(9, "示例路径 A→B→C→D 的代价是 6；动画到这里展示了状态压缩的两个维度：访问集合 + 最后位置。", 15, 3, 2, 6, -1, nil, []string{"答案候选：dp[1111][D]=6"}),
	}
	return concreteTrace("bitmask-state", "状态压缩 DP：集合与最后位置", code, frames...)
}

func redBlueDeepFrame(line int, narration string, red, mid, blue, groups int, feasible bool, scanned, segments []string) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"red": itoa(red), "mid": itoa(mid), "blue": itoa(blue), "groups": itoa(groups)}, State: redBlueState{Numbers: []int{7, 2, 5, 10, 8}, Minimum: 9, Maximum: 32, Red: red, Mid: mid, Blue: blue, Groups: groups, Feasible: feasible, Scanned: append([]string(nil), scanned...), Segments: append([]string(nil), segments...)}}
}

func redesignedBinaryRedBlueTrace() Trace {
	code := []string{"red, blue := max(nums)-1, sum(nums)", "for red+1 < blue {", "    mid := red + (blue-red)/2", "    groups, sum := 1, 0", "    for _, x := range nums {", "        if sum+x > mid { groups++; sum = x } else { sum += x }", "    }", "    if groups <= k { blue = mid } else { red = mid }", "}", "return blue"}
	frames := []Frame{
		redBlueDeepFrame(0, "例题 nums=[7,2,5,10,8]、k=2。red=9 是不可行哨兵，blue=32 是总和这个可行哨兵。", 9, -1, 32, 0, false, nil, nil),
		redBlueDeepFrame(2, "第一次取 mid=20；二分只负责选答案候选，是否可行要交给一次完整的分组扫描。", 9, 20, 32, 0, false, nil, nil),
		redBlueDeepFrame(4, "扫描 7：当前组和从 0 写成 7，没有超过 mid=20。", 9, 20, 32, 1, false, []string{"7→sum=7"}, []string{"[7]"}),
		redBlueDeepFrame(4, "扫描 2：组和 7+2=9 仍不超过 20，继续留在第一组。", 9, 20, 32, 1, false, []string{"7→7", "2→9"}, []string{"[7,2]"}),
		redBlueDeepFrame(4, "扫描 5：组和写成 14，第一组继续扩展。", 9, 20, 32, 1, false, []string{"7→7", "2→9", "5→14"}, []string{"[7,2,5]"}),
		redBlueDeepFrame(5, "扫描 10：14+10>20，必须新开第二组，组和重置为 10。", 9, 20, 32, 2, false, []string{"7→7", "2→9", "5→14", "10→new group"}, []string{"[7,2,5]", "[10]"}),
		redBlueDeepFrame(5, "扫描 8：10+8=18，最终只需 2 组；mid=20 可行，染蓝并令 blue=20。", 9, 20, 20, 2, true, []string{"7→7", "2→9", "5→14", "10→new", "8→18"}, []string{"[7,2,5]", "[10,8]"}),
		redBlueDeepFrame(2, "下一次 mid=14；扫描后 10 与 8 无法放进同一组，所需 3 组，mid 染红。", 9, 14, 20, 3, false, []string{"[7,2,5]", "[10]", "[8]"}, []string{"[7,2,5]", "[10]", "[8]"}),
		redBlueDeepFrame(2, "mid=17；分组为 [7,2,5]、[10]、[8]，仍需 3 组，继续把 red 推到 17。", 17, 17, 20, 3, false, []string{"[7,2,5]", "[10]", "[8]"}, []string{"[7,2,5]", "[10]", "[8]"}),
		redBlueDeepFrame(2, "mid=18；分组为 [7,2,5]、[10,8]，需要 2 组，mid 染蓝。", 17, 18, 18, 2, true, []string{"7→7", "2→9", "5→14", "10→new", "8→18"}, []string{"[7,2,5]", "[10,8]"}),
		redBlueDeepFrame(9, "blue-red=1，区间 (17,18] 只剩 18；它是第一个可行的蓝点，返回 18。", 17, -1, 18, 2, true, nil, []string{"答案上界=18"}),
	}
	return concreteTrace("binary-red-blue", "二分答案：红蓝染色", code, frames...)
}

func sequenceDeepFrame(line int, narration string, numbers, tails []int, tailStates []string, current, probe int, probeState, action string) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"current": itoa(current), "probe": itoa(probe), "LIS length": itoa(len(tails))}, State: sequenceState{Numbers: append([]int{}, numbers...), Current: current, Tails: append([]int{}, tails...), TailStates: append([]string{}, tailStates...), Action: action, Probe: probe, ProbeState: probeState}}
}

func redesignedLISTrace() Trace {
	code := []string{"tails := []int{}", "for _, x := range nums {", "    pos := sort.SearchInts(tails, x)", "    if pos == len(tails) { tails = append(tails, x) }", "    else { tails[pos] = x }", "}", "return len(tails)"}
	numbers := []int{10, 9, 2, 5, 3, 7, 101, 18}
	ready := func(n int) []string {
		result := make([]string, n)
		for i := range result {
			result[i] = "ready"
		}
		return result
	}
	frames := []Frame{
		sequenceDeepFrame(0, "例题 [10,9,2,5,3,7,101,18]。tails[k] 保存长度 k+1 能达到的最小结尾，初始为空。", numbers, nil, nil, -1, -1, "", "初始化"),
		sequenceDeepFrame(1, "读 x=10；tails 为空，二分没有可比较的结尾，pos=0。", numbers, nil, nil, 0, -1, "", "读取 10"),
		sequenceDeepFrame(3, "pos 等于 tails 长度，追加 10；长度 1 的最小结尾建立。", numbers, []int{10}, []string{"current"}, 0, -1, "", "追加 10"),
		sequenceDeepFrame(1, "读 x=9；在 tails=[10] 中检查位置 0，10>=9，应在这里替换。", numbers, []int{10}, []string{"dependency"}, 1, 0, "current", "二分检查"),
		sequenceDeepFrame(4, "pos=0，替换 tails[0]=9；长度仍为 1，但结尾变小。", numbers, []int{9}, []string{"current"}, 1, -1, "", "替换结尾"),
		sequenceDeepFrame(1, "读 x=2，检查 tails[0]=9；9>=2，pos=0。", numbers, []int{9}, []string{"dependency"}, 2, 0, "current", "二分检查"),
		sequenceDeepFrame(4, "替换为 tails=[2]；后续数字更容易接在结尾之后。", numbers, []int{2}, []string{"current"}, 2, -1, "", "替换结尾"),
		sequenceDeepFrame(1, "读 x=5，检查 tails[0]=2；2<5，二分向右。", numbers, []int{2}, []string{"dependency"}, 3, 0, "rejected", "二分向右"),
		sequenceDeepFrame(3, "pos=1 等于长度，追加 5，tails=[2,5]。", numbers, []int{2, 5}, []string{"ready", "current"}, 3, -1, "", "增加长度"),
		sequenceDeepFrame(1, "读 x=3，检查 tails[1]=5；5>=3，锁定替换位置 1。", numbers, []int{2, 5}, []string{"ready", "dependency"}, 4, 1, "current", "二分检查"),
		sequenceDeepFrame(4, "替换 tails[1]=3；长度不变，长度 2 的结尾从 5 降到 3。", numbers, []int{2, 3}, []string{"ready", "current"}, 4, -1, "", "替换结尾"),
		sequenceDeepFrame(1, "读 x=7：确认 2<7、3<7，pos=2；它可以延长当前最优长度。", numbers, []int{2, 3}, []string{"dependency", "dependency"}, 5, 1, "rejected", "二分向右"),
		sequenceDeepFrame(3, "追加 7，tails=[2,3,7]，长度变为 3。", numbers, []int{2, 3, 7}, []string{"ready", "ready", "current"}, 5, -1, "", "增加长度"),
		sequenceDeepFrame(1, "读 x=101；所有结尾都小于它，二分到 pos=3。", numbers, []int{2, 3, 7}, []string{"dependency", "dependency", "dependency"}, 6, 2, "rejected", "二分向右"),
		sequenceDeepFrame(3, "追加 101，tails 长度达到 4。", numbers, []int{2, 3, 7, 101}, ready(4), 6, -1, "", "增加长度"),
		sequenceDeepFrame(1, "读 x=18；检查 7<18 后发现 101>=18，定位 pos=3。", numbers, []int{2, 3, 7, 101}, ready(4), 7, 3, "current", "二分检查"),
		sequenceDeepFrame(4, "替换 101 为 18，tails=[2,3,7,18]；长度保持 4。", numbers, []int{2, 3, 7, 18}, ready(4), 7, -1, "", "替换结尾"),
		sequenceDeepFrame(6, "扫描结束，tails 长度 4，所以 LIS 长度为 4；tails 是长度摘要，不是唯一原序列。", numbers, []int{2, 3, 7, 18}, ready(4), -1, -1, "", "完成")}
	return concreteTrace("sequence-tails", "LIS：用最小结尾的 tails 做二分", code, frames...)
}

func redesignedGravityTrace() Trace {
	code := []string{"write := len(row) - 1", "for i := len(row)-1; i >= 0; i-- {", "    if row[i] == '*' { write = i - 1 }", "    if row[i] == '#' { row[i], row[write] = '.', '#' ; write-- }", "}"}
	frames := []Frame{
		gravityFrame([]string{"#", ".", ".", "*", "#", "."}, -1, 5, 0, "例题一行 #..*#.。从右向左扫描，write=5 是右侧无障碍段的落点。"),
		gravityFrame([]string{"#", ".", ".", "*", "#", "."}, 5, 5, 1, "读取 i=5='.'：空位不会消耗落点，write 仍为 5。"),
		gravityFrame([]string{"#", ".", ".", "*", "#", "."}, 4, 5, 3, "读取 i=4='#'：当前块准备写到蓝色 write=5，而不是逐秒向右移动。"),
		gravityFrame([]string{"#", ".", ".", "*", ".", "#"}, 4, 5, 3, "执行交换：# 从 4 写到 5，原位置变为 '.'；随后 write 左移到 4。"),
		gravityFrame([]string{"#", ".", ".", "*", ".", "#"}, 3, 3, 2, "读取障碍 *：它不移动，并把左侧区间的 write 重置为 2。"),
		gravityFrame([]string{"#", ".", ".", "*", ".", "#"}, 2, 2, 1, "读取 i=2='.'：左侧空位保持，write 仍指向 2。"),
		gravityFrame([]string{"#", ".", ".", "*", ".", "#"}, 1, 2, 1, "读取 i=1='.'：继续跳过空位，未扫描区间不被破坏。"),
		gravityFrame([]string{"#", ".", ".", "*", ".", "#"}, 0, 2, 3, "读取 i=0='#'：当前块准备落到左侧区间 write=2。"),
		gravityFrame([]string{".", ".", "#", "*", ".", "#"}, 0, 1, 3, "执行交换并令 write=1；障碍两侧分别稳定，最终是 ..#* .#。"),
		gravityFrame([]string{".", ".", "#", "*", ".", "#"}, -1, 1, 4, "扫描结束：..#*.#。这个状态可逐行复用到完整旋转盒子。"),
	}
	return concreteTrace("row-gravity", "带障碍的重力模拟：写指针从右向左", code, frames...)
}

func redesignedStartSortedIntervalsTrace() Trace {
	code := []string{"sort.Slice(intervals, byStart)", "merged := [][]int{intervals[0]}", "for _, current := range intervals[1:] {", "    last := merged[len(merged)-1]", "    if current[0] <= last[1] { last[1] = max(last[1], current[1]) }", "    else { merged = append(merged, current) }", "}"}
	intervals := []Interval{{Label: "A", Start: 1, End: 3}, {Label: "B", Start: 2, End: 6}, {Label: "C", Start: 8, End: 10}, {Label: "D", Start: 15, End: 18}}
	input := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("输入区间", greedyRangeIntervals(intervals, states)...)
	}
	merged := func(segments ...greedyRangeSegment) greedyRangeTrack {
		return makeGreedyRangeTrack("合并结果", segments...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 [[1,3],[2,6],[8,10],[15,18]]。排序后未来区间只会从当前段右侧进入。", "区间合并：按开始时间", 0, 18, map[string]string{"merged": "[[1,3]]"}, input(map[int]string{0: "current"}), merged(makeGreedyRangeSegment(1, 3, "[1,3]", "current", "range"))),
		greedyRangeFrame(2, "读取 current=[2,6]；比较 current.start=2 与 last.end=3，2<=3，存在重叠。", "检查第二段", 0, 18, map[string]string{"current": "[2,6]", "last": "[1,3]", "merged": "[[1,3]]"}, input(map[int]string{0: "dependency", 1: "current"}), merged(makeGreedyRangeSegment(1, 3, "last=[1,3]", "dependency", "range"))),
		greedyRangeFrame(4, "只扩展当前段右端：last.end=max(3,6)=6，merged 变为 [[1,6]]。", "合并重叠段", 0, 18, map[string]string{"last.end": "6", "merged": "[[1,6]]"}, input(map[int]string{0: "ready", 1: "ready"}), merged(makeGreedyRangeSegment(1, 6, "[1,6]", "current", "range"))),
		greedyRangeFrame(2, "读取 [8,10]；8>6，当前段已经固定，不能再被它覆盖。", "检查断开区间", 0, 18, map[string]string{"current": "[8,10]", "last.end": "6"}, input(map[int]string{0: "ready", 1: "ready", 2: "current"}), merged(makeGreedyRangeSegment(1, 6, "[1,6]", "dependency", "range"))),
		greedyRangeFrame(5, "断开时复制 current 追加新段，而不是修改 last；merged=[[1,6],[8,10]]。", "追加新段", 0, 18, map[string]string{"merged": "[[1,6],[8,10]]"}, input(map[int]string{0: "ready", 1: "ready", 2: "dependency"}), merged(makeGreedyRangeSegment(1, 6, "[1,6]", "ready", "range"), makeGreedyRangeSegment(8, 10, "[8,10]", "current", "range"))),
		greedyRangeFrame(2, "读取 [15,18]，15>10，再次进入断开分支。", "检查最后区间", 0, 18, map[string]string{"current": "[15,18]", "last.end": "10"}, input(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"}), merged(makeGreedyRangeSegment(1, 6, "[1,6]", "ready", "range"), makeGreedyRangeSegment(8, 10, "[8,10]", "dependency", "range"))),
		greedyRangeFrame(5, "追加后得到三个互不重叠段；排序键决定了只需观察 merged 的最后一段。", "区间合并：完成", 0, 18, map[string]string{"merged": "[[1,6],[8,10],[15,18]]"}, input(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready"}), merged(makeGreedyRangeSegment(1, 6, "[1,6]", "ready", "range"), makeGreedyRangeSegment(8, 10, "[8,10]", "ready", "range"), makeGreedyRangeSegment(15, 18, "[15,18]", "current", "range"))),
	}
	return concreteTrace("greedy-range", "区间：按开始时间合并", code, frames...)
}

func redesignedMeetingRoomsTrace() Trace {
	code := []string{"sort.Ints(starts); sort.Ints(ends)", "rooms, end := 0, 0", "for _, start := range starts {", "    if start < ends[end] { rooms++ } else { end++ }", "}", "return rooms"}
	meetings := []Interval{{Label: "A", Start: 0, End: 30}, {Label: "B", Start: 5, End: 10}, {Label: "C", Start: 15, End: 20}}
	meetingTrack := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("会议区间", greedyRangeIntervals(meetings, states)...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 [[0,30],[5,10],[15,20]]。横轴是统一时间轴；开始与结束指针都在会议区间上移动。", "会议室：双指针扫描", 0, 30, map[string]string{"rooms": "0", "start": "-", "end": "0"}, meetingTrack(nil), makeGreedyRangeTrack("开始时间"), makeGreedyRangeTrack("最早结束")),
		greedyRangeFrame(2, "读取 start=0，比较最早结束 10；0<10，当前会议无法复用房间。", "检查第一个开始时间", 0, 30, map[string]string{"start": "0", "nextEnd": "10", "rooms": "0"}, meetingTrack(map[int]string{0: "current"}), makeGreedyRangeTrack("开始时间"), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "start=0", 0, "current"), makeGreedyRangeMarker("最早结束", "end=10", 10, "dependency")),
		greedyRangeFrame(3, "并发房间数写成 1；end 指针仍停在最早结束的 10。", "开第一间房", 0, 30, map[string]string{"rooms": "1", "end": "0"}, meetingTrack(map[int]string{0: "current"}), makeGreedyRangeTrack("开始时间", makeGreedyRangeSegment(0, 1, "0", "current", "item")), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "start=0", 0, "current"), makeGreedyRangeMarker("最早结束", "end=10", 10, "dependency")),
		greedyRangeFrame(2, "读取 start=5，5<10；第一场仍未结束，不能复用，rooms 将增加。", "第二场产生重叠", 0, 30, map[string]string{"start": "5", "nextEnd": "10", "rooms": "1"}, meetingTrack(map[int]string{0: "dependency", 1: "current"}), makeGreedyRangeTrack("开始时间"), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "start=5", 5, "current"), makeGreedyRangeMarker("最早结束", "end=10", 10, "dependency")),
		greedyRangeFrame(3, "rooms 从 1 写成 2；峰值并发目前是 2，两个重叠的区间仍完整可见。", "开第二间房", 0, 30, map[string]string{"rooms": "2", "end": "0"}, meetingTrack(map[int]string{0: "dependency", 1: "current"}), makeGreedyRangeTrack("开始时间", makeGreedyRangeSegment(0, 1, "0", "ready", "item"), makeGreedyRangeSegment(5, 6, "5", "current", "item")), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "start=5", 5, "current"), makeGreedyRangeMarker("最早结束", "end=10", 10, "dependency")),
		greedyRangeFrame(2, "读取 start=15，15>=最早结束 10；end 指针前进到 20，表示释放一间房。", "找到可复用房间", 0, 30, map[string]string{"start": "15", "ends[end]": "10", "end": "1", "rooms": "2"}, meetingTrack(map[int]string{0: "dependency", 1: "dependency", 2: "current"}), makeGreedyRangeTrack("开始时间"), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "start=15", 15, "current"), makeGreedyRangeMarker("最早结束", "end=20", 20, "current"), makeGreedyRangeMarker("最早结束", "释放=10", 10, "rejected")),
		greedyRangeFrame(3, "复用释放的房间，不增加 rooms；最终峰值仍为 2。", "会议室：完成", 0, 30, map[string]string{"rooms": "2", "answer": "2"}, meetingTrack(map[int]string{0: "ready", 1: "ready", 2: "ready"}), makeGreedyRangeTrack("开始时间", makeGreedyRangeSegment(0, 1, "0", "ready", "item"), makeGreedyRangeSegment(5, 6, "5", "ready", "item"), makeGreedyRangeSegment(15, 16, "15", "current", "item")), makeGreedyRangeTrack("最早结束"), makeGreedyRangeMarker("开始时间", "完成", 15, "current"), makeGreedyRangeMarker("最早结束", "end=20", 20, "ready")),
	}
	return concreteTrace("greedy-range", "区间：最少会议室", code, frames...)
}

func redesignedWeightedIntervalsTrace() Trace {
	code := []string{"sort intervals by end", "dp[0] = 0", "for i := 1; i <= n; i++ {", "    prev := latest compatible prefix", "    skip := dp[i-1]", "    take := weight[i] + dp[prev]", "    dp[i] = max(skip, take)", "}"}
	jobs := []Interval{{Label: "A", Start: 1, End: 3}, {Label: "B", Start: 2, End: 5}, {Label: "C", Start: 4, End: 6}}
	jobTrack := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("工作区间", greedyRangeIntervals(jobs, states)...)
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 A=[1,3],w=5；B=[2,5],w=100；C=[4,6],w=5。区间条保留时间范围，收益显示在下方 DP 轨道。", "带权区间调度：前驱 DP", 0, 6, map[string]string{"dp": "[0,?, ?, ?]"}, jobTrack(nil), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 1, "dp[0]=0", "current", "range"))),
		greedyRangeFrame(2, "处理 A：它之前没有兼容工作，二分得到 prev=0。", "寻找 A 的前驱", 0, 6, map[string]string{"job": "A", "prev": "0"}, jobTrack(map[int]string{0: "current"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 1, "dp[0]", "dependency", "range")), makeGreedyRangeMarker("工作区间", "A", 1, "current")),
		greedyRangeFrame(4, "A 的两种选择：skip=dp[0]=0，take=5+dp[0]=5；橙色候选取更大值。", "比较 A 的两种选择", 0, 6, map[string]string{"skip": "0", "take": "5"}, jobTrack(map[int]string{0: "dependency"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 1, "skip=0", "dependency", "range"), makeGreedyRangeSegment(1, 3, "take=5", "current", "range")), makeGreedyRangeMarker("工作区间", "A", 1, "current")),
		greedyRangeFrame(6, "写 dp[1]=5；这是结束在 3 之前的最优收益。", "写 dp[1]", 0, 6, map[string]string{"dp[1]": "5"}, jobTrack(map[int]string{0: "ready"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 3, "dp[1]=5", "current", "range"))),
		greedyRangeFrame(2, "处理 B：start=2，最近兼容前驱仍为 0；不能因为 A 结束更早就强制选择 A。", "寻找 B 的前驱", 0, 6, map[string]string{"job": "B", "prev": "0"}, jobTrack(map[int]string{0: "ready", 1: "current"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 3, "dp[1]=5", "dependency", "range")), makeGreedyRangeMarker("工作区间", "B", 2, "current")),
		greedyRangeFrame(4, "B 的 skip=dp[1]=5，take=100+dp[0]=100；高收益工作胜过结束更早的 A。", "比较 B 的两种选择", 0, 6, map[string]string{"skip": "5", "take": "100"}, jobTrack(map[int]string{0: "ready", 1: "current"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 3, "skip=5", "dependency", "range"), makeGreedyRangeSegment(2, 5, "take=100", "current", "range")), makeGreedyRangeMarker("工作区间", "B", 2, "current")),
		greedyRangeFrame(6, "写 dp[2]=100；此处正是无权区间调度贪心不能复用的分叉。", "写 dp[2]", 0, 6, map[string]string{"dp[2]": "100"}, jobTrack(map[int]string{0: "ready", 1: "ready"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 5, "dp[2]=100", "current", "range"))),
		greedyRangeFrame(2, "处理 C：start=4，最近兼容前驱是 A，prev=1；读取 dp[1]=5。", "寻找 C 的前驱", 0, 6, map[string]string{"job": "C", "prev": "1", "dp[prev]": "5"}, jobTrack(map[int]string{0: "ready", 1: "ready", 2: "current"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 3, "dp[1]=5", "dependency", "range")), makeGreedyRangeMarker("工作区间", "C", 4, "current")),
		greedyRangeFrame(4, "C 的 take=5+dp[1]=10，skip=dp[2]=100；跳过 C 更优。", "比较 C 的两种选择", 0, 6, map[string]string{"skip": "100", "take": "10"}, jobTrack(map[int]string{0: "ready", 1: "ready", 2: "rejected"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(4, 6, "take=10", "rejected", "range"), makeGreedyRangeSegment(0, 5, "skip=100", "current", "range")), makeGreedyRangeMarker("工作区间", "C", 4, "rejected")),
		greedyRangeFrame(6, "写 dp[3]=100；最终选择 B，最大收益为 100。", "写最终答案", 0, 6, map[string]string{"dp[3]": "100", "answer": "100"}, jobTrack(map[int]string{0: "ready", 1: "ready", 2: "ready"}), makeGreedyRangeTrack("DP收益", makeGreedyRangeSegment(0, 6, "dp[3]=100", "current", "range"))),
	}
	return concreteTrace("greedy-range", "区间：带权调度 DP", code, frames...)
}

func redesignedKadaneTrace() Trace {
	code := []string{"current, best := nums[0], nums[0]", "for _, x := range nums[1:] {", "    extend := current + x", "    current = max(x, extend)", "    best = max(best, current)", "}", "return best"}
	values := []string{"-2", "1", "-3", "4", "-1", "2", "1", "-5", "4"}
	array := func(states map[int]string) greedyRangeTrack {
		return makeGreedyRangeTrack("数组位置", greedyRangeItems(values, states)...)
	}
	currentRange := func(start, end int, label, state string) greedyRangeTrack {
		return makeGreedyRangeTrack("当前子段", makeGreedyRangeSegment(start, end+1, label, state, "range"))
	}
	bestRange := func(start, end int, label, state string) greedyRangeTrack {
		return makeGreedyRangeTrack("全局最优", makeGreedyRangeSegment(start, end+1, label, state, "range"))
	}
	frames := []Frame{
		greedyRangeFrame(0, "例题 [-2,1,-3,4,-1,2,1,-5,4]。横轴始终保留整个数组；当前子段和全局最优分别占固定轨道。", "Kadane：初始化", 0, 9, map[string]string{"current": "-2", "best": "-2", "range": "[0,0]"}, array(map[int]string{0: "current"}), currentRange(0, 0, "sum=-2", "current"), bestRange(0, 0, "best=-2", "ready"), makeGreedyRangeMarker("数组位置", "i=0", 0, "current")),
		greedyRangeFrame(2, "读 x=1；在原数组位置 1 上同时保留延续候选 -1 和重开候选 1。", "比较延续与重开", 0, 9, map[string]string{"x": "1", "extend": "-1", "restart": "1"}, array(map[int]string{0: "dependency", 1: "current"}), currentRange(0, 1, "extend=-1", "dependency"), bestRange(0, 0, "best=-2", "ready"), makeGreedyRangeMarker("数组位置", "i=1", 1, "current")),
		greedyRangeFrame(3, "取较大者，current=1；当前子段移动为单点 [1,1]，数组主体没有换屏。", "写 current=1", 0, 9, map[string]string{"current": "1", "range": "[1,1]"}, array(map[int]string{0: "ready", 1: "current"}), currentRange(1, 1, "sum=1", "current"), bestRange(0, 0, "best=-2", "ready"), makeGreedyRangeMarker("数组位置", "i=1", 1, "current")),
		greedyRangeFrame(4, "更新 best=max(-2,1)=1；全局最优轨道扩展到位置 1。", "更新 best", 0, 9, map[string]string{"current": "1", "best": "1"}, array(map[int]string{0: "ready", 1: "ready"}), currentRange(1, 1, "sum=1", "ready"), bestRange(1, 1, "best=1", "current"), makeGreedyRangeMarker("数组位置", "best", 1, "current")),
		greedyRangeFrame(2, "读 x=-3；extend=-2 比 restart=-3 更大，当前子段延续到位置 2，但和降为 -2。", "负数削弱当前段", 0, 9, map[string]string{"x": "-3", "extend": "-2", "restart": "-3", "current": "-2"}, array(map[int]string{1: "dependency", 2: "current"}), currentRange(1, 2, "sum=-2", "current"), bestRange(1, 1, "best=1", "ready"), makeGreedyRangeMarker("数组位置", "i=2", 2, "current")),
		greedyRangeFrame(2, "读 x=4；extend=2，restart=4，旧 current 为负，重开得到更优。", "从 4 重开", 0, 9, map[string]string{"x": "4", "extend": "2", "restart": "4"}, array(map[int]string{2: "dependency", 3: "current"}), currentRange(3, 3, "restart=4", "current"), bestRange(1, 1, "best=1", "ready"), makeGreedyRangeMarker("数组位置", "i=3", 3, "current")),
		greedyRangeFrame(3, "写 current=4、best=max(1,4)=4；当前子段和全局最优都落在位置 3。", "写入第二个状态", 0, 9, map[string]string{"current": "4", "best": "4", "range": "[3,3]"}, array(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "current"}), currentRange(3, 3, "sum=4", "current"), bestRange(3, 3, "best=4", "current"), makeGreedyRangeMarker("数组位置", "i=3", 3, "current")),
		greedyRangeFrame(2, "依次读 -1、2、1：current 由 3、5、6 逐步沿原数组延伸到位置 6。", "连续延续正贡献", 0, 9, map[string]string{"current": "6", "best": "6", "range": "[3,6]"}, array(map[int]string{3: "dependency", 4: "ready", 5: "ready", 6: "current"}), currentRange(3, 6, "sum=6", "current"), bestRange(3, 6, "best=6", "ready"), makeGreedyRangeMarker("数组位置", "i=6", 6, "current")),
		greedyRangeFrame(2, "读 x=-5：extend=1 胜过 restart=-5，当前子段仍沿位置 7 延伸；best 保持 6。", "短暂回落但不重开", 0, 9, map[string]string{"x": "-5", "extend": "1", "restart": "-5", "current": "1", "best": "6"}, array(map[int]string{6: "dependency", 7: "current"}), currentRange(3, 7, "sum=1", "dependency"), bestRange(3, 6, "best=6", "ready"), makeGreedyRangeMarker("数组位置", "i=7", 7, "current")),
		greedyRangeFrame(2, "读最后 x=4：extend=5，restart=4，current=5，小于保留在原位置的 best=6。", "最后一次比较", 0, 9, map[string]string{"x": "4", "extend": "5", "restart": "4", "current": "5", "best": "6"}, array(map[int]string{7: "dependency", 8: "current"}), currentRange(3, 8, "sum=5", "dependency"), bestRange(3, 6, "best=6", "ready"), makeGreedyRangeMarker("数组位置", "i=8", 8, "current")),
		greedyRangeFrame(6, "扫描结束，返回 best=6，对应原数组中的连续区间 [4,-1,2,1]。", "Kadane：完成", 0, 9, map[string]string{"answer": "6", "range": "[4,-1,2,1]"}, array(map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready", 5: "ready", 6: "ready", 7: "ready", 8: "ready"}), currentRange(3, 8, "current=5", "ready"), bestRange(3, 6, "best=6", "current"), makeGreedyRangeMarker("数组位置", "答案", 6, "current")),
	}
	return concreteTrace("greedy-range", "贪心：最大子数组 Kadane", code, frames...)
}
