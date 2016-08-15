package redis

// The pool package implements a connection pool for redis connections which is
// thread-safe

import (
	"github.com/fzzy/radix/redis"
	"time"
)

// A simple connection pool. It will create a small pool of initial connections,
// and if more connections are needed they will be created on demand. If a
// connection is returned and the pool is full it will be closed.
type Pool struct {
	network,
	addr string
	pool          chan *redisClient
	df            DialFunc
	clientTimeout time.Duration
}

// add createdTime to check timeout
type redisClient struct {
	Conn        *redis.Client
	createdTime time.Time
}

// A function which can be passed into NewCustomPool
type DialFunc func(network, addr string) (*redis.Client, error)

// A Pool whose connections are all created using f(network, addr). The size
// indicates the maximum number of idle connections to have waiting to be used
// at any given moment. f will be used to create all new connections associated
// with this pool.
//
// The following is an example of using NewCustomPool to have all connections
// automatically get AUTH called on them upon creation
//
//  df := func(network, addr string) (*redis.Client, error) {
//      client, err := redis.Dial(network, addr)
//      if err != nil {
//          return nil, err
//      }
//      if err = client.Cmd("AUTH", "SUPERSECRET").Err; err != nil {
//          client.Close()
//          return nil, err
//      }
//      return client, nil
//  }
//  p, _ := pool.NewCustomPool("tcp", "127.0.0.1:6379", 10, 300 * time.Seconde, df)
//
func NewCustomPool(network, addr string, size int, clientTimeout time.Duration, df DialFunc) (*Pool, error) {
	pool := make([]*redisClient, 0, size)
	now := time.Now()
	for i := 0; i < size; i++ {
		client, err := df(network, addr)
		if err != nil {
			for _, rc := range pool {
				rc.Conn.Close()
			}
			return nil, err
		}
		if client != nil {
			pool = append(pool, &redisClient{
				Conn:        client,
				createdTime: now,
			})
		}
	}
	p := Pool{
		network:       network,
		addr:          addr,
		clientTimeout: clientTimeout,
		pool:          make(chan *redisClient, len(pool)),
		df:            df,
	}
	for i := range pool {
		p.pool <- pool[i]
	}

	p.heartbeat()

	return &p, nil
}

// Creates a new Pool whose connections are all created using
// redis.Dial(network, addr). The size indicates the maximum number of idle
// connections to have waiting to be used at any given moment
func NewPool(network, addr string, size int, clientTimeout time.Duration) (*Pool, error) {
	return NewCustomPool(network, addr, size, clientTimeout, redis.Dial)
}

// Calls NewPool, but if there is an error it return a pool of the same size but
// without any connections pre-initialized (can be used the same way, but if
// this happens there might be something wrong with the redis instance you're
// connecting to)
func NewOrEmptyPool(network, addr string, size int) *Pool {
	pool, err := NewPool(network, addr, size, 0)
	if err != nil {
		pool = &Pool{
			network: network,
			addr:    addr,
			pool:    make(chan *redisClient, size),
			df:      redis.Dial,
		}
	}
	return pool
}

// Retrieves an available redis client. If there are none available it will
// create a new one on the fly
func (p *Pool) Get() (*redisClient, error) {
	select {
	case rc := <-p.pool:
		if p.clientTimeout > 0 && time.Now().Sub(rc.createdTime) > p.clientTimeout {
			rc.Conn.Close()
			return p.generate()
		}
		return rc, nil
	default:
		return p.generate()
	}
}

// heartbeat keep the client alive
func (p *Pool) heartbeat() {
	go func() {
		for {
			p.Cmd("PING")
			time.Sleep(1 * time.Second)
		}
	}()
}

// Returns a client back to the pool. If the pool is full the client is closed
// instead. If the client is already closed (due to connection failure or
// what-have-you) it should not be put back in the pool. The pool will create
// more connections as needed.
func (p *Pool) Put(rc *redisClient) {
	select {
	case p.pool <- rc:
	default:
		rc.Conn.Close()
	}
}

// Cmd automatically gets one client from the pool, executes the given command
// (returning its result), and puts the client back in the pool
func (p *Pool) Cmd(cmd string, args ...interface{}) *redis.Reply {
	var (
		reply *redis.Reply
		err   error
	)
	rc, err := p.Get()
	if err != nil {
		return &redis.Reply{
			Err: err,
		}
	}
	reply = rc.Conn.Cmd(cmd, args...)
	err = reply.Err
	defer p.CarefullyPut(rc, &err)
	return reply
}

// A useful helper method which acts as a wrapper around Put. It will only
// actually Put the conn back if potentialErr is not an error or is a
// redis.CmdError. It would be used like the following:
//
//  func doSomeThings(p *Pool) error {
//      rc, redisErr := p.Get()
//      if redisErr != nil {
//          return redisErr
//      }
//      defer p.CarefullyPut(cn, &redisErr)
//
//      var i int
//      i, redisErr = rc.Conn.Cmd("GET", "foo").Int()
//      if redisErr != nil {
//          return redisErr
//      }
//
//      redisErr = rc.Conn.Cmd("SET", "foo", i * 3).Err
//      return redisErr
//  }
//
// If we were just using the normal Put we wouldn't be able to defer it because
// we don't want to Put back a connection which is broken. This method takes
// care of doing that check so we can still use the convenient defer
func (p *Pool) CarefullyPut(rc *redisClient, potentialErr *error) {
	if potentialErr != nil && *potentialErr != nil {
		// We don't care about command errors, they don't indicate anything
		// about the connection integrity
		if _, ok := (*potentialErr).(*redis.CmdError); !ok {
			rc.Conn.Close()
			return
		}
	}
	p.Put(rc)
}

// Removes and calls Close() on all the connections currently in the pool.
// Assuming there are no other connections waiting to be Put back this method
// effectively closes and cleans up the pool.
func (p *Pool) Empty() {
	var rc *redisClient
	for {
		select {
		case rc = <-p.pool:
			rc.Conn.Close()
		default:
			return
		}
	}
}

// generate a client
func (p *Pool) generate() (*redisClient, error) {
	conn, err := p.df(p.network, p.addr)
	if err != nil {
		return nil, err
	}
	rc := &redisClient{
		Conn:        conn,
		createdTime: time.Now(),
	}
	return rc, nil
}
