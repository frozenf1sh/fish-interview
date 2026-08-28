package trace

func concreteFlowTrace(name string) (Trace, bool) {
	factory, ok := concreteFlowFactories[name]
	if !ok {
		return Trace{}, false
	}
	return factory(), true
}

var concreteFlowFactories = map[string]func() Trace{
	"flow-greedy-reachability":       redesignedGreedyReachabilityTrace,
	"flow-greedy-lexicographic":      redesignedGreedyLexicographicTrace,
	"flow-greedy-interval-endpoints": redesignedGreedyEndpointsTrace,
	"flow-bfs-shortest-path":         redesignedBFSShortestTrace,
	"flow-bfs-multi-source":          redesignedBFSMultiSourceTrace,
	"flow-bfs-topological":           redesignedBFSTopologicalTrace,
	"flow-dfs-tree":                  redesignedDFSTreeTrace,
	"flow-dfs-grid":                  redesignedDFSGridTrace,
	"flow-dfs-path":                  redesignedDFSPathTrace,
	"flow-backtracking-choose-skip":  redesignedChooseSkipTrace,
	"flow-backtracking-enumeration":  redesignedEnumerationTrace,
	"flow-list-fast-slow":            redesignedFastSlowTrace,
	"flow-list-merge":                redesignedMergeListsTrace,
	"flow-tree-bst":                  redesignedBSTTrace,
	"flow-tree-lca":                  redesignedLCATrace,
	"flow-tree-path-sum":             redesignedTreePathSumTrace,
	"flow-tree-dp":                   redesignedTreeDPTrace,
	"flow-string-window":             SlidingWindowAtMostTrace,
	"flow-string-golang":             redesignedStringGoTrace,
	"flow-string-palindrome":         PalindromeIntervalDPTrace,
	"flow-list-merge-sort":           ListMergeSortTrace,
	"flow-window-exact":              SlidingWindowExactTrace,
	"flow-window-at-most":            SlidingWindowAtMostTrace,
	"flow-window-minimum":            SlidingWindowMinimumTrace,
	"flow-string-kmp":                redesignedKMPTrace,
	"flow-lcs-space":                 redesignedLCSSpaceTrace,
}

func concreteTrace(kind, title string, pseudocode []string, frames ...Frame) Trace {
	return Trace{Kind: kind, Title: title, Pseudocode: pseudocode, Frames: frames}
}

func tokenRow(label string, values []string, states map[int]string) exampleLane {
	items := make([]exampleItem, len(values))
	for index, value := range values {
		state := "ready"
		if configured, ok := states[index]; ok {
			state = configured
		}
		items[index] = item(value, state)
	}
	return lane(label, items...)
}

func matrixFromRows(rows []string, states map[string]string) []matrixCell {
	result := make([]matrixCell, 0)
	for row, value := range rows {
		for column, character := range []rune(value) {
			state := "ready"
			if configured, ok := states[itoa(row)+":"+itoa(column)]; ok {
				state = configured
			}
			result = append(result, matrixCell{Row: row, Column: column, Label: string(character), State: state})
		}
	}
	return result
}

var graphLinks = []nodeLink{{From: "A", To: "B"}, {From: "A", To: "C"}, {From: "B", To: "D"}, {From: "C", To: "D"}}

func graphNodes(states map[string]string) []nodeVisual {
	base := []nodeVisual{{ID: "A", Label: "A", X: 60, Y: 42}, {ID: "B", Label: "B", X: 145, Y: 108}, {ID: "C", Label: "C", X: 235, Y: 108}, {ID: "D", Label: "D", X: 310, Y: 42}}
	for index := range base {
		base[index].State = "ready"
		if state, ok := states[base[index].ID]; ok {
			base[index].State = state
		}
	}
	return base
}

var treeLinks = []nodeLink{{From: "3", To: "5"}, {From: "3", To: "1"}, {From: "5", To: "6"}, {From: "5", To: "2"}, {From: "1", To: "0"}, {From: "1", To: "8"}}

func treeNodes(states map[string]string) []nodeVisual {
	base := []nodeVisual{{ID: "3", Label: "3", X: 180, Y: 28}, {ID: "5", Label: "5", X: 105, Y: 92}, {ID: "1", Label: "1", X: 255, Y: 92}, {ID: "6", Label: "6", X: 55, Y: 158}, {ID: "2", Label: "2", X: 140, Y: 158}, {ID: "0", Label: "0", X: 220, Y: 158}, {ID: "8", Label: "8", X: 305, Y: 158}}
	for index := range base {
		base[index].State = "ready"
		if state, ok := states[base[index].ID]; ok {
			base[index].State = state
		}
	}
	return base
}

func withNodes(base []nodeVisual, states map[string]string) []nodeVisual {
	for index := range base {
		base[index].State = "ready"
		if state, ok := states[base[index].ID]; ok {
			base[index].State = state
		}
	}
	return base
}

func greedyReachabilityTrace() Trace {
	code := []string{"farthest := 0", "for i, jump := range nums {", "    if i > farthest { return false }", "    farthest = max(farthest, i+jump)", "}", "return true"}
	frames := []Frame{
		exampleFrame(0, "例题 nums=[2,3,1,1,4]。起点可达，farthest=0。", "跳跃游戏：从左到右扫描", tokenRow("nums", []string{"2", "3", "1", "1", "4"}, nil), lane("边界", item("farthest=0", "current"))),
		exampleFrame(2, "读取 i=0 与 farthest=0：0 没有越过蓝色边界，因此当前位置可参与扩展。", "跳跃游戏：检查可达性", tokenRow("nums", []string{"2", "3", "1", "1", "4"}, map[int]string{0: "current"}), lane("边界", item("farthest=0", "dependency"))),
		exampleFrame(3, "计算 0+2=2，写入 farthest=2。", "跳跃游戏：更新最远位置", tokenRow("nums", []string{"2", "3", "1", "1", "4"}, map[int]string{0: "dependency"}), lane("计算", item("0+2=2", "current"), item("farthest=2", "current"))),
		exampleFrame(3, "读取 i=1 的跳跃 3；蓝色 farthest=2 覆盖 i=1，计算 1+3=4。", "跳跃游戏：第二次扩展", tokenRow("nums", []string{"2", "3", "1", "1", "4"}, map[int]string{0: "ready", 1: "current"}), lane("边界", item("farthest=2", "dependency"), item("新 farthest=4", "current"))),
		exampleFrame(5, "farthest=4 已覆盖末尾下标 4，返回 true。", "跳跃游戏：到达终点", tokenRow("nums", []string{"2", "3", "1", "1", "4"}, map[int]string{0: "ready", 1: "ready"}), lane("答案", item("farthest=4", "ready"), item("true", "current"))),
	}
	return concreteTrace("example-state", "可达边界贪心：跳跃游戏", code, frames...)
}

