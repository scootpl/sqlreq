# SQLReq

[![Go Reference](https://pkg.go.dev/badge/github.com/scootpl/sqlreq.svg)](https://pkg.go.dev/github.com/scootpl/sqlreq) [![Go Report Card](https://goreportcard.com/badge/github.com/scootpl/sqlreq)](https://goreportcard.com/report/github.com/scootpl/sqlreq)

SQLReq is a wrapper around standard Go net/http client. You can execute REST requests using SQL-like syntax.

## Quick Start

```go
import "github.com/scootpl/sqlreq"

xauth, status, err := sqlreq.SelectHEADER(`
    X-AUTH FROM http://example.com/api/init
    WHERE login = testlogin AND pass = testpass
    WITH POST`)
```

```go
token := "12345678"

payload := Msg{
    Message: "test",
}

body, status, err := sqlreq.SelectBODY(`
    FROM http://example.com/api/status
    WHERE HEADER x-auth = %s AND
    PAYLOAD = ?
    WITH POST`, token, &payload)
```

```go
resp, status, err := sqlreq.SelectRESPONSE(`
    FROM http://example.com/api/status
    WHERE HEADER x-auth = '123456' AND
    message = 'test'
    WITH TIMEOUT 30
    WITH POST`)

```

See [GoDoc](https://godoc.org/github.com/scootpl/sqlreq) for more details.