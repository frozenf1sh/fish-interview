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
	Caption string       `json:"caption"`
	Nodes   []nodeVisual `json:"nodes"`
	Links   []nodeLink   `json:"links"`
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