func greedyLexicographicTrace() Trace {
	code := []string{"last := lastIndex(s)", "for _, ch := range s {", "    if inStack[ch] { continue }", "    for top > ch && last[top] > i { pop() }", "    push(ch)", "}"}
	input := []string{"c", "b", "a", "c", "d", "c", "b", "c"}
	frames := []Frame{
		exampleFrame(0, "例题 s=cbacdcbc。最后出现位置告诉我们：弹出的字符是否还能在后面补回。", "去重字典序：cbacdcbc", tokenRow("输入", input, nil), lane("栈", item("[]", "current"))),
		exampleFrame(4, "读取 c，栈为空，压入 c。", "去重字典序：读 c", tokenRow("输入", input, map[int]string{0: "current"}), lane("栈", item("c", "current"))),
		exampleFrame(3, "读取 b：栈顶 c>b 且 c 在后面还会出现，c 是蓝色依赖，先弹出。", "去重字典序：c 可补回", tokenRow("输入", input, map[int]string{0: "dependency", 1: "current"}), lane("栈", item("c", "dependency"), item("b", "current"))),
		exampleFrame(4, "把 b 压入，当前前缀为 b。", "去重字典序：写入 b", tokenRow("输入", input, map[int]string{0: "ready", 1: "current"}), lane("栈", item("b", "current"))),
		exampleFrame(4, "继续读 a、c、d 后栈为 acd；读到第二个 c 时已在栈，跳过。", "去重字典序：去重", tokenRow("输入", input, map[int]string{2: "ready", 3: "ready", 4: "ready", 5: "current"}), lane("栈", item("a", "ready"), item("c", "ready"), item("d", "ready"))),
		exampleFrame(5, "读 b 时 d 可在最后一个位置补回，弹 d 后压 b；再读 c 跳过，结果 acdb。", "去重字典序：最终栈", tokenRow("输入", input, map[int]string{6: "ready", 7: "current"}), lane("栈", item("a", "ready"), item("c", "ready"), item("d", "ready"), item("b", "ready")), lane("答案", item("acdb", "current"))),
	}
	return concreteTrace("example-state", "字典序贪心：删除重复字母", code, frames...)
}

func greedyEndpointsTrace() Trace {
	code := []string{"sort intervals by start", "right := first.end", "for _, in := range intervals[1:] {", "    if in.start > right { arrows++; right = in.end }", "    else { right = min(right, in.end) }", "}"}
	intervals := []string{"[1,6]", "[2,8]", "[7,12]", "[10,16]"}
	frames := []Frame{
		exampleFrame(0, "例题：用最少箭射爆 [[1,6],[2,8],[7,12],[10,16]]。输入已按开始时间排序。", "区间端点：最少箭数", tokenRow("区间", intervals, nil), lane("共同右端", item("right=6", "current"), item("arrows=1", "ready"))),
		exampleFrame(4, "读取 [2,8]：开始 2 没越过蓝色 right=6，仍与第一组相交。", "区间端点：检查交集", tokenRow("区间", intervals, map[int]string{0: "ready", 1: "current"}), lane("共同右端", item("right=6", "dependency"))),
		exampleFrame(4, "写 right=min(6,8)=6；箭仍放在 6，可同时命中前两段。", "区间端点：收紧交集", tokenRow("区间", intervals, map[int]string{0: "ready", 1: "ready"}), lane("共同右端", item("right=6", "current"), item("arrows=1", "ready"))),
		exampleFrame(3, "读取 [7,12]：7>6，交集为空。必须新增一支箭，并把 right 写成 12。", "区间端点：新建一组", tokenRow("区间", intervals, map[int]string{2: "current"}), lane("共同右端", item("right=6", "dependency"), item("right=12", "current"), item("arrows=2", "current"))),
		exampleFrame(4, "读取 [10,16]：10<=12，第二支箭保留在 right=min(12,16)=12。", "区间端点：最终", tokenRow("区间", intervals, map[int]string{3: "current"}), lane("答案", item("arrows=2", "current"), item("位置：6,12", "ready"))),
	}
	return concreteTrace("example-state", "区间端点贪心：最少箭数", code, frames...)
}

func bfsShortestTrace() Trace {
	code := []string{"queue := []Node{A}; dist[A] = 0", "for head < len(queue) {", "    cur := queue[head]; head++", "    for _, next := range graph[cur] {", "        if unseen(next) { dist[next] = dist[cur]+1; enqueue(next) }", "    }", "}"}
	frames := []Frame{
		nodeFrame(0, "例题图从 A 到 D。初始化：A 的距离为 0，先入队。", "无权图最短路：A→D", graphNodes(map[string]string{"A": "current"}), graphLinks),
		nodeFrame(2, "出队 A；A 是当前层，B、C 是待检查的邻居。", "队列=[A]，当前=A", graphNodes(map[string]string{"A": "current", "B": "dependency", "C": "dependency"}), graphLinks),
		nodeFrame(4, "首次访问 B，写 dist[B]=1 并入队。", "队列=[B]，dist[B]=1", graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency"}), graphLinks),
		nodeFrame(4, "首次访问 C，写 dist[C]=1 并入队。", "队列=[B,C]，dist[C]=1", graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current"}), graphLinks),
		nodeFrame(4, "出队 B 时首次到达 D，写 dist[D]=2。之后 C 再见到 D 不会覆盖它。", "队列=[C,D]，dist[D]=2", graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "dependency", "D": "current"}), graphLinks),
	}
	return concreteTrace("node-link-state", "BFS：无权图最短路", code, frames...)
}

