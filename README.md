# oops

[![Go Reference](https://pkg.go.dev/badge/github.com/calebcase/oops.svg)](https://pkg.go.dev/github.com/calebcase/oops)
[![Go Report Card](https://goreportcard.com/badge/github.com/calebcase/oops)](https://goreportcard.com/report/github.com/calebcase/oops)

Oops implements the batteries that are missing from [errors][errors]:

* Stack traces: attach a stack trace to any error
* Chaining: combining multiple errors into a single stack of errors (as you
  might want when handling the errors from disposing resources like
  [file_close][file.Close()])
* Namespacing: create an error factory that prefixes errors with a given name
* Shadowing: hide the exact error behind a package level error (as you might
  want when trying to stabilize your API's supported errors)

---

[errors]: https://pkg.go.dev/errors
[file_close]: https://pkg.go.dev/os#File.Close
