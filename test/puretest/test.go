package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const (
	Enone = iota
	Eio
	Einval
)

func main() {
	//AboutNew()
	//AboutMake()
	//AboutSlice()
	//AboutMap()
	//AboutPrint()
	//AboutAppend()
	//AboutConst()
	//AboutMethod()
	//AboutInterface()
	//AboutTypeAssertiong()
	//AboutEmbedding()
	//AboutRecover()

	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadLine()

	fmt.Println(getGOMAXPROCS())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(getGOMAXPROCS())
}

//这段代码很简练，有golang的好味道，返回值是命名的，有默认初始值
func ReadFull(r bufio.Reader, buf []byte) (n int, err error) {
	for len(buf) > 0 && err == nil {
		var nr int
		nr, err = r.Read(buf)
		n += nr
		buf = buf[nr:]
	}
	return
}

func getGOMAXPROCS() int {
	//GOMAXPROCS sets the maximum number of CPUs that can be executing simultaneously and returns the previous setting.
	//If n < 1, it does not change the current setting. The number of logical CPUs on the local machine can be queried with NumCPU.
	//This call will go away when the scheduler improves.
	return runtime.GOMAXPROCS(0)
}

func AboutRecover() {
	//见AboutPanic()
	go func() {
		i := AboutPanic()
		fmt.Println("AboutPanic() returns ", i)
	}()

	//这里有个注意点，recover()只能使得出问题的goroutine自己挂了，不影响整个进程。所以recover后，此goroutine退出(如果主goroutine挂了也不行)，其它goroutine正常运行。
}

func AboutPanic() (ret int) {
	//注册recover()函数的defer通常写在函数的开始处，保证首先被注册掉，如果写在panic下面，则不会执行，因为还没注册就挂了
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("work failed: ", err)
			ret = 1 //修改返回值，异常退出时的返回值。defer可以修改命名的返回变量
		}
	}()

	//通过触发，就是手工调用panic()函数，可以通知有某个不可继续的错误发生了，要退出程序，通常panic的参数是个字符串
	panic("just test panic")
	return 0

	//panic行为说明
	//调用panic将立刻中断当前函数的执行，并展开当前Goroutine的调用栈，依次执行之前注册的defer函数。
	//当栈展开操作达到该Goroutine栈顶端时，程序将终止。

	//recovery()
	//调用 recover() 方法会终止栈展开操作并返回之前传递给 panic 方法的那个参数。
	//由于在栈展开过程中，只有defer型函数会被执行，因此 recover 的调用必须置于defer函数内才有效。
	//见开始那个defer，标准写法就是1）函数开始用defer注册recover 2)recover的写法也跟上面那个差不多
}

func AboutChannel() {
	//没啥新鲜的，有个例子不错
	//LeakyList，构造一个100缓冲的池子，池子里是*Buffer，需要就从里面取，没有就new。用完了池子没满就放到里面去，否则就Free
	/*
		var freeList = make(chan *Buffer, 100)
		var serverChan = make(chan *Buffer)

		func client() {
			for {
				var b *Buffer
				// Grab a buffer if available; allocate if not.
				select {
				case b = <-freeList:
					// Got one; nothing more to do.
				default:
					// None free, so allocate a new one.
					b = new(Buffer)
				}
				load(b) // Read next message from the net.
				serverChan <- b // Send to server.
			}
		}

		func server() {
		for {
			b := <-serverChan // Wait for work.
			process(b)
			// Reuse buffer if there's room.
			select {
			case freeList <- b:
				// Buffer on free list; nothing more to do.
			default:
				// Free list full, just carry on.
			}
			}
		}
	*/
}

//GO中比较标准的内嵌类型，即匿名包含一个类型，类似于OO中的继承。ReadWriter1默认已经有Reader Writer两个接口了，应该也有ReaderWriter这个接口吧？
//这样ReadWriter1，即当我们内嵌一个类型时，该类型的所有方法会变成外部类型的方法，
//但是当这些方法被调用时，其接收的参数仍然是内部类型，而非外部类型
type ReadWriter1 struct {
	bufio.Reader // *bufio.Reader
	bufio.Writer // *bufio.Writer
}

