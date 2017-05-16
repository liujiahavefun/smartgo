package socket

import (
	"sync"
)

type PacketList struct {
	list      []*Packet
	listGuard sync.Mutex
	listCond  *sync.Cond
}

func NewPacketList() *PacketList {
	self := &PacketList{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}

func (self *PacketList) Add(packet *Packet) {
	self.listGuard.Lock()
	self.list = append(self.list, packet)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *PacketList) Reset() {
	self.list = self.list[0:0]
}

func (self *PacketList) BeginPick() []*Packet {
	//condition variable标准用法
	self.listGuard.Lock()
	for len(self.list) == 0 {
		self.listCond.Wait()
	}
	self.listGuard.Unlock()

	self.listGuard.Lock()

	return self.list
}

func (self *PacketList) EndPick() {
	self.Reset()
	self.listGuard.Unlock()
}