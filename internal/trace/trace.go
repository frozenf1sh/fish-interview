package trace

import "sort"

// Trace separates an algorithm's state transitions from its browser renderer.
type Trace struct {
	Kind       string   `json:"kind"`
	Title      string   `json:"title"`
	Pseudocode []string `json:"pseudocode"`
	Frames     []Frame  `json:"frames"`
}

type Frame struct {
	ActiveLine int               `json:"activeLine"`
	Narration  string            `json:"narration"`
	Variables  map[string]string `json:"variables"`
	State      any               `json:"state"`
}

type Interval struct {
	Label string `json:"label"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type intervalState struct {
	Intervals []intervalView `json:"intervals"`
	Cursor    int            `json:"cursor"`
}

type intervalView struct {
	Interval
	Status string `json:"status"`
}

// IntervalScheduling returns a replayable trace for the earliest-finish-time greedy rule.
func IntervalScheduling(input []Interval) Trace {
	intervals := append([]Interval(nil), input...)
	sort.SliceStable(intervals, func(i, j int) bool { return intervals[i].End < intervals[j].End })

	trace := Trace{
		Kind:  "intervals",
		Title: "区间调度：按结束时间选择",
		Pseudocode: []string{
			"按结束时间升序排序 intervals",
			"end = -∞; chosen = 0",
			"for interval in intervals:",
			"    if interval.start >= end:",
			"        选择 interval",
			"        end = interval.end",
			"return chosen",
		},
	}
	end, chosen := -1, 0
	selected := make([]bool, len(intervals))
	trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, -1, end, chosen, "", 0, "先按结束时间排序：结束更早的区间为后续保留更多空间。"))
	for i, interval := range intervals {
		trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, i, end, chosen, "candidate", 2, "考察 "+interval.Label+"：它从 "+itoa(interval.Start)+" 开始，当前可用边界是 "+itoa(end)+"。"))
		if interval.Start >= end {
			trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, i, end, chosen, "candidate", 3, "开始时间满足条件，因此它与已选区间不重叠。"))
			end, chosen = interval.End, chosen+1
			selected[i] = true
			trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, i, end, chosen, "selected", 5, "选择 "+interval.Label+"，把后续可用边界更新为 "+itoa(end)+"。"))
		} else {
			trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, i, end, chosen, "rejected", 3, interval.Label+" 与已选区间重叠；跳过它。"))
		}
	}
	trace.Frames = append(trace.Frames, intervalFrame(intervals, selected, -1, end, chosen, "", 6, "扫描结束：共选择 "+itoa(chosen)+" 个互不重叠区间。"))
	return trace
}

func intervalFrame(intervals []Interval, selected []bool, cursor, end, chosen int, cursorStatus string, line int, narration string) Frame {
	state := intervalState{Cursor: cursor, Intervals: make([]intervalView, len(intervals))}
	for i, interval := range intervals {
		status := "pending"
		if selected[i] {
			status = "selected"
		}
		if i == cursor {
			status = cursorStatus
		}
		state.Intervals[i] = intervalView{Interval: interval, Status: status}
	}
	return Frame{
		ActiveLine: line,
		Narration:  narration,
		Variables: map[string]string{
			"end":    itoa(end),
			"chosen": itoa(chosen),
		},
		State: state,
	}
}

func itoa(value int) string {
	if value == -1 {
		return "-∞"
	}
	if value == 0 {
		return "0"
	}
	digits := [20]byte{}
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = byte(value%10) + '0'
		value /= 10
	}
	return string(digits[i:])
}
