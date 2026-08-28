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

func TestConcretePatternTracesKeepKeyExampleSteps(t *testing.T) {
	tests := []struct {
		name         string
		intermediate string
		final        string
	}{
		{"flow-greedy-reachability", "计算 0+2=2", "返回 true"},
		{"flow-greedy-lexicographic", "c>b", "结果 acdb"},
		{"flow-greedy-interval-endpoints", "right=min(6,8)=6", "right=min(12,16)=12"},
		{"flow-bfs-shortest-path", "写 dist[B]=1", "写 dist[D]=2"},
		{"flow-bfs-multi-source", "写距离 1", "队列层数保证"},
		{"flow-bfs-topological", "B 的入度 1→0", "A,B,C,D"},
		{"flow-dfs-tree", "sum(6)=6", "sum(5)=6+2+5=13"},
		{"flow-dfs-grid", "写成 0", "整座岛都已置为 0"},
		{"flow-dfs-path", "path=[A,B,D]", "第二条路径 A,C,D"},
		{"flow-backtracking-choose-skip", "写 path=[2]", "[1,2]"},
		{"flow-backtracking-enumeration", "used[0]=true", "恢复 used"},
		{"flow-list-fast-slow", "slow 从 1 走到 2", "返回 true"},
		{"flow-list-merge", "tail.Next=1", "剩余 6"},
		{"flow-tree-bst", "允许范围是 (5,+∞)", "错误一侧"},
		{"flow-tree-lca", "直接返回 1", "最近公共祖先"},
		{"flow-tree-path-sum", "remain 从 14 写成 11", "返回 true"},
		{"flow-tree-dp", "take：选择 3", "答案取 max"},
		{"flow-string-window", "频次 a 从 1 变为 2", "最大长度保持 3"},
		{"flow-string-golang", "[]rune", "WriteString(Go)"},
		{"flow-string-palindrome", "a 与 a", "奇数扩张得到 bab"},
		{"flow-string-kmp", "j 回退到 0", "报告起点 2"},
		{"flow-lcs-space", "读取覆盖前 dp[1]=0", "答案是 dp[2]=1"},
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
			if last := trace.Frames[len(trace.Frames)-1]; !strings.Contains(last.Narration, test.final) {
				t.Fatalf("final frame %q misses %q", last.Narration, test.final)
			}
		})
	}
}

func TestAdditionalGreedyTracesKeepTheirDecisionFrames(t *testing.T) {
	tests := []struct {
		trace        Trace
		intermediate string
		final        string
	}{
		{StartSortedIntervalsTrace(), "开始 2<=蓝色 last.end=3", "最终得到 [[1,6],[8,10],[15,18]]"},
		{MeetingRoomsTrace(), "5<10", "rooms 不增加"},
		{WeightedIntervalsTrace(), "take=100", "最终 dp[3]=100"},
		{KadaneTrace(), "直接从 1 重开", "best 写成 6"},
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
	if last.Variables["mask"] != "1111" || last.Variables["cost"] != "18" {
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
	bitmaskState, ok := bitmask.Frames[1].State.(bitmaskState)
	if !ok || bitmaskState.PreviousLast != 0 {
		t.Fatalf("bitmask trace should retain the previous last city: %#v", bitmask.Frames[1].State)
	}
}
