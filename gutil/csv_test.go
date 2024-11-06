// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var csvContent = `
Id,Season,Faction,Region
1,3,1,"{""1"":2,""2"":1,""3"":3}"
2,3,2,"{""1"":5,""2"":6,""3"":4}"
`

type TTStruct struct {
	Id      int32
	Season  int32
	Faction int32
	Region  map[int32]int32
}

func TestReadCSVTableFrom(t *testing.T) {
	var rd = strings.NewReader(csvContent)
	rec, err := ReadCSVTableFrom(rd)
	require.Nil(t, err)
	require.NotNil(t, rec)

	assert.Equal(t, 2, rec.RowSize())
	assert.Equal(t, 4, rec.ColSize())

	assert.Equal(t, "1", rec.GetRowField(0, "Id"))
	assert.Equal(t, "3", rec.GetRowField(1, "Season"))
}

func TestCSVTable_ScanRow(t *testing.T) {
	var rd = strings.NewReader(csvContent)
	rec, err := ReadCSVTableFrom(rd)
	assert.Nil(t, err)
	assert.NotNil(t, rec)

	var aa = new(TTStruct)
	assert.Nil(t, rec.ScanRow(0, aa))
	assert.Equal(t, int32(1), aa.Id)
	assert.Equal(t, int32(3), aa.Season)
	assert.Equal(t, int32(1), aa.Faction)
	assert.Equal(t, 3, len(aa.Region))
}
