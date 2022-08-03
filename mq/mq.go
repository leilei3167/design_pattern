package mq

/*
消息队列抽象为mq接口,以适应不同的消息队列;
需要符合接口隔离原则,一个模块不应该提供用户不需要的接口,接口尽可能的小
但是过度地细化和拆分接口，也会导致系统的接口数量的上涨，从而产生更大的维护成本。接口的粒度需要根据具体的业务场景来定，可以参考单一职责原则，将那些为同一类客户端程序提供服务的接口合并在一起。
 * 例子：
 * 根据消息队列的模型，拆分成Consumable, Producible两个接口，由Mq继承它们
 * 生产者依赖Producible，消费者依赖Consumable，符合ISP

其实从领域的角度看 Mq直接具备两个方法是合理的,但是从使用方来看,对于MemoryMqInput来说他只需要生产消息,而不需要再额外实现
消费消息的行为
*/

type Topic string
type Message struct { //代表一个消息
	topic   Topic
	payload string
}

// Consumable 消费接口,从消息队列中消费数据
type Consumable interface {
	Consume(topic Topic) (*Message, error)
}

// Producible 生产接口
type Producible interface {
	Produce(message *Message) error
}

// Mq 代表消息队列的抽象,具备生产和消费的行为
type Mq interface {
	Consumable
	Producible
}

// NewMessage 构造方法
func NewMessage(topic Topic, payload string) *Message {
	return &Message{
		topic:   topic,
		payload: payload,
	}
}
func (m Message) Topic() Topic {
	return m.topic
}

func (m Message) Payload() string {
	return m.payload
}
