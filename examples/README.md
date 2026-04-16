# Examples

本目录提供 `spectrum-go-sdk` 的最小可运行示例。

## Zap 示例

```bash
go run ./examples/zap
```

## Slog 示例

```bash
go run ./examples/slog
```

## OTLP 联调示例

以下示例面向生产形态联调，默认通过 OTLP gRPC 上报到本机 `localhost:4317`：

```bash
go run ./examples/otlp/zap
go run ./examples/otlp/slog
```

运行前建议先启动本机 Agent / Collector，并监听 `localhost:4317`。

这两份示例还演示了：

- 由上层框架构造标准请求上下文
- 将 `request_id`、`tenant_id`、`traceparent` 写入标准上下文
- SDK 从标准上下文中自动注入日志字段

## 可选环境变量

示例默认以开发模式运行，并输出到本地控制台。

如果希望通过统一环境变量注入应用元数据，可以先设置：

```bash
set STELLAR_APP_NAME=example-service
set STELLAR_APP_NAMESPACE=stellar.examples
set STELLAR_APP_VERSION=1.0.0
set STELLAR_ENV=dev
```
