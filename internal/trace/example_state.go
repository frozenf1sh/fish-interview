package trace

type exampleItem struct {
	Label string `json:"label"`
	State string `json:"state"`
}

type exampleLane struct {
	Label string        `json:"label"`
	Items []exampleItem `json:"items"`
}

type exampleState struct {
	Caption string        `json:"caption"`
	Lanes   []exampleLane `json:"lanes"`
}

// greedyRangeState keeps one stable coordinate system while the greedy
// decision moves through the example. Segments are ranges on the same axis;
// markers identify the current scan position or boundary without replacing
// the whole visual with a new list of values.
type greedyRangeState struct {
	Caption string              `json:"caption"`
	Min     int                 `json:"min"`
	Max     int                 `json:"max"`
	Tracks  []greedyRangeTrack  `json:"tracks"`
	Markers []greedyRangeMarker `json:"markers"`
}

type greedyRangeTrack struct {
	Label    string               `json:"label"`
	Segments []greedyRangeSegment `json:"segments"`
}

type greedyRangeSegment struct {
	Label string `json:"label"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	State string `json:"state"`
	Kind  string `json:"kind"`
}

type greedyRangeMarker struct {
	Track    string `json:"track"`
	Label    string `json:"label"`
	Position int    `json:"position"`
	State    string `json:"state"`
}

type matrixCell struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Label  string `json:"label"`
	State  string `json:"state"`
}

type matrixState struct {
	Caption string       `json:"caption"`
	Rows    int          `json:"rows"`
	Columns int          `json:"columns"`
	Cells   []matrixCell `json:"cells"`
}

type nodeLink struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type nodeVisual struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	State string  `json:"state"`
}

type nodeLinkState struct {
	Caption     string            `json:"caption"`
	Nodes       []nodeVisual      `json:"nodes"`
	Links       []nodeLink        `json:"links"`
	ActiveLinks []nodeLink        `json:"activeLinks,omitempty"`
	CallStack   []string          `json:"callStack,omitempty"`
	Values      map[string]string `json:"values,omitempty"`
	Path        []string          `json:"path,omitempty"`
}

type cycleListState struct {
	Caption   string            `json:"caption"`
	Nodes     []nodeVisual      `json:"nodes"`
	Links     []nodeLink        `json:"links"`
	Pointers  map[string]string `json:"pointers"`
	CallStack []string          `json:"callStack,omitempty"`
}

type mergeListState struct {
	Caption string        `json:"caption"`
	Left    []exampleItem `json:"left"`
	Right   []exampleItem `json:"right"`
	Result  []exampleItem `json:"result"`
	Tail    string        `json:"tail"`
	Chosen  string        `json:"chosen"`
}

type mergeSortListState struct {
	Caption  string        `json:"caption"`
	Original []exampleItem `json:"original"`
	Source   []exampleItem `json:"source"`
	Left     []exampleItem `json:"left"`
	Right    []exampleItem `json:"right"`
	Result   []exampleItem `json:"result"`
	Active   []string      `json:"active"`
	Stack    []string      `json:"stack"`
	Phase    string        `json:"phase"`
}

type kGroupListState struct {
	Caption   string            `json:"caption"`
	Chain     []string          `json:"chain"`
	Detached  []string          `json:"detached"`
	Working   []string          `json:"working"`
	Pointers  map[string]string `json:"pointers"`
	Highlight []string          `json:"highlight"`
	Group     []string          `json:"group"`
	Phase     string            `json:"phase"`
}

func item(label, state string) exampleItem { return exampleItem{Label: label, State: state} }

func lane(label string, items ...exampleItem) exampleLane {
	return exampleLane{Label: label, Items: items}
}

func exampleFrame(line int, narration, caption string, lanes ...exampleLane) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"example": caption}, State: exampleState{Caption: caption, Lanes: lanes}}
}

func matrixFrame(line int, narration, caption string, rows, columns int, cells []matrixCell) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"example": caption}, State: matrixState{Caption: caption, Rows: rows, Columns: columns, Cells: cells}}
}

func nodeFrame(line int, narration, caption string, nodes []nodeVisual, links []nodeLink) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"example": caption}, State: nodeLinkState{Caption: caption, Nodes: nodes, Links: links}}
}

func nodeFrameDetail(line int, narration, caption string, variables map[string]string, nodes []nodeVisual, links, activeLinks []nodeLink, stack, path []string, values map[string]string) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: nodeLinkState{
		Caption: caption, Nodes: nodes, Links: links, ActiveLinks: activeLinks,
		CallStack: append([]string(nil), stack...), Path: append([]string(nil), path...), Values: cloneStringMap(values),
	}}
}

func matrixFrameDetail(line int, narration, caption string, variables map[string]string, rows, columns int, cells []matrixCell) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: matrixState{Caption: caption, Rows: rows, Columns: columns, Cells: cells}}
}

func cycleListFrame(line int, narration, caption string, variables map[string]string, nodes []nodeVisual, links []nodeLink, pointers map[string]string, stack []string) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: cycleListState{Caption: caption, Nodes: nodes, Links: links, Pointers: cloneStringMap(pointers), CallStack: append([]string(nil), stack...)}}
}

func mergeListFrame(line int, narration, caption string, variables map[string]string, left, right, result []exampleItem, tail, chosen string) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: mergeListState{Caption: caption, Left: append([]exampleItem{}, left...), Right: append([]exampleItem{}, right...), Result: append([]exampleItem{}, result...), Tail: tail, Chosen: chosen}}
}

func mergeSortListFrame(line int, narration, caption string, variables map[string]string, source, left, right, result []exampleItem, stack []string, phase string) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: mergeSortListState{
		Caption:  caption,
		Original: mergeSortItems([]string{"4", "2", "1", "3"}, nil),
		Source:   append([]exampleItem{}, source...),
		Left:     append([]exampleItem{}, left...),
		Right:    append([]exampleItem{}, right...),
		Result:   append([]exampleItem{}, result...),
		Active:   mergeSortActive(source, left, right),
		Stack:    append([]string{}, stack...),
		Phase:    phase,
	}}
}

func mergeSortActive(source, left, right []exampleItem) []string {
	active := make([]string, 0, len(source)+len(left)+len(right))
	appendUnique := func(items []exampleItem) {
		for _, value := range items {
			found := false
			for _, current := range active {
				if current == value.Label {
					found = true
					break
				}
			}
			if !found {
				active = append(active, value.Label)
			}
		}
	}
	if len(left)+len(right) > 0 {
		appendUnique(left)
		appendUnique(right)
	} else {
		appendUnique(source)
	}
	return active
}

func kGroupListFrame(line int, narration, caption string, variables map[string]string, chain, detached, working []string, pointers map[string]string, highlight, group []string, phase string) Frame {
	copyVariables := make(map[string]string, len(variables)+1)
	for key, value := range variables {
		copyVariables[key] = value
	}
	copyVariables["example"] = caption
	return Frame{ActiveLine: line, Narration: narration, Variables: copyVariables, State: kGroupListState{
		Caption: caption, Chain: append([]string{}, chain...), Detached: append([]string{}, detached...), Working: append([]string{}, working...),
		Pointers: cloneStringMap(pointers), Highlight: append([]string{}, highlight...), Group: append([]string{}, group...), Phase: phase,
	}}
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
