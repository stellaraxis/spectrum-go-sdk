package config

import "time"

const (
	DefaultProtocol         = "grpc"
	DefaultFormat           = "json"
	DefaultOutput           = "otlp"
	DefaultFallbackFilePath = "logs/stellspec-fallback.log"
	defaultRetryEnabled     = true
	defaultRetryInitial     = 5 * time.Second
	defaultRetryMaxInterval = 30 * time.Second
	defaultRetryMaxElapsed  = time.Minute

	OutputOTLP    = "otlp"
	OutputStdout  = "stdout"
	OutputStderr  = "stderr"
	OutputConsole = "console"

	FormatJSON    = "json"
	FormatConsole = "console"
)

// RetryConfig 定义 OTLP 导出失败后的重试策略。
// 这里的重试最终会映射到 otlploggrpc.WithRetry(...)。
type RetryConfig struct {
	// Enabled 控制是否开启 exporter 级别的瞬时错误重试。
	// 默认值为 true；如果显式设置为 false，则 exporter 遇到可重试错误时会直接失败返回。
	Enabled *bool
	// InitialInterval 表示第一次重试前的等待时间。
	// 默认值为 5 秒，适合本机 log-agent 短暂重启、端口短暂不可用等场景。
	InitialInterval *time.Duration
	// MaxInterval 表示指数退避过程中的最大等待间隔。
	// 默认值为 30 秒，用于避免在 log-agent 长时间不可用时持续高频重试。
	MaxInterval *time.Duration
	// MaxElapsedTime 表示单批日志允许重试的总时长。
	// 默认值为 1 分钟；到达该时长后 exporter 会放弃本次推送并返回错误。
	MaxElapsedTime *time.Duration
}

// DefaultRetryConfig returns the SDK default retry policy used for OTLP exporter.
func DefaultRetryConfig() RetryConfig {
	enabled := defaultRetryEnabled
	initial := defaultRetryInitial
	maxInterval := defaultRetryMaxInterval
	maxElapsed := defaultRetryMaxElapsed

	return RetryConfig{
		Enabled:         &enabled,
		InitialInterval: &initial,
		MaxInterval:     &maxInterval,
		MaxElapsedTime:  &maxElapsed,
	}
}

func (c RetryConfig) normalize() RetryConfig {
	cfg := c
	if cfg.Enabled == nil {
		enabled := defaultRetryEnabled
		cfg.Enabled = &enabled
	}
	if cfg.InitialInterval == nil {
		initial := defaultRetryInitial
		cfg.InitialInterval = &initial
	}
	if cfg.MaxInterval == nil {
		maxInterval := defaultRetryMaxInterval
		cfg.MaxInterval = &maxInterval
	}
	if cfg.MaxElapsedTime == nil {
		maxElapsed := defaultRetryMaxElapsed
		cfg.MaxElapsedTime = &maxElapsed
	}
	return cfg
}

