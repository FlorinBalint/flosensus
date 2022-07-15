# This must be the first line in Makefile
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

# TODO: Make everything relative to the binary

BUILD=${mkfile_dir}build
GOBIN=${BUILD}/bin
GOSRC=${BUILD}/src
BINARY=flosensus
GOPROTO=${GOSRC}/flosensus/raft/proto
GO_PROTO_MODULE="github.com/FlorinBalint/flosensus/raft/proto"

.PHONY: protos deps build run clean test

${GOPROTO}:
	mkdir -p ${GOPROTO}

${GOPROTO}/go.mod: ${GOPROTO}
		cd ${GOPROTO} && go mod init ${GO_PROTO_MODULE} && cd -

deps:
	go get github.com/FlorinBalint/flosensus/raft/proto@v0.1.0
	go mod download google.golang.org/grpc

protos:
	protoc -I=${mkfile_dir}proto  --go_out=${GOPROTO} --go_opt=paths=source_relative \
		${mkfile_dir}proto/config.proto
	protoc -I=${mkfile_dir}proto  --go_out=${GOPROTO} --go_opt=paths=source_relative \
		--go-grpc_out=${GOPROTO} --go-grpc_opt=paths=source_relative \
		${mkfile_dir}proto/raft.proto

build: ${GOPROTO}/go.mod protos deps
	go build ${LDFLAGS} -o ${GOBIN}/${BINARY} cmd/server.go

run:
	cd ${GOBIN} && ./${BINARY}

		
test: ${GOPROTO}/go.mod config_proto
	go test ${mkfile_dir}/...

clean:
	rm -rf ${BUILD}
