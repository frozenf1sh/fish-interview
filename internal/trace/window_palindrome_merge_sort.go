package trace

// ListMergeSortTrace shows top-down merge sort on one concrete linked list.
// Every split, comparison, write, and tail append gets its own frame so the
// result lane remains as stable as the merge-two-lists animation.
func ListMergeSortTrace() Trace {
	code := []string{
		"if head == nil || head.Next == nil { return head }",
		"mid := splitWithFastSlow(head)",
		"left := sort(head)",
		"right := sort(mid)",
		"for left != nil && right != nil {",
		"    choose the smaller head",
		"    tail.Next = chosen; tail = tail.Next",
		"}",
		"tail.Next = remaining; return dummy.Next",
	}
	caption := "链表：归并排序的拆分与合并"
	items := func(values ...string) []exampleItem { return mergeSortItems(values, nil) }
	current := func(values ...string) []exampleItem { return mergeSortItems(values, map[int]string{0: "current"}) }
	ready := func(values ...string) []exampleItem { return mergeSortItems(values, nil) }
	focusedOriginal := func(values ...string) []exampleItem {
		states := make(map[int]string, len(values))
		for _, value := range values {
			for index, original := range []string{"4", "2", "1", "3"} {
				if original == value {
					states[index] = "current"
				}
			}
		}
		return mergeSortItems([]string{"4", "2", "1", "3"}, states)
	}
	detail := func(line int, narration string, variables map[string]string, source, left, right, result []exampleItem, stack []string, phase string, topology mergeSortTopology) Frame {
		return mergeSortListFrameWithTopology(line, narration, caption, variables, source, left, right, result, stack, phase, topology)
	}
	frames := []Frame{
		mergeSortListFrame(0, "例题 4→2→1→3。先固定展示整条原链，fast/slow 只负责定位切口。", caption, map[string]string{"head": "4", "mid": "1", "stack": "sort(4,2,1,3)"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[4,2,1,3]"}, "准备拆分"),
		detail(1, "slow 停在 2；这一步只做一件事：去掉原链上的 2→1 箭头。节点还留在原链，拓扑先发生变化。", map[string]string{"slow": "2", "fast": "nil", "cut": "2→1"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[4,2]", "[1,3]"}, "移除中点箭头", mergeSortTopology{original: focusedOriginal("4", "2"), originalLinks: []bool{true, false, true}}),
		detail(1, "箭头已经消失；现在把 4、2 和 1、3 从原链的位置移动到下方两条子链。", map[string]string{"left": "4→2", "right": "1→3"}, items("4", "2", "1", "3"), items("4", "2"), items("1", "3"), nil, []string{"[4,2]", "[1,3]"}, "节点移动到左右子链", mergeSortTopology{originalLinks: []bool{true, false, true}}),
		detail(1, "第一次拆分完成：原链仍作为总体参照，左、右子链已经成为递归要处理的输入。", map[string]string{"left": "4→2", "right": "1→3"}, items("4", "2", "1", "3"), items("4", "2"), items("1", "3"), nil, []string{"[4,2]", "[1,3]"}, "第一次拆分完成", mergeSortTopology{originalLinks: []bool{true, false, true}}),
		detail(0, "递归进入左子链 4→2；先隐藏上一层左右子链，只保留原链和递归栈作为参照。", map[string]string{"head": "4", "mid": "2", "stack": "sort(4,2)"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[4,2]", "[1,3]"}, "进入左子链", mergeSortTopology{originalLinks: []bool{true, false, true}}),
		detail(1, "高亮原链上的 4、2，去掉它们之间的 4→2 箭头；这是真正的子链断开。", map[string]string{"cut": "4→2", "mid": "2"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[4]", "[2]", "[1,3]"}, "断开左子链箭头", mergeSortTopology{original: focusedOriginal("4", "2"), originalLinks: []bool{false, false, true}}),
		detail(1, "隐藏上一层子链后，4 和 2 分别移动到两个单节点子链；右半 1→3 暂存在递归栈里。", map[string]string{"left": "4", "right": "2"}, items("4", "2", "1", "3"), items("4"), items("2"), nil, []string{"[4]", "[2]", "[1,3]"}, "移动为单节点链", mergeSortTopology{originalLinks: []bool{false, false, true}, leftLabel: "左单链", rightLabel: "右单链"}),
		mergeSortListFrame(0, "4 与 2 都是单节点有序链，递归返回后准备合并。", caption, map[string]string{"left": "4", "right": "2"}, items("4", "2"), ready("4"), ready("2"), nil, []string{"merge([4],[2])", "[1,3]"}, "准备合并 4 与 2"),
		mergeSortListFrame(4, "读取左链头 4 和右链头 2；两条输入链保持在原位，结果从 dummy 开始。", caption, map[string]string{"left": "4", "right": "2", "tail": "dummy"}, items("4", "2"), current("4"), current("2"), nil, []string{"merge([4],[2])", "[1,3]"}, "比较 4 与 2"),
		mergeSortListFrame(5, "2 更小，先记录选择；节点尚未从输入链移动。", caption, map[string]string{"chosen": "2", "tail": "dummy"}, items("4", "2"), current("4"), current("2"), nil, []string{"merge([4],[2])", "[1,3]"}, "选择 2"),
		mergeSortListFrame(6, "把 2 从右链动态接到 dummy 后面；临时结果链第一次出现节点。", caption, map[string]string{"tail.Next": "2", "tail": "2"}, items("4"), current("4"), nil, ready("2"), []string{"merge([4],[2])", "[1,3]"}, "2 移到 dummy 后"),
		mergeSortListFrame(7, "右链为空，剩余 4 不再比较；下一步把它接到结果尾。", caption, map[string]string{"remaining": "4", "tail": "2"}, items("4"), current("4"), nil, ready("2"), []string{"merge([4],[2])", "[1,3]"}, "发现剩余 4"),
		mergeSortListFrame(8, "把 4 接到 2 后面，得到有序临时链 2→4，返回上一层。", caption, map[string]string{"tail.Next": "4", "return": "2→4"}, items("2", "4"), nil, nil, ready("2", "4"), []string{"[1,3]"}, "左半合并完成"),
		detail(0, "递归处理右半 1→3；先隐藏旧的左右子链，再定位 1→3 的切口。", map[string]string{"head": "1", "mid": "3", "stack": "sort(1,3)"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[2,4]", "[1,3]"}, "进入右子链", mergeSortTopology{originalLinks: []bool{true, false, true}}),
		detail(1, "高亮原链上的 1、3，去掉 1→3 箭头；右子链也必须先改变 Next 拓扑。", map[string]string{"cut": "1→3", "mid": "3"}, items("4", "2", "1", "3"), nil, nil, nil, []string{"[1]", "[3]", "[2,4]"}, "断开右子链箭头", mergeSortTopology{original: focusedOriginal("1", "3"), originalLinks: []bool{true, false, false}}),
		detail(1, "1 和 3 分别移动到单节点子链；到这里四个叶子节点都已经独立。", map[string]string{"left": "1", "right": "3"}, items("4", "2", "1", "3"), items("1"), items("3"), nil, []string{"[1]", "[3]", "[2,4]"}, "移动为单节点链", mergeSortTopology{originalLinks: []bool{true, false, false}, leftLabel: "左单链", rightLabel: "右单链"}),
		mergeSortListFrame(4, "读取 1 和 3；两条单节点链都在输入区，dummy 仍固定在结果头。", caption, map[string]string{"left": "1", "right": "3", "tail": "dummy"}, items("1", "3"), current("1"), current("3"), nil, []string{"merge([1],[3])"}, "比较 1 与 3"),
		mergeSortListFrame(5, "1 更小，选择左链节点 1。", caption, map[string]string{"chosen": "1"}, items("1", "3"), current("1"), current("3"), nil, []string{"merge([1],[3])"}, "选择 1"),
		mergeSortListFrame(6, "把 1 动态接到 dummy 后；右链的 3 仍保留在输入区。", caption, map[string]string{"tail.Next": "1", "tail": "1"}, items("3"), nil, current("3"), ready("1"), []string{"merge([1],[3])"}, "1 移到 dummy 后"),
		mergeSortListFrame(7, "左链耗尽，剩余节点只有 3；直接接到结果尾。", caption, map[string]string{"remaining": "3", "tail": "1"}, nil, nil, current("3"), ready("1"), []string{"merge([1],[3])"}, "发现剩余 3"),
		mergeSortListFrame(8, "得到有序右半 1→3；现在两侧输入分别是 2→4 和 1→3。", caption, map[string]string{"return": "1→3"}, items("2", "4", "1", "3"), ready("2", "4"), ready("1", "3"), ready("1", "3"), []string{"merge([2,4],[1,3])"}, "右半合并完成"),
		mergeSortListFrame(4, "顶层合并读取 2 和 1；原链的 2→1 箭头仍然是断开的，输入区保持固定。", caption, map[string]string{"left": "2", "right": "1", "tail": "dummy"}, items("4", "2", "1", "3"), current("2"), current("1"), nil, []string{"merge([2,4],[1,3])"}, "比较 2 与 1"),
		mergeSortListFrame(5, "1 更小，单独记录选择；还没有写入临时链。", caption, map[string]string{"chosen": "1"}, items("4", "2", "1", "3"), current("2"), current("1"), nil, []string{"merge([2,4],[1,3])"}, "选择 1"),
		mergeSortListFrame(6, "把 1 从右链动态移动到 dummy 后，临时链变成 dummy→1。", caption, map[string]string{"tail.Next": "1", "tail": "1", "right": "3"}, items("4", "2", "3"), current("2", "4"), current("3"), ready("1"), []string{"merge([2,4],[1,3])"}, "1 移到 dummy 后"),
		mergeSortListFrame(4, "读取 2 和 3；临时链里的 1 保持绿色，不随输入区来回切换。", caption, map[string]string{"left": "2", "right": "3", "tail": "1"}, items("4", "2", "3"), current("2", "4"), current("3"), ready("1"), []string{"merge([2,4],[1,3])"}, "比较 2 与 3"),
		mergeSortListFrame(5, "2 更小，选择左链头 2。", caption, map[string]string{"chosen": "2"}, items("4", "2", "3"), current("2", "4"), current("3"), ready("1"), []string{"merge([2,4],[1,3])"}, "选择 2"),
		mergeSortListFrame(6, "把 2 动态移动到临时链 1 后，得到 dummy→1→2。", caption, map[string]string{"tail.Next": "2", "tail": "2", "left": "4"}, items("4", "3"), current("4"), current("3"), ready("1", "2"), []string{"merge([2,4],[1,3])"}, "2 移到结果尾"),
		mergeSortListFrame(4, "读取 4 和 3；只高亮这两个当前头，已合并的 1、2 不再被重画。", caption, map[string]string{"left": "4", "right": "3", "tail": "2"}, items("4", "3"), current("4"), current("3"), ready("1", "2"), []string{"merge([2,4],[1,3])"}, "比较 4 与 3"),
		mergeSortListFrame(5, "3 更小，选择右链头 3。", caption, map[string]string{"chosen": "3"}, items("4"), current("4"), current("3"), ready("1", "2"), []string{"merge([2,4],[1,3])"}, "选择 3"),
		mergeSortListFrame(6, "把 3 接到临时链尾，结果变成 dummy→1→2→3；左链只剩 4。", caption, map[string]string{"tail.Next": "3", "tail": "3", "remaining": "4"}, items("4"), current("4"), nil, ready("1", "2", "3"), []string{"merge([2,4],[1,3])"}, "3 移到结果尾"),
		mergeSortListFrame(7, "比较循环结束，剩余链是 4；最后一次接入也单独展示。", caption, map[string]string{"remaining": "4", "tail": "3"}, items("4"), current("4"), nil, ready("1", "2", "3"), []string{"merge([2,4],[1,3])"}, "发现剩余 4"),
		mergeSortListFrame(8, "把 4 接到结果尾，临时链完整变为 dummy→1→2→3→4。", caption, map[string]string{"tail.Next": "4", "answer": "1→2→3→4"}, ready("1", "2", "3", "4"), nil, nil, ready("1", "2", "3", "4"), []string{"merge([2,4],[1,3])"}, "临时链完成"),
		detail(8, "下一步不是瞬间替换：把临时链中的 1、2、3、4 作为同一批节点，动态移动到原链排序后的四个位置。", map[string]string{"overlay": "dummy.Next→原链", "answer": "1→2→3→4"}, ready("1", "2", "3", "4"), nil, nil, ready("1", "2", "3", "4"), nil, "临时链覆盖原链", mergeSortTopology{original: mergeSortItems([]string{"1", "2", "3", "4"}, nil), originalLinks: []bool{true, true, true}, overlay: true}),
		detail(8, "覆盖完成；临时链退场，原链现在就是排序后的 1→2→3→4，返回 dummy.Next。", map[string]string{"answer": "1→2→3→4"}, ready("1", "2", "3", "4"), nil, nil, ready("1", "2", "3", "4"), nil, "排序完成", mergeSortTopology{original: mergeSortItems([]string{"1", "2", "3", "4"}, nil), originalLinks: []bool{true, true, true}}),
	}
	return concreteTrace("linked-list-merge-sort", "链表：归并排序（细粒度拆分与合并）", code, frames...)
}

func mergeSortItems(values []string, states map[int]string) []exampleItem {
	result := make([]exampleItem, len(values))
	for index, value := range values {
		state := "ready"
		if configured, ok := states[index]; ok {
			state = configured
		}
		result[index] = item(value, state)
	}
	return result
}

func windowCharacters(values string, states map[int]string) []greedyRangeSegment {
	result := make([]greedyRangeSegment, 0, len([]rune(values)))
	for index, character := range []rune(values) {
		state := "ready"
		if configured, ok := states[index]; ok {
			state = configured
		}
		result = append(result, makeGreedyRangeSegment(index, index+1, string(character), state, "item"))
	}
	return result
}

func windowFrame(line int, narration, caption, input string, left, right, min, max int, variables map[string]string, inputStates map[int]string, windowState, windowLabel string, windowStatus string, best *greedyRangeSegment) Frame {
	tracks := []any{makeGreedyRangeTrack("字符串位置", windowCharacters(input, inputStates)...)}
	windowEnd := right + 1
	if right < left {
		windowEnd = left
	}
	windowSegment := makeGreedyRangeSegment(left, windowEnd, windowLabel, windowStatus, "range")
	if best != nil {
		tracks = append(tracks, makeGreedyRangeTrack("当前窗口", windowSegment, *best))
	} else {
		tracks = append(tracks, makeGreedyRangeTrack("当前窗口", windowSegment))
	}
	tracks = append(tracks, makeGreedyRangeTrack("状态", makeGreedyRangeSegment(0, 1, windowState, "dependency", "item")))
	markers := []any{
		makeGreedyRangeMarker("当前窗口", "left", left, "current"),
		makeGreedyRangeMarker("当前窗口", "right", right+1, "current"),
	}
	tracks = append(tracks, markers...)
	return greedyRangeFrame(line, narration, caption, min, max, variables, tracks...)
}

func SlidingWindowExactTrace() Trace {
	input := "abcad"
	code := []string{"for right := range s {", "    add s[right]", "    for distinct > k { remove s[left]; left++ }", "    if distinct == k { record window }"}
	best := func() *greedyRangeSegment {
		value := makeGreedyRangeSegment(0, 3, "[0,2] = abc", "ready", "range")
		return &value
	}
	frames := []Frame{
		windowFrame(0, "例题 s=abcad、k=3。窗口是字符串位置轴上的连续区间；先从空窗口开始。", "恰好满足：准备", input, 0, -1, 0, 5, map[string]string{"left": "0", "right": "-1", "distinct": "0", "k": "3"}, nil, "distinct=0", "空", "current", nil),
		windowFrame(1, "右端扩张一步，把 a 纳入窗口；distinct 从 0 变为 1。", "扩张到 a", input, 0, 0, 0, 5, map[string]string{"left": "0", "right": "0", "add": "a", "distinct": "1"}, map[int]string{0: "current"}, "distinct=1", "[0,0] a", "current", nil),
		windowFrame(1, "右端再扩张一步，把 b 纳入；窗口仍连续，distinct=2，还没有达到 k。", "扩张到 ab", input, 0, 1, 0, 5, map[string]string{"left": "0", "right": "1", "add": "b", "distinct": "2"}, map[int]string{0: "dependency", 1: "current"}, "distinct=2", "[0,1] ab", "current", nil),
		windowFrame(1, "右端扩张到 c；distinct=3 恰好满足条件，此刻记录 [0,2]。", "恰好满足 abc", input, 0, 2, 0, 5, map[string]string{"left": "0", "right": "2", "add": "c", "distinct": "3", "best": "[0,2]"}, map[int]string{0: "ready", 1: "ready", 2: "current"}, "distinct=3，命中", "[0,2] abc", "ready", best()),
		windowFrame(1, "继续扩张纳入 a；a 是重复字符，distinct 仍为 3，窗口依然恰好满足。", "重复字符不改变 distinct", input, 0, 3, 0, 5, map[string]string{"left": "0", "right": "3", "add": "a", "distinct": "3"}, map[int]string{2: "dependency", 3: "current"}, "distinct=3，命中", "[0,3] abca", "current", best()),
		windowFrame(1, "右端纳入 d 后 distinct=4，条件被破坏；先停在这里，下一帧才移动 left。", "条件被破坏", input, 0, 4, 0, 5, map[string]string{"left": "0", "right": "4", "add": "d", "distinct": "4"}, map[int]string{4: "rejected"}, "distinct=4 > k", "[0,4] abcad", "rejected", best()),
		windowFrame(2, "收缩一步，移除左端 a；a 仍在窗口中出现一次，distinct 仍为 4。", "收缩移除左端 a", input, 1, 4, 0, 5, map[string]string{"left": "1", "right": "4", "remove": "a", "distinct": "4"}, map[int]string{0: "rejected", 4: "current"}, "distinct=4 > k", "[1,4] bcad", "rejected", best()),
		windowFrame(2, "再次收缩，移除 b；窗口变成 cad，distinct 回到 3，恢复恰好满足。", "收缩到恰好满足", input, 2, 4, 0, 5, map[string]string{"left": "2", "right": "4", "remove": "b", "distinct": "3", "best": "[0,2]"}, map[int]string{0: "rejected", 1: "rejected", 2: "ready", 4: "current"}, "distinct=3，命中", "[2,4] cad", "current", best()),
		windowFrame(3, "扫描结束；每次扩张或收缩都只改变一个边界，恰好满足的最长窗口是 abc，长度 3。", "恰好满足：完成", input, 0, 4, 0, 5, map[string]string{"best": "[0,2]", "length": "3"}, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready"}, "best=3", "[0,2] abc", "ready", best()),
	}
	return concreteTrace("window-range", "滑动窗口：恰好满足 k 个不同字符", code, frames...)
}

func SlidingWindowAtMostTrace() Trace {
	input := "eceba"
	code := []string{"for right := range s {", "    add s[right]", "    for distinct > k { remove s[left]; left++ }", "    best = max(best, right-left+1)"}
	frames := []Frame{
		windowFrame(0, "例题 s=eceba、k=2。目标是满足“最多 2 种字符”的最长连续窗口。", "满足最多：准备", input, 0, -1, 0, 5, map[string]string{"left": "0", "right": "-1", "distinct": "0", "best": "0"}, nil, "distinct=0 ≤ 2", "空", "current", nil),
		windowFrame(1, "扩张纳入 e；窗口 [0,0] 合法，best 更新为 1。", "扩张到 e", input, 0, 0, 0, 5, map[string]string{"left": "0", "right": "0", "distinct": "1", "best": "1"}, map[int]string{0: "current"}, "合法", "[0,0] e", "current", nil),
		windowFrame(1, "扩张纳入 c；distinct=2 仍合法，best 更新为 2。", "扩张到 ec", input, 0, 1, 0, 5, map[string]string{"left": "0", "right": "1", "distinct": "2", "best": "2"}, map[int]string{0: "ready", 1: "current"}, "合法", "[0,1] ec", "current", nil),
		windowFrame(1, "扩张纳入 e；重复字符不增加 distinct，窗口 [0,2] 继续合法，best=3。", "扩张到 ece", input, 0, 2, 0, 5, map[string]string{"left": "0", "right": "2", "distinct": "2", "best": "3"}, map[int]string{1: "dependency", 2: "current"}, "合法", "[0,2] ece", "current", nil),
		windowFrame(1, "扩张纳入 b 后 distinct=3，先停在非法窗口；收缩留给下一步。", "扩张导致超限", input, 0, 3, 0, 5, map[string]string{"left": "0", "right": "3", "distinct": "3", "best": "3"}, map[int]string{3: "rejected"}, "distinct=3 > 2", "[0,3] eceb", "rejected", nil),
		windowFrame(2, "收缩一步移除 e；窗口仍含 c、e、b 三种字符，所以还不合法。", "收缩移除 e", input, 1, 3, 0, 5, map[string]string{"left": "1", "right": "3", "remove": "e", "distinct": "3"}, map[int]string{0: "rejected", 3: "current"}, "distinct=3 > 2", "[1,3] ceb", "rejected", nil),
		windowFrame(2, "再次收缩移除 c；distinct 降为 2，窗口 eb 恢复合法，best 保持 3。", "收缩到合法", input, 2, 3, 0, 5, map[string]string{"left": "2", "right": "3", "remove": "c", "distinct": "2", "best": "3"}, map[int]string{0: "rejected", 1: "rejected", 2: "ready", 3: "current"}, "合法", "[2,3] eb", "current", nil),
		windowFrame(1, "扩张纳入 a；distinct=3，再次进入非法状态。", "扩张到 eba", input, 2, 4, 0, 5, map[string]string{"left": "2", "right": "4", "distinct": "3", "best": "3"}, map[int]string{4: "rejected"}, "distinct=3 > 2", "[2,4] eba", "rejected", nil),
		windowFrame(2, "收缩移除 e；窗口 ba 只剩两种字符，恢复合法，best=3。", "收缩到 ba", input, 3, 4, 0, 5, map[string]string{"left": "3", "right": "4", "remove": "e", "distinct": "2", "best": "3"}, map[int]string{2: "rejected", 3: "ready", 4: "current"}, "合法", "[3,4] ba", "current", nil),
		windowFrame(3, "扫描完成；最长合法区间仍是 [0,2] 的 ece，长度 3。", "满足最多：完成", input, 0, 4, 0, 5, map[string]string{"best": "3", "answer": "ece"}, map[int]string{0: "ready", 1: "ready", 2: "ready", 3: "ready", 4: "ready"}, "best=3", "[0,2] ece", "ready", nil),
	}
	return concreteTrace("window-range", "滑动窗口：最多 k 个不同字符", code, frames...)
}

func SlidingWindowMinimumTrace() Trace {
	input := "ADOBECODEBANC"
	code := []string{"for right := range s {", "    add s[right]", "    for window covers target {", "        record shorter answer", "        remove s[left]; left++", "    }", "}"}
	frames := []Frame{}
	appendFrame := func(line int, narration, caption string, left, right int, variables map[string]string, states map[int]string, status string, label string, best *greedyRangeSegment) {
		frames = append(frames, windowFrame(line, narration, caption, input, left, right, 0, len([]rune(input)), variables, states, "target=ABC", label, status, best))
	}
	appendFrame(0, "例题 s=ADOBECODEBANC、target=ABC。右端扩张寻找覆盖，覆盖后才逐步收缩左端。", "最少覆盖：准备", 0, -1, map[string]string{"left": "0", "right": "-1", "need": "A,B,C", "best": "∞"}, nil, "current", "空", nil)
	chars := []rune(input)
	for index := 0; index <= 5; index++ {
		states := map[int]string{index: "current"}
		appendFrame(1, "右端扩张一步，纳入 "+string(chars[index])+"；窗口还未覆盖 A、B、C。", "扩张到 "+string(chars[index]), 0, index, map[string]string{"left": "0", "right": itoa(index), "add": string(chars[index]), "missing": "继续寻找"}, states, "current", "[0,"+itoa(index)+"]", nil)
	}
	best := func(start, end int, label string) *greedyRangeSegment {
		value := makeGreedyRangeSegment(start, end+1, label, "ready", "range")
		return &value
	}
	appendFrame(3, "纳入 C 后，窗口 [0,5] 首次覆盖 target；先记录它，再开始收缩。", "首次覆盖 ADOBEC", 0, 5, map[string]string{"left": "0", "right": "5", "cover": "true", "best": "[0,5]"}, map[int]string{5: "current"}, "ready", "[0,5] ADOBEC", best(0, 5, "best [0,5]"))
	appendFrame(4, "覆盖成立，尝试移除左端 A；这是收缩动作本身，下一帧再判断覆盖是否仍成立。", "收缩移除 A", 1, 5, map[string]string{"left": "1", "right": "5", "remove": "A", "cover": "false"}, map[int]string{0: "rejected", 1: "dependency"}, "rejected", "[1,5] DOBEC", best(0, 5, "best [0,5]"))
	for index := 6; index <= 10; index++ {
		states := map[int]string{index: "current"}
		appendFrame(1, "左端暂时停在 1，右端扩张纳入 "+string(chars[index])+"；窗口仍未重新覆盖 A。", "继续扩张到 "+string(chars[index]), 1, index, map[string]string{"left": "1", "right": itoa(index), "add": string(chars[index]), "missing": "A"}, states, "current", "[1,"+itoa(index)+"]", best(0, 5, "best [0,5]"))
	}
	appendFrame(3, "纳入 A 后 [1,10] 覆盖 target；记录后，收缩每一步都单独展示。", "第二次覆盖", 1, 10, map[string]string{"left": "1", "right": "10", "cover": "true", "best": "[0,5]"}, map[int]string{10: "current"}, "ready", "[1,10]", best(0, 5, "best [0,5]"))
	appendFrame(4, "移除 D，窗口 [2,10] 仍覆盖；长度缩短，但还不是最短。", "收缩移除 D", 2, 10, map[string]string{"left": "2", "right": "10", "remove": "D", "cover": "true"}, map[int]string{1: "rejected", 2: "dependency"}, "ready", "[2,10]", best(0, 5, "best [0,5]"))
	appendFrame(4, "移除 O，窗口 [3,10] 仍覆盖；继续尝试收缩。", "收缩移除 O", 3, 10, map[string]string{"left": "3", "right": "10", "remove": "O", "cover": "true"}, map[int]string{2: "rejected", 3: "dependency"}, "ready", "[3,10]", best(0, 5, "best [0,5]"))
	appendFrame(4, "移除 B，覆盖失效；此刻把刚才最短的 [3,10] 记为候选。", "收缩移除 B", 4, 10, map[string]string{"left": "4", "right": "10", "remove": "B", "cover": "false", "best": "[3,10]"}, map[int]string{3: "rejected", 4: "dependency"}, "rejected", "[4,10]", best(3, 10, "best [3,10]"))
	for index := 11; index <= 12; index++ {
		states := map[int]string{index: "current"}
		appendFrame(1, "右端扩张纳入 "+string(chars[index])+"；左端保持在 4，仍缺少 A。", "扩张到 "+string(chars[index]), 4, index, map[string]string{"left": "4", "right": itoa(index), "add": string(chars[index]), "missing": "A"}, states, "current", "[4,"+itoa(index)+"]", best(3, 10, "best [3,10]"))
	}
	appendFrame(3, "纳入最后的 C 后，窗口 [4,12] 覆盖 target；进入最后一轮收缩。", "第三次覆盖", 4, 12, map[string]string{"left": "4", "right": "12", "cover": "true", "best": "[3,10]"}, map[int]string{12: "current"}, "ready", "[4,12]", best(3, 10, "best [3,10]"))
	for index := 4; index <= 9; index++ {
		states := map[int]string{index: "rejected", 12: "dependency"}
		appendFrame(4, "覆盖仍成立，移除左端 "+string(chars[index])+"；每次只左移一个位置，继续寻找更短覆盖。", "连续收缩到 "+itoa(index+1), index+1, 12, map[string]string{"left": itoa(index + 1), "right": "12", "remove": string(chars[index]), "cover": "true"}, states, "ready", "["+itoa(index+1)+",12]", best(3, 10, "best [3,10]"))
	}
	appendFrame(4, "移除 B 后，窗口 [10,12] 缺少 B；上一帧的 [9,12] 长度 4 成为新答案。", "收缩移除 B", 10, 12, map[string]string{"left": "10", "right": "12", "remove": "B", "cover": "false", "best": "[9,12]"}, map[int]string{9: "rejected", 10: "dependency"}, "rejected", "[10,12]", best(9, 12, "best [9,12]"))
	appendFrame(6, "扫描完成；最短覆盖区间是 [9,12]，对应字符串 BANC。", "最少覆盖：完成", 9, 12, map[string]string{"answer": "BANC", "length": "4"}, map[int]string{9: "ready", 10: "ready", 11: "ready", 12: "ready"}, "ready", "[9,12] BANC", best(9, 12, "answer BANC"))
	return concreteTrace("window-range", "滑动窗口：覆盖 target 的最小区间", code, frames...)
}

func PalindromeIntervalDPTrace() Trace {
	input := []rune("babad")
	n := len(input)
	values := make([][]int, n)
	done := make([][]bool, n)
	for row := range values {
		values[row] = make([]int, n)
		done[row] = make([]bool, n)
	}
	code := []string{
		"for i := 0; i < n; i++ { dp[i][i] = true }",
		"for length := 2; length <= n; length++ {",
		"    for left := 0; left+length <= n; left++ {",
		"        right := left + length - 1",
		"        dp[left][right] = s[left] == s[right] && (length <= 2 || dp[left+1][right-1])",
		"        if dp[left][right] && length > bestLength { update answer }",
		"    }",
		"}",
	}
	frames := []Frame{palindromeFrame(input, values, done, -1, -1, nil, 0, "矩阵只保留上三角：dp[l][r] 表示闭区间 s[l:r+1] 是否为回文；先准备空表。", map[string]string{"best": ""})}
	for index := 0; index < n; index++ {
		values[index][index] = 1
		done[index][index] = true
		frames = append(frames, palindromeFrame(input, values, done, index, index, nil, 0, "长度 1 的区间 ["+itoa(index)+","+itoa(index)+"] 只有一个字符，必然是回文。", map[string]string{"left": itoa(index), "right": itoa(index), "best": string(input[index])}))
	}
	bestStart, bestLength := 0, 1
	for length := 2; length <= n; length++ {
		for left := 0; left+length <= n; left++ {
			right := left + length - 1
			frames = append(frames, palindromeFrame(input, values, done, left, right, palindromeDependencies(length, left, right), 4, "准备计算区间 ["+itoa(left)+","+itoa(right)+"]；蓝色格子是本次读取的子区间依赖。", map[string]string{"left": itoa(left), "right": itoa(right), "length": itoa(length)}))
			match := input[left] == input[right] && (length <= 2 || values[left+1][right-1] == 1)
			if match {
				values[left][right] = 1
				if length > bestLength {
					bestStart, bestLength = left, length
				}
			}
			done[left][right] = true
			frames = append(frames, palindromeFrame(input, values, done, left, right, nil, 5, "字符 "+string(input[left])+" 与 "+string(input[right])+" "+map[bool]string{true: "相同，且内部区间已是回文，写入 true。", false: "不满足两端相同或内部依赖，写入 false。"}[match], map[string]string{"left": itoa(left), "right": itoa(right), "dp": map[bool]string{true: "true", false: "false"}[match], "best": string(input[bestStart : bestStart+bestLength])}))
		}
	}
	frames = append(frames, palindromeFrame(input, values, done, bestStart, bestStart+bestLength-1, nil, 7, "按长度填表结束；最长回文子串是 s["+itoa(bestStart)+":"+itoa(bestStart+bestLength)+"] = "+string(input[bestStart:bestStart+bestLength])+"。", map[string]string{"answer": string(input[bestStart : bestStart+bestLength]), "length": itoa(bestLength)}))
	return Trace{Kind: "dp-grid", Title: "区间 DP：最长回文子串", Pseudocode: code, Frames: frames}
}

func palindromeDependencies(length, left, right int) []gridPoint {
	if length <= 2 {
		return nil
	}
	return []gridPoint{{Row: left + 1, Column: right - 1}}
}

func palindromeFrame(input []rune, values [][]int, done [][]bool, currentRow, currentColumn int, dependencies []gridPoint, line int, narration string, variables map[string]string) Frame {
	dependencySet := gridPointSet(dependencies)
	state := gridState{Title: "dp[l][r]：区间 s[l:r+1] 是否为回文", Rows: make([]string, len(input)), Columns: make([]string, len(input))}
	for index, character := range input {
		state.Rows[index] = "l=" + itoa(index) + "(" + string(character) + ")"
		state.Columns[index] = "r=" + itoa(index) + "(" + string(character) + ")"
	}
	for row := range values {
		for column := range values[row] {
			cellState := "unused"
			if row <= column {
				cellState = "pending"
				if done[row][column] {
					cellState = "ready"
				}
				if row == currentRow && column == currentColumn {
					cellState = "current"
				}
			}
			state.Cells = append(state.Cells, gridCell{Row: row, Column: column, Value: values[row][column], State: cellState, Dependency: dependencySet[gridPoint{row, column}]})
		}
	}
	return Frame{ActiveLine: line, Narration: narration, Variables: variables, State: state}
}
