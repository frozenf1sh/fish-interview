---
id: algo.string.golang
kind: concept
title: 字符串：Go 的 byte、rune 与构造模板
summary: byte 适合 ASCII 下标访问，rune 适合 Unicode 字符；频繁拼接用 strings.Builder，切片索引按字节而非字符。
parents: [algo.patterns.string]
tags: [string, golang, unicode]
links: [algo.string.window]
trace: flow-string-golang
---

## 先判断题目字符集

小写英文字母可直接按 byte：`s[i]` 与 `[26]int` 最快。可能含中文、emoji 或任意 Unicode 时，转为 rune：

```go
for _, r := range s { // r 是 Unicode code point
	_ = r
}
runes := []rune(s)
answer := string(runes[left:right])
```

`len(s)` 是字节数，不能和 rune 下标混用。

## 常用构造

```go
var b strings.Builder
b.Grow(len(s))        // 已知大致长度时预分配
for i := range s { b.WriteByte(s[i]) }
result := b.String()
```

反复 `result += part` 会多次复制。分割单词用 `strings.Fields`，按分隔符用 `strings.Split`；需要原地排序字符时转 `[]byte` 后 `sort.Slice`。
