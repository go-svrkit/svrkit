// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"gopkg.in/svrkit.v1/conv"
)

var errRowOutOfRange = errors.New("row out of range")

// CSVTable represents a CSV table
type CSVTable struct {
	Header map[string]int
	Rows   [][]string
}

// ColSize returns the number of columns
func (t *CSVTable) ColSize() int {
	return len(t.Header)
}

// RowSize returns the number of rows
func (t *CSVTable) RowSize() int {
	return len(t.Rows)
}

// GetRowField returns the field value of the column
func (t *CSVTable) GetRowField(rowIdx int, name string) string {
	idx, ok := t.Header[name]
	if ok && rowIdx >= 0 && rowIdx < len(t.Rows) {
		var row = t.Rows[rowIdx]
		if idx >= 0 && idx < len(row) {
			return row[idx]
		}
	}
	return ""
}

func (t *CSVTable) ScanRow(row int, target any) error {
	if row < 0 && row >= len(t.Rows) {
		return errRowOutOfRange
	}
	var rval = reflect.ValueOf(target)
	if !rval.IsValid() && !rval.CanSet() {
		return fmt.Errorf("target is not valid or cannot be set")
	}
	if kind := rval.Kind(); kind != reflect.Ptr {
		return fmt.Errorf("target must be a pointer, got %v", kind)
	}
	rval = rval.Elem()
	if kind := rval.Kind(); kind != reflect.Struct {
		return fmt.Errorf("target must be struct type, got %v", kind)
	}
	var rt = rval.Type()
	var numField = rt.NumField()
	for i := 0; i < numField; i++ {
		var ft = rt.Field(i)
		var content = t.GetRowField(row, ft.Name)
		if content != "" {
			val, err := conv.ReflectConvAny(ft.Type, content)
			if err != nil {
				return fmt.Errorf("cannot set field %s: %w", ft.Name, err)
			}
			var field = rval.Field(i)
			field.Set(val)
		}
	}
	return nil
}

func NewCSVRecords() *CSVTable {
	return &CSVTable{
		Header: make(map[string]int),
	}
}

func isRowAllEmpty(row []string) bool {
	for i := 0; i < len(row); i++ {
		row[i] = strings.TrimSpace(row[i])
		if row[i] != "" {
			return false
		}
	}
	return true
}

func ReadCSVTableFrom(rd io.Reader) (*CSVTable, error) {
	var reader = csv.NewReader(rd)
	titles, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var rec = NewCSVRecords()
	for i, title := range titles {
		rec.Header[title] = i
	}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return rec, nil
			}
			return nil, err
		}
		if !isRowAllEmpty(row) {
			rec.Rows = append(rec.Rows, row)
		}
	}
}

// ReadCSVTable 从csv文件读取数据
func ReadCSVTable(filename string) (*CSVTable, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ReadCSVTableFrom(f)
}
