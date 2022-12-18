package services

import (
	"github.com/fitan/mykit/mytemplate/services/hello"
	"github.com/fitan/mykit/mytemplate/services/physicalMachine"
	"github.com/google/wire"
)

type Handlers struct {
	Hello           hello.Handler
	PhysicalMachine physicalMachine.Handler
}

var Set = wire.NewSet(
	hello.HelloSet,
	physicalMachine.Set,
	wire.Struct(new(Handlers), "*"),
)
