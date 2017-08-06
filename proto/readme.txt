请按照哥总结的步骤来

1）安装protcbuf 3.x，网上教程一大堆，自己鼓捣吧。安装完后，确认运行"protoc --version"可以正确运行

2) 安装proto_gen_go，运行"go get github.com/golang/protobuf/protoc-gen-go"，然后去其目录，{GOPATH}/sr/github.com/golang/protobuf下面，然后
   运行"go build" & "go install"，其会被安装到{GOPATH}/bin下面。当然了，最好将{GOPATH}/bin加入到环境变量中去，这样方便你我他。

3）cd protoc_gen_msg目录，go build & go install，编译并安装protoc_gen_msg

4）编写proto文件，当然了，推荐用proto3的语法

5）编写sh脚本，这个脚本做两件事，一件事是调用protoc_gen_go将proto文件编译为.go文件，另一件事是调用protoc_gen_msg获取proto的meta信息并生成对应的msgid.go

6) 生成的msgid.go并不能直接用，会出现循环依赖。其实也不完全是，对于session_event确实会的，因为session_event包含了socket库，但是socket库又依赖了session_event
   最好还是将msgid.go改名为msgid.go_，同时自己在需要的地方，register需要的event就好。