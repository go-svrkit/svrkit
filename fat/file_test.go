// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package fat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestIsFileExist(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"abcdefgxyz", false},
		{"./fs_test.go", true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := IsFileExist(tt.input); got != tt.want {
				t.Errorf("IsFileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadToLines(t *testing.T) {
	tests := []struct {
		input   string
		want    []string
		wantErr bool
	}{
		{"", []string(nil), false},
		{"a\nb\nc", []string{"a", "b", "c"}, false},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var rd = strings.NewReader(tt.input)
			got, err := ReadToLines(rd)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("ReadToLines() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadToLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}
