package oops_test

import (
	"testing"

	"github.com/calebcase/oops"
)

func TestNew(t *testing.T) {
	a := func() error {
		return oops.New("bad stuff")
	}

	b := func() error {
		return a()
	}

	c := func() error {
		return b()
	}

	err := c()

	t.Logf("error: %v\n", err)
}
