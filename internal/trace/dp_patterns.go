package trace

type gridCell struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Value      int    `json:"value"`
	State      string `json:"state"`
	Dependency bool   `json:"dependency"`
}

type gridPoint struct{ Row, Column int }

type gridState struct {
	Title   string     `json:"title"`
	Rows    []string   `json:"rows"`
	Columns []string   `json:"columns"`
	Cells   []gridCell `json:"cells"`
}

// LCSTrace fills the table for the longest common subsequence of abcde and ace.
func LCSTrace() Trace {
	a, b := "abcde", "ace"
	values := make([][]int, len(a)+1)
	done := make([][]bool, len(a)+1)
	for i := range values {
		values[i] = make([]int, len(b)+1)
		done[i] = make([]bool, len(b)+1)
		done[i][0] = true
	}
	for j := range done[0] {
		done[0][j] = true
	}
	result := Trace{
		Kind:  "dp-grid",
		Title: "LCS：两个前缀的 DP 表",
		Pseudocode: []string{
			"for i := 1; i <= len(a); i++ {",
			"    for j := 1; j <= len(b); j++ {",
			"        if a[i-1] == b[j-1] { dp[i][j] = dp[i-1][j-1] + 1 }",
			"        else { dp[i][j] = max(dp[i-1][j], dp[i][j-1]) }",
			"    }",
			"}",
		},
	}
	result.Frames = append(result.Frames, gridFrame("dp[i][j]：两个前缀的 LCS 长度", labels("∅", a), labels("∅", b), values, done, -1, -1, nil, 0, "第 0 行和第 0 列表示一个前缀为空，LCS 长度固定为 0。"))
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				values[i][j] = values[i-1][j-1] + 1
				done[i][j] = true
				result.Frames = append(result.Frames, gridFrame("dp[i][j]：两个前缀的 LCS 长度", labels("∅", a), labels("∅", b), values, done, i, j, []gridPoint{{i - 1, j - 1}}, 2, "字符 "+string(a[i-1])+" 相同：在左上角结果后接上它。"))
				continue
			}
			values[i][j] = max(values[i-1][j], values[i][j-1])
			done[i][j] = true
			result.Frames = append(result.Frames, gridFrame("dp[i][j]：两个前缀的 LCS 长度", labels("∅", a), labels("∅", b), values, done, i, j, []gridPoint{{i - 1, j}, {i, j - 1}}, 3, "字符 "+string(a[i-1])+" 与 "+string(b[j-1])+" 不同：保留上方或左方的较优结果。"))
		}
	}
	result.Frames = append(result.Frames, gridFrame("dp[i][j]：两个前缀的 LCS 长度", labels("∅", a), labels("∅", b), values, done, len(a), len(b), []gridPoint{{len(a) - 1, len(b) - 1}}, 5, "完整前缀的 LCS 长度为 "+itoa(values[len(a)][len(b)])+"。"))
	return result
}

