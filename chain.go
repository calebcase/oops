package oops

import (
	"fmt"
	"strings"

	"github.com/calebcase/oops/lines"
)

// ChainError is a list of errors oldest to newest.
type ChainError []error

var _ error = ChainError{}
var _ unwrapper = ChainError{}

func (ce ChainError) Error() string {
	return fmt.Sprintf("%v", ce)
}

func (ce ChainError) Unwrap() error {
	if len(ce) == 0 || len(ce) == 1 {
		return nil
	}

	if len(ce) == 2 {
		return ce[1]
	}

	return ChainError(ce[1:])
}

// Format implements fmt.Format.
func (ce ChainError) Format(f fmt.State, verb rune) {
	if len(ce) == 0 {
		return
	}

	flag := ""
	if f.Flag(int('+')) {
		flag = "+"
	}

	if flag == "" {
		fmt.Fprintf(f, "%"+string(verb), ce[0])

		return
	}

	errs := make([]string, 0, len(ce))

	fmt.Fprintf(f, "chain(len=%d):\n", len(ce))
	for i, err := range ce {
		lines := lines.Indent(lines.Sprintf("%"+flag+string(verb), err), "··", 1)

		errs = append(errs, fmt.Sprintf("··[%d] %s", i, lines[0]))

		if len(lines) > 1 {
			errs = append(errs, lines[1:]...)
		}
	}

	f.Write([]byte(strings.Join(errs, "\n")))
}

// Chain combines errors into a chain of errors. nil errors are removed.
func Chain(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	ce := ChainError{}

	for _, err := range errs {
		if err == nil {
			continue
		}

		if peer, ok := err.(ChainError); ok {
			ce = append(ce, peer...)

			continue
		}

		ce = append(ce, err)
	}

	return ce
}

// ChainP combines errors into a chain of errors. nil errors are removed.
func ChainP(err *error, errs ...error) {
	if err == nil {
		return
	}

	*err = Chain(append([]error{*err}, errs...)...)
}
