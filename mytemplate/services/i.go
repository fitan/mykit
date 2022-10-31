package services

import (
	"github.com/fitan/mykit/mytemplate/services/hello"
	"github.com/google/wire"
)

type Handlers struct {
	Hello hello.Handler
}

var Iset = wire.NewSet(
	hello.HelloSet,
	wire.Struct(new(Handlers), "*"),
)
