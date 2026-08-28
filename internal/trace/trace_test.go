package trace

import (
	"strings"
	"testing"
)

func TestIntervalSchedulingTrace(t *testing.T) {
	got := IntervalScheduling([]Interval{
		{Label: "B", Start: 2, End: 5},
		{Label: "A", Start: 1, End: 3},
		{Label: "C", Start: 3, End: 6},
		{Label: "D", Start: 5, End: 7},
		{Label: "E", Start: 6, End: 8},
	})
	if got.Kind != "intervals" || len(got.Frames) < 10 {
		t.Fatalf("unexpected trace: %#v", got)
	}
	last := got.Frames[len(got.Frames)-1]
	if last.Variables["chosen"] != "3" {
		t.Fatalf("chosen = %q, want 3", last.Variables["chosen"])
	}
}

func TestLinearDPTrace(t *testing.T) {
	got := LinearDPClimbStairs()
	if got.Kind != "dp-table" || got.Frames[len(got.Frames)-1].Variables["i"] != "5" {
		t.Fatalf("unexpected trace: %#v", got)
	}
}

func TestLinkedListRewireTrace(t *testing.T) {
	trace := LinkedListRewireTrace()
	detached, ok := trace.Frames[4].State.(linkedListState)
	if !ok || strings.Join(detached.Chain, ",") != "D,1,2,4,5" || strings.Join(detached.Detached, ",") != "3" {
		t.Fatalf("trace must expose the detach step: %#v", trace.Frames[4].State)
	}
	last, ok := trace.Frames[len(trace.Frames)-1].State.(linkedListState)
	if !ok || strings.Join(last.Chain, ",") != "D,1,4,3,2,5" {
		t.Fatalf("unexpected linked-list trace: %#v", trace.Frames[len(trace.Frames)-1])
	}
}

func TestLinkedListKGroupTrace(t *testing.T) {
	value := LinkedListKGroupTrace()
	if value.Kind != "linked-list-k-group" || len(value.Frames) < 20 {
		t.Fatalf("unexpected k-group trace: %#v", value)
	}
	detached, ok := value.Frames[5].State.(kGroupListState)
	if !ok || strings.Join(detached.Chain, ",") != "D,4,5" || strings.Join(detached.Detached, ",") != "1,2,3" {
		t.Fatalf("trace must expose the detached group: %#v", value.Frames[5].State)
	}
	last, ok := value.Frames[len(value.Frames)-1].State.(kGroupListState)
	if !ok || strings.Join(last.Chain, ",") != "D,3,2,1,4,5" {
		t.Fatalf("unexpected k-group final chain: %#v", value.Frames[len(value.Frames)-1])
	}
}

func TestFlowTraceHasReplayableSteps(t *testing.T) {
	got, ok := FlowTrace("flow-bfs-shortest-path")
	if !ok || got.Kind != "node-link-state" || len(got.Frames) < 5 {
		t.Fatalf("unexpected flow trace: %#v, ok=%v", got, ok)
	}
	state, ok := got.Frames[1].State.(nodeLinkState)
	if !ok || len(state.Nodes) != 4 || state.Nodes[0].State != "current" {
		t.Fatalf("unexpected flow state: %#v", got.Frames[1].State)
	}
	if _, ok := FlowTrace("flow-not-found"); ok {
		t.Fatal("unknown flow trace should not resolve")
	}
}

func TestSequenceAndSimulationTracesReachExpectedStates(t *testing.T) {
	lis := LISTrace()
	dependency, ok := lis.Frames[3].State.(sequenceState)
	if !ok || len(dependency.TailStates) != 1 || dependency.TailStates[0] != "dependency" {
		t.Fatalf("LIS should mark the inspected tail blue: %#v", lis.Frames[3].State)
	}
	lisState, ok := lis.Frames[len(lis.Frames)-1].State.(sequenceState)
	if !ok || strings.Join([]string{itoa(lisState.Tails[0]), itoa(lisState.Tails[1]), itoa(lisState.Tails[2]), itoa(lisState.Tails[3])}, ",") != "2,3,7,18" {
		t.Fatalf("unexpected LIS trace: %#v", lis.Frames[len(lis.Frames)-1])
	}
	gravity := RowGravityTrace()
	gravityState, ok := gravity.Frames[len(gravity.Frames)-1].State.(gravityState)
	if !ok || strings.Join(gravityState.Cells, "") != "..#*.#" {
		t.Fatalf("unexpected gravity trace: %#v", gravity.Frames[len(gravity.Frames)-1])
	}
}

func TestEveryExposedTraceMeetsPlayerContract(t *testing.T) {
	if failures := ValidateAllPlayerContracts(); len(failures) != 0 {
		t.Fatalf("invalid traces: %#v", failures)
	}
}