func bfsMultiSourceTrace() Trace {
	code := []string{"enqueue all sources with dist=0", "for head < len(queue) {", "    cur := queue[head]; head++", "    for _, next := range neighbors(cur) {", "        if dist[next] == -1 { dist[next] = dist[cur]+1; enqueue(next) }", "    }", "}"}
	rows := []string{"0..", "...", "..0"}
	frames := []Frame{
		matrixFrame(0, "例题两个源在左上与右下，均写入距离 0 并同时入队。", "多源 BFS：0 是源，. 未访问", 3, 3, matrixFromRows(rows, map[string]string{"0:0": "current", "2:2": "current"})),
		matrixFrame(2, "出队左上源；右侧和下侧两个格子是本次参与扩散的蓝色邻居。", "处理源 (0,0)", 3, 3, matrixFromRows(rows, map[string]string{"0:0": "current", "0:1": "dependency", "1:0": "dependency", "2:2": "ready"})),
		matrixFrame(4, "首次访问 (0,1) 与 (1,0)，写距离 1。", "距离 1 入队", 3, 3, matrixFromRows([]string{"01.", "1..", "..0"}, map[string]string{"0:1": "current", "1:0": "current"})),
		matrixFrame(2, "再出队右下源；它的左侧和上侧格子也写距离 1。", "处理源 (2,2)", 3, 3, matrixFromRows([]string{"01.", "1.1", ".10"}, map[string]string{"1:2": "current", "2:1": "current", "2:2": "dependency"})),
		matrixFrame(4, "所有格子的第一次写入来自最近的源；队列层数保证无需再比较多个源。", "最终最近源距离", 3, 3, matrixFromRows([]string{"012", "121", "210"}, map[string]string{"1:1": "current"})),
	}
	return concreteTrace("matrix-state", "多源 BFS：最近源距离", code, frames...)
}

func bfsTopologicalTrace() Trace {
	code := []string{"enqueue every indegree-0 node", "for head < len(queue) {", "    v := queue[head]; head++", "    for _, to := range graph[v] {", "        indegree[to]--", "        if indegree[to] == 0 { enqueue(to) }", "    }", "}"}
	frames := []Frame{
		nodeFrame(0, "例题依赖 A→B、A→C、B→D、C→D。只有 A 入度为 0，先入队。", "Kahn：入度 A=0 B=1 C=1 D=2", graphNodes(map[string]string{"A": "current"}), graphLinks),
		nodeFrame(2, "出队 A，A 加入拓扑序。B、C 的入度是本次要更新的蓝色状态。", "队列=[A]，序列=[]", graphNodes(map[string]string{"A": "current", "B": "dependency", "C": "dependency"}), graphLinks),
		nodeFrame(4, "删除 A→B：B 的入度 1→0，B 入队。", "队列=[B]，序列=[A]", graphNodes(map[string]string{"A": "ready", "B": "current", "C": "dependency"}), graphLinks),
		nodeFrame(4, "删除 A→C：C 的入度 1→0，C 入队。", "队列=[B,C]，序列=[A]", graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "current"}), graphLinks),
		nodeFrame(5, "B、C 都处理后，D 的入度 2→0 才入队；最终序列可为 A,B,C,D。", "队列=[D]，序列=[A,B,C]", graphNodes(map[string]string{"A": "ready", "B": "ready", "C": "ready", "D": "current"}), graphLinks),
	}
	return concreteTrace("node-link-state", "BFS：Kahn 拓扑排序", code, frames...)
}

func dfsTreeTrace() Trace {
	code := []string{"if node == nil { return 0 }", "left := dfs(node.Left)", "right := dfs(node.Right)", "return left + right + node.Val"}
	frames := []Frame{
		nodeFrame(0, "例题求树节点和。根 3 还不能计算，先递归左孩子 5。", "树 DFS：后序求和", treeNodes(map[string]string{"3": "current", "5": "dependency"}), treeLinks),
		nodeFrame(1, "进入 5 后继续到叶子 6；6 的左右空节点返回基础值 0。", "当前递归栈：3→5→6", treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks),
		nodeFrame(3, "6 的左右结果都已得到，写入子树和 sum(6)=6。", "返回 6", treeNodes(map[string]string{"3": "dependency", "5": "dependency", "6": "current"}), treeLinks),
		nodeFrame(2, "继续计算 2 后，5 的两个孩子结果 6 与 2 都成为蓝色依赖。", "准备组合节点 5", treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "dependency", "2": "dependency"}), treeLinks),
		nodeFrame(3, "写入 sum(5)=6+2+5=13，再向根 3 返回。", "后序写入子树答案", treeNodes(map[string]string{"3": "dependency", "5": "current", "6": "ready", "2": "ready"}), treeLinks),
	}
	return concreteTrace("node-link-state", "DFS：树的后序组合", code, frames...)
}

func dfsGridTrace() Trace {
	code := []string{"if out of bounds or grid[r][c] == '0' { return }", "grid[r][c] = '0'", "for _, dir := range dirs {", "    dfs(r+dir[0], c+dir[1])", "}"}
	rows := []string{"110", "010", "011"}
	frames := []Frame{
		matrixFrame(0, "例题统计岛屿。外层扫描到 (0,0)=1，当前格子橙色，右邻居和下邻居蓝色。", "DFS 网格：1 是陆地", 3, 3, matrixFromRows(rows, map[string]string{"0:0": "current", "0:1": "dependency", "1:0": "dependency"})),
		matrixFrame(1, "进入即把 (0,0) 写成 0，避免从邻居绕回它。", "标记访问", 3, 3, matrixFromRows([]string{"010", "010", "011"}, map[string]string{"0:0": "current", "0:1": "dependency"})),
		matrixFrame(3, "递归到右侧 (0,1)，先高亮它，再写成 0。", "继续同一个连通块", 3, 3, matrixFromRows([]string{"010", "010", "011"}, map[string]string{"0:1": "current", "1:1": "dependency"})),
		matrixFrame(1, "(0,1) 标记后，向下进入 (1,1)。", "标记并向下", 3, 3, matrixFromRows([]string{"000", "010", "011"}, map[string]string{"1:1": "current"})),
		matrixFrame(3, "递归继续覆盖 (2,1)、(2,2)。本次 DFS 返回时，整座岛都已置为 0。", "一个连通块完成", 3, 3, matrixFromRows([]string{"000", "000", "000"}, map[string]string{"2:2": "ready"})),
	}
	return concreteTrace("matrix-state", "DFS：网格连通块", code, frames...)
}