//包含一个自类型，即命名的包含一个类型的子字段，这样Reader/Writer的接口，需要自己重新搞一下
type ReadWriter2 struct {
	reader bufio.Reader
	writer bufio.Writer
}

func AboutEmbedding() {
	//这里还需要再看书！！！

	// 1: interface
	//Reader和Writer是两个接口，ReadWriter将两个接口Combine在一起，GO推崇组合，而不是继承。即接口只能被embed而不能被继承
	type Reader interface {
		Read(p []byte) (n int, err error)
	}
	type Writer interface {
		Write(p []byte) (n int, err error)
	}
	// ReadWriter is the interface that combines the Reader and Writer interfaces.
	type ReadWriter interface {
		Reader
		Writer
	}

	// 2: struct
	// 见上面
	s := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	br := bufio.NewReader(s) //return *bufio.Reader

	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b) //return *bufio.Writer

	rw1 := &ReadWriter1{*br, *bw}
	rw2 := &ReadWriter2{*br, *bw}
	fmt.Printf("%T\n", rw2)

	br.WriteTo(rw1)
	//br.WriteTo(rw2)

	data := []byte{1, 2, 3, 4, 5}
	rw1.Read(data)
	fmt.Printf("%v\n", data)
	rw1.Write(data)
	fmt.Printf("%v\n", data)
}

func About_() {
	//_为占位符，大体上的怎么用不再赘述，就说两个比较罕见的用法

	// import _ "os"
	//表面引入包但不能使用包，通常用于仅仅需要包的init()函数，但不使用包中其它的功能

	// import fmt 如果不用fmt的话，会编译报错，如果不想改可以这样，用空白占位符使用包的功能
	// var _ = fmt.Printf // For debugging; delete when done

	// GO中类型转换（尤其是接口类型的），编译器会检查，检查某个类型是否实现了某个接口所需的函数，这样才能转换，否则失败，静态检查
	// 动态检查通常发生在如json.Marshaler，会对接口做动态检查，就是类型断言
	// 在一个包中，如果想确保某个类型实现了某个类型可以这样，如下确保*RawMessage实现json.Marshaler接口
	// var _ json.Marshaler = (*RawMessage)(nil)
	// 编译器编译时，会做这个检查，而且这个检查不影响功能实现
}

func AboutTypeAssertiong() {
	t := "liujia"
	var t1 interface{} = t

	//根据接口的动态类型，获取值，可以这样。如果switch case只有一个，即判断某个接口是否是某个动态类型，并取值，但有更好的办法，就是类型断言
	switch val := t1.(type) {
	case string:
		fmt.Println("string: ", val)
	case int:
		fmt.Println("int: ", val)
	default:
		fmt.Println("I do not know")
	}

	//类型断言, t.(具体要断言的类型，如string int []byte等)，ok表示是不是要断言的类型，如果是的话，val是转换后的类型值
	if val, ok := t1.(string); ok {
		fmt.Println("string:", val)
	} else {
		fmt.Println("not string")
	}
}

//其实可以直接将ByteSize定义为float64，这样下面switch那里，直接b/PB就可以了，不用再转换了，将 1 << 10这样的值赋给float64也没问题
type ByteSize int64

const (
	_           = iota             // ignore first value by assigning to blank identifier，第一个iota应该是0
	KB ByteSize = 1 << (10 * iota) //第二个iota就是1
	MB                             //默认重复上面的表达式，就是说这行其实是MB ByteSize = 1 << (10 * iota)，并且第三行iota应该是2
	GB
	TB
	PB
)

func (b ByteSize) String() string {
	switch {
	case b >= PB:
		return fmt.Sprintf("%.2fPB", float64(b)/float64(PB))
	case b >= TB:
		return fmt.Sprintf("%.2fTB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.2fGB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2fMB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2fKB", float64(b)/float64(KB))
	}
	return fmt.Sprintf("%.2fB", b)
}

type Sequence []int

