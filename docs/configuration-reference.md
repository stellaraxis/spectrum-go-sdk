# Spectrum Go SDK 配置说明

本文档使用中文详细说明 `config.Config` 中每个字段的语义，以及 `STELLAR_*`、`SPECTRUM_*` 环境变量在当前 SDK 中的含义、默认行为与推荐使用方式。

## Config 字段说明

| 字段 | 类型 | 默认值 | 详细说明 |
| :--- | :--- | :--- | :--- |
| `ServiceName` | `string` | 无 | 当前服务的主标识，会写入 OTel `service.name`。这是日志平台做服务维度查询、聚合、告警路由时最核心的字段之一。 |
| `ServiceNamespace` | `string` | 空 | 当前服务所属的逻辑命名空间。它用于表达业务域或产品线边界，不建议直接复用 Kubernetes Namespace。 |
| `ServiceVersion` | `string` | 空 | 当前服务版本，推荐使用发布版本号、镜像 tag 或构建版本。它是发布回溯、版本对比和问题归因的重要字段。 |
| `ServiceInstanceID` | `string` | 空 | 当前运行实例的唯一标识。在容器场景可使用 Pod 名，在物理机/虚拟机场景可使用实例 ID。 |
| `Environment` | `string` | 空 | 部署环境标识，例如 `dev`、`test`、`prod`。它只表达环境，不应混入区域、机房或租户信息。 |
| `Cluster` | `string` | 空 | 集群标识，用于表达当前实例所在的部署集群，便于跨集群检索日志和排查集群级故障。 |
| `Region` | `string` | 空 | 大区或地域标识，通常对应云厂商地域或公司内部的大区概念。 |
| `Zone` | `string` | 空 | 可用区标识，用于排查单可用区故障、路由不均衡和局部抖动问题。 |
| `IDC` | `string` | 空 | 机房、园区或基础设施单元标识，适合在更细粒度设施层面定位问题。 |
| `HostName` | `string` | 空 | 宿主机名。在容器场景中通常是节点主机名，在主机部署场景中就是当前机器名。 |
| `HostIP` | `string` | 空 | 宿主机 IP，用于排查节点网络、宿主机路由与主机级别问题。 |
| `NodeName` | `string` | 空 | Kubernetes 节点名，仅在 K8s 场景下有意义。 |
| `K8sNamespace` | `string` | 空 | Kubernetes Namespace，用于表达运行时资源隔离空间。 |
| `PodName` | `string` | 空 | 当前 Pod 名称，常用于容器化场景下定位单实例。 |
| `PodIP` | `string` | 空 | 当前 Pod IP，用于网络问题、服务发现问题定位。 |
| `ContainerName` | `string` | 空 | 当前容器名称。当一个 Pod 中存在多个容器时可区分具体日志来源。 |
| `Endpoint` | `string` | `localhost:4317` | OTLP 导出目标地址。生产环境通常应指向本机 log-agent。 |
| `Insecure` | `bool` | `true` | 是否关闭 OTLP/gRPC 传输安全。面向本机 agent 或可信内网链路时常设为 `true`。 |
| `Protocol` | `string` | `grpc` | 当前 SDK 支持的导出协议。现阶段仅支持 OTLP/gRPC。 |
| `Format` | `string` | `json` 或 `console` | 本地输出格式。开发环境通常使用 `console` 便于阅读，结构化输出通常使用 `json`。 |
| `Output` | `string` | `otlp` 或 `console` | 日志输出目标。生产环境默认走 `otlp`，开发环境默认走控制台输出。 |
| `Level` | `string` | `info` | 最小日志级别。低于该级别的日志通常不会进入导出链路。 |
| `Development` | `bool` | `false` | 是否按开发环境模式运行。为 `true` 时，SDK 会优先选择对调试友好的默认行为。 |
| `EnableCaller` | `bool` | `true` | 是否写出调用方文件、行号、函数名。开启后更易排查，但会带来额外开销。 |
| `EnableStacktrace` | `bool` | `true` | 是否让错误日志自动带堆栈。利于故障定位，但会显著放大日志体积。 |
| `BatchTimeout` | `time.Duration` | `5s` | BatchProcessor 的批量导出周期。值越小延迟越低，值越大吞吐越高。 |
| `ExportTimeout` | `time.Duration` | `3s` | 单次导出的超时时间。超时后当前批次会失败，并进入失败兜底逻辑。 |
| `MaxBatchSize` | `int` | `512` | 单批次最大日志条数。超过后会拆分导出。 |
| `MaxQueueSize` | `int` | `2048` | BatchProcessor 队列最大容量。队列满时 OTel 会丢弃最老记录，因此需要按峰值流量谨慎设置。 |
| `FallbackFilePath` | `string` | `logs/spectrum-fallback.log` | OTLP 推送最终失败后的本地兜底文件路径。当前实现按 JSON 行顺序追加。 |
| `Retry.Enabled` | `*bool` | `true` | 是否启用 OTLP exporter 的瞬时错误重试。关闭后遇到可重试错误会直接失败。 |
| `Retry.InitialInterval` | `*time.Duration` | `5s` | 第一次重试前等待多久。适合处理 log-agent 短暂重启、端口瞬断等场景。 |
| `Retry.MaxInterval` | `*time.Duration` | `30s` | 指数退避过程中的最大等待时间，避免长时间故障时频繁重试。 |
| `Retry.MaxElapsedTime` | `*time.Duration` | `1m` | 单批日志允许重试的总时长。到达上限后会放弃推送并触发本地兜底。 |
| `Headers` | `map[string]string` | 空 | OTLP 额外请求头，常用于鉴权、租户标识、平台路由等场景。 |
| `ResourceAttributes` | `map[string]string` | 空 | 追加到 OTel Resource 的扩展属性，适合补充业务侧自定义元数据。 |

