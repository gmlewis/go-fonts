// Code generated by protoc-gen-go. DO NOT EDIT.
// source: glyphs.proto

package glyphs

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

type Glyphs struct {
	Glyphs               []*Glyph `protobuf:"bytes,1,rep,name=glyphs,proto3" json:"glyphs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Glyphs) Reset()         { *m = Glyphs{} }
func (m *Glyphs) String() string { return proto.CompactTextString(m) }
func (*Glyphs) ProtoMessage()    {}
func (*Glyphs) Descriptor() ([]byte, []int) {
	return fileDescriptor_glyphs_5b62025fda4d7cef, []int{0}
}
func (m *Glyphs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Glyphs.Unmarshal(m, b)
}
func (m *Glyphs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Glyphs.Marshal(b, m, deterministic)
}
func (dst *Glyphs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Glyphs.Merge(dst, src)
}
func (m *Glyphs) XXX_Size() int {
	return xxx_messageInfo_Glyphs.Size(m)
}
func (m *Glyphs) XXX_DiscardUnknown() {
	xxx_messageInfo_Glyphs.DiscardUnknown(m)
}

var xxx_messageInfo_Glyphs proto.InternalMessageInfo

func (m *Glyphs) GetGlyphs() []*Glyph {
	if m != nil {
		return m.Glyphs
	}
	return nil
}

type Glyph struct {
	HorizAdvX            float64     `protobuf:"fixed64,1,opt,name=horiz_adv_x,json=horizAdvX,proto3" json:"horiz_adv_x,omitempty"`
	Unicode              string      `protobuf:"bytes,2,opt,name=unicode,proto3" json:"unicode,omitempty"`
	GerberLP             string      `protobuf:"bytes,3,opt,name=gerber_l_p,json=gerberLP,proto3" json:"gerber_l_p,omitempty"`
	PathSteps            []*PathStep `protobuf:"bytes,4,rep,name=path_steps,json=pathSteps,proto3" json:"path_steps,omitempty"`
	Mbb                  *MBB        `protobuf:"bytes,5,opt,name=mbb,proto3" json:"mbb,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Glyph) Reset()         { *m = Glyph{} }
func (m *Glyph) String() string { return proto.CompactTextString(m) }
func (*Glyph) ProtoMessage()    {}
func (*Glyph) Descriptor() ([]byte, []int) {
	return fileDescriptor_glyphs_5b62025fda4d7cef, []int{1}
}
func (m *Glyph) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Glyph.Unmarshal(m, b)
}
func (m *Glyph) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Glyph.Marshal(b, m, deterministic)
}
func (dst *Glyph) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Glyph.Merge(dst, src)
}
func (m *Glyph) XXX_Size() int {
	return xxx_messageInfo_Glyph.Size(m)
}
func (m *Glyph) XXX_DiscardUnknown() {
	xxx_messageInfo_Glyph.DiscardUnknown(m)
}

var xxx_messageInfo_Glyph proto.InternalMessageInfo

func (m *Glyph) GetHorizAdvX() float64 {
	if m != nil {
		return m.HorizAdvX
	}
	return 0
}

func (m *Glyph) GetUnicode() string {
	if m != nil {
		return m.Unicode
	}
	return ""
}

func (m *Glyph) GetGerberLP() string {
	if m != nil {
		return m.GerberLP
	}
	return ""
}

func (m *Glyph) GetPathSteps() []*PathStep {
	if m != nil {
		return m.PathSteps
	}
	return nil
}

func (m *Glyph) GetMbb() *MBB {
	if m != nil {
		return m.Mbb
	}
	return nil
}

