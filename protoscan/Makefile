PROTOSCAN_PATH=${GOPATH}/src/github.com/znly/protein/protoscan

all: protobuf

protobuf-clean:
	rm -f ${PROTOSCAN_PATH}/*.pb.go

protobuf: protobuf-clean
	protoc --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ \
		   --proto_path=${GOPATH}/src/github.com/gogo/ \
		   --proto_path=${PROTOSCAN_PATH}/protobuf/ \
		   --gogofaster_out=${PROTOSCAN_PATH}/ \
		   ${PROTOSCAN_PATH}/protobuf/protobuf_schema.proto