// Methods required by sort.Interface.
func (s Sequence) Len() int {
	return len(s)
}
func (s Sequence) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Method for printing - sorts the elements before printing.
func (s Sequence) String() string {
	//sort.Sort(s)
	str := "["
	for i, elem := range s {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprint(elem)
	}
	return str + "]"
}

//Sequence和[]int，底层是同样的类型，仅仅名字区别。做转换时，不生成新值（float和int转换会生成新值）。但是下面这样写，由于Sequence 和[]int还是不同的类型，不对应于同一个String()
//所以这里不会导致String()函数的无限递归。
func (s Sequence) String2() string {
	return fmt.Sprint([]int(s))
}

//在GO中，将表达式的类型进行转换，来访问不同的方法集合，是一种常见用法。如下，将Sequence转为[]int后再排序
//现在， Sequence 可以没有实现多个接口（排序和打印），相反的，我们利用了能够将数据项转换为多个类型
//（ Sequence ， sort.IntSlice 和 []int ）的能力，每个类型完成工作的一部分。这在实际中不常见，但是却可以很有效。
func (s Sequence) sort() string {
	sort.IntSlice(s).Sort()
	return fmt.Sprint([]int(s))
}

func AboutInterface() {
	//Sequence实现了sort接口的三个函数，和String()函数
	s := Sequence([]int{1, 3, 2, 7, 4})
	//sort.Sort(s)
	s.Swap(1, 3)
	fmt.Println(s)

	fmt.Println("...")
	sort.Sort(s)
	fmt.Println(s)

	fmt.Println(s.String2())
}

func AboutConst() {
	//常量一般用这种方式定义，只能是数值（整数浮点等），字符串，bool等编译器能确定的，其它需要运行期才能确定的，就是说要调用函数的都不行
	//const(
	//	xxx := xxx
	//)

	//上面代码，很有代表性
	size := 1<<20 + 1<<10
	fmt.Println(ByteSize(size)) //1.00MB
}

func AboutVar() {
	//变量像常量那样，不过值可以是运行期才能确定的
	//var(
	//	home = os.Getenv("HOME")   //获取环境变量
	//	user = os.Getenv("USER")
	//	gopath = os.Getenv("GOPATH")
	//)
}

func init() {
	fmt.Println("in init() function")
}

func AboutInit() {
	//最后，每个源文件可以定义自己的不带参数的（niladic） init 函数，来设置它所需的状态。（实际上每个
	//文件可以有多个 init 函数。） init 是在程序包中所有变量声明都被初始化，以及所有被导入的程序包中的
	//变量初始化之后才被调用。

	//除了用于无法通过声明来表示的初始化以外， init 函数的一个常用法是在真正执行之前进行验证或者修复
	//程序状态的正确性。
}

type TT []byte

func (tt TT) Append1(data []byte) (n int, err error) {
	fmt.Printf("---Append1 %T\n", tt)
	tt = append(tt, data...)
	return len(data), nil
}

func (tt *TT) Append2(data []byte) (n int, err error) {
	fmt.Printf("---Append2 %T\n", tt)
	*tt = append(*tt, data...)
	return len(data), nil
}

func (tt *TT) Swap1(i, j int) {
	t := *tt
	t[i], t[j] = t[j], t[i]
}

func (tt TT) Swap2(i, j int) {
	tt[i], tt[j] = tt[j], tt[i]
}

//实现io.Write接口，可以用于fmt.Fprint等。在字节切片上使用 Write 的思想，是实现 bytes.Buffer 的核心！！！
func (tt *TT) Write(data []byte) (n int, err error) {
	*tt = append(*tt, data...) //写出这样是编译不过的。。。tt = append(tt, data...)
	return len(data), nil
}

type TTT int

func (self TTT) Increament1() {
	self++
}

func (self *TTT) Increament2() {
	(*self)++
}

