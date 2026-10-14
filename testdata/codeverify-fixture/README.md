# CodeVerify Go Fixture

一个用于测试 Go 代码验证 Agent 的独立、小型仓库。无第三方依赖。

## 运行

```bash
go test ./...
```

现有测试覆盖了常见正常路径和部分错误路径，但**不保证覆盖全部业务约束**。

## 验证任务

请让 CodeVerify 自主分析这个仓库：

- 使用 Repository Overview 选择需要深入检查的函数。
- 使用 `get_code_context` 查询静态调用关系，并按需读取源码。
- 从源码注释、函数签名和调用链推导可检验的业务约束。
- **只新增或修改 `_test.go` 测试文件，不要修改生产代码。**
- 使用 `go test ./...` 生成可执行证据，并给出对应的风险结论。

## 仓库结构

```text
ledger/account.go         账户创建、扣款、入账和余额读取
ledger/account_test.go    账户行为的现有测试
service/transfer.go       跨账户业务操作
service/transfer_test.go  跨账户操作的现有测试
```

这是一个内存模型，不测试并发安全、数据库持久化或外部服务故障。
