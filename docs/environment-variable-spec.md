# Stellar Axis 基础应用环境变量规范

本文档是 `stellspec-go-sdk` 仓库内的实现参考副本。

全体系的权威版本位于：

- `E:\PersonalCode\stellar\docs\environment-variable-spec.md`

当前文档仅用于说明 `stellspec-go-sdk` 如何落地该规范。

设计目标如下：

- 所有中间件使用统一前缀和统一语义
- 应用名称、版本、环境、部署拓扑等元数据由平台统一注入
- 中间件在业务应用无显式配置时，也能自动识别当前应用身份
- 规范同时适用于 Kubernetes 注入与物理机 / 虚拟机装机注入

## 命名原则

- 全局统一前缀使用 `STELLAR_`
- 变量名统一使用大写英文加下划线
- 变量语义必须稳定，不因单个中间件变化而变化
- 中间件私有配置不得污染全局基础元数据命名空间

## 优先级原则

推荐所有中间件统一遵循以下优先级：

1. 代码显式配置
2. 产品级环境变量
3. 全局 `STELLAR_*` 环境变量
4. 中间件自身默认值

说明如下：

- `STELLAR_*` 用于平台统一注入的基础应用元数据
- 产品级环境变量用于单个中间件局部覆盖
- 代码显式配置始终拥有最高优先级

## 必选基础变量

以下变量建议作为大型企业的最小必选集合：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_APP_NAME` | `user-service` | 应用或服务名，全体系唯一识别的核心主键之一 |
| `STELLAR_APP_NAMESPACE` | `stellar.trade` | 应用逻辑命名空间，用于区分业务域或产品域 |
| `STELLAR_APP_VERSION` | `1.4.2` | 当前应用版本 |
| `STELLAR_APP_INSTANCE_ID` | `user-service-7f6d9d6d7b-2xk9p` | 当前运行实例 ID |
| `STELLAR_ENV` | `dev` / `test` / `prod` | 部署环境标识 |

## 推荐拓扑变量

以下变量建议由平台侧统一注入，供日志、链路、指标、配置等中间件共同使用：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_CLUSTER` | `cluster-sh-prod-01` | 集群标识 |
| `STELLAR_REGION` | `cn-east-1` | 大区标识 |
| `STELLAR_ZONE` | `cn-east-1a` | 可用区标识 |
| `STELLAR_IDC` | `sh-a` | 机房或园区标识 |
| `STELLAR_HOST_NAME` | `node-01` | 主机名 |
| `STELLAR_HOST_IP` | `10.10.0.11` | 主机 IP |

## Kubernetes 推荐变量

如果运行在 Kubernetes 中，建议额外统一注入：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_NODE_NAME` | `worker-node-01` | 节点名称 |
| `STELLAR_K8S_NAMESPACE` | `trade` | Kubernetes Namespace |
| `STELLAR_POD_NAME` | `user-service-7f6d9d6d7b-2xk9p` | Pod 名称 |
| `STELLAR_POD_IP` | `172.20.10.23` | Pod IP |
| `STELLAR_CONTAINER_NAME` | `app` | 容器名称 |

## 统一语义约束

为避免平台长期演进后语义漂移，推荐固定如下解释：

- `APP_NAME`
  指业务应用、服务或工作负载的稳定名称，不使用机器名或 Pod 名替代
- `APP_NAMESPACE`
  指逻辑业务域，不等同于 Kubernetes Namespace
- `APP_INSTANCE_ID`
  指单个运行实例的唯一标识，可使用 Pod 名、实例 ID 或平台生成值
- `ENV`
  只用于表达部署环境，不混入区域、租户、集群等信息
- `CLUSTER`
  指部署集群，不限制底层是 Kubernetes 还是其他编排平台
- `REGION / ZONE / IDC`
  用于表达地理与基础设施拓扑，不与环境语义混用

## 产品级覆盖建议

在统一 `STELLAR_*` 的前提下，各中间件可以保留自己的产品级覆盖变量，例如：

- 日志平台：`STELLSPEC_*`
- 链路平台：`STARTRACE_*`
- 配置中心：`NEBULA_*`

但产品级变量应只覆盖本产品关心的局部配置，不应重新定义全局基础元数据的含义。

以日志 SDK 为例：

- `STELLSPEC_SERVICE_NAME` 可以覆盖 `STELLAR_APP_NAME`
- `STELLSPEC_ENDPOINT` 可以定义日志上报地址
- `STELLSPEC_OUTPUT` 可以定义日志输出方式

## Kubernetes 注入建议

推荐由平台在 Pod 启动时统一注入这些变量，而不是由业务侧单独维护。典型来源包括：

- `metadata.name`
- `metadata.namespace`
- `spec.nodeName`
- `status.podIP`
- `status.hostIP`
- 镜像版本或发布系统注入的应用版本号

## 对 Stellspec SDK 的约束

当前 `stellspec-go-sdk` 需要遵循本规范：

- 默认读取 `STELLAR_*` 基础元数据
- 将这些元数据写入 OpenTelemetry Resource
- 允许 `STELLSPEC_*` 对日志 SDK 局部行为进行覆盖
- 在业务代码未显式传入 `ServiceName`、`ServiceVersion` 等字段时，自动从 `STELLAR_*` 中补全

## 架构建议：抽离通用基础元数据 SDK

从平台长期演进角度看，`STELLAR_*` 这批变量更适合作为“全中间件共享的基础应用元数据协议”，而不是由每个 SDK 各自重复解析一遍。推荐后续单独沉淀一个基础 SDK 或基础模块，例如：

- `stellar-appmeta-sdk`
- `stellar-env-sdk`
- `stellar-runtime-meta`

推荐该基础 SDK 统一提供以下能力：

- 统一解析 `STELLAR_*` 环境变量
- 提供强类型配置模型，例如 `AppMeta`、`RuntimeMeta`
- 统一优先级与默认值策略
- 提供转 OTel Resource、日志字段、链路标签、指标标签的适配能力
- 作为各中间件 SDK 的共同依赖

这样做的收益包括：

- 避免每个 SDK 都重复维护一遍相同字段和解析逻辑
- 保证日志、链路、指标、配置等系统对同一组基础元数据的理解一致
- 平台新增或调整基础字段时，只需在公共 SDK 中演进一次
- 业务应用和各中间件 SDK 的接入体验更统一
