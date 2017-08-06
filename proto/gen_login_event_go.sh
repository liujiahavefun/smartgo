#!/usr/bin/env bash

#!/usr/bin/env bash
outdir=login_event
plugindir=/Users/liujia/go/bin

#outdir不存在则创建
if [ ! -d "${outdir}" ]; then
    /bin/mkdir "${outdir}"
fi

protoc --plugin=protoc-gen-go=${GOPATH}/bin//protoc-gen-go --go_out ${outdir} --proto_path "." login_event.proto
protoc --plugin=protoc-gen-msg=${GOPATH}/bin//protoc-gen-msg --msg_out=msgid.go:${outdir} --proto_path "." login_event.proto