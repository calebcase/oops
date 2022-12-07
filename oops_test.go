package oops_test

import (
	"fmt"

	"github.com/calebcase/oops"
)

func ExampleNew() {
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

	fmt.Printf("error: %v\n", err)
	// Output: error: bad stuff
}

func ExampleNew_details() {
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

	fmt.Printf("error: %+v\n", err)
}
