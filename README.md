# CodeVerify

面向 Go 项目的代码验证 Agent 原型，使用 **Go + Responses API + go/packages/AST**，通过工具调用将代码分析、风险记录、测试生成和验证结果串联为一个 Agent Loop。

> 个人实验项目，仍在开发中。当前版本为演示/测试原型，**不应直接处理不可信仓库**。

## 核心流程

1. 扫描 Go 仓库并构建函数索引与静态调用关系。
2. 使用多轮 Agent Loop 调用只读文件工具、调用关系查询和风险记录工具。
3. 在受约束的路径下写入 `_test.go`，运行 `go test ./...`。
4. 关联工具结果、风险 ID 和验证证据，输出裁决和总结。
5. 通过 Transcript、近期轮次与风险摘要维护上下文。

## 目录结构

- `cmd/main.go`：CLI 演示入口；**当前目标仓库路径在代码中固定为测试 fixture**。
- `internal/agent/`：Responses API 接入与多轮循环。
- `internal/analysis/`：Go 包、函数索引与静态调用图。
- `internal/core/`：风险、证据、状态及仓库路径处理。
- `internal/tool/`：文件读取、测试写入、执行和风险判定工具。
- `testdata/`：用于演示代码分析和测试行为的 Go fixture。

## 运行方法

需要与 `go.mod` 匹配的 Go 工具链，以及兼容 Responses API 的模型接口。当前代码使用 `deepseek-flash`，从环境变量 `DEEPSEEK_API_KEY` 读取密钥。**不要将密钥写入仓库或日志**。

```bash
export DEEPSEEK_API_KEY='your-key-from-secure-store'
go run ./cmd
```

当前入口使用 `./testdata/codeverify-callgraph-fixture`，并非支持任意仓库参数的成熟 CLI。API 调用会产生费用；测试运行过程中可能执行生成的代码。

## 安全边界与限制

- 写文件工具只接受 `_test.go`，并执行仓库路径边界校验；这**不等于隔离沙箱**。
- `go test ./...` 当前经由宿主环境的 `os/exec` 执行，未实现容器沙箱、资源/网络限制或统一权限控制。
- 静态调用关系不代表运行时所有动态调用；证据判定仍需人工复核。
- 目前没有发布误报率、缺陷发现率或端到端性能评测，不应把它描述成生产可用的自动修复/审计工具。
- `run.log` 是调试输出，不应包含生产环境或个人数据；后续日志应仅在本地保留。

## 后续方向

测试执行隔离、可复现验证样例、失败恢复、结果可视化，以及对不同 Go 项目的通用入口。