func AboutMethod() {
	//所谓method就是定义在类型上的函数
	//这里的类型通常是1）别名类型，即类似 type ByteSize int64这样的 2）自定义struct，如 type T struct{}
	//并且不能是指针（type出来的除外吧？）和接口

	//方法的接收者，可以是1）值，如上Append1的TT  2)指针，如上Append2的*TT
	//值的接收者不会修改接收者，因为传递的是值的拷贝
	//指针的接收者，可以改变接收者，因为传递的指针，通过指针修改接收者
	//若方法接收者是值，可以在指针或者值上调用，反正都是返回拷贝。若是指针，则只能在指针上调用。（使用拷贝的值来调用它们没有意义，但可以编译，只不过将会导致那些修改会被丢弃）
	//我的理解是，对于值接收者，传递值或者指针都可以，但值不会变。对于指针接收者，传递值或者指针都行，但值会变

	var t1 TT = TT("liujia")
	var t2 *TT = &t1
	//var t3 *TT = &(TT("111")) //这样直接赋给一个临时变量的指针，似乎不行，因为TT("111")实际是一个强转，这里有点奇怪，转完也是TT啊
	t3 := &(TT{1, 2, 3})  //这样可以，直接将TT初始化
	t4 := &([]byte{1, 2}) //这样是可以的，但t3是*[]byte，而不是*TT

	//test Swap，下面这些都能编译过并且运行正常
	//这个跟上说的其实一样，但有代表性，所谓传递指针和值，就是传递的调用者本身还是调用者的拷贝
	//但对于slice而言，由于是引用类型，传值和指针无甚区别，而且是swap操作，都会修改底层数组
	t3.Swap1(0, 2)
	fmt.Println(t3)
	(*t3).Swap1(0, 2)
	fmt.Println(t3)
	t3.Swap2(0, 2)
	fmt.Println(t3)
	(*t3).Swap2(0, 2)
	fmt.Println(t3)
	fmt.Println("...")

	//同样的，参考一下TTT的定义和实现
	i := TTT(1)
	p := &i
	i.Increament1()
	fmt.Println("i ", i) //i 1  对于值接收者，值没变
	i.Increament2()
	fmt.Println("i ", i) //i 2  对于指针接收者，值+1
	p.Increament1()
	fmt.Println("p ", p, " *p ", *p) //p  0xc42000a3c8  *p  2   传递指针，对于值接收者一样，值不变
	p.Increament2()
	fmt.Println("p ", p, " *p ", *p) //p  0xc42000a3c8  *p  3  传递指针，值+1

	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)

	fmt.Printf("%T\n", t1) //main.TT
	fmt.Printf("%T\n", t2) //*main.TT
	fmt.Printf("%T\n", t3) //*main.TT
	fmt.Printf("%T\n", t4) //*[]uint8

	t5 := []byte{10, 11}
	size, err := t1.Append1(t5) //在值上调用值接收者函数，值不会变
	fmt.Printf("%v %v %v\n", size, err, t1)
	fmt.Println("")

	size, err = t1.Append2(t5) //在值上调用指针接收者函数，自动转为指针，值会变
	fmt.Printf("%v %v %v\n", size, err, t1)
	fmt.Println("")

	size, err = t2.Append1(t5) //在指针上调用值接收者函数，自动转为值，值不会变
	fmt.Printf("%v %v %v\n", size, err, t2)
	fmt.Println("")

	size, err = t2.Append2(t5) //在指针上调用指针接收者函数，值会变
	fmt.Printf("%v %v %v\n", size, err, t2)
	fmt.Println("")

	//1) 总结一下，GO其实不怎么区分指针和值，接收指针或者值的时候，传值或指针都ok，编译器会自动转换
	//2) 只有指针接收者版本的method可以修改原值，值接收者绝逼不行，改了也是白改，会丢弃掉

	//将TT作为io.Writer，注意这里绝逼是要修改原值的，因为要向buffer里面写啊，所以必须是*TT，即指针接收者，否则行为是错误的
	fmt.Fprintf(t3, "%s", "liujia")
	fmt.Printf("%v", t3)
}

