package trace

type linkedListState struct {
	Chain    []string          `json:"chain"`
	Detached []string          `json:"detached"`
	Pointers map[string]string `json:"pointers"`
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
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"head": "1"}, 0, "dummy 指向原 head：即使翻转从第一个节点开始，返回头也不需要特判。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "3", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, 2, "left=2，因此 pre 停在节点 1，cur 是翻转段当前头节点 2。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "2", "4", "5"}, []string{"3"}, map[string]string{"pre": "1", "cur": "2", "next": "3"}, 5, "先断开 3：2 的 next 改为 4，3 暂时脱离主链。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "4", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, 7, "把 3 头插到 pre 后：局部顺序从 2→3 变为 3→2。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "3", "2", "5"}, []string{"4"}, map[string]string{"pre": "1", "cur": "2", "next": "4"}, 5, "下一轮断开 4：cur 始终是翻转段尾部，不需要移动。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "4", "3", "2", "5"}, nil, map[string]string{"pre": "1", "cur": "2"}, 7, "4 再头插，区间 [2,4] 完成反转。dummy.Next 仍是最终头。"))
	result.Frames = append(result.Frames, linkedListFrame([]string{"D", "1", "4", "3", "2", "5"}, nil, map[string]string{"head": "1"}, 9, "返回 dummy.Next，得到 1→4→3→2→5。"))
	return result
}

func linkedListFrame(chain, detached []string, pointers map[string]string, line int, narration string) Frame {
	return Frame{ActiveLine: line, Narration: narration, Variables: map[string]string{"operation": "断开 → 头插 → 重连"}, State: linkedListState{Chain: chain, Detached: detached, Pointers: pointers}}
}
