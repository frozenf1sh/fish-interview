package trace

type flowSpec struct {
	Title      string
	Pseudocode []string
	Steps      []string
	Narration  []string
}

type flowState struct {
	Steps   []string `json:"steps"`
	Current int      `json:"current"`
}

// FlowTrace is a lightweight visual replay for template-oriented algorithm cards.
func FlowTrace(name string) (Trace, bool) {
	spec, ok := flowSpecs[name]
	if !ok {
		return Trace{}, false
	}
	result := Trace{Kind: "flow-steps", Title: spec.Title, Pseudocode: spec.Pseudocode}
	for index := range spec.Steps {
		narration := spec.Narration[index]
		result.Frames = append(result.Frames, Frame{
			ActiveLine: min(index, len(spec.Pseudocode)-1),
			Narration:  narration,
			Variables:  map[string]string{"stage": spec.Steps[index]},
			State:      flowState{Steps: spec.Steps, Current: index},
		})
	}
	return result, true
}

var flowSpecs = map[string]flowSpec{
	"flow-greedy-reachability": {
		Title:      "可达边界：扫描时只维护最远位置",
		Pseudocode: []string{"farthest := 0", "for i, jump := range nums {", "    if i > farthest { return false }", "    farthest = max(farthest, i+jump)", "}", "return true"},
		Steps:      []string{"初始化可达边界", "确认当前位置可达", "用当前位置扩展边界", "边界覆盖终点"},
		Narration:  []string{"边界 0 表示起点本身可达。", "扫描到 i 前先比较 i 与 farthest。", "只有可达位置才能贡献新的最远边界。", "farthest 覆盖末尾即可提前结束。"},
	},
	"flow-greedy-lexicographic": {
		Title:      "字典序贪心：弹出仍可补回的较大栈顶",
		Pseudocode: []string{"record last occurrence", "for ch := range s {", "    skip if already in stack", "    pop while larger and reusable", "    push ch", "}"},
		Steps:      []string{"记录最后位置", "读入当前字符", "判断栈顶能否补回", "弹出并压入当前字符"},
		Narration:  []string{"最后出现位置决定一个字符能否安全删除。", "已在栈中则不重复加入。", "更大且后面还会出现的栈顶可以让位。", "压入后保持当前前缀的最小字典序。"},
	},
	"flow-greedy-interval-endpoints": {
		Title:      "区间端点：维护仍可共同命中的右端",
		Pseudocode: []string{"sort by start", "right := first.end", "for interval := range rest {", "    if interval.start > right { start new point }", "    else { right = min(right, interval.end) }", "}"},
		Steps:      []string{"按起点排序", "建立当前交集", "检查新区间是否相交", "收紧交集或新增点"},
		Narration:  []string{"排序后新区间只会向右推进。", "right 是当前所有区间共同可放点的最右位置。", "起点越过 right 时交集为空。", "相交则收紧 right；不相交才增加一次选择。"},
	},
	"flow-bfs-shortest-path": {
		Title:      "BFS：按距离层出队",
		Pseudocode: []string{"enqueue start; dist[start] = 0", "for head < len(queue) {", "    cur := queue[head]", "    for next := range neighbors(cur) {", "        mark and enqueue unseen next", "    }", "}"},
		Steps:      []string{"起点入队", "取出当前层节点", "枚举邻居", "首次访问时写距离并入队"},
		Narration:  []string{"队列初始层距离为 0。", "先入队的状态距离不会更大。", "所有边权相同，邻居距离只加 1。", "入队时标记，保证第一次到达就是最短。"},
	},
	"flow-bfs-multi-source": {
		Title:      "多源 BFS：所有源同时扩散",
		Pseudocode: []string{"enqueue every source", "for head < len(queue) {", "    cur := queue[head]", "    visit unassigned neighbors", "}"},
		Steps:      []string{"全部源入队", "取出当前扩散层", "访问未分配邻居", "写入最近源距离"},
		Narration:  []string{"所有源都是距离 0，不需要逐个重复 BFS。", "队列保证先扩散距离小的层。", "未分配格首次被触及。", "该层数就是它到最近源的距离。"},
	},
	"flow-bfs-topological": {
		Title:      "Kahn：入度为 0 的节点逐层释放",
		Pseudocode: []string{"enqueue every indegree-0 node", "for head < len(queue) {", "    v := queue[head]", "    decrement indegree of successors", "    enqueue successor when indegree is 0", "}"},
		Steps:      []string{"收集零入度", "取出可执行节点", "删除它的出边", "新零入度节点入队"},
		Narration:  []string{"没有前置依赖的节点可以立刻执行。", "出队即加入拓扑序列。", "处理节点等价于删除其所有出边。", "若最终数量不足 n，剩余图含环。"},
	},
	"flow-dfs-tree": {
		Title:      "树 DFS：先递归孩子，再组合子树答案",
		Pseudocode: []string{"if node == nil { return base }", "left := dfs(node.Left)", "right := dfs(node.Right)", "return combine(left, right, node)"},
		Steps:      []string{"处理空树边界", "递归左子树", "递归右子树", "后序组合结果"},
		Narration:  []string{"空节点给出最小规模答案。", "左子树先独立返回。", "右子树同样返回。", "父节点只组合已完成的两个孩子结果。"},
	},
	"flow-dfs-grid": {
		Title:      "网格 DFS：进入一个连通块并淹没",
		Pseudocode: []string{"if out of bounds or blocked { return }", "mark current visited", "for dir := range dirs {", "    dfs(next cell)", "}"},
		Steps:      []string{"检查边界与障碍", "进入格子立刻标记", "枚举四个方向", "递归扩展同一连通块"},
		Narration:  []string{"非法格子不进入。", "进入时标记避免绕回。", "方向数组统一邻居生成。", "外层每启动一次 DFS 就覆盖一个连通块。"},
	},
	"flow-dfs-path": {
		Title:      "路径 DFS：选择、递归、撤销",
		Pseudocode: []string{"if target { copy path; return }", "for next := range neighbors {", "    path = append(path, next)", "    dfs(next)", "    path = path[:len(path)-1]", "}"},
		Steps:      []string{"到达终点复制路径", "选择下一节点", "递归深入", "撤销选择"},
		Narration:  []string{"答案必须复制，否则后续回溯会修改它。", "把候选放进当前路径。", "递归探索该分支。", "恢复现场给下一个兄弟分支。"},
	},
	"flow-backtracking-choose-skip": {
		Title:      "选或不选：每个位置分成两条边",
		Pseudocode: []string{"if index == n { collect; return }", "dfs(index + 1)", "path = append(path, nums[index])", "dfs(index + 1)", "pop path"},
		Steps:      []string{"到达叶子收集", "走不选分支", "加入当前元素", "走选择分支并回退"},
		Narration:  []string{"所有位置决策完成时得到一个子集。", "第一条边不改变 path。", "第二条边把当前元素加入。", "返回时撤销，避免影响上一层。"},
	},
	"flow-backtracking-enumeration": {
		Title:      "枚举候选：循环选择下一步",
		Pseudocode: []string{"if path complete { collect; return }", "for candidate := range choices {", "    if invalid { continue }", "    choose; dfs(); undo", "}"},
		Steps:      []string{"确定递归终点", "枚举候选", "跳过非法或重复", "选择、递归、撤销"},
		Narration:  []string{"路径长度或目标条件决定何时收集。", "每一层由循环选择下一步。", "排列用 used，组合用 start 控制候选范围。", "三步缺任何一步都会造成漏解或串状态。"},
	},
	"flow-list-fast-slow": {
		Title:      "快慢指针：相对速度为 1",
		Pseudocode: []string{"slow, fast := head, head", "for fast != nil && fast.Next != nil {", "    slow = slow.Next", "    fast = fast.Next.Next", "    if slow == fast { return true }", "}"},
		Steps:      []string{"两指针同起点", "慢指针走一步", "快指针走两步", "相遇或快指针到尾"},
		Narration:  []string{"初始同起点不代表有环，比较发生在移动后。", "slow 记录较慢位置。", "fast 每轮多走一步。", "有环必追上；无环则 fast 先到 nil。"},
	},
	"flow-list-merge": {
		Title:      "合并链表：tail 持续尾插较小节点",
		Pseudocode: []string{"dummy, tail := new node, dummy", "while a != nil && b != nil {", "    choose smaller head", "    tail.Next = chosen; tail = tail.Next", "}", "tail.Next = remaining"},
		Steps:      []string{"dummy 固定结果头", "比较两条当前头", "接入较小节点", "tail 前进并接上剩余链"},
		Narration:  []string{"dummy 消除第一个节点特判。", "两个输入头都是未合并前缀。", "较小者一定能安全进入结果尾。", "一条链耗尽后另一条已整体有序。"},
	},
	"flow-tree-bst": {
		Title:      "BST 验证：祖先共同收紧值域",
		Pseudocode: []string{"if node == nil { return true }", "if node outside (low, high) { return false }", "validate left with high=node", "validate right with low=node"},
		Steps:      []string{"传入允许区间", "检查当前节点", "递归左子树", "递归右子树"},
		Narration:  []string{"范围来自所有祖先，不只来自父节点。", "当前值必须严格落在开区间。", "左子树上界收紧为当前值。", "右子树下界收紧为当前值。"},
	},
	"flow-tree-lca": {
		Title:      "LCA：左右子树各带回一个目标",
		Pseudocode: []string{"if root is nil, p, or q { return root }", "left := lca(root.Left)", "right := lca(root.Right)", "if either is nil { return the other }", "return root"},
		Steps:      []string{"命中目标或空节点", "搜索左子树", "搜索右子树", "根据两个返回值确定祖先"},
		Narration:  []string{"目标节点本身向上返回。", "左子树可能带回 p、q 或 LCA。", "右子树同理。", "左右非空说明两目标首次在当前节点汇合。"},
	},
	"flow-tree-path-sum": {
		Title:      "路径和：沿根到叶路径递减目标",
		Pseudocode: []string{"if node == nil { return false }", "remain -= node.Val", "if leaf { return remain == 0 }", "return dfs(left, remain) || dfs(right, remain)"},
		Steps:      []string{"进入节点", "扣除当前值", "叶子检查", "向孩子继续"},
		Narration:  []string{"空节点不是一条完整路径。", "remain 始终表示后续还需凑出的和。", "只有叶子才可判断一条根到叶路径完成。", "任一孩子成功即可返回。"},
	},
	"flow-tree-dp": {
		Title:      "树形 DP：孩子先返回 take / skip",
		Pseudocode: []string{"leftTake, leftSkip := dfs(left)", "rightTake, rightSkip := dfs(right)", "take = value + leftSkip + rightSkip", "skip = max(leftTake,leftSkip)+max(rightTake,rightSkip)", "return take, skip"},
		Steps:      []string{"后序得到左状态", "后序得到右状态", "计算选当前节点", "计算不选当前节点"},
		Narration:  []string{"每个孩子返回两种可选状态。", "父节点等到两侧都完成。", "选父节点时孩子被限制为 skip。", "不选父节点时每个孩子可独立取较大值。"},
	},
	"flow-string-window": {
		Title:      "滑动窗口：右扩张、左收缩、合法时记录",
		Pseudocode: []string{"for right := range s {", "    add s[right]", "    for window invalid { remove s[left]; left++ }", "    record answer", "}"},
		Steps:      []string{"右端纳入字符", "更新频次", "条件破坏时收缩左端", "窗口合法后记录答案"},
		Narration:  []string{"右端每轮只前进一次。", "频次表保存窗口状态。", "左端一直收缩到重新合法。", "此时窗口长度才可以作为候选答案。"},
	},
	"flow-string-golang": {
		Title:      "Go 字符串：先选字节还是 rune，再选构造方式",
		Pseudocode: []string{"if input is ASCII { use s[i] }", "else { chars := []rune(s) }", "var builder strings.Builder", "for each output part { builder.WriteString(part) }", "return builder.String()"},
		Steps:      []string{"判断字符集", "选择 byte 或 rune", "创建 Builder", "循环写入并返回"},
		Narration:  []string{"下标、长度、切片都按 byte；题目若含 Unicode 先明确语义。", "ASCII 的 byte 更直接；按字符处理时转为 []rune。", "连续拼接避免用 + 反复创建中间字符串。", "每次写入一段，最后一次性得到结果。"},
	},
	"flow-string-palindrome": {
		Title:      "中心扩展：每个中心向两侧验证",
		Pseudocode: []string{"for center := range s {", "    expand(center, center)", "    expand(center, center+1)", "}", "update longest while characters match"},
		Steps:      []string{"枚举中心", "处理奇数回文", "处理偶数回文", "匹配时向两边扩张"},
		Narration:  []string{"每个回文都可归属到一个中心。", "单点中心覆盖 aba。", "双点中心覆盖 abba。", "首次失配即当前中心无法再扩张。"},
	},
	"flow-string-kmp": {
		Title:      "KMP：失配时跳到最长可复用前缀",
		Pseudocode: []string{"build pi for pattern", "for text character {", "    while mismatch and j>0 { j = pi[j-1] }", "    if match { j++ }", "    if j == len(pattern) { report match }", "}"},
		Steps:      []string{"预处理 pi", "比较当前字符", "失配回退 j", "匹配推进并报告"},
		Narration:  []string{"pi 记录模式前缀和后缀的最大重叠。", "主串指针不会回退。", "回退的是模式已匹配长度 j。", "j 达到模式长度时找到一次匹配。"},
	},
	"flow-lcs-space": {
		Title:      "LCS 一维压缩：保存左上角旧值",
		Pseudocode: []string{"for i := range a {", "    diagonal := 0", "    for j := range b {", "        up := dp[j]", "        update dp[j] from diagonal, up, dp[j-1]", "        diagonal = up", "    }", "}"},
		Steps:      []string{"新行起始 diagonal=0", "读取覆盖前的 up", "读取左方与左上角", "写 dp[j]", "把旧 up 交给下一列"},
		Narration:  []string{"diagonal 对应上一行第 0 列。", "覆盖前 dp[j] 还是上方。", "dp[j-1] 已更新为左方，diagonal 是左上角。", "写入当前行当前位置。", "先保存的旧 up 成为下一格左上角。"},
	},
}
