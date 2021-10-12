package kafka_consume

import (
	"errors"

	"github.com/Shopify/sarama"
)

const (
	// 默认kafka版本
	defaultKafkaVersion = "2.0.0"
	// 默认读取超时
	defaultReadTimeout = 10000
	// 默认写入超时
	defaultWriteTimeout = 10000
	// 默认连接超时
	defaultDialTimeout = 2000
	// 默认分区平衡策略
	defaultPartitionBalanceStrategy = "range"
	// 默认从分区读取失败后重试间隔时间
	defaultRetryInterval = 2000
	// 默认消息的最大允许大小
	defaultMaxMessageBytes = 1048576
	// 默认消息最大处理时间(毫秒)
	defaultMaxProcessingTime = 100
	// 默认偏移初始化
	defaultOffsetInitial = "oldest"
	// 默认隔离级别
	defaultIsolationLevel = "ReadUncommitted"
	// 默认消费者数量
	defaultConsumeCount = 1
	// 默认通道缓冲数
	defaultChannelBufferSize = 256
	// 默认消费失败后重新消费等待时间
	defaultReConsumeWaitTime = 1000
)

type ServiceConfig struct {
	Address                  string // 地址, 多个地址用半角逗号连接
	KafkaVersion             string // kafka版本, 示例: v2.0.0, 2.0.0
	ReadTimeout              int64  // 超时(毫秒)
	WriteTimeout             int64  // 超时(毫秒)
	DialTimeout              int64  // 超时(毫秒)
	PartitionBalanceStrategy string // 分区平衡策略, range, sticky, round_robin
	RetryInterval            int    // 从分区读取失败后重试间隔时间(毫秒)
	MaxMessageBytes          int32  // 消息的最大允许大小(字节)
	MaxProcessingTime        int    // 消息最大处理时间(毫秒)
	OffsetInitial            string // 找不到消费者组偏移记录时进行偏移初始化, newest 表示新的消费者不会消费以前的数据, oldest 表示新的消费者会从能消费的旧数据开始消费
	IsolationLevel           string // 隔离级别, ReadUncommitted, ReadCommitted
	ConsumeCount             int    // 消费者数量, 会为消费者组创建多个消费者进行消费, 建议设置为topic的分区数
	ChannelBufferSize        int    // 通道缓冲数, 要在内部和外部通道中缓冲的事件数量
	ReConsumeWaitTime        int64  // 消费失败后重新消费等待时间(毫秒)
	kConf                    *sarama.Config
}

func newConfig() *ServiceConfig {
	return &ServiceConfig{
		KafkaVersion:   defaultKafkaVersion,
		OffsetInitial:  defaultOffsetInitial,
		IsolationLevel: defaultIsolationLevel,
	}
}

func (conf *ServiceConfig) Check() error {
	if conf.Address == "" {
		return errors.New("address为空")
	}
	if conf.ReadTimeout <= 0 {
		conf.ReadTimeout = defaultReadTimeout
	}
	if conf.WriteTimeout <= 0 {
		conf.WriteTimeout = defaultWriteTimeout
	}
	if conf.DialTimeout <= 0 {
		conf.DialTimeout = defaultDialTimeout
	}
	if conf.PartitionBalanceStrategy == "" {
		conf.PartitionBalanceStrategy = defaultPartitionBalanceStrategy
	}
	if conf.RetryInterval <= 0 {
		conf.RetryInterval = defaultRetryInterval
	}
	if conf.MaxMessageBytes <= 0 {
		conf.MaxMessageBytes = defaultMaxMessageBytes
	}
	if conf.MaxProcessingTime <= 0 {
		conf.MaxProcessingTime = defaultMaxProcessingTime
	}
	if conf.OffsetInitial == "" {
		conf.OffsetInitial = defaultOffsetInitial
	}
	if conf.IsolationLevel == "" {
		conf.IsolationLevel = defaultIsolationLevel
	}
	if conf.ConsumeCount <= 0 {
		conf.ConsumeCount = defaultConsumeCount
	}
	if conf.ChannelBufferSize <= 0 {
		conf.ChannelBufferSize = defaultChannelBufferSize
	}
	if conf.ReConsumeWaitTime < 1 {
		conf.ReConsumeWaitTime = defaultReConsumeWaitTime
	}
	return nil
}
