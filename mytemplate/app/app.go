package app

import (
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/myinit"
	"github.com/fitan/mykit/myrouter"
	"github.com/fitan/mykit/mytemplate/conf"
	"go.uber.org/zap"
)

type App struct {
	Router *myrouter.Router
	Gorm   *mygorm.DB
	Log    *zap.SugaredLogger
	Cfg    *conf.Conf
	SD     *myinit.SD
}

func (a *App) Run() {
	a.SD.Register()
	a.Router.Run(a.Cfg.App.Addr)
	a.SD.Wait()
}
