package kafka_consume

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"
)

type Context struct {
	core.ILogger
	*sarama.ConsumerMessage
	GroupID string // 消费者组ID
}

type RegistryKafkaConsumerHandlerFunc = func(ctx *Context) error

type ConsumerConfig struct {
	Topics  []string
	GroupID string
	Handler RegistryKafkaConsumerHandlerFunc
	Opts    []ConsumerOption
	*ServiceConfig
}

type consumerCli struct {
	app       core.IApp
	conf      *ConsumerConfig
	consumers []sarama.ConsumerGroup
	*consumerOptions
	done chan struct{}
}

func newConsumer(app core.IApp, conf *ConsumerConfig) *consumerCli {
	c := &consumerCli{
		app:             app,
		conf:            conf,
		consumerOptions: newConsumerOptions(),
		done:            make(chan struct{}),
	}

	for _, o := range conf.Opts {
		o(c.consumerOptions)
	}

	if c.Disable {
		return c
	}

	count := c.ConsumeCount
	if count == 0 {
		count = c.conf.ConsumeCount
	}

	c.consumers = make([]sarama.ConsumerGroup, count)
	for i := 0; i < count; i++ {
		consumer, err := c.makeConsumer()
		if err != nil {
			app.Fatal("生成消费者失败", zap.Error(err))
		}
		c.consumers[i] = consumer
	}

	return c
}

func (c *consumerCli) Start() error {
	if c.Disable {
		return nil
	}

	for _, consume := range c.consumers {
		if c.conf.EnabledErrorsChannel {
			go c.errorsChannelCallback(consume)
		}
		go c.start(consume)
	}

	return nil
}

func (c *consumerCli) makeConsumer() (sarama.ConsumerGroup, error) {
	kConf := sarama.NewConfig()

	kConf.Net.DialTimeout = time.Duration(c.conf.DialTimeout) * time.Millisecond
	kConf.Net.ReadTimeout = time.Duration(c.conf.ReadTimeout) * time.Millisecond
	kConf.Net.WriteTimeout = time.Duration(c.conf.WriteTimeout) * time.Millisecond

	switch strings.ToLower(c.conf.PartitionBalanceStrategy) {
	case "sticky":
		kConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "round_robin":
		kConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	default:
		kConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}
	kConf.Consumer.Retry.Backoff = time.Duration(c.conf.RetryInterval) * time.Millisecond
	kConf.Consumer.Fetch.Max = c.conf.MaxMessageBytes
	kConf.Consumer.MaxProcessingTime = time.Duration(c.conf.MaxProcessingTime)
	kConf.Consumer.Return.Errors = c.conf.EnabledErrorsChannel // 如果启用了该选项，未交付的消息将在Errors通道上返回，包括error
	switch strings.ToLower(c.conf.OffsetInitial) {
	case "oldest":
		kConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		kConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	switch c.conf.IsolationLevel {
	case "ReadCommitted":
		kConf.Consumer.IsolationLevel = sarama.ReadCommitted
	default:
		kConf.Consumer.IsolationLevel = sarama.ReadUncommitted
	}

	kConf.ChannelBufferSize = c.conf.ChannelBufferSize
	if strings.HasPrefix(c.conf.KafkaVersion, "v") {
		c.conf.KafkaVersion = c.conf.KafkaVersion[1:]
	}
	var err error
	kConf.Version, err = sarama.ParseKafkaVersion(c.conf.KafkaVersion)
	if err != nil {
		return nil, fmt.Errorf("无法解析版本号: %v", err)
	}
	if strings.HasPrefix(c.conf.KafkaVersion, "v") {
		c.conf.KafkaVersion = c.conf.KafkaVersion[1:]
	}

	consumer, err := sarama.NewConsumerGroup(strings.Split(c.conf.Address, ","), c.conf.GroupID, kConf)
	if err != nil {
		return nil, fmt.Errorf("创建kafka消费者失败, topics:%s, group_id:%s, err:%s", strings.Join(c.conf.Topics, ","), c.conf.GroupID, err)
	}
	return consumer, nil
}

// 开始消费
func (c *consumerCli) start(consume sarama.ConsumerGroup) {
	for {
		select {
		case <-c.app.BaseContext().Done():
			return
		case <-c.done:
			return
		default:
		}

		err := consume.Consume(c.app.BaseContext(), c.conf.Topics, c)
		if err != nil {
			c.app.Error("kafka消费者运行时错误", zap.Error(err))
		}
	}
}

// 错误通道回调
func (c *consumerCli) errorsChannelCallback(consume sarama.ConsumerGroup) {
	for {
		select {
		case <-c.app.BaseContext().Done():
			return
		case <-c.done:
			return
		case err := <-consume.Errors():
			c.app.Error("kafka消费者收到错误消息", zap.Error(err))
			if c.ErrorsChannelCallback != nil {
				err = utils.Recover.WrapCall(func() error {
					c.ErrorsChannelCallback(err)
					return nil
				})
				if err != nil {
					panicErrDetail := utils.Recover.GetRecoverErrorDetail(err)
					panicErrInfos := strings.Split(panicErrDetail, "\n")
					c.app.Error("kafka消费者收到错误消息, 错误回调时产生panic",
						zap.String("error", panicErrInfos[0]), zap.Strings("detail", panicErrInfos[1:]))
				}
			}
		}
	}
}

func (c *consumerCli) Close() error {
	close(c.done)

	if c.Disable {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(c.consumers))
	for _, consumer := range c.consumers {
		go func(consumer sarama.ConsumerGroup) {
			_ = consumer.Close()
			wg.Done()
		}(consumer)
	}
	wg.Wait()
	return nil
}

func (c *consumerCli) Setup(session sarama.ConsumerGroupSession) error   { return nil }
func (c *consumerCli) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
func (c *consumerCli) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	for msg := range claim.Messages() {
		if err = c.process(session, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c *consumerCli) process(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	ctx := &Context{
		ILogger: c.app.NewSessionLogger(
			zap.String("kafka_topic", msg.Topic),
			zap.String("kafka_group_id", c.conf.GroupID),
			zap.Int32("kafka_partition", msg.Partition),
			zap.Int64("kafka_offset", msg.Offset),
		),
		ConsumerMessage: msg,
		GroupID:         c.conf.GroupID,
	}

	ctx.Debug("kafkaConsumer.receive")
	err := utils.Recover.WrapCall(func() error {
		return c.conf.Handler(ctx)
	})
	if err == nil {
		sess.MarkMessage(msg, "")
		ctx.Debug("kafkaConsumer.success")
		return nil
	}

	errDetail := utils.Recover.GetRecoverErrorDetail(err)
	ctx.Error("kafkaConsumer.error!", zap.String("error", errDetail))
	return err
}
