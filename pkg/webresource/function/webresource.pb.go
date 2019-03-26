// Code generated by protoc-gen-go. DO NOT EDIT.
// source: webresource.proto

package function

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WebResource struct {
	OrgID                string            `protobuf:"bytes,1,opt,name=OrgID" json:"OrgID,omitempty"`
	Spec                 string            `protobuf:"bytes,2,opt,name=Spec" json:"Spec,omitempty"`
	RelativePath         string            `protobuf:"bytes,3,opt,name=RelativePath" json:"RelativePath,omitempty"`
	NormalisedKey        string            `protobuf:"bytes,4,opt,name=NormalisedKey" json:"NormalisedKey,omitempty"`
	MD5Hash              string            `protobuf:"bytes,6,opt,name=MD5Hash" json:"MD5Hash,omitempty"`
	SourcePath           string            `protobuf:"bytes,7,opt,name=SourcePath" json:"SourcePath,omitempty"`
	Derivatives          map[string]string `protobuf:"bytes,5,rep,name=Derivatives" json:"Derivatives,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *WebResource) Reset()         { *m = WebResource{} }
func (m *WebResource) String() string { return proto.CompactTextString(m) }
func (*WebResource) ProtoMessage()    {}
func (*WebResource) Descriptor() ([]byte, []int) {
	return fileDescriptor_webresource_5d1c44529e16a0bb, []int{0}
}
func (m *WebResource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebResource.Unmarshal(m, b)
}
func (m *WebResource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebResource.Marshal(b, m, deterministic)
}
func (dst *WebResource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebResource.Merge(dst, src)
}
func (m *WebResource) XXX_Size() int {
	return xxx_messageInfo_WebResource.Size(m)
}
func (m *WebResource) XXX_DiscardUnknown() {
	xxx_messageInfo_WebResource.DiscardUnknown(m)
}

var xxx_messageInfo_WebResource proto.InternalMessageInfo

func (m *WebResource) GetOrgID() string {
	if m != nil {
		return m.OrgID
	}
	return ""
}

func (m *WebResource) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *WebResource) GetRelativePath() string {
	if m != nil {
		return m.RelativePath
	}
	return ""
}

func (m *WebResource) GetNormalisedKey() string {
	if m != nil {
		return m.NormalisedKey
	}
	return ""
}

func (m *WebResource) GetMD5Hash() string {
	if m != nil {
		return m.MD5Hash
	}
	return ""
}

func (m *WebResource) GetSourcePath() string {
	if m != nil {
		return m.SourcePath
	}
	return ""
}

func (m *WebResource) GetDerivatives() map[string]string {
	if m != nil {
		return m.Derivatives
	}
	return nil
}

func init() {
	proto.RegisterType((*WebResource)(nil), "protos.WebResource")
	proto.RegisterMapType((map[string]string)(nil), "protos.WebResource.DerivativesEntry")
}

func init() { proto.RegisterFile("webresource.proto", fileDescriptor_webresource_5d1c44529e16a0bb) }

var fileDescriptor_webresource_5d1c44529e16a0bb = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x4f, 0x4d, 0x2a,
	0x4a, 0x2d, 0xce, 0x2f, 0x2d, 0x4a, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03,
	0x53, 0xc5, 0x4a, 0x87, 0x99, 0xb8, 0xb8, 0xc3, 0x53, 0x93, 0x82, 0xa0, 0xb2, 0x42, 0x22, 0x5c,
	0xac, 0xfe, 0x45, 0xe9, 0x9e, 0x2e, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90,
	0x10, 0x17, 0x4b, 0x70, 0x41, 0x6a, 0xb2, 0x04, 0x13, 0x58, 0x10, 0xcc, 0x16, 0x52, 0xe2, 0xe2,
	0x09, 0x4a, 0xcd, 0x49, 0x2c, 0xc9, 0x2c, 0x4b, 0x0d, 0x48, 0x2c, 0xc9, 0x90, 0x60, 0x06, 0xcb,
	0xa1, 0x88, 0x09, 0xa9, 0x70, 0xf1, 0xfa, 0xe5, 0x17, 0xe5, 0x26, 0xe6, 0x64, 0x16, 0xa7, 0xa6,
	0x78, 0xa7, 0x56, 0x4a, 0xb0, 0x80, 0x15, 0xa1, 0x0a, 0x0a, 0x49, 0x70, 0xb1, 0xfb, 0xba, 0x98,
	0x7a, 0x24, 0x16, 0x67, 0x48, 0xb0, 0x81, 0xe5, 0x61, 0x5c, 0x21, 0x39, 0x2e, 0xae, 0x60, 0xb0,
	0xbb, 0xc0, 0x36, 0xb0, 0x83, 0x25, 0x91, 0x44, 0x84, 0xdc, 0xb8, 0xb8, 0x5d, 0x52, 0x8b, 0x32,
	0xcb, 0xc0, 0x36, 0x16, 0x4b, 0xb0, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0xa9, 0x40, 0xbc, 0x58, 0xac,
	0x87, 0xe4, 0x2f, 0x3d, 0x24, 0x65, 0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0xc8, 0x1a, 0xa5, 0xec,
	0xb8, 0x04, 0xd0, 0x15, 0x08, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x42, 0xc3, 0x01, 0xc4, 0x04,
	0x85, 0x4d, 0x59, 0x62, 0x4e, 0x69, 0x2a, 0x34, 0x18, 0x20, 0x1c, 0x2b, 0x26, 0x0b, 0xc6, 0x24,
	0x48, 0x68, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x20, 0x76, 0x0d, 0xc8, 0x69, 0x01, 0x00,
	0x00,
}