// IntervalMergeTrace fills intervals in increasing length order for stone merging.
func IntervalMergeTrace() Trace {
	stones := []int{3, 5, 1, 2}
	n := len(stones)
	prefix := make([]int, n+1)
	values := make([][]int, n)
	done := make([][]bool, n)
	for i := range values {
		values[i] = make([]int, n)
		done[i] = make([]bool, n)
		done[i][i] = true
		prefix[i+1] = prefix[i] + stones[i]
	}
	result := Trace{
		Kind:  "dp-grid",
		Title: "区间 DP：按长度填充上三角",
		Pseudocode: []string{
			"for length := 2; length <= n; length++ {",
			"    for l := 0; l+length <= n; l++ {",
			"        r := l + length - 1",
			"        dp[l][r] = inf",
			"        for k := l; k < r; k++ {",
			"            mergeCost := prefix[r+1] - prefix[l]",
			"            candidate := dp[l][k] + dp[k+1][r] + mergeCost",
			"            dp[l][r] = min(dp[l][r], candidate)",
			"        }",
			"    }",
			"}",
		},
	}
	result.Frames = append(result.Frames, intervalFrameGrid(values, done, -1, -1, nil, 0, "长度为 1 的区间无需合并，主对角线初始化为 0。"))
	for length := 2; length <= n; length++ {
		for left := 0; left+length <= n; left++ {
			right := left + length - 1
			total := prefix[right+1] - prefix[left]
			best := int(^uint(0) >> 1)
			result.Frames = append(result.Frames, intervalFrameGrid(values, done, left, right, nil, 3, "准备计算区间 ["+itoa(left)+","+itoa(right)+"]：先把 dp[l][r] 设为无穷大。"))
			for split := left; split < right; split++ {
				candidate := values[left][split] + values[split+1][right] + total
				if candidate < best {
					best = candidate
				}
				values[left][right] = best
				result.Frames = append(result.Frames, intervalFrameGrid(values, done, left, right, []gridPoint{{left, split}, {split + 1, right}}, 7, "长度 "+itoa(length)+"，区间 ["+itoa(left)+","+itoa(right)+"]：枚举 k="+itoa(split)+"，候选代价为 "+itoa(candidate)+"，当前最小值为 "+itoa(best)+"。"))
			}
			done[left][right] = true
			result.Frames = append(result.Frames, intervalFrameGrid(values, done, left, right, nil, 8, "区间 ["+itoa(left)+","+itoa(right)+"] 完成，最小合并代价为 "+itoa(best)+"。"))
		}
	}
	result.Frames = append(result.Frames, intervalFrameGrid(values, done, 0, n-1, nil, 10, "完整区间的最小合并代价为 "+itoa(values[0][n-1])+"。"))
	return result
}

func intervalFrameGrid(values [][]int, done [][]bool, currentRow, currentColumn int, dependencies []gridPoint, line int, narration string) Frame {
	state := gridState{Title: "dp[l][r]：闭区间 [l,r] 的最小合并代价", Rows: []string{"l=0", "l=1", "l=2", "l=3"}, Columns: []string{"r=0", "r=1", "r=2", "r=3"}}
	dependencySet := gridPointSet(dependencies)
	for row := range values {
		for column := range values[row] {
			cellState := "pending"
			if column < row {
				cellState = "unused"
			} else if done[row][column] {
				cellState = "ready"
			}
			if row == currentRow && column == currentColumn {
				cellState = "current"
			}
			state.Cells = append(state.Cells, gridCell{Row: row, Column: column, Value: values[row][column], State: cellState, Dependency: dependencySet[gridPoint{row, column}]})
		}
	}
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"order": "区间长度从短到长", "current": intervalName(currentRow, currentColumn)}, State: state}
}

// StockTrace tracks the two end-of-day states for unlimited stock transactions.
func StockTrace() Trace {
	prices := []int{7, 1, 5, 3, 6, 4}
	hold, cash := -prices[0], 0
	state := stockState{Days: []stockDay{{Day: 0, Price: prices[0], Hold: hold, Cash: cash, State: "current"}}, Current: 0}
	result := Trace{
		Kind:  "stock-state",
		Title: "股票状态 DP：持仓与空仓",
		Pseudocode: []string{
			"hold, cash = -prices[0], 0",
			"for _, price := range prices[1:] {",
			"    prevHold, prevCash := hold, cash",
			"    hold = max(prevHold, prevCash-price)",
			"    cash = max(prevCash, prevHold+price)",
			"}",
			"return cash",
		},
	}
	result.Frames = append(result.Frames, stockFrame(state, 0, "第 0 天：买入后持仓收益为 -7，什么也不做的空仓收益为 0。"))
	for day, price := range prices[1:] {
		prevHold, prevCash := hold, cash
		hold = max(prevHold, prevCash-price)
		cash = max(prevCash, prevHold+price)
		state.Days = append(state.Days, stockDay{Day: day + 1, Price: price, Hold: hold, Cash: cash, State: "current"})
		state.Current = day + 1
		result.Frames = append(result.Frames, stockFrame(state, 4, "第 "+itoa(day+1)+" 天价格为 "+itoa(price)+"：两个新状态都只读取昨天的 hold 与 cash。"))
	}
	result.Frames = append(result.Frames, stockFrame(state, 6, "最后一天空仓收益为 "+itoa(cash)+"，这就是可实现的最大收益。"))
	return result
}

type stockState struct {
	Days    []stockDay `json:"days"`
	Current int        `json:"current"`
}

