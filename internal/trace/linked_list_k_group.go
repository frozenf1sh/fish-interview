package trace

// LinkedListKGroupTrace demonstrates reverseKGroup on one chain.  The main
// chain never disappears: locating a group, detaching it, reversing one link
// at a time, and reconnecting it are separate observable states.
func LinkedListKGroupTrace() Trace {
	code := []string{
		"groupPrev := dummy",
		"for {",
		"    kth := getKth(groupPrev, k)",
		"    if kth == nil { break }",
		"    groupNext := kth.Next",
		"    prev, cur := groupNext, groupPrev.Next",
		"    for cur != groupNext {",
		"        next := cur.Next",
		"        cur.Next = prev",
		"        prev, cur = cur, next",
		"    }",
		"    oldHead := groupPrev.Next",
		"    groupPrev.Next = kth",
		"    groupPrev = oldHead",
		"}",
		"return dummy.Next",
	}
	caption := "链表：每 k 个节点一组翻转"
	frame := func(line int, narration, phase string, chain, detached, working []string, pointers map[string]string, highlight, group []string) Frame {
		return kGroupListFrame(line, narration, caption, map[string]string{
			"k": "3", "groupPrev": pointers["groupPrev"], "kth": pointers["kth"],
			"groupNext": pointers["groupNext"], "cur": pointers["cur"], "prev": pointers["prev"],
		}, chain, detached, working, pointers, highlight, group, phase)
	}
	main := []string{"D", "1", "2", "3", "4", "5"}
	frames := []Frame{
		frame(0, "例题 D→1→2→3→4→5，k=3。dummy 固定在头部；目标是每三个节点翻转，余下不足一组的节点保持原样。", "准备", main, nil, nil, map[string]string{"groupPrev": "D", "kth": "-", "cur": "-", "prev": "-"}, []string{"D"}, []string{"1", "2", "3"}),
		frame(2, "从 groupPrev=D 出发定位第 3 个节点；此时 kth 还没有走到目标位置。", "定位第 1 组", main, nil, nil, map[string]string{"groupPrev": "D", "kth": "D", "cur": "1", "prev": "-"}, []string{"D", "1"}, []string{"1", "2", "3"}),
		frame(2, "kth 前进到 1，已经数到第 1 个组内节点；主链保持不变。", "定位第 1 组", main, nil, nil, map[string]string{"groupPrev": "D", "kth": "1", "cur": "1", "prev": "-"}, []string{"1"}, []string{"1", "2", "3"}),
		frame(2, "kth 前进到 2，继续寻找完整的三节点组。", "定位第 1 组", main, nil, nil, map[string]string{"groupPrev": "D", "kth": "2", "cur": "1", "prev": "-"}, []string{"2"}, []string{"1", "2", "3"}),
		frame(2, "kth 到达 3，第一组完整；下一条 groupNext=4 也被保存，防止翻转时丢掉后半段。", "找到完整组", main, nil, nil, map[string]string{"groupPrev": "D", "kth": "3", "groupNext": "4", "cur": "1", "prev": "4"}, []string{"3", "4"}, []string{"1", "2", "3"}),
		frame(4, "保存 groupNext=4 后，从主链断开 1→2→3；后半段 4→5 仍留在主链上。", "断开第 1 组", []string{"D", "4", "5"}, []string{"1", "2", "3"}, nil, map[string]string{"groupPrev": "D", "kth": "3", "groupNext": "4", "prev": "4", "cur": "1"}, []string{"1", "3", "4"}, []string{"1", "2", "3"}),
		frame(5, "初始化局部翻转：prev 指向 groupNext=4，cur 指向组头 1；临时链先显示待反转的组内节点。", "初始化翻转", []string{"D", "4", "5"}, []string{"1", "2", "3"}, []string{"1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "4", "cur": "1"}, []string{"1", "4"}, []string{"1", "2", "3"}),
		frame(7, "读取 cur=1 的 next=2；先把 next 保存下来，才能让 cur 前进。", "读取 next", []string{"D", "4", "5"}, []string{"1", "2", "3"}, []string{"1", "4"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "4", "cur": "1", "next": "2"}, []string{"1", "2"}, []string{"1", "2", "3"}),
		frame(8, "执行 1.Next=4：第一个节点反向接到 groupNext 前面；临时链只绘制真实的组内节点，4 仍留在主链避免重复。", "改写一条指针", []string{"D", "4", "5"}, []string{"2", "3"}, []string{"1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "1", "cur": "2"}, []string{"1", "4"}, []string{"1", "2", "3"}),
		frame(9, "prev 前进到 1，cur 前进到 2；主链、待处理组和临时链仍然同时可见。", "推进指针", []string{"D", "4", "5"}, []string{"2", "3"}, []string{"1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "1", "cur": "2"}, []string{"1", "2"}, []string{"1", "2", "3"}),
		frame(7, "读取 cur=2 的 next=3；这次仍先读后写。", "读取 next", []string{"D", "4", "5"}, []string{"2", "3"}, []string{"1", "4"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "1", "cur": "2", "next": "3"}, []string{"2", "3"}, []string{"1", "2", "3"}),
		frame(8, "执行 2.Next=1，第二个节点插到临时链头，变成 2→1；临时链的尾部仍指向主链中的 4。", "改写一条指针", []string{"D", "4", "5"}, []string{"3"}, []string{"2", "1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "2", "cur": "3"}, []string{"2", "1"}, []string{"1", "2", "3"}),
		frame(9, "prev 前进到 2，cur 前进到 3；最后一个组内节点待处理。", "推进指针", []string{"D", "4", "5"}, []string{"3"}, []string{"2", "1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "2", "cur": "3"}, []string{"2", "3"}, []string{"1", "2", "3"}),
		frame(7, "读取 cur=3 的 next=4；4 是原组之后的节点，不能被重新翻转。", "读取 next", []string{"D", "4", "5"}, []string{"3"}, []string{"2", "1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "2", "cur": "3", "next": "4"}, []string{"3", "4"}, []string{"1", "2", "3"}),
		frame(8, "执行 3.Next=2；三次改写完成，临时链为 3→2→1，尾部仍接着 groupNext=4。", "完成组内翻转", []string{"D", "4", "5"}, nil, []string{"3", "2", "1"}, map[string]string{"groupPrev": "D", "groupNext": "4", "prev": "3", "cur": "4"}, []string{"3", "2", "1"}, []string{"1", "2", "3"}),
		frame(12, "oldHead=1 是反转后这一组的尾部；执行 D.Next=3，把临时链覆盖回主链。", "重新接回主链", []string{"D", "3", "2", "1", "4", "5"}, nil, nil, map[string]string{"groupPrev": "D", "kth": "3", "groupNext": "4", "oldHead": "1", "cur": "4"}, []string{"D", "3", "1", "4"}, []string{"1", "2", "3"}),
		frame(13, "groupPrev 移到旧头 1；下一组从 1 后面开始定位。", "准备下一组", []string{"D", "3", "2", "1", "4", "5"}, nil, nil, map[string]string{"groupPrev": "1", "kth": "1", "cur": "4"}, []string{"1", "4"}, []string{"4", "5"}),
		frame(2, "定位第二组时 kth 走到 4，再走到 5；只有两个节点，不足 k=3。", "定位剩余节点", []string{"D", "3", "2", "1", "4", "5"}, nil, nil, map[string]string{"groupPrev": "1", "kth": "5", "cur": "4"}, []string{"4", "5"}, []string{"4", "5"}),
		frame(3, "kth=nil 表示剩余链长度不足 3；不翻转 4→5，直接返回 dummy.Next。", "保留不足一组", []string{"D", "3", "2", "1", "4", "5"}, nil, nil, map[string]string{"groupPrev": "1", "kth": "nil", "cur": "-"}, []string{"4", "5"}, nil),
		frame(15, "完成：D→3→2→1→4→5。每个完整组都在断开、逐指针反转、重新连接后才进入最终链。", "完成", []string{"D", "3", "2", "1", "4", "5"}, nil, nil, map[string]string{"groupPrev": "1", "kth": "nil"}, []string{"D", "3", "2", "1", "4", "5"}, nil),
	}
	return concreteTrace("linked-list-k-group", "链表：k 个一组翻转（定位、断开与重连）", code, frames...)
}
