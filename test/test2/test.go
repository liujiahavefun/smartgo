package main

import (
    "fmt"
)

type A interface {
    GetId() int64  
}

type B interface {
    A
    GetName() string
}

type AImpl struct {
    id int64
}

func NewA(id int64) *AImpl {
    return &AImpl{
        id: id,
    }
}

func (self *AImpl) GetId() int64 {
    return self.id
}

type BImpl struct {
    *AImpl
    name string
}

func NewB(id int64, name string) B {
    return &BImpl{
        AImpl: NewA(id),
        name: name,
    }
}

func (self *BImpl) GetName() string {
    return self.name
}

func main() {
    b := NewB(1, "liujia")
    if a, ok := b.(interface {GetId2() int64}); ok {
        fmt.Println("b impl GetId2()")
        id := a.GetId2()
        fmt.Printf("%T %v \n", a, a)
        fmt.Println(id)
    }else {
        fmt.Println("b NOT impl GetId2()")
    }

    if a, ok := b.(interface {GetId() int64}); ok {
        fmt.Println("b impl GetId()")
        id := a.GetId()
        fmt.Printf("%T %v \n", a, a)
        fmt.Println(id)
    }else {
        fmt.Println("b NOT impl GetId()")
    }
}