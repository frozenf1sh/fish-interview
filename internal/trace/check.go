package trace

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidatePlayerContract checks every property the browser needs before the first frame can render.
func ValidatePlayerContract(value Trace) error {
	if value.Kind == "" || value.Title == "" {
		return fmt.Errorf("trace misses kind or title")
	}
	if len(value.Pseudocode) == 0 || len(value.Frames) == 0 {
		return fmt.Errorf("trace %q misses pseudocode or frames", value.Title)
	}
	lessEntity, greaterEntity := "&"+"lt;", "&"+"gt;"
	for line, code := range value.Pseudocode {
		if strings.Contains(code, lessEntity) || strings.Contains(code, greaterEntity) {
			return fmt.Errorf("trace %q pseudocode line %d contains HTML entity", value.Title, line)
		}
	}
	for index, frame := range value.Frames {
		if frame.State == nil || frame.Narration == "" || frame.Variables == nil {
			return fmt.Errorf("trace %q frame %d misses state, narration, or variables", value.Title, index)
		}
		if frame.ActiveLine < 0 || frame.ActiveLine >= len(value.Pseudocode) {
			return fmt.Errorf("trace %q frame %d has invalid active line %d", value.Title, index, frame.ActiveLine)
		}
		if err := validateRendererState(value.Kind, frame.State); err != nil {
			return fmt.Errorf("trace %q frame %d: %w", value.Title, index, err)
		}
	}
	return nil
}

func validateRendererState(kind string, state any) error {
	requiredSlices := map[string][]string{
		"intervals":              {"intervals"},
		"dp-table":               {"cells"},
		"dp-grid":                {"rows", "columns", "cells"},
		"rolling-dependency":     {},
		"bitmask-state":          {"names"},
		"linked-list":            {"chain"},
		"binary-red-blue":        {"numbers"},
		"sequence-tails":         {"numbers", "tails", "tailStates"},
		"row-gravity":            {"cells"},
		"example-state":          {"lanes"},
		"greedy-range":           {"tracks", "markers"},
		"matrix-state":           {"cells"},
		"node-link-state":        {"nodes", "links"},
		"cycle-list-state":       {"nodes", "links"},
		"linked-list-merge":      {"left", "right", "result"},
		"linked-list-merge-sort": {"original", "originalLinks", "source", "left", "leftLinks", "right", "rightLinks", "result", "stack"},
		"linked-list-k-group":    {"chain", "detached", "working", "group"},
		"window-range":           {"tracks", "markers"},
		"flow-steps":             {"steps"},
	}
	fields, ok := requiredSlices[kind]
	if !ok {
		return fmt.Errorf("unknown renderer kind %q", kind)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return err
	}
	for _, field := range fields {
		value, exists := object[field]
		if !exists || string(value) == "null" {
			return fmt.Errorf("state field %q must be an array, not null", field)
		}
	}
	return nil
}

// ValidateJSONRoundTrip protects the initial browser render, which sees decoded JSON rather than Go values.
func ValidateJSONRoundTrip(value Trace) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	var decoded Trace
	if err := json.Unmarshal(data, &decoded); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return ValidatePlayerContract(decoded)
}

// AllTraces returns every trace exposed by the HTTP server, including concrete card traces and flow fallbacks.
func AllTraces() map[string]Trace {
	result := map[string]Trace{
		"interval-scheduling":    IntervalScheduling([]Interval{{Label: "A", Start: 1, End: 3}, {Label: "B", Start: 2, End: 5}, {Label: "C", Start: 3, End: 6}, {Label: "D", Start: 5, End: 7}, {Label: "E", Start: 6, End: 8}}),
		"linear-dp":              LinearDPClimbStairs(),
		"space-rolling":          SpaceOptimizationTrace(),
		"lcs-dp":                 LCSTrace(),
		"interval-dp":            IntervalMergeTrace(),
		"stock-dp":               StockTrace(),
		"bitmask-dp":             BitmaskTrace(),
		"linked-list-rewire":     LinkedListRewireTrace(),
		"list-merge-sort":        ListMergeSortTrace(),
		"list-k-group":           LinkedListKGroupTrace(),
		"path-dp":                PathTrace(),
		"reverse-path-dp":        ReversePathTrace(),
		"binary-red-blue":        BinaryRedBluePartition(),
		"lis":                    LISTrace(),
		"row-gravity":            RowGravityTrace(),
		"interval-start-merge":   StartSortedIntervalsTrace(),
		"meeting-rooms":          MeetingRoomsTrace(),
		"weighted-intervals":     WeightedIntervalsTrace(),
		"kadane":                 KadaneTrace(),
		"sliding-window-exact":   SlidingWindowExactTrace(),
		"sliding-window-at-most": SlidingWindowAtMostTrace(),
		"sliding-window-minimum": SlidingWindowMinimumTrace(),
		"palindrome-interval-dp": PalindromeIntervalDPTrace(),
	}
	for name := range flowSpecs {
		if value, ok := FlowTrace(name); ok {
			result[name] = value
		}
	}
	return result
}

// ValidateAllPlayerContracts returns one actionable error per invalid named trace.
func ValidateAllPlayerContracts() map[string]error {
	failures := map[string]error{}
	for name, value := range AllTraces() {
		if err := ValidatePlayerContract(value); err != nil {
			failures[name] = err
			continue
		}
		if err := ValidateJSONRoundTrip(value); err != nil {
			failures[name] = err
		}
	}
	return failures
}