type stockDay struct {
	Day   int    `json:"day"`
	Price int    `json:"price"`
	Hold  int    `json:"hold"`
	Cash  int    `json:"cash"`
	State string `json:"state"`
}

func stockFrame(state stockState, line int, narration string) Frame {
	view := stockState{Current: state.Current, Days: append([]stockDay(nil), state.Days...)}
	for i := range view.Days {
		view.Days[i].State = "ready"
	}
	view.Days[len(view.Days)-1].State = "current"
	if len(view.Days) > 1 {
		view.Days[len(view.Days)-2].State = "dependency"
	}
	last := view.Days[len(view.Days)-1]
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"price": itoa(last.Price), "hold": itoa(last.Hold), "cash": itoa(last.Cash)}, State: view}
}

// BitmaskTrace demonstrates expanding a visit set while retaining the last city.
func BitmaskTrace() Trace {
	steps := []bitmaskStep{
		{Mask: 1, Last: 0, Cost: 0, Narration: "从城市 0 出发：mask=0001，只访问了城市 0。"},
		{Mask: 3, Last: 1, Cost: 2, Narration: "访问城市 1：把第 1 位设为 1，状态变为 mask=0011，最后位置为 1。"},
		{Mask: 11, Last: 3, Cost: 6, Narration: "访问城市 3：新状态保留集合和最后位置，后续代价只需从城市 3 继续扩展。"},
		{Mask: 15, Last: 2, Cost: 18, Narration: "访问城市 2：mask=1111，所有城市都已在集合中。"},
	}
	result := Trace{
		Kind:  "bitmask-state",
		Title: "状态压缩 DP：集合与最后位置",
		Pseudocode: []string{
			"for mask := range dp {",
			"    for last := range dp[mask] { dp[mask][last] = inf }",
			"}",
			"dp[1<<0][0] = 0",
			"for mask := 1; mask < 1<<n; mask++ {",
			"    for last := 0; last < n; last++ {",
			"        if dp[mask][last] == inf { continue }",
			"        for next := 0; next < n; next++ {",
			"            if mask&(1<<next) != 0 { continue }",
			"            nextMask := mask | (1 << next)",
			"            candidate := dp[mask][last] + cost[last][next]",
			"            dp[nextMask][next] = min(dp[nextMask][next], candidate)",
			"        }",
			"    }",
			"}",
		},
	}
	for index, step := range steps {
		line := 3
		previousLast := -1
		if index > 0 {
			line = 11
			previousLast = steps[index-1].Last
		}
		result.Frames = append(result.Frames, Frame{ActiveLine: line, Narration: step.Narration, Variables: map[string]string{"mask": binaryMask(step.Mask, 4), "last": "city " + itoa(step.Last), "cost": itoa(step.Cost)}, State: bitmaskState{Names: []string{"0", "1", "2", "3"}, Mask: step.Mask, Last: step.Last, PreviousLast: previousLast, Cost: step.Cost}})
	}
	return result
}

type bitmaskStep struct {
	Mask      int
	Last      int
	Cost      int
	Narration string
}

type bitmaskState struct {
	Names        []string `json:"names"`
	Mask         int      `json:"mask"`
	Last         int      `json:"last"`
	PreviousLast int      `json:"previousLast"`
	Cost         int      `json:"cost"`
}

