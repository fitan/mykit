package app

import (
	"context"
	"github.com/fitan/mykit/myconf"
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/myinit"
	"github.com/fitan/mykit/mylog"
	"github.com/fitan/mykit/myrouter"
	"github.com/fitan/mykit/mytemplate/conf"
	"github.com/fitan/mykit/mytemplate/handlers"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormLog "log"
	"net/http"
	"os"
	"time"
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

	newLogger := logger.New(
		gormLog.New(os.Stdout, "\r\n", gormLog.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(conf.Mysql.DSN), &gorm.Config{
		Logger: newLogger,
	})
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

func initConsul(conf *conf.Conf, log *zap.SugaredLogger) (*api.Client, error) {
	log.Infow("init consul staring...")
	defer log.Infow("init consul end...")
	cfg := api.DefaultConfig()
	cfg.Address = conf.Consul.Addr
	cfg.Token = conf.Consul.Token
	return api.NewClient(cfg)
}

type zapSugarLogger func(msg string, keysAndValues ...interface{})

func (l zapSugarLogger) Log(kv ...interface{}) error {
	l("", kv...)
	return nil
}

func initSD(conf *conf.Conf, client *api.Client, log *zap.SugaredLogger) (*myinit.SD, error) {
	return myinit.InitSD(conf.App.Name, conf.App.Addr, conf.App.Port, client, log)
}

func initTracer(conf *conf.Conf, log *zap.SugaredLogger) *sdktrace.TracerProvider {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(conf.App.Name),
		),
	)
	if err != nil {
		log.Fatalw("resource.New", "err", err.Error())
	}
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Trace.TracerProviderAddr)))
	if err != nil {
		log.Fatalw("jaeger.New", "err", err.Error())
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	return tp
}
