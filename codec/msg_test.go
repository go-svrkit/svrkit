package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/codec/testdata"
)

func TestMarshal(t *testing.T) {
	var req = &testdata.BuildReq{
		Type: testdata.BuildingType_Hospital,
		PosX: 1234,
		PosY: 5678,
	}
	data, err := Marshal(req)
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)
}