func TestFormerFlowCardsUseConcreteExampleStates(t *testing.T) {
	for name := range flowSpecs {
		trace, ok := FlowTrace(name)
		if !ok || trace.Kind == "flow-steps" || len(trace.Frames) < 4 {
			t.Fatalf("trace %s is not a concrete replay: %#v", name, trace)
		}
	}
}

func TestRedesignedFlowTracesKeepConcreteTransitions(t *testing.T) {
	tests := []struct {
		name         string
		intermediate string
		final        string
		minFrames    int
	}{
		{"flow-greedy-reachability", "0+2=2", "返回 true", 8},
		{"flow-greedy-lexicographic", "执行 pop", "最后的 c 再次跳过", 12},
		{"flow-greedy-interval-endpoints", "第二组右端", "扫描结束", 8},
		{"flow-bfs-shortest-path", "写 dist[B]=1", "队列为空", 10},
		{"flow-bfs-multi-source", "写入距离 1", "所有格子的", 10},
		{"flow-bfs-topological", "入度从 1 减到 0", "所有 4 个", 8},
		{"flow-dfs-tree", "sum(6)", "根拿到", 15},
		{"flow-dfs-grid", "写 (0,1)=0", "递归逐层返回", 12},
		{"flow-dfs-path", "[A,B,D]", "所有邻居", 10},
		{"flow-backtracking-choose-skip", "选择 2", "现场没有", 12},
		{"flow-backtracking-enumeration", "写 used[0]", "所有分支结束", 14},
		{"flow-list-fast-slow", "同时移动", "slow==fast", 9},
		{"flow-list-merge", "选择 A 的 1", "完整有序链", 15},
		{"flow-tree-bst", "必须满足根传下来的下界", "返回 false", 6},
		{"flow-tree-lca", "命中 q", "最近公共祖先", 6},
		{"flow-tree-path-sum", "remain=11", "父调用短路", 7},
		{"flow-tree-dp", "take(3)", "最终答案", 14},
		{"flow-string-window", "扩张纳入 b", "完成", 10},
		{"flow-string-golang", "[]rune", "得到 Go中", 7},
		{"flow-string-palindrome", "写入 true", "最长回文子串", 20},
		{"flow-string-kmp", "pi[2]=1", "报告匹配", 9},
		{"flow-lcs-space", "读取覆盖前 dp[1]=0", "一维数组最终", 11},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trace, ok := FlowTrace(test.name)
			if !ok {
				t.Fatal("trace not found")
			}
			intermediateFound := false
			for _, frame := range trace.Frames[:len(trace.Frames)-1] {
				intermediateFound = intermediateFound || strings.Contains(frame.Narration, test.intermediate)
			}
			if !intermediateFound {
				t.Fatalf("missing intermediate step %q", test.intermediate)
			}
			if len(trace.Frames) < test.minFrames {
				t.Fatalf("trace has %d frames, want at least %d", len(trace.Frames), test.minFrames)
			}
			if last := trace.Frames[len(trace.Frames)-1]; !strings.Contains(last.Narration, test.final) {
				t.Fatalf("final frame %q misses %q", last.Narration, test.final)
			}
		})
	}
}

func TestRedesignedAdditionalTracesKeepDecisionFrames(t *testing.T) {
	tests := []struct {
		trace        Trace
		intermediate string
		final        string
	}{
		{StartSortedIntervalsTrace(), "2<=3", "追加后得到"},
		{MeetingRoomsTrace(), "5<10", "复用释放"},
		{WeightedIntervalsTrace(), "take=100", "写 dp[3]=100"},
		{KadaneTrace(), "重开得到更优", "扫描结束"},
	}
	for _, test := range tests {
		t.Run(test.trace.Title, func(t *testing.T) {
			intermediateFound := false
			for _, frame := range test.trace.Frames[:len(test.trace.Frames)-1] {
				intermediateFound = intermediateFound || strings.Contains(frame.Narration, test.intermediate)
			}
			if !intermediateFound {
				t.Fatalf("missing decision frame %q", test.intermediate)
			}
			if last := test.trace.Frames[len(test.trace.Frames)-1]; !strings.Contains(last.Narration, test.final) {
				t.Fatalf("final frame %q misses %q", last.Narration, test.final)
			}
		})
	}
}