## 环境变量说明

### 全局基础变量 `STELLAR_*`

这些变量属于全体系的基础元数据协议，不只服务于日志 SDK，而是推荐供日志、链路、指标、配置等所有中间件共同消费。

| 环境变量 | 映射字段 | 详细说明 |
| :--- | :--- | :--- |
| `STELLAR_APP_NAME` | `ServiceName` | 当前业务应用或服务名称，是平台识别服务身份的基础主键之一。 |
| `STELLAR_APP_NAMESPACE` | `ServiceNamespace` | 当前业务逻辑命名空间，用于表达业务域、产品域或组织边界。 |
| `STELLAR_APP_VERSION` | `ServiceVersion` | 当前应用版本，推荐由发布系统统一注入。 |
| `STELLAR_APP_INSTANCE_ID` | `ServiceInstanceID` | 当前运行实例唯一标识，推荐由运行平台统一生成。 |
| `STELLAR_ENV` | `Environment` | 当前部署环境，例如 `dev`、`test`、`prod`。 |
| `STELLAR_CLUSTER` | `Cluster` | 当前实例所属集群标识。 |
| `STELLAR_REGION` | `Region` | 当前实例所属地域或大区。 |
| `STELLAR_ZONE` | `Zone` | 当前实例所属可用区。 |
| `STELLAR_IDC` | `IDC` | 当前实例所属机房或园区。 |
| `STELLAR_HOST_NAME` | `HostName` | 当前宿主机主机名。 |
| `STELLAR_HOST_IP` | `HostIP` | 当前宿主机 IP。 |
| `STELLAR_NODE_NAME` | `NodeName` | 当前 Kubernetes 节点名称。 |
| `STELLAR_K8S_NAMESPACE` | `K8sNamespace` | 当前 Kubernetes Namespace。 |
| `STELLAR_POD_NAME` | `PodName` | 当前 Pod 名称。 |
| `STELLAR_POD_IP` | `PodIP` | 当前 Pod IP。 |
| `STELLAR_CONTAINER_NAME` | `ContainerName` | 当前容器名称。 |

### Spectrum 产品级变量 `SPECTRUM_*`

这些变量只作用于 `spectrum-go-sdk`，用于覆盖日志 SDK 自身的行为。推荐优先级如下：

