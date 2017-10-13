package socket

type EventDispatcher interface {
	//注册事件回调
	AddCallback(id uint32, f func(interface{})) *CallbackContext

	//注册默认处理
	AddDefaultCallback(f func(interface{})) *CallbackContext

	RemoveCallback(id uint32)

	//设置事件截获钩子, 在CallData中调用钩子
	InjectData(func(interface{}) bool)

	//直接调用消费者端的handler
	CallData(data interface{})

	//清除所有回调
	Clear()

	Count() int

	CountByID(id uint32) int

	VisitCallback(callback func(uint32, *CallbackContext) VisitOperation)
}

type CallbackContext struct {
	ID   	uint32
	Tag 	interface{}
	Func 	func(interface{})
}

func NewEventDispatcher() EventDispatcher {
	self := &eventDispatcher{
		handlerByMsgPeer: make(map[uint32][]*CallbackContext),
	}

	return self
}

type eventDispatcher struct {
	//保证注册发生在初始化, 读取发生在之后可以不用锁，并且对一个msg，对应一堆handler
	handlerByMsgPeer 	map[uint32][]*CallbackContext
	handlerDefault 		*CallbackContext
	inject 				func(interface{}) bool
}

//注册事件回调
func (self *eventDispatcher) AddCallback(id uint32, f func(interface{})) *CallbackContext {
	//事件
	ctxList, ok := self.handlerByMsgPeer[id]
	if !ok {
		ctxList = make([]*CallbackContext, 0)
	}

	newCtx := &CallbackContext{
		ID:   id,
		Func: f,
	}

	ctxList = append(ctxList, newCtx)
	self.handlerByMsgPeer[id] = ctxList

	return newCtx
}

func (self *eventDispatcher) AddDefaultCallback(f func(interface{})) *CallbackContext {
	newCtx := &CallbackContext{
		ID:   0,
		Func: f,
	}
	self.handlerDefault = newCtx
	return newCtx
}

func (self *eventDispatcher) RemoveCallback(id uint32) {
	delete(self.handlerByMsgPeer, id)
}

//注入回调, 返回false时表示不再投递
func (self *eventDispatcher) InjectData(f func(interface{}) bool) {
	self.inject = f
}

type VisitOperation int

const (
	VISIT_OPERATION_CONTINUE = iota // 循环下一个
	VISIT_OPERATION_REMOVE          // 删除当前元素
	VISIT_OPERATION_EXIT            // 退出循环
)

func (self *eventDispatcher) VisitCallback(callback func(uint32, *CallbackContext) VisitOperation) {
	var needDelete []uint32

	for id, ctxList := range self.handlerByMsgPeer {
		var needRefresh bool
		var index = 0
		for {
			if index >= len(ctxList) {
				break
			}

			ctx := ctxList[index]
			op := callback(id, ctx)
			switch op {
				case VISIT_OPERATION_CONTINUE:
					index++
				case VISIT_OPERATION_REMOVE:
					if len(ctxList) == 1 {
						needDelete = append(needDelete, id)
					}

					ctxList = append(ctxList[:index], ctxList[index+1:]...)
					needRefresh = true
				case VISIT_OPERATION_EXIT:
					goto END_LOOP
			}
		}

		if needRefresh {
			self.handlerByMsgPeer[id] = ctxList
		}
	}

END_LOOP:
	if len(needDelete) > 0 {
		for _, id := range needDelete {
			delete(self.handlerByMsgPeer, id)
		}
	}
}

func (self *eventDispatcher) Clear() {
	self.handlerByMsgPeer = make(map[uint32][]*CallbackContext)
}

func (self *eventDispatcher) Exists(id uint32) bool {
	_, ok := self.handlerByMsgPeer[id]
	return ok
}

func (self *eventDispatcher) Count() int {
	return len(self.handlerByMsgPeer)
}

func (self *eventDispatcher) CountByID(id uint32) int {
	if v, ok := self.handlerByMsgPeer[id]; ok {
		return len(v)
	}

	return 0
}

type contentIndexer interface {
	ContextID() uint32
}

//通过数据接口调用
func (self *eventDispatcher) CallData(data interface{}) {
	switch d := data.(type) {
		//ID索引的消息
		case contentIndexer:
			if self == nil {
				logErrorln("recv indexed event, but event dispatcher nil, id: %d", d.ContextID())
				return
			}

			//先处理注入
			if self.inject != nil && !self.inject(data) {
				return
			}

			if ctxList, ok := self.handlerByMsgPeer[d.ContextID()]; ok {
				for _, ctx := range ctxList {
					ctx.Func(data)
				}
			}else {
				if self.handlerDefault != nil {
					self.handlerDefault.Func(d)
				}
			}
		//直接回调
		case func():
			d()
		default:
			logErrorln("unknown queue data: ", data)
	}
}