// PathTrace fills the minimum-path table from the top left to the bottom right.
func PathTrace() Trace {
	grid := [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}
	rows, columns := len(grid), len(grid[0])
	values := make([][]int, rows)
	done := make([][]bool, rows)
	for row := range values {
		values[row] = make([]int, columns)
		done[row] = make([]bool, columns)
	}
	values[0][0], done[0][0] = grid[0][0], true
	result := Trace{
		Kind:  "dp-grid",
		Title: "寻路 DP：网格最小路径和",
		Pseudocode: []string{
			"dp[0][0] = grid[0][0]",
			"for c := 1; c < n; c++ { dp[0][c] = dp[0][c-1] + grid[0][c] }",
			"for r := 1; r < m; r++ {",
			"    dp[r][0] = dp[r-1][0] + grid[r][0]",
			"    for c := 1; c < n; c++ {",
			"        dp[r][c] = min(dp[r-1][c], dp[r][c-1]) + grid[r][c]",
			"    }",
			"}",
		},
	}
	result.Frames = append(result.Frames, pathFrame(grid, values, done, 0, 0, nil, 0, "起点的最小路径和就是它自己的权重 1。"))
	for column := 1; column < columns; column++ {
		values[0][column] = values[0][column-1] + grid[0][column]
		done[0][column] = true
		result.Frames = append(result.Frames, pathFrame(grid, values, done, 0, column, []gridPoint{{0, column - 1}}, 1, "首行只能从左侧进入，蓝色格子是本次读取的左侧状态。"))
	}
	for row := 1; row < rows; row++ {
		values[row][0] = values[row-1][0] + grid[row][0]
		done[row][0] = true
		result.Frames = append(result.Frames, pathFrame(grid, values, done, row, 0, []gridPoint{{row - 1, 0}}, 3, "首列只能从上方进入，蓝色格子是本次读取的上方状态。"))
		for column := 1; column < columns; column++ {
			values[row][column] = min(values[row-1][column], values[row][column-1]) + grid[row][column]
			done[row][column] = true
			result.Frames = append(result.Frames, pathFrame(grid, values, done, row, column, []gridPoint{{row - 1, column}, {row, column - 1}}, 5, "蓝色格子是本次参与比较的上方与左方状态，再加当前格子权重 "+itoa(grid[row][column])+"。"))
		}
	}
	result.Frames = append(result.Frames, pathFrame(grid, values, done, rows-1, columns-1, []gridPoint{{rows - 2, columns - 1}, {rows - 1, columns - 2}}, 7, "右下角的最小路径和为 "+itoa(values[rows-1][columns-1])+"。"))
	return result
}

func pathFrame(grid, values [][]int, done [][]bool, currentRow, currentColumn int, dependencies []gridPoint, line int, narration string) Frame {
	state := gridState{Title: "dp[r][c]：到达该格的最小路径和", Rows: []string{"r=0", "r=1", "r=2"}, Columns: []string{"c=0", "c=1", "c=2"}}
	dependencySet := gridPointSet(dependencies)
	for row := range values {
		for column := range values[row] {
			cellState := "pending"
			if done[row][column] {
				cellState = "ready"
			}
			if row == currentRow && column == currentColumn {
				cellState = "current"
			}
			state.Cells = append(state.Cells, gridCell{Row: row, Column: column, Value: values[row][column], State: cellState, Dependency: dependencySet[gridPoint{row, column}]})
		}
	}
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"cell": "(" + itoa(currentRow) + "," + itoa(currentColumn) + ")", "weight": itoa(grid[currentRow][currentColumn])}, State: state}
}

func gridFrame(title string, rows, columns []string, values [][]int, done [][]bool, currentRow, currentColumn int, dependencies []gridPoint, line int, narration string) Frame {
	state := gridState{Title: title, Rows: rows, Columns: columns}
	dependencySet := gridPointSet(dependencies)
	for row := range values {
		for column := range values[row] {
			cellState := "pending"
			if done[row][column] {
				cellState = "ready"
			}
			if row == currentRow && column == currentColumn {
				cellState = "current"
			}
			state.Cells = append(state.Cells, gridCell{Row: row, Column: column, Value: values[row][column], State: cellState, Dependency: dependencySet[gridPoint{row, column}]})
		}
	}
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"current": intervalName(currentRow, currentColumn)}, State: state}
}

func gridPointSet(points []gridPoint) map[gridPoint]bool {
	result := make(map[gridPoint]bool, len(points))
	for _, point := range points {
		result[point] = true
	}
	return result
}

func labels(prefix, input string) []string {
	result := []string{prefix}
	for _, character := range input {
		result = append(result, string(character))
	}
	return result
}

func intervalName(row, column int) string {
	if row < 0 || column < 0 {
		return "初始化"
	}
	return "[" + itoa(row) + "," + itoa(column) + "]"
}

func binaryMask(value, width int) string {
	result := make([]byte, width)
	for index := width - 1; index >= 0; index-- {
		result[index] = byte(value&1) + '0'
		value >>= 1
	}
	return string(result)
}
