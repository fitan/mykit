package myhttp

import (
	"context"
	"github.com/go-resty/resty/v2"
	"net/http"
)

// CreateRequestFunc creates an outgoing HTTP request based on the passed
// request object. It's designed to be used in HTTP clients, for client-side
// endpoints. It's a more powerful version of EncodeRequestFunc, and can be used
// if more fine-grained control of the HTTP request is required.
type CreateRequestFunc func(context.Context, interface{}) (*http.Request, error)

type RestyEncodeRequestFunc func(context.Context, *resty.Request, interface{}) error

type RestyCreateRequestFunc func(context.Context, interface{}) (*resty.Request, error)

type RestyDecodeResponseFunc func(context.Context, *resty.Response) (response interface{}, err error)

type RestyRequestFunc func(context.Context, *resty.Request) context.Context
