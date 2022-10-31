//go:build wireinject
// +build wireinject

package app

import (
	"github.com/fitan/mykit/mytemplate/services"
	"github.com/google/wire"
)

var confSet = wire.NewSet(initConf)
var gormSet = wire.NewSet(initGORM)
var logSet = wire.NewSet(initLog)
var routerSet = wire.NewSet(initRouter)
var atomicLevelSet = wire.NewSet(initAtomicLevel)
var consulSet = wire.NewSet(initConsul)
var sdSet = wire.NewSet(initSD)
var myGORMSet = wire.NewSet(initMyGORM)
var mwsSet = wire.NewSet(initEndpointMiddleware)
var optsSet = wire.NewSet(initHttpServiceOptions)
var muxSet = wire.NewSet(initMux)

var initSet = wire.NewSet(
	consulSet,
	confSet,
	gormSet,
	logSet,
	routerSet,
	atomicLevelSet,
	sdSet,
	myGORMSet,
	mwsSet,
	optsSet,
	muxSet,
	services.Iset,
)

func InitApp() (App, error) {
	panic(wire.Build(
		initSet,
		wire.Struct(new(App), "*")))
	return App{}, nil
}
