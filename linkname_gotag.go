// Code generated by "linkname-gen -symbol github.com/gogo/protobuf/protoc-gen-gogo/generator.(*Generator).goTag -def func goTag(*generator.Generator, *generator.Descriptor, *descriptor.FieldDescriptorProto, string) string"; DO NOT EDIT.

package protein

import (
	_ "unsafe"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

// This avoids some weird goimports issues that I've been facing on linux/amd64
// since Go1.9+.
var (
	_ *generator.Generator
	_ *generator.Descriptor
	_ *descriptor.FieldDescriptorProto
)

//go:linkname goTag github.com/gogo/protobuf/protoc-gen-gogo/generator.(*Generator).goTag
func goTag(*generator.Generator, *generator.Descriptor, *descriptor.FieldDescriptorProto, string) string