// Config defines the runtime configuration used by the SDK.
type Config struct {
	// ServiceName 是当前服务的主标识。
	// 它会写入 OTel Resource 的 service.name，也是日志平台区分服务的核心字段之一。
	ServiceName string
	// ServiceNamespace 是服务所属的逻辑命名空间。
	// 它通常用于区分业务域、产品线或组织边界，不等同于 Kubernetes Namespace。
	ServiceNamespace string
	// ServiceVersion 表示当前进程运行的应用版本。
	// 推荐使用发布版本号、镜像 tag 或构建版本，便于发布比对和回溯问题。
	ServiceVersion string
	// ServiceInstanceID 表示单个运行实例的唯一标识。
	// 在 Kubernetes 中通常可使用 Pod 名，在虚拟机/物理机场景中可使用实例 ID 或进程级实例标识。
	ServiceInstanceID string
	// Environment 表示部署环境，例如 dev、test、prod。
	// 它只表达环境，不建议混入区域、租户、机房等其他语义。
	Environment string
	// Cluster 表示当前实例所在的集群标识。
	// 该字段通常由平台统一注入，便于跨集群检索日志和排查拓扑问题。
	Cluster string
	// Region 表示当前实例所在的大区或地域。
	// 它通常对应云区域或公司内部的大区概念。
	Region string
	// Zone 表示当前实例所在的可用区。
	// 当平台需要分析单可用区故障、限流或抖动问题时，该字段非常关键。
	Zone string
	// IDC 表示当前实例所在的机房、园区或基础设施单元。
	// 当 Region/Zone 颗粒度不足时，可用 IDC 补充更细粒度定位信息。
	IDC string
	// HostName 表示宿主机名称。
	// 在容器部署场景下通常是节点主机名，在物理机/虚拟机场景下就是当前机器名。
	HostName string
	// HostIP 表示宿主机 IP。
	// 该字段可用于排查节点网络、主机级别故障或基础设施路由问题。
	HostIP string
	// NodeName 表示 Kubernetes 节点名称。
	// 如果当前应用不运行在 Kubernetes 中，该字段可以为空。
	NodeName string
	// K8sNamespace 表示 Kubernetes Namespace。
	// 它描述的是运行时资源隔离空间，不等同于业务逻辑命名空间 ServiceNamespace。
	K8sNamespace string
	// PodName 表示当前实例所在 Pod 名称。
	// 它通常可直接作为容器化场景中的实例 ID 参考值。
	PodName string
	// PodIP 表示当前 Pod 的 IP 地址。
	// 它在排查 Service 发现、Sidecar 通信和 Pod 网络问题时很有帮助。
	PodIP string
	// ContainerName 表示当前业务容器名称。
	// 当一个 Pod 内存在多个容器时，该字段可以帮助区分具体日志来源。
	ContainerName string
	// Endpoint 是 OTLP 导出的目标地址。
	// 在生产环境中通常指向本机 log-agent，例如 localhost:4317。
	Endpoint string
	// Insecure 控制 OTLP/gRPC 连接是否关闭传输层安全。
	// 当日志发送目标是本机 agent 或可信内网链路时，通常可设置为 true。
	Insecure bool
	// Protocol 表示 SDK 与日志接收端通信使用的协议。
	// 当前仓库仅支持 grpc，对应 OTLP/gRPC 导出链路。
	Protocol string
	// Format 表示日志在本地输出时采用的格式。
	// 开发环境通常使用 console 便于阅读，结构化链路通常使用 json。
	Format string
	// Output 表示日志最终的输出目标。
	// 可选值包括 otlp、stdout、stderr、console，用于区分生产链路和本地调试链路。
	Output string
	// Level 表示日志最小输出级别。
	// 该字段会影响 bridge 层是否允许对应级别的日志进入导出链路。
	Level string
	// Development 表示是否按开发环境行为运行。
	// 为 true 时，SDK 会优先选择本地控制台输出和更友好的调试体验。
	Development bool
	// EnableCaller 控制是否收集并写出调用方文件名、行号和函数名。
	// 该能力有助于排查问题，但会带来一定性能开销。
	EnableCaller bool
	// EnableStacktrace 控制错误级别日志是否自动携带堆栈。
	// 对生产故障定位很有帮助，但会显著增加单条日志体积。
	EnableStacktrace bool
	// BatchTimeout 表示 BatchProcessor 的批量导出间隔。
	// 值越小延迟越低，值越大吞吐越高，但在异常时堆积时间也会更长。
	BatchTimeout time.Duration
	// ExportTimeout 表示单次批量导出的最长执行时间。
	// 超过该时间后，本次导出会被取消，并交由失败处理逻辑兜底。
	ExportTimeout time.Duration
	// MaxBatchSize 表示单次导出的最大日志条数。
	// 超过该数量后，BatchProcessor 会拆分为多批导出。
	MaxBatchSize int
	// MaxQueueSize 表示 BatchProcessor 内部队列的最大容量。
	// 当队列写满后，OTel 会丢弃最老的日志记录，因此该值需要根据峰值流量合理配置。
	MaxQueueSize int
	// FallbackFilePath 表示 OTLP 推送最终失败后的本地兜底文件路径。
	// 当前实现会把失败批次按 JSON 行追加到该文件，便于后续补偿与排查。
	FallbackFilePath string
	// Retry 定义 exporter 级别的瞬时错误重试策略。
	// 该重试只负责网络/接收端短暂异常场景；重试耗尽后会进入本地文件兜底。
	Retry RetryConfig
	// Headers 表示发送给 OTLP 接收端的额外请求头。
	// 常用于租户标识、鉴权令牌或平台侧路由标识透传。
	Headers map[string]string
	// ResourceAttributes 表示额外附加到 OTel Resource 的自定义属性。
	// 适合写入业务域扩展元数据，但不建议覆盖标准 service.* 或基础设施字段。
	ResourceAttributes map[string]string
}

// Default returns a baseline configuration with sensible defaults.
func Default() Config {
	return Config{
		Protocol:         DefaultProtocol,
		Format:           DefaultFormat,
		Output:           DefaultOutput,
		Level:            "info",
		Insecure:         true,
		BatchTimeout:     5 * time.Second,
		ExportTimeout:    3 * time.Second,
		MaxBatchSize:     512,
		MaxQueueSize:     2048,
		FallbackFilePath: DefaultFallbackFilePath,
		Retry:            DefaultRetryConfig(),
		EnableCaller:     true,
		EnableStacktrace: true,
	}
}

// LoadFromEnv returns a config populated from environment variables.
func LoadFromEnv() (Config, error) {
	cfg := Default()
	if err := cfg.ApplyEnv(); err != nil {
		return Config{}, err
	}
	return cfg.Normalize()
}

// IsDevelopment reports whether the config should run in development mode.
func (c Config) IsDevelopment() bool {
	if c.Development {
		return true
	}
	switch c.Environment {
	case "", "dev", "local", "development":
		return true
	default:
		return false
	}
}

// Normalize applies derived defaults and validates the configuration.
func (c Config) Normalize() (Config, error) {
	cfg := c
	cfg.applyGlobalMetadataDefaults()

	if cfg.Protocol == "" {
		cfg.Protocol = DefaultProtocol
	}
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 5 * time.Second
	}
	if cfg.ExportTimeout <= 0 {
		cfg.ExportTimeout = 3 * time.Second
	}
	if cfg.MaxBatchSize <= 0 {
		cfg.MaxBatchSize = 512
	}
	if cfg.MaxQueueSize <= 0 {
		cfg.MaxQueueSize = 2048
	}
	if cfg.FallbackFilePath == "" {
		cfg.FallbackFilePath = DefaultFallbackFilePath
	}
	cfg.Retry = cfg.Retry.normalize()
	if cfg.Output == "" {
		if cfg.IsDevelopment() {
			cfg.Output = OutputConsole
		} else {
			cfg.Output = OutputOTLP
		}
	}
	if cfg.Format == "" {
		if cfg.Output == OutputConsole {
			cfg.Format = FormatConsole
		} else {
			cfg.Format = FormatJSON
		}
	}
	if cfg.IsDevelopment() && cfg.Output == OutputOTLP && cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:4317"
	}
	if !cfg.IsDevelopment() && cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:4317"
	}
	if cfg.Headers == nil {
		cfg.Headers = map[string]string{}
	}
	if cfg.ResourceAttributes == nil {
		cfg.ResourceAttributes = map[string]string{}
	}

	return cfg, cfg.Validate()
}
