package main

import (
	"fmt"
)

//go tool vet -shadow xxx.go用于检查如下的隐藏变量。直接在源文件目录下：go tool vet -shadow *.go
func go_tool_vet() {
	x := 1
	fmt.Println(x) //prints 1
	{
		fmt.Println(x) //prints 1
		x := 2
		fmt.Println(x) //prints 2，在当前作用域覆盖掉了之前的x
	}
	fmt.Println(x) //prints 1 (bad if you need 2)
}

//slice可以是nil，并且可以在nil的slice上append。在nil的map上会panic
//对map不能使用cap
//所谓空的字符串，即 var s string，不是nil（nil只能是slice map channel interface或者函数），而是“”
//所以判断空字符串是 if s == ""

//数组不是引用的，数组作为函数参数，是值传递，想修改只能传数组指针
//数组 x := [3]int{1,2,3}
//数组指针 *[]int
//更好选择是传递slice，这个是引用的