func dfsPathTrace() Trace {
	code := []string{"if node == target { collect copy(path); return }", "for _, next := range graph[node] {", "    path = append(path, next)", "    dfs(next)", "    path = path[:len(path)-1]", "}"}
	frames := []Frame{
		nodeFrame(1, "例题枚举 A 到 D 的所有路径。起始 path=[A]，选择邻居 B。", "路径 DFS：target=D", graphNodes(map[string]string{"A": "dependency", "B": "current"}), graphLinks),
		nodeFrame(2, "选择 B 后先写 path=[A,B]，再递归进入 B。", "path=[A,B]", graphNodes(map[string]string{"A": "dependency", "B": "current"}), graphLinks),
		nodeFrame(2, "从 B 选择 D，写 path=[A,B,D]。", "path=[A,B,D]", graphNodes(map[string]string{"A": "dependency", "B": "dependency", "D": "current"}), graphLinks),
		nodeFrame(0, "到达 D，复制 path 到答案，不能直接保存共享切片。", "答案=[[A,B,D]]", graphNodes(map[string]string{"A": "ready", "B": "ready", "D": "current"}), graphLinks),
		nodeFrame(4, "返回 B 后撤销 D，再撤销 B；接着可选择 C 得到第二条路径 A,C,D。", "回溯恢复 path=[A]", graphNodes(map[string]string{"A": "current", "B": "ready", "C": "dependency", "D": "ready"}), graphLinks),
	}
	return concreteTrace("node-link-state", "DFS：路径枚举与回溯", code, frames...)
}

func chooseSkipTrace() Trace {
	code := []string{"if index == len(nums) { collect copy(path); return }", "dfs(index + 1)", "path = append(path, nums[index])", "dfs(index + 1)", "path = path[:len(path)-1]"}
	frames := []Frame{
		exampleFrame(1, "例题 nums=[1,2]。在 index=0 先走不选 1 的分支，path 仍为空。", "选或不选：构造所有子集", tokenRow("nums", []string{"1", "2"}, map[int]string{0: "current"}), lane("path", item("[]", "ready"))),
		exampleFrame(1, "index=1 再走不选 2，抵达叶子，复制 []。", "第一条叶子", tokenRow("nums", []string{"1", "2"}, map[int]string{1: "current"}), lane("答案", item("[]", "current"))),
		exampleFrame(2, "回到 index=1，选择 2：写 path=[2]。", "选择 2", tokenRow("nums", []string{"1", "2"}, map[int]string{1: "current"}), lane("path", item("2", "current"))),
		exampleFrame(3, "递归到叶子，复制 [2]；返回时执行 pop，path 恢复为空。", "收集并回退", lane("答案", item("[]", "ready"), item("[2]", "current")), lane("path", item("[]", "dependency"))),
		exampleFrame(2, "回到 index=0 后选择 1，后续同样枚举得到 [1] 与 [1,2]。", "另一条选择边", tokenRow("nums", []string{"1", "2"}, map[int]string{0: "current"}), lane("path", item("1", "current")), lane("答案", item("[] [2] [1] [1,2]", "ready"))),
	}
	return concreteTrace("example-state", "回溯：选或不选子集", code, frames...)
}

func enumerationTrace() Trace {
	code := []string{"if len(path) == len(nums) { collect copy(path); return }", "for i, value := range nums {", "    if used[i] { continue }", "    used[i] = true; path = append(path, value)", "    dfs(); path = path[:len(path)-1]; used[i] = false", "}"}
	frames := []Frame{
		exampleFrame(1, "例题 nums=[1,2,3] 求排列。第 0 层枚举候选 1、2、3。", "排列：used 防止重复选择", tokenRow("候选", []string{"1", "2", "3"}, map[int]string{0: "current"}), lane("used", item("F", "ready"), item("F", "ready"), item("F", "ready"))),
		exampleFrame(3, "选择 1：写 used[0]=true 与 path=[1]。", "第一层选择 1", tokenRow("候选", []string{"1", "2", "3"}, map[int]string{0: "current"}), lane("used", item("T", "current"), item("F", "ready"), item("F", "ready")), lane("path", item("1", "current"))),
		exampleFrame(2, "下一层枚举时 1 已被使用，红色跳过；2、3 仍可选。", "跳过已使用元素", tokenRow("候选", []string{"1", "2", "3"}, map[int]string{0: "rejected", 1: "dependency", 2: "dependency"}), lane("path", item("1", "ready"))),
		exampleFrame(3, "选择 2，写 path=[1,2]，再选择 3 即得到一个完整排列。", "深入一条分支", tokenRow("候选", []string{"1", "2", "3"}, map[int]string{1: "current", 2: "dependency"}), lane("path", item("1", "ready"), item("2", "current"))),
		exampleFrame(4, "收集 [1,2,3] 后逐层 pop 并恢复 used；恢复后才能尝试兄弟候选。", "撤销现场", lane("答案", item("[1,2,3]", "current")), lane("path", item("[1]", "dependency")), lane("used", item("T", "ready"), item("F", "current"), item("F", "ready"))),
	}
	return concreteTrace("example-state", "回溯：枚举下一个候选", code, frames...)
}

func listFrame(chain []string, pointers map[string]string, highlights []string, line int, narration string) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"example": "链表 1→2→3→4→2（有环）"}, State: linkedListState{Chain: chain, Pointers: pointers, Highlight: highlights}}
}

func fastSlowTrace() Trace {
	code := []string{"slow, fast := head, head", "for fast != nil && fast.Next != nil {", "    slow = slow.Next", "    fast = fast.Next.Next", "    if slow == fast { return true }", "}"}
	frames := []Frame{
		listFrame([]string{"1", "2", "3", "4", "2"}, map[string]string{"slow": "1", "fast": "1"}, []string{"1"}, 0, "例题链表 1→2→3→4→2。slow 与 fast 都从 head=1 起步，尚未比较。"),
		listFrame([]string{"1", "2", "3", "4", "2"}, map[string]string{"slow": "2", "fast": "1"}, []string{"1", "2"}, 2, "第一轮先让 slow 从 1 走到 2；1 与 2 是本次参与的节点。"),
		listFrame([]string{"1", "2", "3", "4", "2"}, map[string]string{"slow": "2", "fast": "3"}, []string{"1", "2", "3"}, 3, "再让 fast 从 1 连走两步到 3。此时 slow=2、fast=3，未相遇。"),
		listFrame([]string{"1", "2", "3", "4", "2"}, map[string]string{"slow": "3", "fast": "2"}, []string{"2", "3", "4"}, 3, "第二轮 slow 到 3，fast 从 3 经 4 回到 2；环让 fast 不会到 nil。"),
		listFrame([]string{"1", "2", "3", "4", "2"}, map[string]string{"slow": "4", "fast": "4"}, []string{"4"}, 4, "第三轮二者都在 4，相遇后返回 true。"),
	}
	return concreteTrace("linked-list", "链表：快慢指针判环", code, frames...)
}

