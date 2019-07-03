// Code generated by protoc-gen-go. DO NOT EDIT.
// source: genbank.proto

package gbparse

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Genbank struct {
	LOCUS      string       `protobuf:"bytes,1,opt,name=LOCUS" json:"LOCUS,omitempty"`
	ACCESSION  string       `protobuf:"bytes,2,opt,name=ACCESSION" json:"ACCESSION,omitempty"`
	DEFINITION string       `protobuf:"bytes,3,opt,name=DEFINITION" json:"DEFINITION,omitempty"`
	VERSION    string       `protobuf:"bytes,4,opt,name=VERSION" json:"VERSION,omitempty"`
	DBLINK     string       `protobuf:"bytes,5,opt,name=DBLINK" json:"DBLINK,omitempty"`
	KEYWORDS   string       `protobuf:"bytes,6,opt,name=KEYWORDS" json:"KEYWORDS,omitempty"`
	SOURCE     string       `protobuf:"bytes,7,opt,name=SOURCE" json:"SOURCE,omitempty"`
	ORGANISM   string       `protobuf:"bytes,8,opt,name=ORGANISM" json:"ORGANISM,omitempty"`
	COMMENT    string       `protobuf:"bytes,9,opt,name=COMMENT" json:"COMMENT,omitempty"`
	SEQUENCE   string       `protobuf:"bytes,10,opt,name=SEQUENCE" json:"SEQUENCE,omitempty"`
	CONTIG     string       `protobuf:"bytes,11,opt,name=CONTIG" json:"CONTIG,omitempty"`
	REFERENCES []*Reference `protobuf:"bytes,12,rep,name=REFERENCES" json:"REFERENCES,omitempty"`
	FEATURES   []*Feature   `protobuf:"bytes,13,rep,name=FEATURES" json:"FEATURES,omitempty"`
}

func (m *Genbank) Reset()                    { *m = Genbank{} }
func (m *Genbank) String() string            { return proto.CompactTextString(m) }
func (*Genbank) ProtoMessage()               {}
func (*Genbank) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Genbank) GetLOCUS() string {
	if m != nil {
		return m.LOCUS
	}
	return ""
}

func (m *Genbank) GetACCESSION() string {
	if m != nil {
		return m.ACCESSION
	}
	return ""
}

func (m *Genbank) GetDEFINITION() string {
	if m != nil {
		return m.DEFINITION
	}
	return ""
}

func (m *Genbank) GetVERSION() string {
	if m != nil {
		return m.VERSION
	}
	return ""
}

func (m *Genbank) GetDBLINK() string {
	if m != nil {
		return m.DBLINK
	}
	return ""
}

func (m *Genbank) GetKEYWORDS() string {
	if m != nil {
		return m.KEYWORDS
	}
	return ""
}

func (m *Genbank) GetSOURCE() string {
	if m != nil {
		return m.SOURCE
	}
	return ""
}

func (m *Genbank) GetORGANISM() string {
	if m != nil {
		return m.ORGANISM
	}
	return ""
}

func (m *Genbank) GetCOMMENT() string {
	if m != nil {
		return m.COMMENT
	}
	return ""
}

func (m *Genbank) GetSEQUENCE() string {
	if m != nil {
		return m.SEQUENCE
	}
	return ""
}

func (m *Genbank) GetCONTIG() string {
	if m != nil {
		return m.CONTIG
	}
	return ""
}

func (m *Genbank) GetREFERENCES() []*Reference {
	if m != nil {
		return m.REFERENCES
	}
	return nil
}

func (m *Genbank) GetFEATURES() []*Feature {
	if m != nil {
		return m.FEATURES
	}
	return nil
}

type Reference struct {
	Number  int32  `protobuf:"varint,1,opt,name=Number" json:"Number,omitempty"`
	ORIGIN  string `protobuf:"bytes,2,opt,name=ORIGIN" json:"ORIGIN,omitempty"`
	AUTHORS string `protobuf:"bytes,3,opt,name=AUTHORS" json:"AUTHORS,omitempty"`
	TITLE   string `protobuf:"bytes,4,opt,name=TITLE" json:"TITLE,omitempty"`
	JOURNAL string `protobuf:"bytes,5,opt,name=JOURNAL" json:"JOURNAL,omitempty"`
	PUBMED  string `protobuf:"bytes,6,opt,name=PUBMED" json:"PUBMED,omitempty"`
}

