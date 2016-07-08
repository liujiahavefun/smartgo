package net

import (
//"bytes"
//"encoding"
//"encoding/binary"
//"io"
//"log"
)

/*
const (
	NTYPE  = 4
	NLEN   = 4
	MAXLEN = 1 << 23 // 8M
)

func init() {
	MessageMap = make(MessageMapType)
	HandlerMap = make(HandlerMapType)
	buf = new(bytes.Buffer)
	//messageCodec = TypeLengthValueCodec{}
}

var (
	MessageMap MessageMapType
	HandlerMap HandlerMapType
	buf        *bytes.Buffer
	//messageCodec Codec
)

//接受一个Message，返回此消息的处理函数
type NewHandlerFunctionType func(Message) MessageHandler

//接受字节流，返回解包后的消息体
type UnmarshalFunctionType func([]byte) (Message, error)

//liujia: Message实际上就是encoding.BinaryMarshaler的子类，并且有一个消息id
type Message interface {
	MessageNumber() int32
	encoding.BinaryMarshaler
}

//消息处理函数就是一个接口，参数是TCPConnection，返回处理成功还是失败
type MessageHandler interface {
	Process(client Connection) bool
}

//map，映射消息号到消息解包函数
type MessageMapType map[int32]UnmarshalFunctionType

//liujia: 注册和获取对message解包的函数
func (mm *MessageMapType) Register(msgType int32, unmarshaler func([]byte) (Message, error)) {
	(*mm)[msgType] = UnmarshalFunctionType(unmarshaler)
}

func (mm *MessageMapType) get(msgType int32) UnmarshalFunctionType {
	if unmarshaler, ok := (*mm)[msgType]; ok {
		return unmarshaler
	}
	return nil
}
*/

//liujia: 注册和获取对message的处理函数，处理时，根据消息号，获取factory，然后得到真正的处理函数，handler = factory(msg)
//如下面的DefaultHeartBeatMessageHandler，将消息和其处理函数包装其它，形成这个factory
//将这个factory注册给HandlerMap，处理里获取factory，然后通过factory(msg)获取绑定的真正handler
//这样将真实的消息，和通用的handler，通过闭包绑定在一起了
//绑定的目的，是要统一发往handleLoop去处理
/*
type HandlerMapType map[int32]NewHandlerFunctionType
func (hm *HandlerMapType) Register(msgType int32, factory func(Message) MessageHandler) {
	(*hm)[msgType] = NewHandlerFunctionType(factory)
}

func (hm *HandlerMapType) get(msgType int32) NewHandlerFunctionType {
	if fn, ok := (*hm)[msgType]; ok {
		return fn
	}
	return nil
}
*/

/* Message number 0 is the preserved message
for long-term connection keeping alive */
/*
type DefaultHeartBeatMessage struct {
	Timestamp int64
}

//liujia:buf是全局变量，MarshalBinary()和UnmarshalDefaultHeartBeatMessage() 都用了，是不是有问题？
func (dhbm DefaultHeartBeatMessage) MarshalBinary() ([]byte, error) {
	buf.Reset()
	err := binary.Write(buf, binary.BigEndian, dhbm.Timestamp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (dhbm DefaultHeartBeatMessage) MessageNumber() int32 {
	return 0
}

func UnmarshalDefaultHeartBeatMessage(data []byte) (message Message, err error) {
	var timestamp int64
	if data == nil {
		return nil, ErrorNilData
	}

	//liujia: 这个buf是上面的全局变量么？
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.BigEndian, &timestamp)
	if err != nil {
		return nil, err
	}
	return DefaultHeartBeatMessage{
		Timestamp: timestamp,
	}, nil
}

type DefaultHeartBeatMessageHandler struct {
	message Message
}

func NewDefaultHeartBeatMessageHandler(msg Message) MessageHandler {
	return DefaultHeartBeatMessageHandler{
		message: msg,
	}
}

func (handler DefaultHeartBeatMessageHandler) Process(client *TCPConnection) bool {
	heartBeatMessage := handler.message.(DefaultHeartBeatMessage)
	log.Printf("Receiving heart beat at %d, updating\n", heartBeatMessage.Timestamp)
	client.HeartBeat = heartBeatMessage.Timestamp
	return true
}

func SetMessageCodec(codec Codec) {
	messageCodec = codec
}
*/

/* Application programmer can define a custom codec themselves */
/*
type Codec interface {
	Decode(*TCPConnection) (Message, error)
	Encode(Message) ([]byte, error)
}

// use type-length-value format: |4 bytes|4 bytes|n bytes <= 8M|
type TypeLengthValueCodec struct{}

func (codec TypeLengthValueCodec) Decode(c *TCPConnection) (Message, error) {
	typeBytes := make([]byte, NTYPE)
	lengthBytes := make([]byte, NLEN)

	_, err := io.ReadFull(c.RawConn(), typeBytes)
	if err != nil {
		return nil, err
	}
	typeBuf := bytes.NewReader(typeBytes)
	var msgType int32
	if err = binary.Read(typeBuf, binary.BigEndian, &msgType); err != nil {
		return nil, err
	}

	_, err = io.ReadFull(c.RawConn(), lengthBytes)
	if err != nil {
		return nil, err
	}
	lengthBuf := bytes.NewReader(lengthBytes)
	var msgLen uint32
	if err = binary.Read(lengthBuf, binary.BigEndian, &msgLen); err != nil {
		return nil, err
	}
	if msgLen > MAXLEN {
		return nil, ErrorIllegalData
	}

	// read real application message
	msgBytes := make([]byte, msgLen)
	_, err = io.ReadFull(c.RawConn(), msgBytes)
	if err != nil {
		return nil, err
	}

	// deserialize message from bytes
	unmarshaler := MessageMap.get(msgType)
	if unmarshaler == nil {
		return nil, ErrorUndefind
	}
	return unmarshaler(msgBytes)
}

func (codec TypeLengthValueCodec) Encode(msg Message) ([]byte, error) {
	data, err := msg.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, msg.MessageNumber())
	binary.Write(buf, binary.BigEndian, int32(len(data)))
	binary.Write(buf, binary.BigEndian, data)
	packet := buf.Bytes()
	return packet, nil
}
*/
