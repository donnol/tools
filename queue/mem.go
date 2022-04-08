package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type memQueue[T Message] struct {
	topics []Topic

	indexMap map[string]int
}

func NewQueue[T Message]() Queue {
	return &memQueue[T]{
		topics:   make([]Topic, 0),
		indexMap: make(map[string]int),
	}
}

func (que *memQueue[T]) AddTopic(name string) Topic {
	if i, ok := que.indexMap[name]; ok {
		return que.topics[i]
	}

	l := len(que.topics)

	topic := &memTopic[T]{
		name:        name,
		list:        make([]T, 0),
		unreadIndex: -1,
		newIndex:    -1,
	}
	que.topics = append(que.topics, topic)
	que.indexMap[name] = l

	return topic
}

func (que *memQueue[T]) TopicByName(name string) (Topic, bool) {
	if i, ok := que.indexMap[name]; ok {
		return que.topics[i], ok
	}
	return &memTopic[T]{}, false
}

func (que *memQueue[T]) Topics() []Topic {
	return que.topics
}

// Produce 逐条生产
func (que *memQueue[T]) Produce(topic Topic, msg Message) (Receipt, error) {
	var ok bool
	name := topic.Name()
	topic, ok = que.TopicByName(name)
	if !ok {
		return &MemReceipt{}, fmt.Errorf("not found topic: %s", name)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("generate id failed: %w", err)
	}

	top, ok := topic.(*memTopic[T])
	if !ok {
		return nil, fmt.Errorf("assert topic to memTopic[T] failed")
	}
	{
		top.mutex.Lock()
		l := len(top.list)
		top.list = append(top.list, msg.(T))
		top.newIndex = l
		if top.unreadIndex == -1 {
			top.unreadIndex = 0
		}
		top.mutex.Unlock()
	}

	return &MemReceipt{Id: id.String()}, nil
}

// Consume 循环消费，会阻塞住
func (que *memQueue[T]) Consume(topic Topic, f func(msg Message) error) error {
	var ok bool
	name := topic.Name()
	topic, ok = que.TopicByName(name)
	if !ok {
		return fmt.Errorf("not found topic: %s", name)
	}

	top, ok := topic.(*memTopic[T])
	if !ok {
		return fmt.Errorf("assert topic to memTopic[T] failed")
	}

	for {
		var msg Message
		{
			top.mutex.RLock()

			if top.unreadIndex == -1 || top.unreadIndex > top.newIndex {
				top.mutex.RUnlock()
				time.Sleep(200 * time.Millisecond)
				continue
			}

			msg = top.list[top.unreadIndex]
			if top.unreadIndex <= top.newIndex {
				top.unreadIndex++
			}

			top.mutex.RUnlock()
		}

		if err := f(msg); err != nil {
			return fmt.Errorf("exec f failed: %w", err)
		}
	}
}

type memTopic[T Message] struct {
	name string

	mutex       sync.RWMutex
	list        []T
	unreadIndex int // 未读信息下标
	newIndex    int // 最新信息下标
}

func (topic *memTopic[T]) Name() string {
	return topic.name
}

type Element interface {
	~int | ~string | ~struct{ Name string }
}

type MemMessage[T Element] struct {
	Id   string       `json:"id"`
	Data map[string]T `json:"data"`
}

func (msg *MemMessage[T]) Marshal() ([]byte, error) {
	return json.Marshal(*msg)
}

func (msg *MemMessage[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, msg)
}

type MemReceipt struct {
	Id string `json:"id"`
}

func (receipt *MemReceipt) String() string {
	return "receipt: " + receipt.Id
}
