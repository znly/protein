PATH_PROTOS=${GOPATH}/src/github.com/znly/protein/protobuf
PATH_SCHEMAS=${PATH_PROTOS}/..
PATH_SCHEMAS_TEST=${PATH_PROTOS}/test

SCHEMAS=                   \
	${PATH_PROTOS}/protobuf_payload.proto \
	${PATH_PROTOS}/protobuf_schema.proto

SCHEMAS_TEST=             \
	${PATH_PROTOS}/test_schema.proto     \
	${PATH_PROTOS}/test_schema_xxx.proto

all: protobuf protobuf-test

protobuf-clean:
	rm -f ${PATH_SCHEMAS}/*.pb.go

protobuf: protobuf-clean
	mkdir -p ${PATH_SCHEMAS}
	protoc --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ \
		   --proto_path=${GOPATH}/src/github.com/gogo/                   \
		   --proto_path=${PATH_PROTOS}/                                 \
		   --gogofaster_out=${PATH_SCHEMAS}/                             \
		   ${SCHEMAS}
	sed -i '' 's#google_protobuf "google/protobuf"#google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"#g' ${PATH_SCHEMAS}/*.pb.go
	sed -i '' 's#"protobuf/gogoproto"#"github.com/gogo/protobuf/gogoproto"#g' ${PATH_SCHEMAS}/*.pb.go

protobuf-test-clean:
	rm -f ${PATH_SCHEMAS_TEST}/*.pb.go

protobuf-test: protobuf-test-clean
	mkdir -p ${PATH_SCHEMAS_TEST}
	protoc --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ \
		   --proto_path=${GOPATH}/src/github.com/gogo/                   \
		   --proto_path=${PATH_PROTOS}/                          \
		   --gogofaster_out=${PATH_SCHEMAS_TEST}/                      \
		   ${SCHEMAS_TEST}
	sed -i '' 's#google_protobuf "google/protobuf"#google_protobuf "github.com/gogo/protobuf/types"#g' ${PATH_SCHEMAS_TEST}/*.pb.go
	sed -i '' 's#"protobuf/gogoproto"#"github.com/gogo/protobuf/gogoproto"#g' ${PATH_SCHEMAS_TEST}/*.pb.go