func mergeListsTrace() Trace {
	code := []string{"dummy := &ListNode{}; tail := dummy", "for a != nil && b != nil {", "    if a.Val <= b.Val { chosen = a; a = a.Next } else { chosen = b; b = b.Next }", "    tail.Next = chosen; tail = tail.Next", "}", "tail.Next = remaining"}
	frames := []Frame{
		exampleFrame(0, "例题 A=1→3→5，B=2→4→6。dummy 固定结果头，tail=dummy。", "合并两个有序链表", tokenRow("A", []string{"1", "3", "5"}, map[int]string{0: "dependency"}), tokenRow("B", []string{"2", "4", "6"}, map[int]string{0: "dependency"}), lane("结果", item("dummy", "current"))),
		exampleFrame(2, "比较两个蓝色头节点 1 与 2，1 更小，选择 A 的 1。", "比较当前头", tokenRow("A", []string{"1", "3", "5"}, map[int]string{0: "current"}), tokenRow("B", []string{"2", "4", "6"}, map[int]string{0: "dependency"}), lane("tail", item("dummy", "dependency"))),
		exampleFrame(3, "执行 tail.Next=1，tail 前进到 1；A 的头移动到 3。", "接入 1", tokenRow("A", []string{"3", "5"}, map[int]string{0: "dependency"}), tokenRow("B", []string{"2", "4", "6"}, map[int]string{0: "dependency"}), lane("结果", item("dummy→1", "current"))),
		exampleFrame(2, "比较 3 与 2，选择 B 的 2，再把它接在 tail=1 后。", "接入 2", tokenRow("A", []string{"3", "5"}, map[int]string{0: "dependency"}), tokenRow("B", []string{"2", "4", "6"}, map[int]string{0: "current"}), lane("结果", item("dummy→1→2", "current"))),
		exampleFrame(5, "重复比较后得到 1→2→3→4→5；A 耗尽时直接接上 B 的剩余 6。", "接上剩余链", lane("结果", item("dummy→1→2→3→4→5→6", "current"))),
	}
	return concreteTrace("example-state", "链表：合并两个有序链表", code, frames...)
}

func bstTrace() Trace {
	code := []string{"if node == nil { return true }", "if node.Val <= low || node.Val >= high { return false }", "return valid(node.Left, low, node.Val) && valid(node.Right, node.Val, high)"}
	links := []nodeLink{{From: "5", To: "1"}, {From: "5", To: "4"}, {From: "4", To: "3"}, {From: "4", To: "6"}}
	nodes := func(states map[string]string) []nodeVisual {
		return withNodes([]nodeVisual{{ID: "5", Label: "5", X: 180, Y: 28}, {ID: "1", Label: "1", X: 90, Y: 95}, {ID: "4", Label: "4", X: 270, Y: 95}, {ID: "3", Label: "3", X: 230, Y: 160}, {ID: "6", Label: "6", X: 315, Y: 160}}, states)
	}
	frames := []Frame{
		nodeFrame(1, "例题树根 5，初始范围 (-∞,+∞)，5 合法。", "BST：范围来自全部祖先", nodes(map[string]string{"5": "current"}), links),
		nodeFrame(2, "递归左孩子 1，允许范围变为 (-∞,5)，1 合法。", "检查左子树范围", nodes(map[string]string{"5": "dependency", "1": "current"}), links),
		nodeFrame(2, "递归右孩子 4，允许范围是 (5,+∞)，4 立刻违反下界 5。", "祖先下界仍生效", nodes(map[string]string{"5": "dependency", "4": "rejected"}), links),
		nodeFrame(1, "不能只比较 4 与父节点 5 的左右关系；即使检查 3，它也同样落在根 5 的错误一侧。", "返回 false", nodes(map[string]string{"5": "dependency", "4": "rejected", "3": "rejected"}), links),
	}
	return concreteTrace("node-link-state", "BST：祖先范围验证", code, frames...)
}

func lcaTrace() Trace {
	code := []string{"if root == nil || root == p || root == q { return root }", "left := lca(root.Left, p, q)", "right := lca(root.Right, p, q)", "if left == nil { return right }", "if right == nil { return left }", "return root"}
	frames := []Frame{
		nodeFrame(1, "例题 p=5、q=1。根 3 不是目标，先向左递归。", "LCA：目标 5 与 1", treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks),
		nodeFrame(0, "到达节点 5，命中 p，直接把 5 向上返回。", "左子树返回 5", treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks),
		nodeFrame(2, "根 3 再递归右子树；1 是 q，直接返回 1。", "右子树返回 1", treeNodes(map[string]string{"3": "dependency", "1": "current"}), treeLinks),
		nodeFrame(5, "根 3 收到 left=5、right=1，两个蓝色返回值都非空，因此当前 3 是最近公共祖先。", "左右返回同时存在", treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks),
	}
	return concreteTrace("node-link-state", "树题：最近公共祖先", code, frames...)
}

func treePathSumTrace() Trace {
	code := []string{"if node == nil { return false }", "remain -= node.Val", "if leaf { return remain == 0 }", "return dfs(node.Left, remain) || dfs(node.Right, remain)"}
	frames := []Frame{
		nodeFrame(1, "例题 target=14，根 3 参与计算：remain 从 14 写成 11。", "路径和：target=14", treeNodes(map[string]string{"3": "current"}), treeLinks),
		nodeFrame(3, "选择左孩子 5；蓝色 remain=11 减去节点值 5，写入 remain=6。", "路径 3→5", treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks),
		nodeFrame(3, "节点 5 不是叶子，remain=6 还不能结算；继续进入孩子 6 或 2。", "只有叶子才结算", treeNodes(map[string]string{"5": "current", "6": "dependency", "2": "dependency"}), treeLinks),
		nodeFrame(2, "路径 3→5→6 到叶子时 remain 从 6 变为 0，才返回 true。", "叶子检查", treeNodes(map[string]string{"3": "ready", "5": "ready", "6": "current"}), treeLinks),
	}
	return concreteTrace("node-link-state", "树题：根到叶路径和", code, frames...)
}

