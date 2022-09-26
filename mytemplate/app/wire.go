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

var initSet = wire.NewSet(
	confSet,
	gormSet,
	logSet,
	routerSet,
	handlerSet,
	atomicLevelSet,
)

func InitApp() (App, error) {
	panic(wire.Build(
		initSet,
		wire.Struct(new(App), "*")))
	return App{}, nil
}