func AboutAppend() {
	//append是编译器内建函数，类似于 append(slice []T, element... T)，但显然GO没有模板，T无法作为模板参数，所以这个函数只能是编译器提供喽

	//append两种用法
	a := []int{1, 2, 3}
	a = append(a, 4, 5, 6) //append会自动处理扩容，后面可以跟任意多个元素
	fmt.Println(a)

	b := []int{9, 8, 7}
	a = append(a, b...) //切片后再添加切片，必须用这种语法
	fmt.Println(a)
}

func AboutPrint() {
	//fmt.Print/Printf/Println  fmt.Sprint/Sprintln/Sprintf  fmt.Fprint/Fprintln/Fprintf
	//Print、Println、Printf 对应普通打印，打印一行，格式化打印，就是向标准输出打印，即os.Stdout
	//Sprintxxx返回字符串，就是讲打印内容扔到字符串里去
	//Fprintxxx第一个参数是io.Writer，即向某个可写的缓冲区打印，可以是os.Stdout,可以是文件，可以是网络等，或者就是一个缓冲区

	i := 1
	s := "liujia"
	fmt.Print(i, " ", s)
	fmt.Println(i, " ", s)

	s1 := fmt.Sprint(i, " ", s)
	s2 := fmt.Sprintf("%d %s, %v %v", i, s, i, s)

	fmt.Fprintln(os.Stdout, s1)
	fmt.Fprint(os.Stderr, s2)

	fmt.Println("\n")

	//注意，像%d，是根据数值的类型来决定格式的，这里不像C一样
	var x uint64 = 1<<64 - 1
	fmt.Printf("%d %x; %d %x\n", x, x, int64(x), int64(x)) //输出：18446744073709551615 ffffffffffffffff; -1 -1

	//%v是比较常见的选项，是按照通用格式打印，并且除了整数浮点bool字符串等，还可以支持数组，切片，map，结构体
	//当打印一个结构体时，带修饰的格式 %+v 会将结构体的域使用它们的名字进行注解，对于任意的值，格式 %#v 会按照完整的Go语法打印出该值
	type T struct {
		a int
		b float64
		c string
	}

	//对于结构体来说%+v或者%#v，都非常好！
	t := &T{7, -2.35, "abc\tdef"}
	fmt.Printf("%v\n", t)  //&{7 -2.35 abc	def}
	fmt.Printf("%+v\n", t) //&{a:7 b:-2.35 c:abc	def}
	fmt.Printf("%#v\n", t) //& main.T{a: 7, b: -2.35, c: "abc\tdef"}

	//%q用于输出带双引号的字符串
	s3 := "haha"
	fmt.Printf("%q \n", s3) //"haha"

	//%T用于打印值的类型，我的建议就是多用%T,，少用%t
	fmt.Printf("%T \n", t) //*main.T
	fmt.Printf("%t \n", t) //&{%!t(int=7) %!t(float64=-2.35) %!t(string=abc	def)}

	//对于结构体可以自定义String() string这个方法，用于%v打印时自定义格式，否则就是系统默认的，不过默认的挺好的
	/*
		func (t *T) String() string {
			return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
		}
	*/

	//可变参数列表
	//Printf类型为：func Printf(format string, v ...interface{}) (n int, err error)
	//假如我们的日志函数，类型为 func LogPrint(format string, v ...interface{}) error
	//内部调用Printf,则方式是 fmt.Printf(format, v...)
	//函数 LogPrint 内部， v 就像是一个类型为 []interface{} 的变量，但是如果其被传递给另一个可变参数的函数，
	//其就像是一个正常的参数列表。这里有一个对我们上面用到的函数 log.Println 的实现。
	//其将参数直接传递给 fmt.Sprintln 来做实际的格式化。
	//传递这种 v ...interface{} 要用 v...告诉编译器这是一个参数列表，否则会当做切片的

	//下面是个用可变参数列表的例子
	/*
			func Min(a ...int) int {
			min := int(^uint(0) >> 1) // largest int
			for _, i := range a {
				if i < min {
					min = i
				}
			}
			return min
		}
	*/
}