func treeDPTrace() Trace {
	code := []string{"leftTake, leftSkip := dfs(node.Left)", "rightTake, rightSkip := dfs(node.Right)", "take := node.Val + leftSkip + rightSkip", "skip := max(leftTake,leftSkip) + max(rightTake,rightSkip)", "return take, skip"}
	frames := []Frame{
		nodeFrame(0, "例题打家劫舍 III。先后序计算左孩子 5 的 (take,skip)。", "树形 DP：节点值即偷取收益", treeNodes(map[string]string{"3": "dependency", "5": "current"}), treeLinks),
		nodeFrame(1, "再计算右孩子 1 的 (take,skip)。孩子结果返回后保持绿色。", "孩子状态已完成", treeNodes(map[string]string{"3": "dependency", "5": "ready", "1": "current"}), treeLinks),
		nodeFrame(2, "计算根的 take：选择 3 时只能读取两个孩子的 skip。蓝色状态参与加法。", "take(3)=3+skip(5)+skip(1)", treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks),
		nodeFrame(3, "计算根的 skip：不选 3 时，左右孩子各自取 take/skip 的较大值。", "skip(3)=max(5)+max(1)", treeNodes(map[string]string{"3": "current", "5": "dependency", "1": "dependency"}), treeLinks),
		nodeFrame(4, "根返回两个状态；答案取 max(take(3),skip(3))。", "最终只在根比较两种状态", treeNodes(map[string]string{"3": "current", "5": "ready", "1": "ready"}), treeLinks),
	}
	return concreteTrace("node-link-state", "树形 DP：选与不选", code, frames...)
}

func stringWindowTrace() Trace {
	code := []string{"for right := range s {", "    count[s[right]]++", "    for count[s[right]] > 1 { count[s[left]]--; left++ }", "    answer = max(answer, right-left+1)", "}"}
	input := []string{"a", "b", "c", "a", "b", "c", "b", "b"}
	frames := []Frame{
		exampleFrame(0, "例题 s=abcabcbb，求无重复最长子串。right=0 读入 a。", "滑动窗口：左闭右闭 [left,right]", tokenRow("字符串", input, map[int]string{0: "current"}), lane("窗口", item("a", "current")), lane("答案", item("len=1", "ready"))),
		exampleFrame(3, "继续读 b、c，窗口 abc 合法，写答案 len=3。", "窗口 [0,2]=abc", tokenRow("字符串", input, map[int]string{0: "ready", 1: "ready", 2: "current"}), lane("窗口", item("a", "ready"), item("b", "ready"), item("c", "current")), lane("答案", item("len=3", "current"))),
		exampleFrame(1, "right=3 再读 a，频次 a 从 1 变为 2；新 a 与旧 a 都参与判断。", "窗口不再合法", tokenRow("字符串", input, map[int]string{0: "dependency", 3: "current"}), lane("频次", item("a=2", "current"))),
		exampleFrame(2, "收缩 left：移除旧 a，left 从 0 写成 1，窗口 bca 恢复合法。", "恢复窗口条件", tokenRow("字符串", input, map[int]string{0: "rejected", 1: "ready", 2: "ready", 3: "current"}), lane("窗口", item("b", "ready"), item("c", "ready"), item("a", "current"))),
		exampleFrame(3, "后续重复 b 时重复同一收缩；最大长度保持 3。", "最终答案", lane("答案", item("len=3", "current")), lane("窗口", item("abc", "ready"))),
	}
	return concreteTrace("example-state", "字符串：无重复滑动窗口", code, frames...)
}

func stringGoTrace() Trace {
	code := []string{"if asciiOnly(s) { use s[i] }", "chars := []rune(s)", "var builder strings.Builder", "for _, part := range parts { builder.WriteString(part) }", "return builder.String()"}
	frames := []Frame{
		exampleFrame(0, "例题 s=Go中。题目若按字符计数，先判断 byte 下标不能直接代表第三个字符。", "Go 字符串：byte 与 rune", tokenRow("bytes", []string{"G", "o", "中(3 bytes)"}, map[int]string{2: "current"})),
		exampleFrame(1, "把 s 转成 []rune 后得到三个可按字符索引的元素 G、o、中。", "rune 切片", tokenRow("runes", []string{"G", "o", "中"}, map[int]string{2: "current"})),
		exampleFrame(2, "例题输出拼接 parts=[Go,中]；创建 Builder，不在循环内用 + 反复创建中间串。", "构造输出", tokenRow("parts", []string{"Go", "中"}, map[int]string{0: "dependency"}), lane("builder", item("\"\"", "current"))),
		exampleFrame(3, "依次 WriteString(Go)、WriteString(中)，新写入部分为橙色。", "Builder 追加", tokenRow("parts", []string{"Go", "中"}, map[int]string{1: "current"}), lane("builder", item("Go中", "current"))),
	}
	return concreteTrace("example-state", "Go 字符串：byte、rune 与 Builder", code, frames...)
}

func palindromeTrace() Trace {
	code := []string{"for center := range s {", "    expand(center, center)", "    expand(center, center+1)", "}", "for l >= 0 && r < len(s) && s[l] == s[r] { l--; r++ }"}
	input := []string{"b", "a", "b", "a", "d"}
	frames := []Frame{
		exampleFrame(0, "例题 s=babad。先枚举单点中心 center=2，即中间 b。", "中心扩展：寻找最长回文", tokenRow("字符串", input, map[int]string{2: "current"})),
		exampleFrame(4, "比较左右相邻 a 与 a，两个蓝色字符相同，向外扩张。", "中心 2：l=1 r=3", tokenRow("字符串", input, map[int]string{1: "dependency", 2: "current", 3: "dependency"}), lane("候选", item("aba", "current"))),
		exampleFrame(4, "再向外比较 b 与 d，失配，中心 2 的最大回文停在 aba。", "首次失配停止", tokenRow("字符串", input, map[int]string{0: "dependency", 4: "rejected"}), lane("最长", item("aba", "ready"))),
		exampleFrame(1, "center=1 的奇数扩张得到 bab，长度同为 3；偶数中心 (1,2) 先比较 a 与 b 即停止。", "奇偶中心都必须枚举", tokenRow("字符串", input, map[int]string{0: "ready", 1: "current", 2: "ready"}), lane("答案", item("bab 或 aba", "current"))),
	}
	return concreteTrace("example-state", "字符串：回文中心扩展", code, frames...)
}

