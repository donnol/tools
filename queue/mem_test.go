package queue

import (
	"log"
	"testing"
)

func TestMemQueue(t *testing.T) {
	que := NewQueue[*MemMessage[string]]()

	topic := que.AddTopic("test")
	t.Logf("topic: %+v", topic)

	finish := make(chan struct{})
	go func() {
		if err := que.Consume(topic, func(msg Message) error {
			t.Logf("msg: %+v", msg)
			finish <- struct{}{}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	receipt, err := que.Produce(topic, &MemMessage[string]{
		Id: "1",
		Data: map[string]string{
			"key": "value",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("receipt: %v", receipt.String())

	<-finish
}
