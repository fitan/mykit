package myhttp

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type RestyResponseFunc func(context.Context, *resty.Response) context.Context

type ClientFinalizerFunc func(ctx context.Context, err error)