func JustTest() {
	//a是array，s是slice，m是map
	//这样初始化array和slice也是可以的，忽略掉前面的EnoneEioEinval等，直接用后面的string类型的value
	a := [...]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	s := []string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
	m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

	fmt.Printf("%T\n", a)
	fmt.Printf("%T\n", s)
	fmt.Printf("%T\n", m)

	//a = append(a, "liujia") //append只能往slice后面追加元素
	s = append(s, "liujai") //这里容易出错的地方就是append返回新的slice，新的slice才是append之后的。不过如果 直接写append(s, "liujai")没有s =，会编译报错

	fmt.Println("--------------------")
	for _, a_ := range a {
		fmt.Println(a_)
	}

	fmt.Println("--------------------")
	for _, s_ := range s {
		fmt.Println(s_)
	}

	for k, v := range m {
		fmt.Println(k, " : ", v)
	}
}

func AboutNew() {
	//new(T)，为T分配一块内存，并置0(意思是设置默认的零值，对int来说是0，对string来说是“”，对sync.Mutex是没有上锁的mutex)
	//并返回*T，即T的指针
	a := new(int) //a是*int
	var b *int = new(int)
	fmt.Printf("%t\n", a)
	fmt.Printf("%T\n", a)

	fmt.Printf("%t\n", *a)
	fmt.Printf("%T\n", *b)
}

type T struct {
	index int
	path  string
	lock  sync.Mutex
	//其它的成员，或者称之为field，类型可以是基本类型或者自定义类型
}

func NewT1(index int, path string) *T {
	t := new(T)
	t.index = index
	t.path = path
	return t
}

func NewT2(index int, path string) *T {
	if index < 0 {
		return nil
	}

	//第一种方式，T的成员顺序并且全部列出，并赋予初始值，T{param1, param2, ...} -- 这种方式称为复合文字 composite literal，GO中返回局部变量的地址没有问题，不像c++哦
	//t := T{index, path, sync.Mutex{}}
	//return &t

	//第二种方式，用filed:value方式显式的初始化T，这种方式不容易出错，并且没有列出来的就是直接用初值好了，不像上面那个必须都列出来
	//t := T{index: index, path: path}
	//return &t

	//更直接的方式
	return &T{index: index, path: path}
}

func AboutConstructor() {
	//new(T)只能返回零值的对象，GO中的初始化构造器为了更为精确的控制对象的构造
	//通常对于T来说，其初始化构造器名为 "NewT(一些参数) *T"
	//直接看上面的T NewT1 NewT2吧
}

func AboutMake() {
	//make只用于分配 slice  map  和 channel，并且返回的是对象，而不是对象的指针，这点和new是不一样的！

	//slice
	s := make([]int, 10, 100) //[]int类型的slice，长度为10，容量100，指向后面对应的数组的前10项。此时s就是有10个0的slice
	fmt.Println("--------------------")
	for _, s_ := range s {
		fmt.Println(s_)
	}

	s1 := make([]int, 5) //[]int类型的slice，长度为5，容量5
	fmt.Println("--------------------")
	for _, s_ := range s1 {
		fmt.Println(s_)
	}

	fmt.Println(cap(s1)) //cap(s1)为5，即容量是5

	//map
	fmt.Println("--------------------")
	m := make(map[int]string)
	m[1] = "liu"
	m[2] = "jia"
	for k, v := range m {
		fmt.Println(k, " : ", v)
	}

	//chan
	c := make(chan bool, 5) //5是缓冲区的大小，最多容纳5个，超过之后往chan里塞就会阻塞
	//c := make(chan bool)    //创建一个无缓冲的bool型Channel，无缓冲区的chan可以兼做通信和同步
	//c <- x      //向一个Channel发送一个值
	//<-c         //从一个Channel中接收一个值
	//x = <-c     //从Channel c接收一个值并将其存储到x中
	//x, ok = <-c //从Channel接收一个值，如果channel关闭了，那么ok将被置为false

	c <- false
	c <- true
	close(c)

	//这种for循环，从range chan里接收，如果chan里没有值就会一直等待，for退出的条件的是chan被close

	for x := range c {
		fmt.Println(x)
	}

	//上面的for循环相当于下面的代码
	/*
		for {
			x, ok := <-c
			if !ok {
				break
			}
			fmt.Println(x)
		}*/
}

