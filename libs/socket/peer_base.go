package socket

//Peer间的共享数据
type peerBase struct {
	//EventDispatcher
	//EventQueue

	name          string
	address       string
	maxPacketSize int
}

func newPeerBase() *peerBase {
	self := &peerBase{
	}

	return self
}

func (self *peerBase) SetName(name string) {
	self.name = name
}

func (self *peerBase) Name() string {
	return self.name
}

func (self *peerBase) SetMaxPacketSize(size int) {
	self.maxPacketSize = size
}

func (self *peerBase) MaxPacketSize() int {
	return self.maxPacketSize
}
