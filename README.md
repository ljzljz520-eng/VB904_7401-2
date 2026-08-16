# VideoLab

VideoLab 是一个纯 Go 命令行视频制作工作流示例。它从本地固定夹具加载海边、城市和夜景片段，在内存仓储中保存调色方案，按确定顺序生成每个片段的抽帧比较计划，并输出可直接粘贴到 Windows 命令行的 ffmpeg 命令。

## 环境

- Go 1.25.13
- `go.mod` 使用语言版本 `go 1.25`
- 使用前设置 `GOTOOLCHAIN=local`

```bash
export GOTOOLCHAIN=local
go run ./cmd/videolab compare
```

默认入口会输出 3 个片段与 3 个调色方案的 9 条抽帧命令。路径刻意包含空格，命令使用 Windows 双引号保护路径。

保存一个运行期内存调色方案：

```bash
go run ./cmd/videolab preset --name warm --filter 'eq=saturation=1.2' --description '暖色'
```

对固定片段执行 ffmpeg：

```bash
go run ./cmd/videolab extract seaside clean
```

该命令需要本机可执行的 `ffmpeg`。仓储和测试使用固定夹具及内存替身，不依赖视频文件、随机数、睡眠、外部 API 或数据库服务。

业务链路测试命令：

```bash
go test -count=1 ./...
```

其中 `TestCompareReportsDecodeFailureDetails` 会稳定暴露当前视频命令执行器的错误包装缺失：ffmpeg 解码失败时页面层只收到通用的命令失败信息，路径和具体错误输出没有被保留。