func AboutArray() {
	//数组
	//1:数组是slice的底层构件
	//2:数组是值，而不是引用，就是说传递数组到函数参数，得到的是拷贝，并且数组赋值也是拷贝了原数组的所有值
	//3:数组的大小是其类型的一部分，就是说[10]int和[20]int是不同的类型
	//4:多用切片吧
}

func AboutSlice() {
	//切片
	//切边是对底层数组的封装，GO中对数组的编程通常是通过slice完成的
	//切片是引用的(可以看做是底层数组的一个指针)，本质上slice是个结构是按值传递的，但是对底层数组的是基于引用的，就是还是引用了同一个底层数组
	//cap(s) 返回容量
	//len(s) 返回长度
	//s = append(s, ...) 添加元素
	//s = s[n:m] 返回 s[n] ... s[m-1]组成的新切片
	//copy(newslice, slice) 拷贝slice的值到另一个newslice中

	//apppend函数的一个示意实现
	/*
		func Append(slice, data []byte) []byte {
			l := len(slice)
			if l+len(data) > cap(slice) { // reallocate
				// Allocate double what's needed, for future growth.
				newSlice := make([]byte, (l+len(data))*2)
				// The copy function is predeclared and works for any slice type.
				copy(newSlice, slice)
				slice = newSlice
			}
			slice = slice[0 : l+len(data)]
			for i, c := range data {
				slice[l+i] = c
			}
			return slice
		}
	*/

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3, 4, 5, 6}
	s3 := make([]int, 5, 10)
	copy(s3, s1)
	fmt.Println(s3) //s3为1 2 3 0 0，因为len(s1) < len(s3)，所以只覆盖了前面的元素，后面的保持不变

	copy(s3, s2)
	fmt.Println(s3) //s3为1 2 3 4 5，因为len(s2) > len(s3)，所以只拷贝了len(s3)个元素

	//二维数组和切片
	type Transform [3][3]float64 // A 3x3 array, really an array of arrays.
	type LinesOfText [][]byte    // A slice of byte slices.

	//构造二维切片的一种方式，先分配一个二维的（行数那么多个），再分配每一个行的
	/*
		// Allocate the top-level slice.
		picture := make([][]uint8, YSize) // One row per unit of y.
		// Loop over the rows, allocating the slice for each row.
		for i := range picture {
			picture[i] = make([]uint8, XSize)
		}
	*/

	//第二种方式，先分配一个二维的（行数那么多的），再分配一个总的，然后每一个是一个子切片，这样效率高，但是不灵活（列数不能变）
	/*
		// Allocate the top-level slice, the same as before.
		picture := make([][]uint8, YSize) // One row per unit of y.
		// Allocate one large slice to hold all the pixels.
		pixels := make([]uint8, XSize*YSize) // Has type []uint8 even though picture is [][]uint8.
		// Loop over the rows, slicing each row from the front of the remaining pixels slice.
		for i := range picture {
			picture[i], pixels = pixels[:XSize], pixels[XSize:] //这样写法很好，就是每次从前面slice(动词)出XSize个，给picture[i]，然后重新设置pixels
		}
	*/
}

func AboutDefer() {
	//如果一个函数注册了多个defer，按照FIFO的顺序执行，即最先注册的最先执行
	//并且defer注册的函数，参数是在defer注册的时候求值的，而不是真正执行时求值的
	//下面的j，会打印出defer注册时的100，而不是defer真正执行时的200

	//FIFO, 打印出 4 3 2 1 0
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)
	}

	//打印出100
	j := 100
	defer fmt.Println(j)
	j = 200

	//defer注册一个函数，打印出200，因为defer注册的时候j已经是200了
	defer func() {
		fmt.Println(j)
	}()
}

