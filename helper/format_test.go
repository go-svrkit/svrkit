package helper

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
)

type Coord struct {
	X int32 `protobuf:"varint,1,opt,name=X,proto3" json:"X,omitempty"`
	Z int32 `protobuf:"varint,2,opt,name=Z,proto3" json:"Z,omitempty"`
}

func (m *Coord) Reset()         { *m = Coord{} }
func (m *Coord) String() string { return proto.CompactTextString(m) }
func (*Coord) ProtoMessage()    {}

func TestPrettyBytes(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0B"},
		{KB, "1.00KB"},
		{-KB, "-1.00KB"},
		{KB + 100, "1.10KB"},
		{MB, "1.00MB"},
		{MB + 10*KB, "1.01MB"},
		{GB, "1.00GB"},
		{GB + 100*MB, "1.10GB"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := PrettyBytes(tt.input); got != tt.want {
				t.Errorf("PrettyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONParse(t *testing.T) {
	tests := []struct {
		input   string
		want    interface{}
		wantErr bool
	}{
		{"1", 1, false},
		{"true", true, false},
		{"127", int8(127), false},
		{"65535", uint16(65535), false},
		{"2147483647", int32(math.MaxInt32), false},
		{"9223372036854775807", int64(math.MaxInt64), false},
		{"3.14", float32(3.14), false},
		{`{"X":12, "Y":34}`, image.Point{X: 12, Y: 34}, false},
		{`{"X":12,"Y":34}`, map[string]int{"X": 12, "Y": 34}, false},
		{"[4294967295,18446744073709551615]", []uint64{math.MaxUint32, math.MaxUint64}, false},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var rval = reflect.New(reflect.TypeOf(tt.want))
			var err = JSONParse([]byte(tt.input), rval.Interface())
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("JSONParse() error = %v, wantErr %v", err, tt.wantErr)
			}
			var got = rval.Elem().Interface()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrettyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONStringify(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		input interface{}
		want  string
	}{
		{1, "1"},
		{true, "true"},
		{math.MaxInt64, "9223372036854775807"},
		{image.Point{X: 12, Y: 34}, `{"X":12,"Y":34}`},
		{map[string]int{"X": 12, "Y": 34}, `{"X":12,"Y":34}`},
		{[]uint64{math.MaxUint32, math.MaxUint64}, `[4294967295,18446744073709551615]`},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := JSONStringify(tt.input); got != tt.want {
				t.Errorf("JSONStringify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatProtoToJSON(t *testing.T) {
	tests := []struct {
		input proto.Message
		want  string
	}{
		{nil, ""},
		{&Coord{}, "{}"},
		{&Coord{X: 12, Z: 34}, `{"X":12,"Z":34}`},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := FormatProtoToJSON(tt.input); got != tt.want {
				t.Errorf("FormatProtoToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSONToProto(t *testing.T) {
	tests := []struct {
		input   string
		want    proto.Message
		wantErr bool
	}{
		{"{}", &Coord{}, true},
		{`{"X":12, "Z":34}`, &Coord{X: 12, Z: 34}, true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var dst = reflect.New(reflect.TypeOf(tt.want).Elem()).Interface()
			var err = ParseJSONToProto(tt.input, dst.(proto.Message))
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("ParseJSONToProto() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, dst) {
				t.Errorf("ParseJSONToProto() %v want %T:%v, got %T:%v", err, tt.want, tt.want, dst, dst)
			}
		})
	}
}
