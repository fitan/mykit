package myhttp

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/transport/http"
	"net/url"
	"strings"
)

type KitHttpClient struct {
	host string
	opts []http.ClientOption
}

func NewKitHttpClient() *KitHttpClient {
	return &KitHttpClient{}
}

func (k *KitHttpClient) SetHost(host string) *KitHttpClient {
	k.host = host
	return k
}

func (k *KitHttpClient) SetOpt(option http.ClientOption) *KitHttpClient {
	k.opts = append(k.opts, option)
	return k
}

func (k *KitHttpClient) R(method, path string, enc http.EncodeRequestFunc, dec http.DecodeResponseFunc) *KitHttpRequest {
	return &KitHttpRequest{
		host:   k.host,
		path:   path,
		params: map[string]string{},
		method: method,
		query:  map[string]string{},
		body:   nil,
		enc:    enc,
		dec:    dec,
		opts:   k.opts,
	}
}

type KitHttpRequest struct {
	host   string
	path   string
	params map[string]string
	method string
	query  map[string]string
	body   interface{}
	enc    http.EncodeRequestFunc
	dec    http.DecodeResponseFunc
	opts   []http.ClientOption
}

func (k *KitHttpRequest) SetHost(host string) *KitHttpRequest {
	k.host = host
	return k
}

func (k *KitHttpRequest) SetParams(key, value string) *KitHttpRequest {
	k.params[key] = value
	return k
}

func (k *KitHttpRequest) SetQuery(key, value string) *KitHttpRequest {
	k.query[key] = value
	return k
}

func (k *KitHttpRequest) SetBody(body interface{}) *KitHttpRequest {
	k.body = body
	return k
}

func (k *KitHttpRequest) SetOption(option http.ClientOption) *KitHttpRequest {
	k.opts = append(k.opts, option)
	return k
}

func (k *KitHttpRequest) Exec(ctx context.Context) (res interface{}, err error) {
	urlStr := fmt.Sprintf("%s/%s", k.host, k.path)

	for key, value := range k.params {
		urlStr = strings.ReplaceAll(urlStr, fmt.Sprintf("{%s}", key), value)
	}

	query := url.Values{}
	for key, value := range k.query {
		query.Add(key, value)
	}

	u, err := url.Parse(fmt.Sprintf("%s?%s", urlStr, query.Encode()))
	if err != nil {
		err = fmt.Errorf("parse url error: %w", err)
		return
	}

	return http.NewClient(k.method, u, k.enc, k.dec, k.opts...).Endpoint()(ctx, k.body)
}
