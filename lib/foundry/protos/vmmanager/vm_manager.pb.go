// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: vm_manager.proto

package vmmanager

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cpu    int32 `protobuf:"varint,1,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory int32 `protobuf:"varint,2,opt,name=memory,proto3" json:"memory,omitempty"`
}

func (x *ClusterRequest) Reset() {
	*x = ClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterRequest) ProtoMessage() {}

func (x *ClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterRequest.ProtoReflect.Descriptor instead.
func (*ClusterRequest) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{0}
}

func (x *ClusterRequest) GetCpu() int32 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *ClusterRequest) GetMemory() int32 {
	if x != nil {
		return x.Memory
	}
	return 0
}

type ClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ClusterResponse) Reset() {
	*x = ClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterResponse) ProtoMessage() {}

func (x *ClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterResponse.ProtoReflect.Descriptor instead.
func (*ClusterResponse) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{1}
}

func (x *ClusterResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ClusterID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ClusterID) Reset() {
	*x = ClusterID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterID) ProtoMessage() {}

func (x *ClusterID) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterID.ProtoReflect.Descriptor instead.
func (*ClusterID) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{2}
}

func (x *ClusterID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{3}
}

func (x *StatusResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type VMID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *VMID) Reset() {
	*x = VMID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VMID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VMID) ProtoMessage() {}

func (x *VMID) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VMID.ProtoReflect.Descriptor instead.
func (*VMID) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{4}
}

func (x *VMID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type VMResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VMResponse) Reset() {
	*x = VMResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vm_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VMResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VMResponse) ProtoMessage() {}

func (x *VMResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vm_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VMResponse.ProtoReflect.Descriptor instead.
func (*VMResponse) Descriptor() ([]byte, []int) {
	return file_vm_manager_proto_rawDescGZIP(), []int{5}
}

func (x *VMResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_vm_manager_proto protoreflect.FileDescriptor

var file_vm_manager_proto_rawDesc = []byte{
	0x0a, 0x10, 0x76, 0x6d, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x09, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0x3a, 0x0a,
	0x0e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x63, 0x70,
	0x75, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x22, 0x25, 0x0a, 0x0f, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x1b, 0x0a, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x2a, 0x0a,
	0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x16, 0x0a, 0x04, 0x56, 0x4d, 0x49,
	0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x20, 0x0a, 0x0a, 0x56, 0x4d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x32, 0x86, 0x04, 0x0a, 0x09, 0x56, 0x4d, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x72,
	0x79, 0x12, 0x48, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x0d, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x76,
	0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x49, 0x44, 0x1a, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x41, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x14, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x40, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x70, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x14, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x72, 0x74, 0x56, 0x4d, 0x12,
	0x0f, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x56, 0x4d, 0x49, 0x44,
	0x1a, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a,
	0x06, 0x53, 0x74, 0x6f, 0x70, 0x56, 0x4d, 0x12, 0x0f, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x56, 0x4d, 0x49, 0x44, 0x1a, 0x19, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x1a, 0x2e, 0x76, 0x6d, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x56, 0x4d, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x0f, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x56,
	0x4d, 0x49, 0x44, 0x1a, 0x15, 0x2e, 0x76, 0x6d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x56, 0x4d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x38, 0x5a, 0x36,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x64, 0x73, 0x6c, 0x61,
	0x62, 0x73, 0x2f, 0x6b, 0x61, 0x74, 0x61, 0x6e, 0x61, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x66, 0x6f,
	0x75, 0x6e, 0x64, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x6d, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vm_manager_proto_rawDescOnce sync.Once
	file_vm_manager_proto_rawDescData = file_vm_manager_proto_rawDesc
)

func file_vm_manager_proto_rawDescGZIP() []byte {
	file_vm_manager_proto_rawDescOnce.Do(func() {
		file_vm_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_vm_manager_proto_rawDescData)
	})
	return file_vm_manager_proto_rawDescData
}

var file_vm_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_vm_manager_proto_goTypes = []interface{}{
	(*ClusterRequest)(nil),  // 0: vmmanager.ClusterRequest
	(*ClusterResponse)(nil), // 1: vmmanager.ClusterResponse
	(*ClusterID)(nil),       // 2: vmmanager.ClusterID
	(*StatusResponse)(nil),  // 3: vmmanager.StatusResponse
	(*VMID)(nil),            // 4: vmmanager.VMID
	(*VMResponse)(nil),      // 5: vmmanager.VMResponse
}
var file_vm_manager_proto_depIdxs = []int32{
	0, // 0: vmmanager.VMFoundry.CreateCluster:input_type -> vmmanager.ClusterRequest
	2, // 1: vmmanager.VMFoundry.DeleteCluster:input_type -> vmmanager.ClusterID
	2, // 2: vmmanager.VMFoundry.StartCluster:input_type -> vmmanager.ClusterID
	2, // 3: vmmanager.VMFoundry.StopCluster:input_type -> vmmanager.ClusterID
	4, // 4: vmmanager.VMFoundry.StartVM:input_type -> vmmanager.VMID
	4, // 5: vmmanager.VMFoundry.StopVM:input_type -> vmmanager.VMID
	2, // 6: vmmanager.VMFoundry.ClusterInfo:input_type -> vmmanager.ClusterID
	4, // 7: vmmanager.VMFoundry.VMInfo:input_type -> vmmanager.VMID
	1, // 8: vmmanager.VMFoundry.CreateCluster:output_type -> vmmanager.ClusterResponse
	3, // 9: vmmanager.VMFoundry.DeleteCluster:output_type -> vmmanager.StatusResponse
	3, // 10: vmmanager.VMFoundry.StartCluster:output_type -> vmmanager.StatusResponse
	3, // 11: vmmanager.VMFoundry.StopCluster:output_type -> vmmanager.StatusResponse
	3, // 12: vmmanager.VMFoundry.StartVM:output_type -> vmmanager.StatusResponse
	3, // 13: vmmanager.VMFoundry.StopVM:output_type -> vmmanager.StatusResponse
	1, // 14: vmmanager.VMFoundry.ClusterInfo:output_type -> vmmanager.ClusterResponse
	5, // 15: vmmanager.VMFoundry.VMInfo:output_type -> vmmanager.VMResponse
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_vm_manager_proto_init() }
func file_vm_manager_proto_init() {
	if File_vm_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vm_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vm_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vm_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vm_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vm_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VMID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vm_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VMResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vm_manager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vm_manager_proto_goTypes,
		DependencyIndexes: file_vm_manager_proto_depIdxs,
		MessageInfos:      file_vm_manager_proto_msgTypes,
	}.Build()
	File_vm_manager_proto = out.File
	file_vm_manager_proto_rawDesc = nil
	file_vm_manager_proto_goTypes = nil
	file_vm_manager_proto_depIdxs = nil
}
