package socket

import (
	"sync"
	"sync/atomic"
)

const (
	TOTAL_TRY_COUNT = 100
)

//TODO: 用带Hash的map进行优化，sesIDAcc用Atomic int64类型
type sessionMgr struct {
	sesMap map[int64]Session

	sesIDAcc    int64
	sesMapGuard sync.RWMutex
}

func newSessionManager() *sessionMgr {
	return &sessionMgr{
		sesMap: make(map[int64]Session),
	}
}

func (self *sessionMgr) Add(ses Session) {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	var tryCount int = TOTAL_TRY_COUNT
	var id int64

	//id翻越处理
	for tryCount > 0 {
		id = atomic.AddInt64(&self.sesIDAcc, 1)
		if _, ok := self.sesMap[id]; !ok {
			break
		}
		tryCount--
	}

	if tryCount == 0 {
		logWarningf("sessionID override! %v", id)
	}

	//ltvses := ses.(*tcpServerSession)
	//ltvses.id = id
	ses.SetID(id)
	self.sesMap[id] = ses
}

func (self *sessionMgr) Remove(ses Session) {
	self.sesMapGuard.Lock()
	delete(self.sesMap, ses.GetID())
	self.sesMapGuard.Unlock()
}

//根据ID获得一个session
func (self *sessionMgr) GetSession(id int64) Session {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	v, ok := self.sesMap[id]
	if ok {
		return v
	}

	return nil
}

//遍历访问所有的session
func (self *sessionMgr) VisitSession(callback func(Session) bool) {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	for _, ses := range self.sesMap {
		if !callback(ses) {
			break
		}
	}
}

func (self *sessionMgr) SessionCount() int {
	//这里加读锁好还是写锁好？
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	return len(self.sesMap)
}
