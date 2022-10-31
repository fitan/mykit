package myhttp

import (
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"runtime"
	"time"
)

type ClientOption func(*Client)

func WithTrace() ClientOption {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return func(client *Client) {
		client.SetTransport(otelhttp.NewTransport(t))
	}
}

type RequestOption func(*Request)

func RequestBefore(before ...BeforeFunc) RequestOption {
	return func(r *Request) {
		r.before = append(r.before, before...)
	}
}

func RequestAfter(after ...AfterFunc) RequestOption {
	return func(r *Request) {
		r.after = append(r.after, after...)
	}
}

type BeforeFunc func(*resty.Request)

type AfterFunc func(*resty.Response, error)
