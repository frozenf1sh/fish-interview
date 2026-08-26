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
