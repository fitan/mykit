package myrouter

import (
	"context"
	"github.com/arl/statsviz"
	"github.com/fitan/mykit/myhttp"
	"github.com/fitan/mykit/myhttpmid"
	"github.com/google/gops/agent"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/pterm/pterm"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMid "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Router struct {
	*mux.Router
	log    *zap.SugaredLogger
	debugM *mux.Router
}

func (r *Router) gops() {
	err := agent.Listen(agent.Options{})
	if err != nil {
		r.log.Panic(err)
	}
}

func (r *Router) walk(m *mux.Router) {
	v := make([][]string, 0, 0)
	v = append(v, []string{"id", "name", "path", "method"})
	id := 0
	m.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, _ := route.GetPathTemplate()
			method, _ := route.GetMethods()
			id = id + 1
			v = append(v, []string{strconv.Itoa(id), route.GetName(), path, strings.Join(method, ",")})

			return nil
		})
	pterm.DefaultTable.WithHasHeader().WithData(v).Render()
}

func (r *Router) health() {
	r.Router.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		myhttp.ResponseJsonEncode(writer, "ok")
	}).Methods(http.MethodGet)
}

func (r *Router) statsviz(m *mux.Router) {
	m.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(statsviz.Ws)
	m.Methods("GET").Path("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(statsviz.Index)
}

func (r *Router) metric(m *mux.Router) {
	mdlw := metricsMid.New(metricsMid.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	r.Router.Use(
		func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				path, _ := mux.CurrentRoute(r).GetPathTemplate()
				std.Handler(path, mdlw, next).ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		})

	m.Methods(http.MethodGet).Name("prom metrics").Path("/metrics").Handler(promhttp.Handler())
}

func (r *Router) debugRouterRun() {
	r.metric(r.debugM)
	r.statsviz(r.debugM)
	r.switchDebug(r.debugM)
	r.walk(r.debugM)
	r.log.Infof("debug listening at: %s", ":8090")
	if err := http.ListenAndServe(":8090", r.debugM); err != nil {
		r.log.Panicf("error while serving metrics: %s", err)
	}
}

func (r *Router) Run(addr string) {
	go r.debugRouterRun()

	if addr == "" {
		addr = "localhost:8080"
	}

	r.walk(r.Router)

	server := &http.Server{
		Addr:    addr,
		Handler: handlers.LoggingHandler(os.Stdout, r.Router),
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				r.log.Panic("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			r.log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	r.log.Infof("server listening at %s", addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		r.log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	r.log.Infow("server done...")
}

func (r *Router) switchDebug(m *mux.Router) {
	debug := myhttpmid.NewDebugSwitch()
	r.Use(debug.Middleware)
	debug.Handlers(m)
}

func (r *Router) SetLog(log *zap.SugaredLogger) {
	r.log = log
}

// RecoveryHandler log
func (r *Router) Println(args ...interface{}) {
	r.log.Errorw("recovery", args...)
}

func New(r *mux.Router) *Router {
	zaplog, _ := zap.NewProduction()
	router := &Router{
		Router: r,
		log:    zaplog.Sugar().WithOptions(zap.AddCallerSkip(1)),
		debugM: mux.NewRouter(),
	}

	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true), handlers.RecoveryLogger(router)))
	router.health()
	router.gops()
	return router
}