func kmpTrace() Trace {
	code := []string{"pi := buildPrefix(pattern)", "for _, ch := range text {", "    for j > 0 && ch != pattern[j] { j = pi[j-1] }", "    if ch == pattern[j] { j++ }", "    if j == len(pattern) { report match }", "}"}
	frames := []Frame{
		exampleFrame(0, "例题 text=ababc，pattern=abc。先构造 pi=[0,0,0]。", "KMP：文本不回退", tokenRow("text", []string{"a", "b", "a", "b", "c"}, nil), tokenRow("pattern", []string{"a", "b", "c"}, nil), lane("pi", item("0", "ready"), item("0", "ready"), item("0", "ready"))),
		exampleFrame(3, "读 text[0]=a，匹配 pattern[0]=a，写 j=1。", "匹配第一个字符", tokenRow("text", []string{"a", "b", "a", "b", "c"}, map[int]string{0: "current"}), tokenRow("pattern", []string{"a", "b", "c"}, map[int]string{0: "dependency"}), lane("j", item("1", "current"))),
		exampleFrame(3, "读 text[1]=b，匹配 pattern[1]=b，写 j=2。", "继续匹配", tokenRow("text", []string{"a", "b", "a", "b", "c"}, map[int]string{1: "current"}), tokenRow("pattern", []string{"a", "b", "c"}, map[int]string{1: "dependency"}), lane("j", item("2", "current"))),
		exampleFrame(2, "读 text[2]=a，与 pattern[2]=c 失配；pi[1]=0 是蓝色依赖，把 j 回退到 0，文本指针不动。", "只回退模式长度", tokenRow("text", []string{"a", "b", "a", "b", "c"}, map[int]string{2: "current"}), tokenRow("pattern", []string{"a", "b", "c"}, map[int]string{2: "rejected"}), lane("j", item("2→0", "current"))),
		exampleFrame(4, "从 text[2] 重新匹配 a、b、c，j 达到 3，报告起点 2。", "找到匹配", tokenRow("text", []string{"a", "b", "a", "b", "c"}, map[int]string{2: "ready", 3: "ready", 4: "current"}), lane("答案", item("match at 2", "current"))),
	}
	return concreteTrace("example-state", "字符串：KMP 模式匹配", code, frames...)
}

func lcsSpaceTrace() Trace {
	code := []string{"dp := make([]int, len(b)+1)", "for i := 1; i <= len(a); i++ {", "    diagonal := 0", "    for j := 1; j <= len(b); j++ {", "        up := dp[j]", "        dp[j] = transition(diagonal, up, dp[j-1])", "        diagonal = up", "    }", "}"}
	frames := []Frame{
		exampleFrame(0, "例题 a=ab、b=ac。dp=[0,0,0] 对应 b 的空前缀、a、c。", "LCS 一维压缩：第 0 行", tokenRow("dp", []string{"0", "0", "0"}, map[int]string{0: "ready", 1: "ready", 2: "ready"}), lane("diagonal", item("0", "current"))),
		exampleFrame(4, "处理 a[0]=a、b[0]=a：读取覆盖前 dp[1]=0 作为蓝色 up，同时读取 diagonal=0。", "当前格 (1,1)", tokenRow("dp", []string{"0", "0", "0"}, map[int]string{0: "dependency", 1: "dependency", 2: "ready"}), lane("diagonal", item("0", "dependency"))),
		exampleFrame(5, "字符相同，写 dp[1]=diagonal+1=1；新 dp[1] 橙色，未动 dp[2] 保持绿色。", "写入当前行", tokenRow("dp", []string{"0", "1", "0"}, map[int]string{0: "ready", 1: "current", 2: "ready"}), lane("旧 up", item("0", "dependency"))),
		exampleFrame(6, "把旧 up=0 写给 diagonal，下一列才能读取上一行左上角。", "更新 diagonal", tokenRow("dp", []string{"0", "1", "0"}, map[int]string{0: "ready", 1: "ready", 2: "dependency"}), lane("diagonal", item("0", "current"))),
		exampleFrame(5, "处理第二行 b 时，同样先保存 up；最终 dp=[0,1,1]，答案是 dp[2]=1。", "最终一维数组", tokenRow("dp", []string{"0", "1", "1"}, map[int]string{0: "ready", 1: "ready", 2: "current"}), lane("答案", item("LCS=1", "current"))),
	}
	return concreteTrace("example-state", "LCS：一维数组空间优化", code, frames...)
}

// StartSortedIntervalsTrace shows why interval merging scans by increasing start time.
func legacyStartSortedIntervalsTrace() Trace {
	code := []string{"sort.Slice(intervals, byStart)", "merged := [][]int{intervals[0]}", "for _, current := range intervals[1:] {", "    last := merged[len(merged)-1]", "    if current[0] <= last[1] { last[1] = max(last[1], current[1]) }", "    else { merged = append(merged, current) }", "}"}
	frames := []Frame{
		exampleFrame(0, "例题合并 [[1,3],[2,6],[8,10],[15,18]]。开始时间排序后，后来的区间不会再出现在当前区间左侧。", "区间合并：按开始时间", tokenRow("输入", []string{"[1,3]", "[2,6]", "[8,10]", "[15,18]"}, nil), lane("merged", item("[1,3]", "current"))),
		exampleFrame(4, "读取 [2,6]：开始 2<=蓝色 last.end=3，两个区间重叠。", "检查相交", tokenRow("输入", []string{"[1,3]", "[2,6]", "[8,10]", "[15,18]"}, map[int]string{0: "dependency", 1: "current"}), lane("last", item("[1,3]", "dependency"))),
		exampleFrame(4, "写 last.end=max(3,6)=6，merged 更新为 [1,6]。", "扩展当前段", lane("merged", item("[1,6]", "current"))),
		exampleFrame(5, "读取 [8,10]：8>6，不相交，直接追加为下一段。", "开启新段", tokenRow("输入", []string{"[1,3]", "[2,6]", "[8,10]", "[15,18]"}, map[int]string{2: "current"}), lane("merged", item("[1,6]", "ready"), item("[8,10]", "current"))),
		exampleFrame(5, "[15,18] 同理追加，最终得到 [[1,6],[8,10],[15,18]]。", "最终合并结果", lane("merged", item("[1,6]", "ready"), item("[8,10]", "ready"), item("[15,18]", "current"))),
	}
	return concreteTrace("example-state", "区间：按开始时间合并", code, frames...)
}

