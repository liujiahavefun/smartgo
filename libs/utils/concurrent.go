/*
Data types and structues for concurrent use:
AtomicInt32
AtomicInt64
AtomicBoolean
ConcurrentMap
*/

package utils

import (
	"sync"
)

const INITIAL_SHARD_SIZE = 16

type ConcurrentMap struct {
	shards []*syncMap
}

func NewConcurrentMap() *ConcurrentMap {
	cm := &ConcurrentMap{
		shards: make([]*syncMap, INITIAL_SHARD_SIZE),
	}
	for i, _ := range cm.shards {
		cm.shards[i] = newSyncMap()
	}
	return cm
}

func (cm *ConcurrentMap) Put(k, v interface{}) error {
	if IsNil(k) {
		return ErrorNilKey
	}
	if IsNil(v) {
		return ErrorNilValue
	}
	if shard, err := cm.shardFor(k); err != nil {
		return err
	} else {
		shard.put(k, v)
	}
	return nil
}

func (cm *ConcurrentMap) PutIfAbsent(k, v interface{}) error {
	if IsNil(k) {
		return ErrorNilKey
	}
	if IsNil(v) {
		return ErrorNilValue
	}
	if shard, err := cm.shardFor(k); err != nil {
		return err
	} else {
		if _, ok := shard.get(k); !ok {
			shard.put(k, v)
		}
	}
	return nil
}

func (cm *ConcurrentMap) Get(k interface{}) (interface{}, bool) {
	if IsNil(k) {
		return nil, false
	}
	if shard, err := cm.shardFor(k); err != nil {
		return nil, false
	} else {
		return shard.get(k)
	}
}

func (cm *ConcurrentMap) ContainsKey(k interface{}) (bool, error) {
	if IsNil(k) {
		return false, ErrorNilKey
	}
	if shard, err := cm.shardFor(k); err != nil {
		return false, err
	} else {
		_, ok := shard.get(k)
		return ok, nil
	}
}

func (cm *ConcurrentMap) Remove(k interface{}) bool {
	if IsNil(k) {
		return false
	}
	if shard, err := cm.shardFor(k); err != nil {
		return false
	} else {
		return shard.remove(k)
	}
}

func (cm *ConcurrentMap) IsEmpty() bool {
	return cm.Size() <= 0
}

func (cm *ConcurrentMap) Clear() {
	for _, s := range cm.shards {
		s.clear()
	}
}

func (cm *ConcurrentMap) Size() int {
	var size int = 0
	for _, s := range cm.shards {
		size += s.size()
	}
	return size
}

func (cm *ConcurrentMap) shardFor(k interface{}) (*syncMap, error) {
	if code, err := Hash(k); err != nil {
		return nil, err
	} else {
		return cm.shards[code&uint32(INITIAL_SHARD_SIZE-1)], nil
	}
}

func (cm *ConcurrentMap) IterKeys() <-chan interface{} {
	kch := make(chan interface{})
	go func() {
		for _, s := range cm.shards {
			s.RLock()
			defer s.RUnlock()
			for k, _ := range s.shard {
				kch <- k
			}
		}
		close(kch)
	}()
	return kch
}

func (cm *ConcurrentMap) IterValues() <-chan interface{} {
	vch := make(chan interface{})
	go func() {
		for _, s := range cm.shards {
			s.RLock()
			defer s.RUnlock()
			for _, v := range s.shard {
				vch <- v
			}
		}
		close(vch)
	}()
	return vch
}

func (cm *ConcurrentMap) IterItems() <-chan Item {
	ich := make(chan Item)
	go func() {
		for _, s := range cm.shards {
			s.RLock()
			defer s.RUnlock()
			for k, v := range s.shard {
				ich <- Item{k, v}
			}
		}
		close(ich)
	}()
	return ich
}

type Item struct {
	Key, Value interface{}
}

type syncMap struct {
	shard map[interface{}]interface{}
	sync.RWMutex
}

func newSyncMap() *syncMap {
	return &syncMap{
		shard: make(map[interface{}]interface{}, 1024),
	}
}

func (sm *syncMap) put(k, v interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.shard[k] = v
}

func (sm *syncMap) get(k interface{}) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	v, ok := sm.shard[k]
	return v, ok
}

func (sm *syncMap) size() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.shard)
}

func (sm *syncMap) remove(k interface{}) bool {
	sm.Lock()
	defer sm.Unlock()
	_, ok := sm.shard[k]
	delete(sm.shard, k)
	return ok
}

func (sm *syncMap) clear() {
	sm.Lock()
	defer sm.Unlock()
	sm.shard = make(map[interface{}]interface{})
}
