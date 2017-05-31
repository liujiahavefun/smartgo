package utils_test

import (
	"fmt"
	"testing"

	"smartgo/libs/utils"
)

func TestNewConcurrentMap(t *testing.T) {
	cm := utils.NewConcurrentMap()
	if cm == nil {
		t.Error("NewConcurrentMap() = nil, want non-nil")
	}
	if !cm.IsEmpty() {
		t.Error("ConcurrentMap::IsEmpty() = false")
	}
	if size := cm.Size(); size != 0 {
		t.Errorf("ConcurrentMap::Size() = %d, want 0", size)
	}
}

func TestConcurrentMapPutAndRemove(t *testing.T) {
	cm := utils.NewConcurrentMap()
	var i int64
	for i = 0; i < 10000; i++ {
		cm.Put(i, fmt.Sprintf("%d", i))
	}

	if size := cm.Size(); size != 10000 {
		t.Errorf("ConcurrentMap::Size() = %d, want 10000", size)
	}

	for i = 0; i < 10000; i++ {
		if ok := cm.Remove(i); !ok {
			t.Errorf("ConcurrentMap::Remove(%d) = %v", i, ok)
		}
	}

	if !cm.IsEmpty() {
		t.Errorf("ConcurrentMap::IsEmpty() = false, size %d", cm.Size())
	}
}

func TestConcurrentMapInt(t *testing.T) {
	key, val := 1, 10
	cm := utils.NewConcurrentMap()
	cm.Put(key, val)
	if size := cm.Size(); size != 1 {
		t.Errorf("ConcurrentMap::Size() = %d, want 1", size)
	}

	var ret interface{}
	var ok bool
	if ret, ok = cm.Get(key); !ok {
		t.Errorf("ConcurrentMap::Get(%d) not ok", key)
	}
	if ret.(int) != val {
		t.Errorf("ConcurrentMap::Get(%d) = %d, want %d", key, ret.(int), val)
	}
}

func TestConcurrentMapString(t *testing.T) {
	keysMap := map[string]bool{
		"Lucy":  false,
		"Lily":  false,
		"Kathy": false,
		"Joana": false,
		"Belle": false,
		"Fiona": false,
	}

	valuesMap := map[string]bool{
		"Product Manager": false,
		"Rust Programmer": false,
		"Python":          false,
		"Golang":          false,
		"Java":            false,
		"Javascript":      false,
	}

	keyValueMap := map[string]string{
		"Lucy":  "Product Manager",
		"Lily":  "Rust Programmer",
		"Kathy": "Python",
		"Joana": "Golang",
		"Belle": "Java",
		"Fiona": "Javascript",
	}

	cm := utils.NewConcurrentMap()
	cm.Put("Lucy", "Product Manager")
	cm.Put("Lily", "C++")
	cm.Put("Kathy", "Python")
	cm.Put("Joana", "Golang")
	cm.Put("Belle", "Java")
	cm.PutIfAbsent("Joana", "Objective-C")
	cm.PutIfAbsent("Fiona", "Javascript")

	if size := cm.Size(); size != 6 {
		t.Errorf("ConcurrentMap::Size() = %d, want 6", size)
	}
	cm.Put("Lily", "Rust Programmer")

	for key := range cm.IterKeys() {
		keysMap[key.(string)] = true
	}

	for k, v := range keysMap {
		if !v {
			t.Errorf("Key %s not in ConcurrentMap", k)
		}
	}

	for value := range cm.IterValues() {
		valuesMap[value.(string)] = true
	}

	for k, v := range valuesMap {
		if !v {
			t.Errorf("Value %s not in ConcurrentMap", k)
		}
	}

	for item := range cm.IterItems() {
		mapKey := item.Key.(string)
		mapVal := item.Value.(string)
		if keyValueMap[mapKey] != mapVal {
			t.Errorf("ConcurrentMap[%s] != %s", mapKey, keyValueMap[mapKey])
		}
	}

	if ok := cm.Remove("Lucy"); !ok {
		t.Error(`ConcurrentMap::Remove("Lucy") = false`)
	}

	if ok, _ := cm.ContainsKey("Lucy"); ok {
		t.Error(`ConcurrentMap::ContainsKey("Lucy") = true`)
	}

	if ok, _ := cm.ContainsKey("Joana"); !ok {
		t.Error(`ConcurrentMap::ContainsKey("Joana") = false`)
	}

	cm.Clear()
	if !cm.IsEmpty() {
		t.Errorf("ConcurrentMap::IsEmpty() = false, size %d", cm.Size())
	}
}
