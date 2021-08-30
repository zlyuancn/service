package memory

import (
	"container/list"
	"fmt"
	"sync"

	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/service/crawler/config"
	"github.com/zly-app/service/crawler/core"
	"github.com/zly-app/service/crawler/seed"
)

type MemoryQueue struct {
	queues map[string]*list.List
	mx     sync.Mutex
}

func (m *MemoryQueue) getQueue(queueName string) *list.List {
	queue, ok := m.queues[queueName]
	if !ok {
		queue = list.New()
		m.queues[queueName] = queue
	}
	return queue
}

func (m *MemoryQueue) Put(queueName string, seed core.ISeed, front bool) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	data, err := seed.Encode()
	if err != nil {
		return fmt.Errorf("seed编码失败: %v", err)
	}

	queue := m.getQueue(queueName)
	if front {
		queue.PushFront(data)
		return nil
	}

	queue.PushBack(data)
	return nil
}

func (m *MemoryQueue) Pop(queueName string, front bool) (core.ISeed, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	queue := m.getQueue(queueName)
	if queue.Len() == 0 {
		return nil, nil
	}

	var element *list.Element
	if front {
		element = queue.Front()
	} else {
		element = queue.Back()
	}

	raw := queue.Remove(element).(string)
	return seed.MakeSeedOfRaw(raw)
}

func (m *MemoryQueue) CheckQueueIsEmpty() (bool, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, suffix := range config.Conf.Frame.QueueSuffixes {
		if config.Conf.Frame.CheckEmptyQueueIgnoreErrorQueue {
			if suffix == config.Conf.Frame.ErrorSeedQueueSuffix || suffix == config.Conf.Frame.ParserErrorSeedQueueSuffix {
				continue
			}
		}
		queueName := config.Conf.Spider.Name + suffix
		if _, ok := m.queues[queueName]; ok {
			return false, nil
		}
	}
	return true, nil
}

func (m *MemoryQueue) Close() error {
	return nil
}

func NewMemoryQueue(app zapp_core.IApp) core.IQueue {
	return &MemoryQueue{
		queues: make(map[string]*list.List),
	}
}
