package net

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"time"

	"github.com/golang/glog"

	. "smartgo/libs/utils"
)

func LoadTLSConfig(certFile, keyFile string, isSkipVerify bool) (tls.Config, error) {
	var config tls.Config
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return config, err
	}
	config = tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: isSkipVerify,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}
	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader
	return config, nil
}

type TLSTCPServer struct {
	certFile string
	keyFile  string
	*TCPServer
}

func NewTLSTCPServer(addr, cert, key string) Server {
	server := &TLSTCPServer{
		certFile:  cert,
		keyFile:   key,
		TCPServer: NewTCPServer(addr).(*TCPServer),
	}

	config, err := LoadTLSConfig(server.certFile, server.keyFile, false)
	if err != nil {
		glog.Fatalln(err)
	}

	setTLSWrapper(func(conn net.Conn) net.Conn {
		return tls.Server(conn, &config)
	})

	return server
}

func (server *TLSTCPServer) IsRunning() bool {
	return server.TCPServer.IsRunning()
}

func (server *TLSTCPServer) GetAllConnections() *ConcurrentMap {
	return server.TCPServer.GetAllConnections()
}

func (server *TLSTCPServer) GetTimingWheel() *TimingWheel {
	return server.TCPServer.GetTimingWheel()
}

func (server *TLSTCPServer) GetWorkerPool() *WorkerPool {
	return server.TCPServer.GetWorkerPool()
}

func (server *TLSTCPServer) GetServerAddress() string {
	return server.TCPServer.GetServerAddress()
}

func (server *TLSTCPServer) Start() {
	server.TCPServer.Start()
}

func (server *TLSTCPServer) Close() {
	server.TCPServer.Close()
}

func (server *TLSTCPServer) SetOnScheduleCallback(duration time.Duration, callback func(time.Time, interface{})) {
	server.TCPServer.SetOnScheduleCallback(duration, callback)
}

func (server *TLSTCPServer) GetOnScheduleCallback() (time.Duration, onScheduleFunc) {
	return server.TCPServer.GetOnScheduleCallback()
}

func (server *TLSTCPServer) SetOnConnectCallback(callback func(Connection) bool) {
	server.TCPServer.SetOnConnectCallback(callback)
}

func (server *TLSTCPServer) GetOnConnectCallback() onConnectFunc {
	return server.TCPServer.GetOnConnectCallback()
}

func (server *TLSTCPServer) SetOnMessageCallback(callback func(Message, Connection)) {
	server.TCPServer.SetOnMessageCallback(callback)
}

func (server *TLSTCPServer) GetOnMessageCallback() onMessageFunc {
	return server.TCPServer.GetOnMessageCallback()
}

func (server *TLSTCPServer) SetOnCloseCallback(callback func(Connection)) {
	server.TCPServer.SetOnCloseCallback(callback)
}

func (server *TLSTCPServer) GetOnCloseCallback() onCloseFunc {
	return server.TCPServer.GetOnCloseCallback()
}

func (server *TLSTCPServer) SetOnErrorCallback(callback func()) {
	server.TCPServer.SetOnErrorCallback(callback)
}

func (server *TLSTCPServer) GetOnErrorCallback() onErrorFunc {
	return server.TCPServer.GetOnErrorCallback()
}

func (server *TLSTCPServer) SetOnPacketRecvCallback(callback onPacketRecvFunc) {
	server.onPacket = onPacketRecvFunc(callback)
}

func (server *TLSTCPServer) GetOnPacketRecvCallback() onPacketRecvFunc {
	return server.onPacket
}
