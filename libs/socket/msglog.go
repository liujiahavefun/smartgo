package socket

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

type MessageLogInfo struct {
	Dir string
	ses Session
	pkt *Packet

	meta *MessageMeta
}

func (self *MessageLogInfo) PeerName() string {
	return self.ses.FromPeer().Name()
}

func (self *MessageLogInfo) SessionID() int64 {
	return self.ses.GetID()
}

func (self *MessageLogInfo) MsgName() string {

	if self.meta == nil {
		return ""
	}

	return self.meta.Name
}

func (self *MessageLogInfo) MsgID() uint32 {
	return self.pkt.MsgID
}

func (self *MessageLogInfo) MsgSize() int {
	return len(self.pkt.Data)
}

func (self *MessageLogInfo) MsgString() string {
	if self.meta == nil {
		return fmt.Sprintf("%v", self.pkt.Data)
	}

	rawMsg, err := ParsePacket(self.pkt, self.meta.Type)
	if err != nil {
		return err.Error()
	}

	return rawMsg.(proto.Message).String()
}

// 是否启用消息日志
var EnableMessageLog bool = true

func msgLog(dir string, ses Session, pkt *Packet) {
	if !EnableMessageLog {
		return
	}

	if pkt == nil {
		//fmt.Println("pkt is nill")
		return
	}

	info := &MessageLogInfo{
		Dir:  dir,
		ses:  ses,
		pkt:  pkt,
		meta: MessageMetaByID(pkt.MsgID),
	}

	// 找到消息需要屏蔽
	if _, ok := msgMetaByID[info.MsgID()]; ok {
		return
	}

	if msgLogHook == nil || (msgLogHook != nil && msgLogHook(info)) {
		logDebugf("#%s(%s) sid: %d %s size: %d | %s", info.Dir, info.PeerName(), info.SessionID(), info.MsgName(), info.MsgSize(), info.MsgString())
	}
}

var msgLogHook func(*MessageLogInfo) bool

func HookMessageLog(hook func(*MessageLogInfo) bool) {
	msgLogHook = hook
}

var msgMetaByID = make(map[uint32]*MessageMeta)

func BlockMessageLog(msgName string) {
	meta := MessageMetaByName(msgName)

	if meta == nil {
		logErrorf("msg log block not found: %s", msgName)
		return
	}

	msgMetaByID[meta.ID] = meta
}
