package pulsar_consume

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/zly-app/zapp/core"
)

type Consume struct {
	app     core.IApp
	consume pulsar.Consumer
	handle  func(message Message) error
}

func (c *Consume) Start() {
	panic("未实现")
}
func (c *Consume) Close() {
	panic("未实现")
}

func NewConsume(app core.IApp, client pulsar.Client, conf *Config, handle func(Message) error) (*Consume, error) {
	co := pulsar.ConsumerOptions{
		Topic:                          "",
		Topics:                         nil,
		TopicsPattern:                  "",
		AutoDiscoveryPeriod:            0,
		SubscriptionName:               "",
		Properties:                     nil,
		SubscriptionProperties:         nil,
		Type:                           0,
		SubscriptionInitialPosition:    0,
		DLQ:                            nil,
		KeySharedPolicy:                nil,
		RetryEnable:                    false,
		MessageChannel:                 nil,
		ReceiverQueueSize:              0,
		NackRedeliveryDelay:            0,
		Name:                           "",
		ReadCompacted:                  false,
		ReplicateSubscriptionState:     false,
		Interceptors:                   nil,
		Schema:                         nil,
		MaxReconnectToBroker:           nil,
		Decryption:                     nil,
		EnableDefaultNackBackoffPolicy: false,
		NackBackoffPolicy:              nil,
	}

	consume, err := client.Subscribe(co)
	if err != nil {
		return nil, fmt.Errorf("订阅失败: %v", err)
	}

	c := &Consume{
		app:     app,
		consume: consume,
		handle:  handle,
	}
	return c, nil
}
