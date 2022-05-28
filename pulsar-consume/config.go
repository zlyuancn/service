package pulsar_consume

const (
	// pulsar地址
	defUrl = "pulsar://localhost:6650"
	// 连接超时
	defConnectionTimeout = 5000
	// 操作超时
	defOperationTimeout = 30000

	// 消费者数量
	defConsumeCount = 1
	// 每个消费者协程数
	defConsumeThreadCount = 1
)

type Config struct {
	Url               string // pulsar地址, 示例: pulsar://localhost:6600,localhost:6650
	ListenerName      string // pulsar使用的监听器名, 示例: external
	ConnectionTimeout int    // 连接超时, 单位毫秒
	OperationTimeout  int    // 操作超时, 单位毫秒

	ConsumeCount       int // 消费者数量
	ConsumeThreadCount int // 每个消费者协程数
}

func NewConfig() *Config {
	return &Config{}
}

func (conf *Config) Check() error {
	if conf.Url == "" {
		conf.Url = defUrl
	}
	if conf.ConnectionTimeout < 1 {
		conf.ConnectionTimeout = defConnectionTimeout
	}
	if conf.OperationTimeout < 1 {
		conf.OperationTimeout = defOperationTimeout
	}

	if conf.ConsumeCount < 1 {
		conf.ConsumeCount = defConsumeCount
	}
	if conf.ConsumeThreadCount < 1 {
		conf.ConsumeThreadCount = defConsumeThreadCount
	}
	return nil
}
