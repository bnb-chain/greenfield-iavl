package iavl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffLayer_SaveRoot(t *testing.T) {
	c := newNdbCache(0)
	err := c.SaveRoot(0, []byte{'a', 'b', 'c'})
	require.NoError(t, err)

	err = c.SaveRoot(1, []byte{'d', 'e', 'f'})
	require.NoError(t, err)

	root := c.GetRoot(0)
	fmt.Println(root)
	root = c.GetRoot(1)
	fmt.Println(root)
}
