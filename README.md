# Fish Interview

一个本地优先的秋招知识查询系统：从知识树定位概念，从算法模式理解题型识别，并用逐帧动画观察算法执行。

它不是进度 Dashboard，也不以公司专区为主。公司考情只是带来源和时间边界的可聚合标签。

## Demo scope

- 知识树与关键词搜索
- 可引用的小知识卡片
- 贪心、DP、Kafka 的种子内容
- Go 生成的算法执行 Trace 与浏览器逐步播放

## Run

```sh
mise exec go@1.26 -- go run ./cmd/fish-interview
```

Open <http://localhost:8080>.

The server listens on `0.0.0.0:8080` by default, so devices on the same LAN can use the Mac's LAN IP and port `8080`. Use `-addr 127.0.0.1:8080` when local-only binding is required.

## Quality checks

```sh
mise exec go@1.26 -- go test ./...
mise exec go@1.26 -- go vet ./...
git diff --check
```

See [architecture.md](docs/architecture.md) for content contracts and extension boundaries.
