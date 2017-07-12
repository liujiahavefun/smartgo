#!/usr/bin/env bash
outdir=test
plugindir=/Users/liujia/go/bin

#outdir不存在则创建
if [ ! -d "${outdir}" ]; then
    /bin/mkdir "${outdir}"
fi

protoc --plugin=protoc-gen-go=${plugindir}/protoc-gen-go --go_out ${outdir} --proto_path "." test.proto
protoc --plugin=protoc-gen-msg=../libs/protoc-gen-msg/protoc-gen-msg --msg_out=msgid.go:${outdir} --proto_path "." test.proto