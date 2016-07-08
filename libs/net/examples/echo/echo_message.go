package echo

import (
//"github.com/golang/glog"
//"smartgo/libs/net"
)

type EchoMessage struct {
	EchoString string
}

/*
func (em EchoMessage) Serialize() ([]byte, error) {
	return []byte(em.Message), nil
}
*/

func (em EchoMessage) Major() uint16 {
	return 1
}

func (em EchoMessage) Minor() uint16 {
	return 1
}

func (em EchoMessage) MarshalBinary() (data []byte, err error) {
	msg := []byte(em.EchoString)
	l := len(msg)
	data = make([]byte, 4+l)
	data = data[0:0]
	data = append(data, byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	data = append(data, msg...)
	return data, nil
}

/*
func DeserializeEchoMessage(data []byte) (message net.Message, err error) {
	if data == nil {
		return nil, net.ErrorNilData
	}
	msg := string(data)
	echo := EchoMessage{
		Message: msg,
	}

	return echo, nil
}
*/

/*
type EchoMessageHandler struct {
	netid   int64
	message net.Message
}

func NewEchoMessageHandler(net int64, msg net.Message) net.MessageHandler {
	return EchoMessageHandler{
		netid:   net,
		message: msg,
	}
}

func (handler EchoMessageHandler) Process(client net.Connection) bool {
	echoMessage := handler.message.(EchoMessage)
	glog.Infof("Receving message %s\n", echoMessage.Message)
	client.Write(handler.message)
	return true
}
*/
