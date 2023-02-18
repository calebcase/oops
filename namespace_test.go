package oops_test

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

var Error = oops.Namespace("oops_test")

func namespaceWrap() (err error) {
	return Error.Wrap(oops.New("bad stuff"))
}

func namespaceWrapP() (err error) {
	defer Error.WrapP(&err)

	return oops.New("bad stuff")
}

func namespaceNew() (err error) {
	return Error.New("bad stuff")
}

func namespaceTrace() (err error) {
	return Error.Trace(errors.New("bad stuff"))
}

func namespaceTraceN() (err error) {
	return Error.TraceN(errors.New("bad stuff"), oops.TraceSkipInternal)
}

func namespaceTraceWithOptions() (err error) {
	return Error.TraceWithOptions(errors.New("bad stuff"), oops.TraceOptions{
		Skip: oops.TraceSkipInternal,
	})
}

func namespaceChain() (err error) {
	return Error.Chain(
		oops.New("bad stuff"),
		oops.New("worse stuff"),
	)
}

func namespaceChainP() (err error) {
	defer func() {
		Error.ChainP(&err, oops.New("worse stuff"))
	}()

	return oops.New("bad stuff")
}

func namespaceShadow() (err error) {
	return Error.Shadow(
		oops.New("internal bad stuff"),
		oops.New("external bad stuff"),
	)
}

func namespaceShadowP() (err error) {
	err = oops.New("internal bad stuff")
	Error.ShadowP(&err, oops.New("external bad stuff"))

	return err
}

func namespaceShadowF() (err error) {
	defer Error.ShadowF(&err, oops.New("external bad stuff"))()

	return oops.New("internal bad stuff")
}

func TestNamespace(t *testing.T) {
	tcs := []struct {
		Fn func() error
	}{
		{Fn: namespaceWrap},
		{Fn: namespaceWrapP},
		{Fn: namespaceNew},
		{Fn: namespaceTrace},
		{Fn: namespaceTraceN},
		{Fn: namespaceTraceWithOptions},
		{Fn: namespaceChain},
		{Fn: namespaceChainP},
		{Fn: namespaceShadow},
		{Fn: namespaceShadowP},
		{Fn: namespaceShadowF},
	}
	for i, tc := range tcs {
		name := runtime.FuncForPC(reflect.ValueOf(tc.Fn).Pointer()).Name()
		t.Run(fmt.Sprintf("%02d/%s", i, name), func(t *testing.T) {
			err := tc.Fn()
			require.Error(t, err)

			t.Logf("%v\n", err)
			t.Logf("%+v\n", err)

			ne := err.(*oops.NamespaceError)
			require.Equal(t, string(Error), ne.Name)

			var te *oops.TraceError
			require.Equal(t, true, errors.As(err, &te))

			fs := te.Data.(oops.Frames)
			require.Equal(t, name, fs[0].Function)
		})
	}
}
