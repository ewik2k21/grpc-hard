// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: spot_instrument_service.proto

package pkg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_spot_instrument_service_proto protoreflect.FileDescriptor

const file_spot_instrument_service_proto_rawDesc = "" +
	"\n" +
	"\x1dspot_instrument_service.proto\x12\x1aspot_instrument_service_v1\x1a&spot_instrument_service_messages.proto2\x87\x01\n" +
	"\x15SpotInstrumentService\x12n\n" +
	"\vViewMarkets\x12..spot_instrument_service_v1.ViewMarketsRequest\x1a/.spot_instrument_service_v1.ViewMarketsResponseB#Z!github.com/ewik2k21/grpc-hard/pkgb\x06proto3"

var file_spot_instrument_service_proto_goTypes = []any{
	(*ViewMarketsRequest)(nil),  // 0: spot_instrument_service_v1.ViewMarketsRequest
	(*ViewMarketsResponse)(nil), // 1: spot_instrument_service_v1.ViewMarketsResponse
}
var file_spot_instrument_service_proto_depIdxs = []int32{
	0, // 0: spot_instrument_service_v1.SpotInstrumentService.ViewMarkets:input_type -> spot_instrument_service_v1.ViewMarketsRequest
	1, // 1: spot_instrument_service_v1.SpotInstrumentService.ViewMarkets:output_type -> spot_instrument_service_v1.ViewMarketsResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_spot_instrument_service_proto_init() }
func file_spot_instrument_service_proto_init() {
	if File_spot_instrument_service_proto != nil {
		return
	}
	file_spot_instrument_service_messages_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_spot_instrument_service_proto_rawDesc), len(file_spot_instrument_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_spot_instrument_service_proto_goTypes,
		DependencyIndexes: file_spot_instrument_service_proto_depIdxs,
	}.Build()
	File_spot_instrument_service_proto = out.File
	file_spot_instrument_service_proto_goTypes = nil
	file_spot_instrument_service_proto_depIdxs = nil
}
