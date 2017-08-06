package socket

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

//封包
type Packet struct {
	MsgID uint32 // 消息ID
	Data  []byte
}

func (self Packet) ContextID() uint32 {
	return self.MsgID
}

//消息到封包
//注意这里要求data一定是proto.Message类型的，参考RegisterMessageMeta()，那里也要求注册的消息都是proto.Message
//注意MessageMeta里的ID就是消息的MsgId
func BuildPacket(data interface{}) (*Packet, *MessageMeta) {
	msg := data.(proto.Message)
	rawdata, err := proto.Marshal(msg)
	if err != nil {
		logErrorln(err)
	}

	meta := MessageMetaByName(MessageFullName(reflect.TypeOf(msg)))
	if meta == nil {
		fmt.Println("invalid or unregistered msg:", reflect.TypeOf(data))
	}

	return &Packet{
		MsgID: uint32(meta.ID),
		Data:  rawdata,
	}, meta
}

//封包到消息
//根据msgType，去掉ptr，New出对应的reflect.Value的"零值"的pointer(类型是reflect.Value)，Interface()返回reflect.Value对应的interface
//说了一堆，就是根据msgType，返回一个空的具体的proto.Message对象，例如传入的msgType是gamedef.TestEchoACK，
//则rawMsg就是一个interface，底层类型就是对应的TestEchoACK的proto.TestEchoACK
func ParsePacket(pkt *Packet, msgType reflect.Type) (interface{}, error) {
	//msgType为ptr类型, new时需要非ptr型
	rawMsg := reflect.New(msgType.Elem()).Interface()

	err := proto.Unmarshal(pkt.Data, rawMsg.(proto.Message))
	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}
