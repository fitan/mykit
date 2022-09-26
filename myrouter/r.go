package myrouter

import (
	"context"
	"fmt"
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
	"log"
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
	log *zap.SugaredLogger
}

func (r *Router) gops() {
	err := agent.Listen(agent.Options{})
	if err != nil {
		r.log.Panic(err)
	}
}

func (r *Router) walk() {
	v := make([][]string, 0, 0)
	v = append(v, []string{"id", "name", "path", "method"})
	id := 0
	r.Router.Walk(
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

func (r *Router) metric() {
	mdlw := metricsMid.New(metricsMid.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	r.Use(
		func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				path, _ := mux.CurrentRoute(r).GetPathTemplate()
				std.Handler(path, mdlw, next).ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		})

	go func() {
		r.log.Infof("metrics listening at: %s", ":8090")
		if err := http.ListenAndServe(":8090", promhttp.Handler()); err != nil {
			r.log.Panicf("error while serving metrics: %s", err)
		}
	}()
}

func (r *Router) Run(addr string) {
	if addr == "" {
		addr = "localhost:8080"
	}

	fmt.Println("start r.walk")
	r.walk()
	fmt.Println("r.walk")

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
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func (r *Router) switchDebug() {
	debug := myhttpmid.NewDebugSwitch()
	r.Use(debug.Middleware)
	debug.Handlers("/debug", r.Router)
}

func (r *Router) Setlog(log *zap.SugaredLogger) {
	r.log = log
}

func New() *Router {
	r := mux.NewRouter()
	zaplog, _ := zap.NewProduction()
	router := &Router{
		Router: r,
		log:    zaplog.Sugar(),
	}

	r.Use(handlers.RecoveryHandler())
	router.health()
	router.metric()
	router.switchDebug()
	router.gops()
	fmt.Println("return router")
	return router
}
