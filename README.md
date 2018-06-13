# Session middleware for [Beego](https://github.com/astaxie/beego)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v gopkg.in/go-session/beego-session.v2
```

### Create file `server.go`

```go
package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gopkg.in/go-session/beego-session.v2"
	"gopkg.in/session.v2"
)

func main() {
	app := beego.NewApp()

	app.Handlers.InsertFilter("*", beego.BeforeRouter,
		beegosession.New(
			session.SetCookieName("session_id"),
			session.SetSign([]byte("sign")),
		),
	)

	app.Handlers.Get("/", func(ctx *context.Context) {
		store := beegosession.FromContext(ctx)
		store.Set("foo", "bar")
		err := store.Save()
		if err != nil {
			ctx.Abort(500, err.Error())
			return
		}
		ctx.Redirect(302, "/foo")
	})

	app.Handlers.Get("/foo", func(ctx *context.Context) {
		store := beegosession.FromContext(ctx)
		foo, ok := store.Get("foo")
		if !ok {
			ctx.Abort(404, "not found")
			return
		}
		ctx.WriteString(fmt.Sprintf("foo:%s", foo))
	})

	beego.BConfig.Listen.HTTPPort = 8080
	app.Run()
}
```

### Build and run

```bash
$ go build server.go
$ ./server
```

### Open in your web browser

<http://localhost:8080>

    foo:bar


## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/go-session/beego-session
[Build-Status-Image]: https://travis-ci.org/go-session/beego-session.svg?branch=master
[codecov-url]: https://codecov.io/gh/go-session/beego-session
[codecov-image]: https://codecov.io/gh/go-session/beego-session/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/go-session/beego-session
[reportcard-image]: https://goreportcard.com/badge/github.com/go-session/beego-session
[godoc-url]: https://godoc.org/github.com/go-session/beego-session
[godoc-image]: https://godoc.org/github.com/go-session/beego-session?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg