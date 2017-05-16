package socket

/*
TODO: 要实现以下几种回调
//start()时，无论是server启动listen还是client启动connect，都去给个回调
type onConnectFunc func(Connection) bool

//一个connection关闭时给回调，无论是server还是client
type onCloseFunc func(Connection)
//出错时给回调，时机？我倾向于server accept错误，client收发包的错误
type onErrorFunc func()

//下面这两个可以合二为一，一个是收到raw data包(byte[])，另一个是解成具体的message对象(做完反序列化之后)
type onPacketRecvFunc func(Connection, *pool.Buffer) (HandlerProc, bool)
type HandlerProc func()

type onMessageFunc func(Message, Connection)

//定时器回调
type onScheduleFunc func(time.Time, interface{})
*/

//某个peer连接过来并建立了session
type OnSessionConnectedFunc func(session Session)

//某个peer连接过来并建立了session
type OnSessionClosedFunc func(session Session)

//出现错误时回调，无论是server还是client
type OnSessionErrorFunc func(session Session, err error)

//收到packet
type OnSessionRecvPacketFunc func(session Session, packet *Packet)

type SessionCallback struct {
	OnConnected 	OnSessionConnectedFunc
	OnClosed    	OnSessionClosedFunc
	OnError     	OnSessionErrorFunc
	OnRecvPacket 	OnSessionRecvPacketFunc
}

func NewSessionCallback(onClosed OnSessionClosedFunc,
						onError OnSessionErrorFunc,
						onPacket OnSessionRecvPacketFunc) *SessionCallback {
	return &SessionCallback{
		OnClosed: onClosed,
		OnError: onError,
		OnRecvPacket: onPacket,
	}
}