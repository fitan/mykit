// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() (App, error) {
	conf, err := initConf()
	if err != nil {
		return App{}, err
	}
	atomicLevel := initAtomicLevel(conf)
	sugaredLogger := initLog(conf, atomicLevel)
	router := initRouter(sugaredLogger, atomicLevel)
	db, err := initGORM(conf, sugaredLogger)
	if err != nil {
		return App{}, err
	}
	app := App{
		Router: router,
		Gorm:   db,
		Log:    sugaredLogger,
		Cfg:    conf,
	}
	return app, nil
}

// wire.go:

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