//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

var confSet = wire.NewSet(initConf)
var gormSet = wire.NewSet(initGORM)
var logSet = wire.NewSet(initLog)
var routerSet = wire.NewSet(initRouter)
var handlerSet = wire.NewSet(initHandler)
var atomicLevelSet = wire.NewSet(initAtomicLevel)
var consulSet = wire.NewSet(initConsul)
var SdSet = wire.NewSet(initSD)

var initSet = wire.NewSet(
	consulSet,
	confSet,
	gormSet,
	logSet,
	routerSet,
	handlerSet,
	atomicLevelSet,
	SdSet,
)

func InitApp() (App, error) {
	panic(wire.Build(
		initSet,
		wire.Struct(new(App), "*")))
	return App{}, nil
}
