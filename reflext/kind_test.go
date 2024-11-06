package reflext

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TTT struct {
	a int
}

func (t *TTT) String() string {
	return fmt.Sprintf("<%d>", t.a)
}

func TestIsInterfaceNil(t *testing.T) {
	var tt *TTT
	var iface fmt.Stringer
	iface = tt
	var mm map[string]string
	var ch chan int

	require.True(t, IsInterfaceNil(nil))
	require.True(t, IsInterfaceNil(tt))
	require.True(t, IsInterfaceNil(mm))
	require.True(t, IsInterfaceNil(ch))
	require.True(t, IsInterfaceNil(iface))
}