func (m *Reference) Reset()                    { *m = Reference{} }
func (m *Reference) String() string            { return proto.CompactTextString(m) }
func (*Reference) ProtoMessage()               {}
func (*Reference) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Reference) GetNumber() int32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Reference) GetORIGIN() string {
	if m != nil {
		return m.ORIGIN
	}
	return ""
}

func (m *Reference) GetAUTHORS() string {
	if m != nil {
		return m.AUTHORS
	}
	return ""
}

func (m *Reference) GetTITLE() string {
	if m != nil {
		return m.TITLE
	}
	return ""
}

func (m *Reference) GetJOURNAL() string {
	if m != nil {
		return m.JOURNAL
	}
	return ""
}

func (m *Reference) GetPUBMED() string {
	if m != nil {
		return m.PUBMED
	}
	return ""
}

type Feature struct {
	TYPE         string            `protobuf:"bytes,1,opt,name=TYPE" json:"TYPE,omitempty"`
	IsCompliment bool              `protobuf:"varint,2,opt,name=isCompliment" json:"isCompliment,omitempty"`
	START        string            `protobuf:"bytes,3,opt,name=START" json:"START,omitempty"`
	STOP         string            `protobuf:"bytes,4,opt,name=STOP" json:"STOP,omitempty"`
	QUALIFIERS   map[string]string `protobuf:"bytes,5,rep,name=QUALIFIERS" json:"QUALIFIERS,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Feature) Reset()                    { *m = Feature{} }
func (m *Feature) String() string            { return proto.CompactTextString(m) }
func (*Feature) ProtoMessage()               {}
func (*Feature) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Feature) GetTYPE() string {
	if m != nil {
		return m.TYPE
	}
	return ""
}

func (m *Feature) GetIsCompliment() bool {
	if m != nil {
		return m.IsCompliment
	}
	return false
}

func (m *Feature) GetSTART() string {
	if m != nil {
		return m.START
	}
	return ""
}

func (m *Feature) GetSTOP() string {
	if m != nil {
		return m.STOP
	}
	return ""
}

func (m *Feature) GetQUALIFIERS() map[string]string {
	if m != nil {
		return m.QUALIFIERS
	}
	return nil
}

func init() {
	proto.RegisterType((*Genbank)(nil), "gbparse.Genbank")
	proto.RegisterType((*Reference)(nil), "gbparse.Reference")
	proto.RegisterType((*Feature)(nil), "gbparse.Feature")
}

func init() { proto.RegisterFile("genbank.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 469 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x93, 0xcf, 0x6a, 0xdb, 0x40,
	0x10, 0xc6, 0x51, 0x1c, 0x5b, 0xf6, 0x24, 0xa1, 0x61, 0x29, 0x65, 0x09, 0xa5, 0x18, 0x9f, 0x72,
	0x28, 0x3e, 0xa4, 0x97, 0x52, 0x28, 0x54, 0x91, 0x47, 0xee, 0x36, 0xb2, 0xd6, 0x99, 0x5d, 0xb5,
	0xe4, 0x68, 0x97, 0x6d, 0x08, 0x89, 0x65, 0xa3, 0xd8, 0x85, 0x3c, 0x45, 0x9f, 0xa0, 0x6f, 0xd8,
	0x87, 0x28, 0xfb, 0x47, 0x8e, 0xe3, 0xdb, 0xfe, 0xe6, 0x9b, 0x6f, 0xf7, 0xf3, 0x8c, 0x05, 0x27,
	0xb7, 0xa6, 0x9a, 0xcf, 0xaa, 0xfb, 0xe1, 0xaa, 0x5e, 0xae, 0x97, 0x2c, 0xbe, 0x9d, 0xaf, 0x66,
	0xf5, 0xa3, 0x19, 0xfc, 0x69, 0x41, 0x3c, 0xf6, 0x12, 0x7b, 0x0d, 0xed, 0x5c, 0xa6, 0xa5, 0xe2,
	0x51, 0x3f, 0x3a, 0xef, 0x91, 0x07, 0xf6, 0x16, 0x7a, 0x49, 0x9a, 0xa2, 0x52, 0x42, 0x16, 0xfc,
	0xc0, 0x29, 0xcf, 0x05, 0xf6, 0x0e, 0x60, 0x84, 0x99, 0x28, 0x84, 0xb6, 0x72, 0xcb, 0xc9, 0x3b,
	0x15, 0xc6, 0x21, 0xfe, 0x8e, 0xe4, 0xbc, 0x87, 0x4e, 0x6c, 0x90, 0xbd, 0x81, 0xce, 0xe8, 0x32,
	0x17, 0xc5, 0x15, 0x6f, 0x3b, 0x21, 0x10, 0x3b, 0x83, 0xee, 0x15, 0xde, 0xfc, 0x90, 0x34, 0x52,
	0xbc, 0xe3, 0x94, 0x2d, 0x5b, 0x8f, 0x92, 0x25, 0xa5, 0xc8, 0x63, 0xef, 0xf1, 0x64, 0x3d, 0x92,
	0xc6, 0x49, 0x21, 0xd4, 0x84, 0x77, 0xbd, 0xa7, 0x61, 0x9b, 0x20, 0x95, 0x93, 0x09, 0x16, 0x9a,
	0xf7, 0x7c, 0x82, 0x80, 0xd6, 0xa5, 0xf0, 0xba, 0xc4, 0x22, 0x45, 0x0e, 0xde, 0xd5, 0xb0, 0x7d,
	0x29, 0x95, 0x85, 0x16, 0x63, 0x7e, 0xe4, 0x5f, 0xf2, 0xc4, 0x2e, 0x00, 0x08, 0x33, 0x24, 0xdb,
	0xa4, 0xf8, 0x71, 0xbf, 0x75, 0x7e, 0x74, 0xc1, 0x86, 0x61, 0x9a, 0x43, 0x32, 0xbf, 0x4c, 0x6d,
	0xaa, 0x9f, 0x86, 0x76, 0xba, 0xd8, 0x7b, 0xe8, 0x66, 0x98, 0xe8, 0x92, 0x50, 0xf1, 0x13, 0xe7,
	0x38, 0xdd, 0x3a, 0x32, 0x33, 0x5b, 0x6f, 0x6a, 0x43, 0xdb, 0x8e, 0xc1, 0xdf, 0x08, 0x7a, 0xdb,
	0x7b, 0x6c, 0x8e, 0x62, 0xb3, 0x98, 0x9b, 0xda, 0x2d, 0xa5, 0x4d, 0x81, 0x6c, 0x5d, 0x92, 0x18,
	0x8b, 0x66, 0x25, 0x81, 0xec, 0xaf, 0x4d, 0x4a, 0xfd, 0x55, 0x92, 0x0a, 0xcb, 0x68, 0xd0, 0x6e,
	0x57, 0x0b, 0x9d, 0x63, 0xd8, 0x83, 0x07, 0xdb, 0xff, 0x4d, 0x96, 0x54, 0x24, 0x79, 0x58, 0x43,
	0x83, 0xf6, 0x85, 0x69, 0x79, 0x39, 0xc1, 0x51, 0xd8, 0x42, 0xa0, 0xc1, 0xbf, 0x08, 0xe2, 0x90,
	0x9a, 0x31, 0x38, 0xd4, 0x37, 0x53, 0x0c, 0x7f, 0x18, 0x77, 0x66, 0x03, 0x38, 0xbe, 0x7b, 0x4c,
	0x97, 0x8b, 0xd5, 0xc3, 0xdd, 0xc2, 0x54, 0x6b, 0x97, 0xaf, 0x4b, 0x2f, 0x6a, 0x36, 0x8b, 0xd2,
	0x09, 0xe9, 0x90, 0xd1, 0x83, 0xbd, 0x4d, 0x69, 0x39, 0x0d, 0x01, 0xdd, 0x99, 0x7d, 0x01, 0xb8,
	0x2e, 0x93, 0x5c, 0x64, 0x02, 0x49, 0xf1, 0xb6, 0x9b, 0x5e, 0x7f, 0x7f, 0x7a, 0xc3, 0xe7, 0x16,
	0xac, 0xd6, 0xf5, 0x13, 0xed, 0x78, 0xce, 0x3e, 0xc3, 0xab, 0x3d, 0x99, 0x9d, 0x42, 0xeb, 0xde,
	0x3c, 0x85, 0xd4, 0xf6, 0x68, 0x03, 0xfd, 0x9e, 0x3d, 0x6c, 0x4c, 0x98, 0xa6, 0x87, 0x4f, 0x07,
	0x1f, 0xa3, 0x79, 0xc7, 0x7d, 0x30, 0x1f, 0xfe, 0x07, 0x00, 0x00, 0xff, 0xff, 0x5d, 0xb2, 0xab,
	0xf3, 0x41, 0x03, 0x00, 0x00,
}
