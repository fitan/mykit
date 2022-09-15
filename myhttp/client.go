package myhttp

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-resty/resty/v2"
	"io"
)

type Client struct {
	client *resty.Client
}

func (c *Client) Endpoint(url UrlOption, enc RestyEncodeRequestFunc, dec RestyDecodeResponseFunc, options ...RequestOption) endpoint.Endpoint {
	r := c.client.R()
	url(r)
	request := &Request{
		req:            makeCreateRequestFunc(r, enc),
		dec:            dec,
		before:         make([]RestyRequestFunc, 0),
		after:          make([]RestyResponseFunc, 0),
		finalizer:      make([]ClientFinalizerFunc, 0),
		bufferedStream: false,
	}
	for _, option := range options {
		option(request)
	}

	return request.Endpoint()
}

type UrlOption func(c *resty.Request)

type RequestOption func(request *Request)

type Request struct {
	req            RestyCreateRequestFunc
	dec            RestyDecodeResponseFunc
	before         []RestyRequestFunc
	after          []RestyResponseFunc
	finalizer      []ClientFinalizerFunc
	bufferedStream bool
}

func (r *Request) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx, cancel := context.WithCancel(ctx)

		var (
			resp *resty.Response
			err  error
		)
		if r.finalizer != nil {
			defer func() {
				//if resp != nil {
				//	ctx = context.WithValue(ctx, ContextKeyResponseHeaders, resp.Header)
				//	ctx = context.WithValue(ctx, ContextKeyResponseSize, resp.Size())
				//}
				for _, f := range r.finalizer {
					f(ctx, err)
				}
			}()
		}

		req, err := r.req(ctx, request)
		if err != nil {
			cancel()
			return nil, err
		}

		for _, f := range r.before {
			ctx = f(ctx, req)
		}

		resp, err = req.SetContext(ctx).Send()

		if err != nil {
			cancel()
			return nil, err
		}

		if r.bufferedStream {
			resp.RawResponse.Body = bodyWithCancel{ReadCloser: resp.RawResponse.Body, cancel: cancel}
		} else {
			defer resp.RawResponse.Body.Close()
			defer cancel()
		}

		for _, f := range r.after {
			ctx = f(ctx, resp)
		}

		response, err := r.dec(ctx, resp)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

func makeCreateRequestFunc(req *resty.Request, enc RestyEncodeRequestFunc) RestyCreateRequestFunc {
	return func(ctx context.Context, i interface{}) (*resty.Request, error) {
		err := enc(ctx, req, i)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}

type bodyWithCancel struct {
	io.ReadCloser

	cancel context.CancelFunc
}

func (bwc bodyWithCancel) Close() error {
	bwc.ReadCloser.Close()
	bwc.cancel()
	return nil
}
