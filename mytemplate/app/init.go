package app

import (
	"context"
	"github.com/fitan/mykit/myconf"
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/myhttp"
	"github.com/fitan/mykit/myinit"
	"github.com/fitan/mykit/mylog"
	"github.com/fitan/mykit/myrouter"
	"github.com/fitan/mykit/mytemplate/conf"
	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
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
	"net/http"
)

type ConfName string

func initConf(confName ConfName) (*conf.Conf, error) {
	v := conf.Conf{}
	err := myconf.ReadFile(string(confName)+".yaml", []string{"./"}, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func initGORM(conf *conf.Conf, log *zap.SugaredLogger) (*gorm.DB, error) {
	log.Infow("init gorm staring...")
	defer log.Infow("init gorm end...")
	db, err := gorm.Open(mysql.Open(conf.Mysql.DSN))
	if err != nil {
		return db, err
	}

	return db, nil
}

func initMyGORM(db *gorm.DB, log *zap.SugaredLogger) *mygorm.DB {
	log.Infow("init gorm staring...")
	defer log.Infow("init gorm end...")
	return mygorm.New(db)
}

func initAtomicLevel(conf *conf.Conf) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(zapcore.Level(conf.Log.Level))
}

func initLog(conf *conf.Conf, level zap.AtomicLevel) *zap.SugaredLogger {
	log := mylog.New(conf.App.Name, conf.Log.Dir, level)
	return log
}

func initMux(conf *conf.Conf) *mux.Router {
	m := mux.NewRouter()
	m.Use(otelmux.Middleware(conf.App.Name))
	return m
}

func initRouter(r *mux.Router, log *zap.SugaredLogger, level zap.AtomicLevel) *myrouter.Router {
	log.Infow("init router staring...")
	defer log.Infow("init router end...")

	myR := myrouter.New(r)
	myR.Setlog(log)
	myR.Methods(http.MethodPut).Path("/log").Handler(level)
	return myR
}

func initConsul(conf *conf.Conf, log *zap.SugaredLogger) (*api.Client, error) {
	log.Infow("init consul staring...")
	defer log.Infow("init consul end...")
	cfg := api.DefaultConfig()
	cfg.Address = conf.Consul.Addr
	cfg.Token = conf.Consul.Token
	return api.NewClient(cfg)
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

func initHttpServiceOptions(log *zap.SugaredLogger) []httpkit.ServerOption {
	options := make([]httpkit.ServerOption, 0)
	options = append(options,
		httpkit.ServerErrorHandler(myhttp.NewErrorHandler(log)),
		httpkit.ServerErrorEncoder(myhttp.ErrorEncoder),
		httpkit.ServerBefore(
			httpkit.PopulateRequestContext,
		),
	)

	return options
}

func initEndpointMiddleware() (res []endpoint.Middleware) {
	return
}