type PathStep struct {
	C                    uint32    `protobuf:"varint,1,opt,name=c,proto3" json:"c,omitempty"`
	P                    []float64 `protobuf:"fixed64,2,rep,packed,name=p,proto3" json:"p,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PathStep) Reset()         { *m = PathStep{} }
func (m *PathStep) String() string { return proto.CompactTextString(m) }
func (*PathStep) ProtoMessage()    {}
func (*PathStep) Descriptor() ([]byte, []int) {
	return fileDescriptor_glyphs_5b62025fda4d7cef, []int{2}
}
func (m *PathStep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PathStep.Unmarshal(m, b)
}
func (m *PathStep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PathStep.Marshal(b, m, deterministic)
}
func (dst *PathStep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PathStep.Merge(dst, src)
}
func (m *PathStep) XXX_Size() int {
	return xxx_messageInfo_PathStep.Size(m)
}
func (m *PathStep) XXX_DiscardUnknown() {
	xxx_messageInfo_PathStep.DiscardUnknown(m)
}

var xxx_messageInfo_PathStep proto.InternalMessageInfo

func (m *PathStep) GetC() uint32 {
	if m != nil {
		return m.C
	}
	return 0
}

func (m *PathStep) GetP() []float64 {
	if m != nil {
		return m.P
	}
	return nil
}

type MBB struct {
	Xmin                 float64  `protobuf:"fixed64,1,opt,name=xmin,proto3" json:"xmin,omitempty"`
	Ymin                 float64  `protobuf:"fixed64,2,opt,name=ymin,proto3" json:"ymin,omitempty"`
	Xmax                 float64  `protobuf:"fixed64,3,opt,name=xmax,proto3" json:"xmax,omitempty"`
	Ymax                 float64  `protobuf:"fixed64,4,opt,name=ymax,proto3" json:"ymax,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MBB) Reset()         { *m = MBB{} }
func (m *MBB) String() string { return proto.CompactTextString(m) }
func (*MBB) ProtoMessage()    {}
func (*MBB) Descriptor() ([]byte, []int) {
	return fileDescriptor_glyphs_5b62025fda4d7cef, []int{3}
}
func (m *MBB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MBB.Unmarshal(m, b)
}
func (m *MBB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MBB.Marshal(b, m, deterministic)
}
func (dst *MBB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MBB.Merge(dst, src)
}
func (m *MBB) XXX_Size() int {
	return xxx_messageInfo_MBB.Size(m)
}
func (m *MBB) XXX_DiscardUnknown() {
	xxx_messageInfo_MBB.DiscardUnknown(m)
}

var xxx_messageInfo_MBB proto.InternalMessageInfo

func (m *MBB) GetXmin() float64 {
	if m != nil {
		return m.Xmin
	}
	return 0
}

func (m *MBB) GetYmin() float64 {
	if m != nil {
		return m.Ymin
	}
	return 0
}

func (m *MBB) GetXmax() float64 {
	if m != nil {
		return m.Xmax
	}
	return 0
}

func (m *MBB) GetYmax() float64 {
	if m != nil {
		return m.Ymax
	}
	return 0
}

func init() {
	proto.RegisterType((*Glyphs)(nil), "glyphs.Glyphs")
	proto.RegisterType((*Glyph)(nil), "glyphs.Glyph")
	proto.RegisterType((*PathStep)(nil), "glyphs.PathStep")
	proto.RegisterType((*MBB)(nil), "glyphs.MBB")
}

func init() { proto.RegisterFile("glyphs.proto", fileDescriptor_glyphs_5b62025fda4d7cef) }

var fileDescriptor_glyphs_5b62025fda4d7cef = []byte{
	// 271 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x90, 0x4f, 0x6b, 0x83, 0x40,
	0x14, 0xc4, 0x79, 0x6a, 0x6c, 0x7c, 0x26, 0x50, 0xf6, 0xb4, 0x87, 0xb6, 0x88, 0xd0, 0xe2, 0x29,
	0x81, 0xf4, 0x13, 0xd4, 0x4b, 0x2f, 0x0d, 0x84, 0x2d, 0x85, 0xde, 0x64, 0xfd, 0x43, 0x14, 0x62,
	0x5c, 0xd4, 0x86, 0xb5, 0x5f, 0xaa, 0x5f, 0xb1, 0xf8, 0xdc, 0xbd, 0xcd, 0xfc, 0xde, 0x0c, 0xcc,
	0x2e, 0x6e, 0xce, 0x97, 0x49, 0xd5, 0xc3, 0x4e, 0xf5, 0xdd, 0xd8, 0x31, 0x7f, 0x71, 0xf1, 0x1e,
	0xfd, 0x77, 0x52, 0xec, 0x19, 0x0d, 0xe3, 0x10, 0xb9, 0x49, 0x78, 0xd8, 0xee, 0x4c, 0x81, 0xee,
	0xc2, 0x16, 0xfe, 0x00, 0x57, 0x44, 0xd8, 0x13, 0x86, 0x75, 0xd7, 0x37, 0xbf, 0x99, 0x2c, 0x6f,
	0x99, 0xe6, 0x10, 0x41, 0x02, 0x22, 0x20, 0xf4, 0x56, 0xde, 0xbe, 0x19, 0xc7, 0xbb, 0x9f, 0x6b,
	0x53, 0x74, 0x65, 0xc5, 0x9d, 0x08, 0x92, 0x40, 0x58, 0xcb, 0x1e, 0x10, 0xcf, 0x55, 0x9f, 0x57,
	0x7d, 0x76, 0xc9, 0x14, 0x77, 0xe9, 0xb8, 0x5e, 0xc8, 0xc7, 0x89, 0xed, 0x11, 0x95, 0x1c, 0xeb,
	0x6c, 0x18, 0x2b, 0x35, 0x70, 0x8f, 0xc6, 0xdc, 0xdb, 0x31, 0x27, 0x39, 0xd6, 0x9f, 0x63, 0xa5,
	0x44, 0xa0, 0x8c, 0x1a, 0xd8, 0x23, 0xba, 0x6d, 0x9e, 0xf3, 0x55, 0x04, 0x49, 0x78, 0x08, 0x6d,
	0xf2, 0x98, 0xa6, 0x62, 0xe6, 0xf1, 0x0b, 0xae, 0x6d, 0x8b, 0x6d, 0x10, 0x0a, 0x5a, 0xba, 0x15,
	0x50, 0xcc, 0x4e, 0x71, 0x27, 0x72, 0x13, 0x10, 0xa0, 0xe2, 0x2f, 0x74, 0x8f, 0x69, 0xca, 0x18,
	0x7a, 0xba, 0x6d, 0xae, 0xe6, 0x3d, 0xa4, 0x67, 0x36, 0xcd, 0xcc, 0x59, 0xd8, 0x64, 0x98, 0x6e,
	0xa5, 0xa6, 0xf9, 0x94, 0x93, 0x7a, 0xc9, 0x49, 0xcd, 0x3d, 0x9b, 0x93, 0x3a, 0xf7, 0xe9, 0xc3,
	0x5f, 0xff, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xc9, 0x4f, 0xc9, 0x80, 0x01, 0x00, 0x00,
}
