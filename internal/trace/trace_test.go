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
	last, ok := trace.Frames[len(trace.Frames)-1].State.(linkedListState)
	if !ok || strings.Join(last.Chain, ",") != "D,1,4,3,2,5" {
		t.Fatalf("unexpected linked-list trace: %#v", trace.Frames[len(trace.Frames)-1])
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
