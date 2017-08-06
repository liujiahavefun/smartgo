#!/usr/bin/env bash
outdir=session_event
plugindir=${GOPATH}/bin/

#outdir不存在则创建
if [ ! -d "${outdir}" ]; then
    /bin/mkdir "${outdir}"
fi

protoc --plugin=protoc-gen-go=${plugindir}/protoc-gen-go --go_out ${outdir} --proto_path "." session_event.proto
protoc --plugin=protoc-gen-msg=${plugindir}/protoc-gen-msg --msg_out=msgid.go:${outdir} --proto_path "." session_event.proto