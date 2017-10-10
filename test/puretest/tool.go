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

//string不是[]byte，底层存储可以看做不可变的[]byte
//string <-> []byte 时，应该是得到一份新的拷贝（不一定，有些地方不会拷贝，看做是特定的优化吧）
//如果修改string的特定字符，不能直接string[0]。要先转成[]byte，改完再转为string，有点蛋疼，没办法。或者string转为[]rune
// s := "123", fmt.Println(s[0]) ok，但s[0] = '2'会编译错误
//字符串编码可以用“unicode/utf8”这个包，例如 utf8.ValidString 判断是否是utf8字符串
//len(s)字符串，返回的是utf8后的字符数。如果想看有多少个额字符，用unicode/utf8包的RuneCountInString
//对字符串做 for range,返回的第一个是第二个unicode字符在utf8底层的索引，第二个是unicode字符(或者是0xfffd，如果有错误编码)。所以如果有中文，第一个参数绝逼不连续
//正确做法是将字符串转为[]byte，然后依次打印每一个byte

//带缓存的channel关闭后，仍然可以读，返回的第一个是值，第二个ok表明是否还有数据，有就是true没有就false
//想关闭的channel发送，会panic
