#!/usr/bin/env bash
outdir=test_event
plugindir=/Users/liujia/go/bin

#outdir不存在则创建
if [ ! -d "${outdir}" ]; then
    /bin/mkdir "${outdir}"
fi

protoc --plugin=protoc-gen-go=${GOPATH}/bin//protoc-gen-go --go_out ${outdir} --proto_path "." test_event.proto
protoc --plugin=protoc-gen-msg=${GOPATH}/bin//protoc-gen-msg --msg_out=msgid.go:${outdir} --proto_path "." test_event.proto