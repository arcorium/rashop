// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: mailer/v1/command.proto

package mailerv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_mailer_v1_command_proto protoreflect.FileDescriptor

var file_mailer_v1_command_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x72, 0x61, 0x73, 0x68, 0x6f,
	0x70, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x6d, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xba, 0x01, 0x0a,
	0x14, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x21, 0x2e,
	0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x22, 0x2e, 0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x23,
	0x2e, 0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x6d, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61, 0x69,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xc2, 0x01, 0x0a, 0x14, 0x63, 0x6f,
	0x6d, 0x2e, 0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x42, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x72, 0x63, 0x6f, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x72, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x52, 0x4d, 0x58, 0xaa, 0x02, 0x10, 0x52, 0x61, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x4d, 0x61,
	0x69, 0x6c, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x10, 0x52, 0x61, 0x73, 0x68, 0x6f, 0x70,
	0x5c, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1c, 0x52, 0x61, 0x73,
	0x68, 0x6f, 0x70, 0x5c, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12, 0x52, 0x61, 0x73, 0x68,
	0x6f, 0x70, 0x3a, 0x3a, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_mailer_v1_command_proto_goTypes = []interface{}{
	(*SendMailRequest)(nil),    // 0: rashop.mailer.v1.SendMailRequest
	(*DeleteMailRequest)(nil),  // 1: rashop.mailer.v1.DeleteMailRequest
	(*SendMailResponse)(nil),   // 2: rashop.mailer.v1.SendMailResponse
	(*DeleteMailResponse)(nil), // 3: rashop.mailer.v1.DeleteMailResponse
}
var file_mailer_v1_command_proto_depIdxs = []int32{
	0, // 0: rashop.mailer.v1.MailerCommandService.Send:input_type -> rashop.mailer.v1.SendMailRequest
	1, // 1: rashop.mailer.v1.MailerCommandService.Delete:input_type -> rashop.mailer.v1.DeleteMailRequest
	2, // 2: rashop.mailer.v1.MailerCommandService.Send:output_type -> rashop.mailer.v1.SendMailResponse
	3, // 3: rashop.mailer.v1.MailerCommandService.Delete:output_type -> rashop.mailer.v1.DeleteMailResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mailer_v1_command_proto_init() }
func file_mailer_v1_command_proto_init() {
	if File_mailer_v1_command_proto != nil {
		return
	}
	file_mailer_v1_message_command_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mailer_v1_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mailer_v1_command_proto_goTypes,
		DependencyIndexes: file_mailer_v1_command_proto_depIdxs,
	}.Build()
	File_mailer_v1_command_proto = out.File
	file_mailer_v1_command_proto_rawDesc = nil
	file_mailer_v1_command_proto_goTypes = nil
	file_mailer_v1_command_proto_depIdxs = nil
}