func AboutString() {
	//字符串是utf8编码，用for range循环变量，得到的不是byte，而是rune，即unicode码点，并且index不是连续的
	s := "刘佳is中国人"
	for i, u := range s {
		fmt.Printf("%#U --- %d\n", u, i)
	}

	//s[0], s[1] = s[1], s[0] //这样不行，不能写
	fmt.Println(s[1]) //这样可以，可以读
}

func AboutArrayAndSlice() {
	//切片，声明并赋值
	s1 := []int{1, 2, 3, 4, 5, 6}
	//s1 := "liujia" //字符串不能这么赋值，就是说字符串不能直接s[i]这么访问
	for i, j := 0, len(s1)-1; i < j; i, j = i+1, j-1 {
		s1[i], s1[j] = s1[j], s1[i]
	}
	fmt.Println(s1)
}

func AboutMap() {
	//map，直接声明并赋值
	m := map[int]string{1: "liu", 2: "jia"}
	for k, v := range m {
		fmt.Printf("%d --- %s\n", k, v)
	}

	//new出一个map，并且添加元素
	m2 := make(map[string]int)
	m2["liu"] = 1
	m2["jia"] = 2

	//用map[key]来获取key对应的value
	//如果key不存在，map[key]返回value的零值
	fmt.Println(m2["liu"])  // 打印 1
	fmt.Println(m2["haha"]) //打印 0，因为key不存在，返回int的零值

	//判断map是否有某个key
	if v, ok := m2["haha"]; ok {
		fmt.Printf("key exist, %b \n", ok)
	} else {
		fmt.Printf("key not exist, %b \n", v)
	}

	//更精炼的代码
	_, present := m2["haha"] //present 是bool，返回key是否在map中
	fmt.Println(present)

	//删除某一个key、value
	delete(m2, "liu")

	//遍历map
	for k, v := range m2 {
		fmt.Printf("%s --- %d\n", k, v)
	}
}

func AboutSwitch() {
	var c byte = 'z'
	if c < 'z' {
		fmt.Println("here")
	}

	//%T和%t不一样
	fmt.Printf("-----  %T \n", c) // uint8  类型
	fmt.Printf("-----  %t \n", c) // %!t(uint8=122)  类型和值都有

	/*
		var c1 byte

		//switch默认是break的，不用显式去写，如果不想break，想几个合并，用fallthrough
		//switch后面不跟特定变量时，case语句是boo表达式
		switch {
		case '0' <= c && c <= '9':
			c1 = c - '0'
			fallthrough
		case 'a' <= c && c <= 'f':
			c1 = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c1 = c - 'A' + 10
		default:
			fmt.Println("default case")
		}
	*/

	//switch后面跟特定变量时，case语句单独的或者,分开的一些变量取值

	//这里顺便演示一下break到标号的用法，switch不能直接跳到标号，外面必须有个作用域，如这里的for{}，
	//相当于break LOOP的作用在于直接从for作用域内部的switch直接跳出for作用域
	//常见的编译错误就是1）LOOP标号和下面的for之间有其它语句 2）switch外面就是LOOP，而没有其它的作用域包着
LOOP:
	for {
		fmt.Println("enter for")
		switch c {
		case 'w':
			fmt.Println("'w'")
		case 'x':
			fmt.Println("'x'")
			//continue LOOP //也可以continue加标号，相当于继续循环，但这里例子这里绝对的死循环
		case 'z':
			fmt.Println("'z'")
			break LOOP
		default:
			fmt.Println("default case")
			break LOOP
		}
		fmt.Println("leave for")
	}

	fmt.Println("after switch")

	//一个接口变量的动态类型
	//switch也可以用来判断interface变量的具体类型 x.(type)
	z := 'z' // z is byte type
	z1 := interface{}(z)
	switch t := z1.(type) { //这个t变量，在每一个case实际上是不同的变量，而且转为了对应的类型
	case bool:
		fmt.Println("bool", t)
	case int:
		fmt.Println("int", t)
	case *byte:
		fmt.Println("byte", *t)
	case byte:
		fmt.Println("byte", t)
	case rune:
		fmt.Println("rune", t)
	}
}
