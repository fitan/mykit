package app

import (
	"github.com/fitan/mykit/myconf"
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/mylog"
	"github.com/fitan/mykit/myrouter"
	"github.com/fitan/mykit/mytemplate/conf"
	"github.com/fitan/mykit/mytemplate/handlers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func initConf() (*conf.Conf, error) {
	v := conf.Conf{}
	err := myconf.ReadFile("conf_dev.yaml", []string{"./"}, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func initGORM(conf *conf.Conf, log *zap.SugaredLogger) (*mygorm.DB, error) {
	log.Infow("init gorm staring...")
	defer log.Infow("init gorm end...")
	db, err := gorm.Open(mysql.Open(conf.Mysql.DSN))
	return mygorm.New(db), err
}

func initAtomicLevel(conf *conf.Conf) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(zapcore.Level(conf.Log.Level))
}

func initLog(conf *conf.Conf, level zap.AtomicLevel) *zap.SugaredLogger {
	log := mylog.New(conf.App.Name, conf.Log.Dir, level)
	return log
}

func initRouter(log *zap.SugaredLogger, level zap.AtomicLevel) *myrouter.Router {
	log.Infow("init router staring...")
	defer log.Infow("init router end...")

	r := myrouter.New()
	r.Setlog(log)
	r.Methods(http.MethodPut).Path("/log").Handler(level)
	return r
}

type initHandlerWire struct {
}

func initHandler(router *myrouter.Router) initHandlerWire {
	handlers.Handlers(router.Router)
	return initHandlerWire{}
}