// MeetingRoomsTrace replays the sorted-start/sorted-end two-pointer method.
func legacyMeetingRoomsTrace() Trace {
	code := []string{"sort.Ints(starts); sort.Ints(ends)", "rooms, end := 0, 0", "for _, start := range starts {", "    if start < ends[end] { rooms++ } else { end++ }", "}", "return rooms"}
	frames := []Frame{
		exampleFrame(0, "例题会议 [[0,30],[5,10],[15,20]]。拆成开始数组 [0,5,15]、结束数组 [10,20,30]。", "会议室：两个排序指针", tokenRow("starts", []string{"0", "5", "15"}, nil), tokenRow("ends", []string{"10", "20", "30"}, nil), lane("rooms", item("0", "current"))),
		exampleFrame(3, "读取 start=0，0<蓝色 end=10，第一场尚未结束，rooms 写成 1。", "开始 0 占用房间", tokenRow("starts", []string{"0", "5", "15"}, map[int]string{0: "current"}), tokenRow("ends", []string{"10", "20", "30"}, map[int]string{0: "dependency"}), lane("rooms", item("1", "current"))),
		exampleFrame(3, "读取 start=5，5<10，第二场仍与第一场重叠，rooms 写成 2。", "开始 5 仍重叠", tokenRow("starts", []string{"0", "5", "15"}, map[int]string{1: "current"}), tokenRow("ends", []string{"10", "20", "30"}, map[int]string{0: "dependency"}), lane("rooms", item("2", "current"))),
		exampleFrame(3, "读取 start=15，15>=10，最早结束的会议释放一间房，end 指针前进而 rooms 不增加。", "开始 15 复用房间", tokenRow("starts", []string{"0", "5", "15"}, map[int]string{2: "current"}), tokenRow("ends", []string{"10", "20", "30"}, map[int]string{0: "dependency", 1: "current"}), lane("答案", item("rooms=2", "current"))),
	}
	return concreteTrace("example-state", "区间：最少会议室", code, frames...)
}

// WeightedIntervalsTrace makes the incompatibility with unweighted greedy visible.
func legacyWeightedIntervalsTrace() Trace {
	code := []string{"sort intervals by end", "dp[0] = 0", "for i := 1; i <= n; i++ {", "    skip := dp[i-1]", "    take := weight[i] + dp[prev(i)]", "    dp[i] = max(skip, take)", "}"}
	frames := []Frame{
		exampleFrame(0, "例题区间 A=[1,3] 权重 5，B=[2,5] 权重 100，C=[4,6] 权重 5。按结束排序后不能只选结束最早的 A。", "带权区间调度：按结束时间 DP", tokenRow("区间", []string{"A:3/5", "B:5/100", "C:6/5"}, nil), tokenRow("dp", []string{"0", "?", "?", "?"}, map[int]string{0: "ready"})),
		exampleFrame(4, "处理 A：前驱为空，take=5+dp[0]=5，写 dp[1]=5。", "写 dp[1]", tokenRow("dp", []string{"0", "5", "?", "?"}, map[int]string{0: "dependency", 1: "current"})),
		exampleFrame(4, "处理 B：前驱仍为空，蓝色 dp[0] 参与 take=100；同时比较 skip=dp[1]=5。", "比较选 B 与跳过 B", tokenRow("dp", []string{"0", "5", "?", "?"}, map[int]string{0: "dependency", 1: "dependency", 2: "current"}), lane("候选", item("skip=5", "dependency"), item("take=100", "current"))),
		exampleFrame(5, "写 dp[2]=100。处理 C 时前驱是 A，因此 take=5+dp[1]=10，小于跳过 C 的 dp[2]=100。", "写 dp[2] 后处理 C", tokenRow("dp", []string{"0", "5", "100", "100"}, map[int]string{1: "dependency", 2: "dependency", 3: "current"})),
		exampleFrame(5, "最终 dp[3]=100，选择 B。带权版本的最优子结构是 DP，不是无权区间调度的贪心。", "最终最优权重", lane("答案", item("100", "current"))),
	}
	return concreteTrace("example-state", "区间：带权调度 DP", code, frames...)
}

// KadaneTrace replays the restart-or-extend decision for maximum subarray.
func legacyKadaneTrace() Trace {
	code := []string{"current, best := nums[0], nums[0]", "for _, x := range nums[1:] {", "    current = max(x, current+x)", "    best = max(best, current)", "}", "return best"}
	values := []string{"-2", "1", "-3", "4", "-1", "2", "1", "-5", "4"}
	frames := []Frame{
		exampleFrame(0, "例题 [-2,1,-3,4,-1,2,1,-5,4]。current 与 best 都从 -2 开始。", "最大子数组：Kadane", tokenRow("nums", values, map[int]string{0: "current"}), lane("状态", item("current=-2", "current"), item("best=-2", "current"))),
		exampleFrame(2, "读 x=1：蓝色 current=-2 加 1 得 -1，比直接从 1 重开差；写 current=1。", "重开还是延续", tokenRow("nums", values, map[int]string{0: "dependency", 1: "current"}), lane("候选", item("extend=-1", "dependency"), item("restart=1", "current"))),
		exampleFrame(3, "current=1 大于旧 best=-2，写 best=1。", "更新全局最优", lane("状态", item("current=1", "ready"), item("best=1", "current"))),
		exampleFrame(2, "读 x=4 前，经过 -3 后 current 已重置为 -2；蓝色 extend=2 小于 restart=4，写 current=4。", "第二次重开", tokenRow("nums", values, map[int]string{3: "current"}), lane("候选", item("extend=2", "dependency"), item("restart=4", "current"))),
		exampleFrame(3, "继续读 -1、2、1，current 依次为 3、5、6，best 写成 6。", "连续正贡献延续", tokenRow("nums", values, map[int]string{4: "ready", 5: "ready", 6: "current"}), lane("答案", item("best=6", "current"), item("子数组 [4,-1,2,1]", "ready"))),
	}
	return concreteTrace("example-state", "贪心：最大子数组 Kadane", code, frames...)
}
