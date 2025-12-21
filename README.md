# go-lint-cascade

[![codecov](https://codecov.io/gh/nikolaydubina/go-lint-cascade/graph/badge.svg?token=D3ww4rcZ0N)](https://codecov.io/gh/nikolaydubina/go-lint-cascade)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-lint-cascade)](https://goreportcard.com/report/github.com/nikolaydubina/go-lint-cascade)

Detect missing cascading calls in Go.

```bash
go install github.com/nikolaydubina/go-lint-cascade@latest
```

```bash
go-lint-cascade ./...
```

For example, if you have cascade calls for `WithDefaults()` in nested config structs, this linter will detect if you missed any `WithDefaults()` calls.

```go
type Config struct {
    DB DBConfig
}

func (s Config) WithDefaults() Config {
    // ERROR: Missing s.DB = s.DB.WithDefaults()
    return s
}

type DBConfig struct {
    Port int
}

func (s DBConfig) WithDefaults() DBConfig { 
    if s.Port == 0 {
        s.Port = 42
    }
    return s
}
```
