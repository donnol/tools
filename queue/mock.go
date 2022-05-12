package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/donnol/tools/inject"
)

var (
	_gen_customCtxMap = make(map[string]inject.CtxFunc)
)

func RegisterProxyMethod(pctx inject.ProxyContext, cf inject.CtxFunc) {
	_gen_customCtxMap[pctx.Uniq()] = cf
}

type ConsumerMock struct {
	ConsumeFunc func(topic Topic, f func(msg Message) error) error
}

var (
	_ Consumer = &ConsumerMock{}

	consumerMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/queue",
		InterfaceName: "Consumer",
	}
	ConsumerMockConsumeProxyContext = func() (pctx inject.ProxyContext) {
		pctx = consumerMockCommonProxyContext
		pctx.MethodName = "Consume"
		return
	}()

	_ = getConsumerProxy
)

func getConsumerProxy(base Consumer) *ConsumerMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ConsumerMock{
		ConsumeFunc: func(topic Topic, f func(msg Message) error) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := ConsumerMockConsumeProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, topic)

				_gen_params = append(_gen_params, f)

				_gen_res := _gen_cf(_gen_ctx, base.Consume, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Consume(topic, f)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *ConsumerMock) Consume(topic Topic, f func(msg Message) error) error {
	return mockRecv.ConsumeFunc(topic, f)
}

type MessageMock struct {
	MarshalFunc func() ([]byte, error)

	UnmarshalFunc func([]byte) error
}

var (
	_ Message = &MessageMock{}

	messageMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/queue",
		InterfaceName: "Message",
	}
	MessageMockMarshalProxyContext = func() (pctx inject.ProxyContext) {
		pctx = messageMockCommonProxyContext
		pctx.MethodName = "Marshal"
		return
	}()
	MessageMockUnmarshalProxyContext = func() (pctx inject.ProxyContext) {
		pctx = messageMockCommonProxyContext
		pctx.MethodName = "Unmarshal"
		return
	}()

	_ = getMessageProxy
)

func getMessageProxy(base Message) *MessageMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &MessageMock{
		MarshalFunc: func() ([]byte, error) {
			_gen_begin := time.Now()

			var _gen_r0 []byte

			var _gen_r1 error

			_gen_ctx := MessageMockMarshalProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Marshal, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].([]byte)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Marshal()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},

		UnmarshalFunc: func(p0 []byte) error {
			_gen_begin := time.Now()

			var _gen_r0 error

			_gen_ctx := MessageMockUnmarshalProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, p0)

				_gen_res := _gen_cf(_gen_ctx, base.Unmarshal, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(error)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Unmarshal(p0)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *MessageMock) Marshal() ([]byte, error) {
	return mockRecv.MarshalFunc()
}

func (mockRecv *MessageMock) Unmarshal(p0 []byte) error {
	return mockRecv.UnmarshalFunc(p0)
}

type ProducerMock struct {
	ProduceFunc func(topic Topic, msg Message) (Receipt, error)
}

var (
	_ Producer = &ProducerMock{}

	producerMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/queue",
		InterfaceName: "Producer",
	}
	ProducerMockProduceProxyContext = func() (pctx inject.ProxyContext) {
		pctx = producerMockCommonProxyContext
		pctx.MethodName = "Produce"
		return
	}()

	_ = getProducerProxy
)

func getProducerProxy(base Producer) *ProducerMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ProducerMock{
		ProduceFunc: func(topic Topic, msg Message) (Receipt, error) {
			_gen_begin := time.Now()

			var _gen_r0 Receipt

			var _gen_r1 error

			_gen_ctx := ProducerMockProduceProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_params = append(_gen_params, topic)

				_gen_params = append(_gen_params, msg)

				_gen_res := _gen_cf(_gen_ctx, base.Produce, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(Receipt)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

				_gen_tmpr1, _gen_exist := _gen_res[1].(error)
				if _gen_exist {
					_gen_r1 = _gen_tmpr1
				}

			} else {
				_gen_r0, _gen_r1 = base.Produce(topic, msg)
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0, _gen_r1
		},
	}
}

func (mockRecv *ProducerMock) Produce(topic Topic, msg Message) (Receipt, error) {
	return mockRecv.ProduceFunc(topic, msg)
}

type ReceiptMock struct {
	StringFunc func() string
}

var (
	_ Receipt = &ReceiptMock{}

	receiptMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/queue",
		InterfaceName: "Receipt",
	}
	ReceiptMockStringProxyContext = func() (pctx inject.ProxyContext) {
		pctx = receiptMockCommonProxyContext
		pctx.MethodName = "String"
		return
	}()

	_ = getReceiptProxy
)

func getReceiptProxy(base Receipt) *ReceiptMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &ReceiptMock{
		StringFunc: func() string {
			_gen_begin := time.Now()

			var _gen_r0 string

			_gen_ctx := ReceiptMockStringProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.String, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.String()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *ReceiptMock) String() string {
	return mockRecv.StringFunc()
}

type TopicMock struct {
	NameFunc func() string
}

var (
	_ Topic = &TopicMock{}

	topicMockCommonProxyContext = inject.ProxyContext{
		PkgPath:       "github.com/donnol/tools/queue",
		InterfaceName: "Topic",
	}
	TopicMockNameProxyContext = func() (pctx inject.ProxyContext) {
		pctx = topicMockCommonProxyContext
		pctx.MethodName = "Name"
		return
	}()

	_ = getTopicProxy
)

func getTopicProxy(base Topic) *TopicMock {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &TopicMock{
		NameFunc: func() string {
			_gen_begin := time.Now()

			var _gen_r0 string

			_gen_ctx := TopicMockNameProxyContext
			_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
			if _gen_ok {
				_gen_params := []any{}

				_gen_res := _gen_cf(_gen_ctx, base.Name, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(string)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.Name()
			}

			log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

			return _gen_r0
		},
	}
}

func (mockRecv *TopicMock) Name() string {
	return mockRecv.NameFunc()
}
