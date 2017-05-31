package utils_test

import (
    "testing"
    "sync"

    "smartgo/libs/utils"
)

func TestAtomicInt64(t *testing.T) {
    ai64 := utils.NewAtomicInt64(0)
    wg := &sync.WaitGroup{}
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            ai64.GetAndIncrement()
            wg.Done()
        }()
    }
    wg.Wait()

    if cnt := ai64.Get(); cnt != 3 {
        t.Errorf("AtomicInt64::Get() = %d, want 3", cnt)
    }

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            ai64.GetAndDecrement()
            wg.Done()
        }()
    }
    wg.Wait()

    if cnt := ai64.Get(); cnt != 0 {
        t.Errorf("AtomicInt64::Get() = %d, want 0", cnt)
    }
}