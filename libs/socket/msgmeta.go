package socket

import (
	"path"
	"reflect"

	"github.com/golang/protobuf/proto"
)

//liujia: Type是proto.Message的实际类型，Name是类似testproto.TestEchoACK这样的可读名字，ID对应实际的消息号MsgID
type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   uint32
}

var (
	metaByName = map[string]*MessageMeta{}
	metaByID   = map[uint32]*MessageMeta{}
)

//注册消息元信息(代码生成专用)
func RegisterMessageMeta(name string, msg proto.Message, id uint32) {
	meta := &MessageMeta{
		Type: reflect.TypeOf(msg),
		Name: name,
		ID:   id,
	}

	metaByName[name] = meta
	metaByID[meta.ID] = meta
}

//根据名字查找消息元信息
func MessageMetaByName(name string) *MessageMeta {
	if v, ok := metaByName[name]; ok {
		return v
	}

	return nil
}

//消息全名，类似gamedef.TestEchoACK，rtype.PkgPath()返回完整package名字类似cellnet/proto/gamedef，path.Base()后变为gamedef
//所以这里要求注册时都要用gamedef.TestEchoACK这种形式？
func MessageFullName(rtype reflect.Type) string {
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	return path.Base(rtype.PkgPath()) + "." + rtype.Name()
}

//根据id查找消息元信息
func MessageMetaByID(id uint32) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

//根据id查找消息名, 没找到返回空
func MessageNameByID(id uint32) string {
	if meta := MessageMetaByID(id); meta != nil {
		return meta.Name
	}

	return ""
}

//遍历消息元信息
func VisitMessageMeta(callback func(*MessageMeta)) {
	for _, meta := range metaByName {
		callback(meta)
	}
}
