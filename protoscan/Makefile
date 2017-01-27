PROTOSCAN_PATH=${GOPATH}/src/github.com/znly/protein/protoscan
PROTOBUF_SCHEMA=protobuf_schema

all: protobuf

protobuf-clean:
	rm -f ${PROTOSCAN_PATH}/*.pb.go

protobuf: protobuf-clean
	protoc --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ \
		   --proto_path=${GOPATH}/src/github.com/gogo/ \
		   --proto_path=${PROTOSCAN_PATH}/protobuf/ \
		   --gogofaster_out=${PROTOSCAN_PATH}/ \
		   ${PROTOSCAN_PATH}/protobuf/${PROTOBUF_SCHEMA}.proto
	sed -i '' 's#"google/protobuf"#"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"#g' ${PROTOSCAN_PATH}/${PROTOBUF_SCHEMA}.pb.go
	sed -i '' 's#"protobuf/gogoproto"#"github.com/gogo/protobuf/gogoproto"#g' ${PROTOSCAN_PATH}/${PROTOBUF_SCHEMA}.pb.go
