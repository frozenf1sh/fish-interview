package trace

import "testing"

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

	stock := StockTrace()
	if len(stock.Frames) < 4 || stock.Frames[len(stock.Frames)-1].Variables["cash"] != "7" {
		t.Fatalf("unexpected stock trace: %#v", stock)
	}
	bitmask := BitmaskTrace()
	last := bitmask.Frames[len(bitmask.Frames)-1]
	if last.Variables["mask"] != "1111" || last.Variables["cost"] != "18" {
		t.Fatalf("unexpected bitmask trace: %#v", bitmask)
	}
}
