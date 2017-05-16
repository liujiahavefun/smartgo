package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
)

const (
	PACKET_HEADER_SIZE = 8    // MsgID(uint32) + Ser(uint16) + Size(uint16)
	TOTAL_SEND_TRY_TIMES = 100
)

//封包流
type PacketStream interface {
	Read() (*Packet, error)
	Write(pkt *Packet) error
	Flush() error
	Close() error
	Raw() net.Conn
}

var (
	packageTagNotMatch     = errors.New("ReadPacket: package tag not match")
	packageDataSizeInvalid = errors.New("ReadPacket: package crack, invalid size")
	packageTooBig          = errors.New("ReadPacket: package too big")
)

type ltvStream struct {
	recvser      uint16
	sendser      uint16
	conn         net.Conn
	sendtagGuard sync.RWMutex

	outputWriter     *bufio.Writer
	outputHeadBuffer *bytes.Buffer

	inputHeadBuffer []byte
	headReader      *bytes.Reader

	maxPacketSize int
}

//将net.Conn包装成一个PacketStream relay模式: 在封包头有clientid信息
func NewPacketStream(conn net.Conn) *ltvStream {
	s := &ltvStream{
		conn:             conn,
		recvser:          1,
		sendser:          1,
		outputWriter:     bufio.NewWriter(conn),
		outputHeadBuffer: bytes.NewBuffer([]byte{}),
		inputHeadBuffer:  make([]byte, PACKET_HEADER_SIZE),
	}
	s.headReader = bytes.NewReader(s.inputHeadBuffer)

	return s
}

//从socket读取1个packet,并返回
func (self *ltvStream) Read() (p *Packet, err error) {
	if _, err = self.headReader.Seek(0, 0); err != nil {
		return nil, err
	}

	if _, err = io.ReadFull(self.conn, self.inputHeadBuffer); err != nil {
		return nil, err
	}

	p = &Packet{}

	//读取ID，4字节
	if err = binary.Read(self.headReader, binary.LittleEndian, &p.MsgID); err != nil {
		return nil, err
	}

	//读取序号，2字节
	var ser uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &ser); err != nil {
		return nil, err
	}

	//读取包大小，2字节
	var fullsize uint16
	if err = binary.Read(self.headReader, binary.LittleEndian, &fullsize); err != nil {
		return nil, err
	}

	//包大小太大
	if self.maxPacketSize > 0 && int(fullsize) > self.maxPacketSize {
		return nil, packageTooBig
	}

	//序列号不匹配
	if self.recvser != ser {
		return nil, packageTagNotMatch
	}

	//packet body的大小
	dataSize := fullsize - PACKET_HEADER_SIZE
	if dataSize < 0 {
		return nil, packageDataSizeInvalid
	}

	//读取数据
	//TODO: 内存池
	p.Data = make([]byte, dataSize)
	if _, err = io.ReadFull(self.conn, p.Data); err != nil {
		return nil, err
	}

	//增加序列号值
	self.recvser++

	return p, nil
}

//发送一个Packet
func (self *ltvStream) Write(pkt *Packet) (err error) {
	//防止将Send放在go内造成的多线程冲突问题
	self.sendtagGuard.Lock()
	defer self.sendtagGuard.Unlock()

	self.outputHeadBuffer.Reset()

	//写ID
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, pkt.MsgID); err != nil {
		return err
	}

	//写序号
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, self.sendser); err != nil {
		return err
	}

	//写包大小
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, uint16(len(pkt.Data)+PACKET_HEADER_SIZE)); err != nil {
		return err
	}

	//发包头
	if err = self.writeFull(self.outputHeadBuffer.Bytes()); err != nil {
		return err
	}

	//发包内容
	if err = self.writeFull(pkt.Data); err != nil {
		return err
	}

	//增加序号值
	self.sendser++

	return
}

//完整发送所有封包
func (self *ltvStream) writeFull(p []byte) error {
	sizeToWrite := len(p)
	for {
		n, err := self.outputWriter.Write(p)
		if err != nil {
			return err
		}

		if n >= sizeToWrite {
			break
		}

		copy(p[0:sizeToWrite-n], p[n:sizeToWrite])
		sizeToWrite -= n
	}

	return nil
}

func (self *ltvStream) Flush() error {
	var err error
	for tryTimes := 0; tryTimes < TOTAL_SEND_TRY_TIMES; tryTimes++ {

		err = self.outputWriter.Flush()

		//如果没写完, flush底层会将没发完的buff准备好, 我们只需要重新调一次flush
		if err != io.ErrShortWrite {
			break
		}
	}

	return err
}

//关闭
func (self *ltvStream) Close() error {
	return self.conn.Close()
}

//返回Raw Socket
func (self *ltvStream) Raw() net.Conn {
	return self.conn
}

