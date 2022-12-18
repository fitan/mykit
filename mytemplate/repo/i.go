package repo

import (
	"github.com/fitan/mykit/mytemplate/repo/physicalMachine"
	"github.com/fitan/mykit/mytemplate/repo/user"
	"github.com/google/wire"
)

type Services struct {
	PhysicalMachine physicalMachine.Service
	User            user.Service
}

var Set = wire.NewSet(
	physicalMachine.Set,
	wire.Struct(new(Services), "*"),
)