func TestGreedyTracesKeepOneStableRangeLayout(t *testing.T) {
	traces := []Trace{
		FlowTraceMust("flow-greedy-reachability"),
		FlowTraceMust("flow-greedy-lexicographic"),
		FlowTraceMust("flow-greedy-interval-endpoints"),
		StartSortedIntervalsTrace(), MeetingRoomsTrace(), WeightedIntervalsTrace(), KadaneTrace(),
	}
	for _, trace := range traces {
		t.Run(trace.Title, func(t *testing.T) {
			if trace.Kind != "greedy-range" {
				t.Fatalf("kind = %q, want greedy-range", trace.Kind)
			}
			first, ok := trace.Frames[0].State.(greedyRangeState)
			if !ok || len(first.Tracks) == 0 {
				t.Fatalf("missing range tracks: %#v", trace.Frames[0].State)
			}
			for index, frame := range trace.Frames[1:] {
				state, ok := frame.State.(greedyRangeState)
				if !ok || len(state.Tracks) != len(first.Tracks) {
					t.Fatalf("frame %d replaces the main range layout: %#v", index+1, frame.State)
				}
				for trackIndex, track := range state.Tracks {
					if track.Label != first.Tracks[trackIndex].Label {
						t.Fatalf("frame %d changes track %d from %q to %q", index+1, trackIndex, first.Tracks[trackIndex].Label, track.Label)
					}
				}
			}
		})
	}
}

func FlowTraceMust(name string) Trace {
	trace, ok := FlowTrace(name)
	if !ok {
		panic("missing flow trace " + name)
	}
	return trace
}

func TestBinaryRedBlueTrace(t *testing.T) {
	got := BinaryRedBluePartition()
	last := got.Frames[len(got.Frames)-1]
	if got.Kind != "binary-red-blue" || last.Variables["blue"] != "18" || last.Variables["red"] != "17" {
		t.Fatalf("unexpected trace: %#v", got)
	}
	firstState, ok := got.Frames[0].State.(redBlueState)
	if !ok || firstState.Minimum != 9 || firstState.Maximum != 32 {
		t.Fatalf("unexpected binary range: %#v", got.Frames[0].State)
	}
}

func TestDPPatternTracesReachExpectedResults(t *testing.T) {
	assertGridValue := func(t *testing.T, trace Trace, row, column, want int) {
		t.Helper()
		if len(trace.Frames) < 4 {
			t.Fatalf("too few frames: %#v", trace)
		}
		state, ok := trace.Frames[len(trace.Frames)-1].State.(gridState)
		if !ok {
			t.Fatalf("state type = %T, want gridState", trace.Frames[len(trace.Frames)-1].State)
		}
		for _, cell := range state.Cells {
			if cell.Row == row && cell.Column == column && cell.Value == want {
				return
			}
		}
		t.Fatalf("cell (%d,%d) = not %d: %#v", row, column, want, state.Cells)
	}
	assertGridValue(t, LCSTrace(), 5, 3, 3)
	assertGridValue(t, IntervalMergeTrace(), 0, 3, 22)
	assertGridValue(t, PathTrace(), 2, 2, 7)
	assertGridValue(t, ReversePathTrace(), 0, 0, 7)
	assertGridValue(t, StockTrace(), 1, 5, 7)

	stock := StockTrace()
	if len(stock.Frames) < 4 || stock.Frames[len(stock.Frames)-1].Variables["cash"] != "7" {
		t.Fatalf("unexpected stock trace: %#v", stock)
	}
	bitmask := BitmaskTrace()
	last := bitmask.Frames[len(bitmask.Frames)-1]
	if last.Variables["mask"] != "1111" || last.Variables["last"] != "D" || last.Variables["cost"] != "6" {
		t.Fatalf("unexpected bitmask trace: %#v", bitmask)
	}
	rolling := SpaceOptimizationTrace()
	rollingState, ok := rolling.Frames[len(rolling.Frames)-1].State.(rollingState)
	if !ok || rollingState.PreviousOne != 8 {
		t.Fatalf("unexpected rolling trace: %#v", rolling.Frames[len(rolling.Frames)-1])
	}
}

func TestIntervalTraceShowsEverySplit(t *testing.T) {
	trace := IntervalMergeTrace()
	splits := map[string]bool{}
	for _, frame := range trace.Frames {
		if strings.Contains(frame.Narration, "枚举 k=") {
			splits[frame.Narration] = true
		}
	}
	if len(splits) != 10 {
		t.Fatalf("split frames = %d, want 10: %#v", len(splits), splits)
	}
}

func TestDPTracesMarkDependencies(t *testing.T) {
	linear := LinearDPClimbStairs()
	linearState, ok := linear.Frames[3].State.(dpTableState)
	if !ok || !linearState.Cells[0].Dependency || !linearState.Cells[1].Dependency {
		t.Fatalf("fibonacci trace should mark dp[0] and dp[1]: %#v", linear.Frames[3].State)
	}
	for _, trace := range []Trace{LCSTrace(), IntervalMergeTrace(), PathTrace(), ReversePathTrace(), StockTrace()} {
		found := false
		for _, frame := range trace.Frames {
			state, ok := frame.State.(gridState)
			if !ok {
				continue
			}
			for _, cell := range state.Cells {
				found = found || cell.Dependency
			}
		}
		if !found {
			t.Fatalf("trace %q never marks a dependency", trace.Title)
		}
	}
	bitmask := BitmaskTrace()
	bitmaskState, ok := bitmask.Frames[3].State.(bitmaskState)
	if !ok || bitmaskState.PreviousLast != 0 {
		t.Fatalf("bitmask trace should retain the previous last city: %#v", bitmask.Frames[3].State)
	}
}

