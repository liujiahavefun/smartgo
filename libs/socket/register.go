package socket

type RegisterMessageContext struct {
	*MessageMeta
	*CallbackContext
}

func MessageRegisteredCount(evd EventDispatcher, msgName string) int {
	msgMeta := MessageMetaByName(msgName)
	if msgMeta == nil {
		return 0
	}

	return evd.CountByID(msgMeta.ID)
}

/*
 * 注册消息处理逻辑
 */
func RegisterMessage(evd EventDispatcher, msgName string, userHandler func(interface{}, Session)) *RegisterMessageContext {
	msgMeta := MessageMetaByName(msgName)
	if msgMeta == nil {
		logErrorf("message register failed, %s", msgName)
		return nil
	}

	ctx := evd.AddCallback(msgMeta.ID, func(data interface{}) {
		if ev, ok := data.(*SessionEvent); ok {
			rawMsg, err := ParsePacket(ev.Packet, msgMeta.Type)
			if err != nil {
				logErrorf("unmarshaling error: %v, raw: %v", err, ev.Packet)
				return
			}

			userHandler(rawMsg, ev.Ses)
		}
	})

	return &RegisterMessageContext{MessageMeta: msgMeta, CallbackContext: ctx}
}

/*
 * 注册默认消息处理逻辑，即如果此消息未被RegisterMessage所注册，都用此回调处理
 */
func RegisterDefault(evd EventDispatcher, defaultHandler func(uint32, []byte, Session)) *RegisterMessageContext {
	ctx := evd.AddDefaultCallback(func(data interface{}) {
		if ev, ok := data.(*SessionEvent); ok {
			defaultHandler(ev.Packet.MsgID, ev.Packet.Data, ev.Ses)
		}
	})

	return &RegisterMessageContext{MessageMeta: nil, CallbackContext: ctx}
}
