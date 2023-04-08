package oops_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

func TestVerbose(t *testing.T) {
	err := oops.New("bad thing")

	lines := len(strings.Split(err.Error(), "\n"))
	require.Equal(t, 1, lines)

	verr := oops.Verbose(err)
	vlines := len(strings.Split(verr.Error(), "\n"))
	require.Greater(t, vlines, 1)

	vflines := len(strings.Split(fmt.Sprintf("%v", verr), "\n"))
	require.Greater(t, vflines, 1)

	require.Nil(t, oops.Verbose(nil))
}
