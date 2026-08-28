package trace

type linkedListState struct {
	Chain     []string          `json:"chain"`
	Detached  []string          `json:"detached"`
	Pointers  map[string]string `json:"pointers"`
	Highlight []string          `json:"highlight"`
}

// LinkedListRewireTrace demonstrates dummy-head based head insertion for reversing a sublist.
func LinkedListRewireTrace() Trace {
	result := Trace{
		Kind:  "linked-list",
		Title: "链表重连：dummy、断开与头插",
		Pseudocode: []string{
			"dummy.Next = head",
			"pre := dummy // 走到待翻转区间前一个节点",
			"cur := pre.Next",
			"for i := 0; i < right-left; i++ {",
			"    next := cur.Next",
			"    cur.Next = next.Next // 从原位置断开 next",
			"    next.Next = pre.Next",
			"    pre.Next = next // next 头插到翻转段前端",
			"}",
			"return dummy.Next",
		},
	}
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"head": "1"}, []string{"D", "1"}, 0, "建立 dummy.Next=head。D 与 head 先高亮，返回头从这里统一取得。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, []string{"1", "2"}, 2, "定位完成：pre=1 在翻转区间前，cur=2 是当前尾节点。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, []string{"2", "3"}, 4, "读取 cur.Next：先高亮 2→3，下一步把 next 保存为节点 3。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2", "next": "3"}, []string{"3"}, 4, "next=3 已保存；后续可以修改 2 的 Next 而不丢失节点 3。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "4", "5"}, []string{"3"}, map[string]string{"pre": "1", "cur": "2", "next": "3"}, []string{"2", "3", "4"}, 5, "执行 cur.Next=next.Next：2 直接连到 4，3 暂时脱离主链。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "4", "5"}, []string{"3"}, map[string]string{"pre": "1", "cur": "2", "next": "3"}, []string{"1", "2", "3"}, 6, "执行 next.Next=pre.Next：脱离的 3 先指向原翻转段头 2。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, []string{"1", "3"}, 7, "执行 pre.Next=next：3 插到 1 后，第一轮头插完成。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, []string{"2", "4"}, 4, "第二轮读取 cur.Next：cur 仍是 2，next 候选是 4。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2", "next": "4"}, []string{"4"}, 4, "next=4 已保存。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "5"}, []string{"4"}, map[string]string{"pre": "1", "cur": "2", "next": "4"}, []string{"2", "4", "5"}, 5, "执行 cur.Next=next.Next：2 直接连到 5，4 脱离主链。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "5"}, []string{"4"}, map[string]string{"pre": "1", "cur": "2", "next": "4"}, []string{"1", "3", "4"}, 6, "执行 next.Next=pre.Next：4 指向当前翻转段头 3。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "4", "3", "2", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, []string{"1", "4"}, 7, "执行 pre.Next=next：4 插到 1 后，区间 [2,4] 已反转。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "4", "3", "2", "5"}, nil, map[string]string{"head": "1"}, []string{"D", "1"}, 9, "返回 dummy.Next，得到 1→4→3→2→5。"))
	return result
}

func linkedListFrame(chain, detached []string, pointers map[string]string, highlight []string, line int, narration string) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"operation": "保存 next → 断开 → 重连"}, State: linkedListState{Chain: chain, Detached: detached, Pointers: pointers, Highlight: highlight}}
}
