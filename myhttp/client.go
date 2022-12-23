package myhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
)

type Client struct {
	*req.Client
}

func NewClient(client *req.Client, options ...ClientOption) *Client {
	c := &Client{Client: client}
	for _, option := range options {
		option(c)
	}
	return c
}

type contextKey int

const HttpRequestName contextKey = 1

type Request struct {
	*req.Request
	before []BeforeFunc
	after  []AfterFunc
}

func (r *Request) WithName(name string) {
	ctx := context.WithValue(r.Context(), HttpRequestName, name)
	r.SetContext(ctx)
}

func (r *Request) DecodeEx(method, url string, decode Decode, options ...RequestOption) (*req.Response, error) {
	for _, option := range options {
		option(r)
	}

	for _, before := range r.before {
		before(r.Request)
	}
	resp, err := r.Request.Send(method, url)

	for _, after := range r.after {
		after(resp, err)
	}

	return decode(resp, err)
}

func (c *Client) R() *Request {
	r := c.Client.R()

	return &Request{
		Request: r,
	}
}

type Decode func(res *req.Response, err error) (*req.Response, error)

func DecodeJsonData(i interface{}) Decode {
	return func(resp *req.Response, err error) (*req.Response, error) {
		if resp.GetStatusCode() != http.StatusOK {
			return resp, fmt.Errorf("unexpected status code %d", resp.GetStatusCode())
		}

		if i == nil {
			return resp, nil
		}

		if _, ok := i.(*string); ok {
			i = resp.String()
			return resp, nil
		}

		body := resp.Bytes()
		result := gjson.GetBytes(body, "code")
		if !result.Exists() {
			err := json.Unmarshal(body, i)
			if err != nil {
				err = errors.Wrap(err, "unmarshal response")
				return resp, err
			}
			return resp, nil
		}

		if result.Int() != http.StatusOK {
			s := gjson.Get(resp.String(), "err").String()
			return resp, fmt.Errorf("response err: %s", s)
		}

		dataResult := gjson.GetBytes(body, "data")
		err = json.Unmarshal(body[dataResult.Index:result.Index+len(result.Raw)], i)
		if err != nil {
			err = errors.Wrap(err, "unmarshal response data")
			return resp, err
		}

		return resp, nil
	}
}