1. 代码显式配置
2. `SPECTRUM_*`
3. `STELLAR_*`
4. SDK 默认值

| 环境变量 | 映射字段 | 详细说明 |
| :--- | :--- | :--- |
| `SPECTRUM_SERVICE_NAME` | `ServiceName` | 覆盖全局 `STELLAR_APP_NAME`，用于日志 SDK 局部重命名服务标识。 |
| `SPECTRUM_SERVICE_NAMESPACE` | `ServiceNamespace` | 覆盖全局逻辑命名空间。 |
| `SPECTRUM_SERVICE_VERSION` | `ServiceVersion` | 覆盖全局服务版本。 |
| `SPECTRUM_SERVICE_INSTANCE_ID` | `ServiceInstanceID` | 覆盖全局实例 ID。 |
| `SPECTRUM_ENVIRONMENT` | `Environment` | 覆盖全局环境标识。 |
| `SPECTRUM_ENDPOINT` | `Endpoint` | 指定 OTLP 推送目标地址，生产环境通常应为本机 log-agent。 |
| `SPECTRUM_PROTOCOL` | `Protocol` | 指定导出协议，当前仅支持 `grpc`。 |
| `SPECTRUM_OUTPUT` | `Output` | 指定输出目标，例如 `otlp`、`stdout`、`stderr`、`console`。 |
| `SPECTRUM_FORMAT` | `Format` | 指定本地输出格式，例如 `json` 或 `console`。 |
| `SPECTRUM_LEVEL` | `Level` | 指定最小日志级别。 |
| `SPECTRUM_INSECURE` | `Insecure` | 指定是否关闭 OTLP/gRPC 传输安全。 |
| `SPECTRUM_DEVELOPMENT` | `Development` | 指定是否启用开发模式。 |
| `SPECTRUM_ENABLE_CALLER` | `EnableCaller` | 指定是否收集调用方位置。 |
| `SPECTRUM_ENABLE_STACKTRACE` | `EnableStacktrace` | 指定是否自动写出错误堆栈。 |
| `SPECTRUM_BATCH_TIMEOUT` | `BatchTimeout` | 指定 BatchProcessor 的导出周期。 |
| `SPECTRUM_EXPORT_TIMEOUT` | `ExportTimeout` | 指定单次导出的超时时间。 |
| `SPECTRUM_MAX_BATCH_SIZE` | `MaxBatchSize` | 指定单次导出的最大日志条数。 |
| `SPECTRUM_MAX_QUEUE_SIZE` | `MaxQueueSize` | 指定 BatchProcessor 最大队列容量。 |
| `SPECTRUM_FALLBACK_FILE_PATH` | `FallbackFilePath` | 指定 OTLP 推送最终失败后的本地兜底文件路径。 |
| `SPECTRUM_RETRY_ENABLED` | `Retry.Enabled` | 指定是否开启 exporter 级别重试。 |
| `SPECTRUM_RETRY_INITIAL_INTERVAL` | `Retry.InitialInterval` | 指定首次重试等待时间，例如 `5s`。 |
| `SPECTRUM_RETRY_MAX_INTERVAL` | `Retry.MaxInterval` | 指定指数退避过程中的最大等待时间，例如 `30s`。 |
| `SPECTRUM_RETRY_MAX_ELAPSED_TIME` | `Retry.MaxElapsedTime` | 指定单批日志总重试时长，例如 `1m`。 |

## 超长日志截断说明

当前 SDK 内部已经实现固定长度截断：

- 最大正文长度固定为 `32 KiB`
- 该值当前不对业务侧开放，目的是把消息体大小控制前移到 SDK，而不是压给低频更新的 log-agent
- 截断时保证 UTF-8 合法
- 截断后会自动补充以下属性：
  - `log.body_truncated=true`
  - `log.body_original_size`
  - `log.body_max_size`

这样做的收益是：

- 控制 OTLP 传输体积
- 降低 BatchProcessor 队列内存压力
- 避免极端长日志把本机 log-agent 压成“策略执行器”
- 让日志体积策略在各语言 SDK 中统一实现