func TestRedesignedTracesExposeStateSpecificModels(t *testing.T) {
	checks := []struct {
		name string
		kind string
		min  int
	}{
		{"flow-dfs-tree", "node-link-state", 15},
		{"flow-list-fast-slow", "cycle-list-state", 9},
		{"flow-list-merge", "linked-list-merge", 15},
		{"list-merge-sort", "linked-list-merge-sort", 25},
		{"list-k-group", "linked-list-k-group", 20},
		{"sliding-window-exact", "window-range", 9},
		{"sliding-window-at-most", "window-range", 10},
		{"sliding-window-minimum", "window-range", 25},
		{"palindrome-interval-dp", "dp-grid", 20},
		{"flow-bfs-shortest-path", "node-link-state", 10},
		{"flow-tree-dp", "node-link-state", 14},
		{"flow-string-kmp", "example-state", 9},
		{"bitmask-dp", "bitmask-state", 12},
		{"binary-red-blue", "binary-red-blue", 11},
		{"lis", "sequence-tails", 18},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			value, ok := AllTraces()[check.name]
			if !ok || value.Kind != check.kind || len(value.Frames) < check.min {
				t.Fatalf("trace=%#v", value)
			}
			for index, frame := range value.Frames {
				if frame.Narration == "" || frame.Variables == nil || frame.State == nil {
					t.Fatalf("frame %d is not observable: %#v", index, frame)
				}
			}
		})
	}
}

func TestNewAnimationTracesExposeStableIntermediateAndFinalStates(t *testing.T) {
	merge := ListMergeSortTrace()
	first, ok := merge.Frames[0].State.(mergeSortListState)
	if !ok || len(first.Source) != 4 || len(merge.Frames) < 25 {
		t.Fatalf("merge sort should expose the full source and many small steps: %#v", merge)
	}
	last, ok := merge.Frames[len(merge.Frames)-1].State.(mergeSortListState)
	if !ok || strings.Join(exampleLabels(last.Result), "→") != "1→2→3→4" {
		t.Fatalf("merge sort final result = %#v", last)
	}
	var cut, moved, overlay mergeSortListState
	for _, frame := range merge.Frames {
		state, ok := frame.State.(mergeSortListState)
		if !ok {
			continue
		}
		switch state.Phase {
		case "移除中点箭头":
			cut = state
		case "节点移动到左右子链":
			moved = state
		case "临时链覆盖原链":
			overlay = state
		}
	}
	if len(cut.OriginalLinks) != 3 || cut.OriginalLinks[1] || len(moved.Left) != 2 || len(moved.Right) != 2 {
		t.Fatalf("merge sort should cut the original arrow before moving nodes: cut=%#v moved=%#v", cut, moved)
	}
	if !overlay.Overlay || strings.Join(exampleLabels(overlay.Original), "→") != "1→2→3→4" || strings.Join(exampleLabels(overlay.Result), "→") != "1→2→3→4" {
		t.Fatalf("merge sort should animate the completed temporary chain over the original chain: %#v", overlay)
	}

	window := SlidingWindowMinimumTrace()
	initial, ok := window.Frames[0].State.(greedyRangeState)
	if !ok || len(initial.Tracks) != 3 || len(initial.Markers) != 2 {
		t.Fatalf("window should start with a fixed interval layout: %#v", initial)
	}
	if window.Kind != "window-range" || len(window.Frames) < 25 {
		t.Fatalf("window trace should retain expansion and shrink frames: %#v", window)
	}

	palindrome := PalindromeIntervalDPTrace()
	palindromeState, ok := palindrome.Frames[len(palindrome.Frames)-1].State.(gridState)
	if !ok || !gridCellValue(palindromeState.Cells, 0, 2) {
		t.Fatalf("palindrome interval dp should finish with dp[0][2]=true: %#v", palindromeState)
	}
}

func exampleLabels(items []exampleItem) []string {
	labels := make([]string, len(items))
	for index, value := range items {
		labels[index] = value.Label
	}
	return labels
}

func gridCellValue(cells []gridCell, row, column int) bool {
	for _, cell := range cells {
		if cell.Row == row && cell.Column == column {
			return cell.Value == 1 && (cell.State == "ready" || cell.State == "current")
		}
	}
	return false
}
