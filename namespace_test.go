package oops_test

import (
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

var Error = oops.Namespace{"oops_test"}

func namespaceFunc() (err error) {
	defer Error.NewP(&err)

	return oops.New("bad stuff")
}

func TestNamespace(t *testing.T) {
	err := namespaceFunc()
	require.Error(t, err)

	t.Logf("%v\n", err)
	t.Logf("%+v\n", err)

	ne := err.(*oops.NamespaceError)
	require.Equal(t, Error.Name, ne.Name)
}
