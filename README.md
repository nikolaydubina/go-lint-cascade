# go-lint-cascade

[![codecov](https://codecov.io/gh/nikolaydubina/go-lint-cascade/graph/badge.svg?token=D3ww4rcZ0N)](https://codecov.io/gh/nikolaydubina/go-lint-cascade)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-lint-cascade)](https://goreportcard.com/report/github.com/nikolaydubina/go-lint-cascade)

Detect missing cascading calls in Go.

For example, if you have cascade calls for `WithDefaults()` config definitions, this linter will detect if you missed any `WithDefaults()` calls.

```bash
go install github.com/nikolaydubina/go-lint-cascade@latest
```

```bash
go-lint-cascade ./...
```

This code would be flagged:

```go
type Outer struct {
    Inner Inner
}

func (s Outer) WithDefaults() Outer {
    // ERROR: Missing s.Inner = s.Inner.WithDefaults()
    return s
}

type Inner struct { Value int }

func (s Inner) WithDefaults() Inner { 
    if s.Value == 0 {
        s.Value = 42
    }
    return s
}
```

Fixed version:
```go
func (s Outer) WithDefaults() Outer {
    s.Inner = s.Inner.WithDefaults()
    return s
}
```
