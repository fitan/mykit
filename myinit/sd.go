package myinit

import (
	"context"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"go-micro.dev/v4/util/addr"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type zapSugarLogger func(msg string, keysAndValues ...interface{})

func (l zapSugarLogger) Log(kv ...interface{}) error {
	l("", kv...)
	return nil
}

type SD struct {
	id           string
	consulClient *api.Client
	reg          *consulsd.Registrar
	sig          chan os.Signal
	log          *zap.SugaredLogger
	cancel       context.CancelFunc
	ctx          context.Context
}

func (s *SD) Register() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s.sig = sig
	s.log.Infof("registering service %s", s.id)
	s.reg.Register()
	s.log.Infow("register service success")
	go func() {
		for {
			select {
			case <-sig:
				s.log.Infow("consul deregister ing...")
				s.Deregister()
				s.log.Infow("consul deregister end...")
				s.cancel()
			case <-time.After(time.Second * 10):
				err := s.consulClient.Agent().PassTTL(s.id, "pass")
				if err != nil {
					s.log.Errorw("consul pass ttl error", "err", err)
				} else {
					s.log.Infow("consul pass ttl success")
				}
			}
		}
	}()
}

func (s *SD) Deregister() {
	s.reg.Deregister()
}

func (s *SD) Wait() {
	<-s.ctx.Done()
	s.log.Infow("sd done...")
	return
}

func InitSD(name, ip string, port int, client *api.Client, log *zap.SugaredLogger) (*SD, error) {
	log.Infow("init consul sd staring...")
	defer log.Infow("init consul sd end...")
	sd := consulsd.NewClient(client)

	if ip == "" {
		extract, err := addr.Extract("")
		if err != nil {
			err = errors.Wrap(err, "addr.Extract")
			return nil, err
		}

		ip = extract
	}
	id := ip + ":" + strconv.Itoa(port)

	reg := consulsd.NewRegistrar(sd, &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    []string{},
		Port:    port,
		Address: ip,
		Check: &api.AgentServiceCheck{
			CheckID:                        id,
			TTL:                            "15s",
			DeregisterCriticalServiceAfter: "15s",
		},
	}, zapSugarLogger(log.Infow))

	ctx, cancel := context.WithCancel(context.Background())

	return &SD{
		id:           id,
		consulClient: client,
		reg:          reg,
		sig:          nil,
		log:          log,
		cancel:       cancel,
		ctx:          ctx,
	}, nil
}
