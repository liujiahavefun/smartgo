package redispool

import (
	"github.com/fzzy/radix/redis"
	"time"
)

type redisClient struct {
	Conn        *redis.Client
	createdTime time.Time
}

type RedisConnFunc func(network, addr string) (*redis.Client, error)

type RedisPool struct {
	network       string
	addr          string
	pool          chan *redisClient
	cf            RedisConnFunc
	clientTimeout time.Duration
}

func NewCustomRedisPool(network, addr string, size int, clientTimeout time.Duration, cf RedisConnFunc) (*RedisPool, error) {
	rp := &RedisPool{
		network:       network,
		addr:          addr,
		clientTimeout: clientTimeout,
		pool:          make(chan *redisClient, size),
		cf:            cf,
	}

	pool := make([]*redisClient, 0, size)
	for i := 0; i < size; i++ {
		client, err := rp.createClient()
		if err != nil || client == nil {
			for _, rc := range pool {
				rc.Conn.Close()
			}
			return nil, err
		}
		pool = append(pool, client)
	}

	for i := range pool {
		rp.pool <- pool[i]
	}

	rp.heartbeat()

	return rp, nil
}

func NewRedisPool(network, addr string, size int, clientTimeout time.Duration) (*RedisPool, error) {
	return NewCustomRedisPool(network, addr, size, clientTimeout, redis.Dial)
}

func (self *RedisPool) Get() (*redisClient, error) {
	select {
	case rc := <-self.pool:
		/*
		   //这里不检查创建时间是否超时，倾向于一直使用此链接，除非执行redis命令出现错误
		   if self.clientTimeout > 0 && time.Now().Sub(rc.createdTime) > self.clientTimeout {
		       rc.Conn.Close()
		       return self.createClient()
		   }
		*/
		return rc, nil
	default:
		return self.createClient()
	}
}

func (self *RedisPool) Put(rc *redisClient) {
	select {
	case self.pool <- rc:
	default:
		rc.Conn.Close()
	}
}

func (self *RedisPool) Cmd(cmd string, args ...interface{}) *redis.Reply {
	var (
		reply *redis.Reply
		err   error
	)

	rc, err := self.Get()
	if err != nil {
		return &redis.Reply{
			Err: err,
		}
	}

	reply = rc.Conn.Cmd(cmd, args...)
	err = reply.Err
	defer self.internalPut(rc, err)
	return reply
}

func (self *RedisPool) Close() {
	var rc *redisClient
	for {
		select {
		case rc = <-self.pool:
			rc.Conn.Close()
		default:
			return
		}
	}
}

func (self *RedisPool) internalPut(rc *redisClient, replyErr error) {
	if replyErr != nil {
		if _, ok := replyErr.(*redis.CmdError); !ok {
			rc.Conn.Close()
			return
		}
	}
	self.Put(rc)
}

func (self *RedisPool) heartbeat() {
	go func() {
		for {
			self.Cmd("PING")
			time.Sleep(1 * time.Second)
		}
	}()
}

func (self *RedisPool) createClient() (*redisClient, error) {
	conn, err := self.cf(self.network, self.addr)
	if err != nil {
		return nil, err
	}
	rc := &redisClient{
		Conn:        conn,
		createdTime: time.Now(),
	}
	return rc, nil
}
