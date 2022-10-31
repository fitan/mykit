package myhttp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
)

type Client struct {
	*resty.Client
}

func NewClient(client *resty.Client, options ...ClientOption) *Client {
	c := &Client{Client: client}
	for _, option := range options {
		option(c)
	}
	return c
}

type Request struct {
	*resty.Request
	before []BeforeFunc
	after  []AfterFunc
}

func (r *Request) DecodeEx(method, url string, decode Decode, options ...RequestOption) (*resty.Response, error) {
	for _, option := range options {
		option(r)
	}

	for _, before := range r.before {
		before(r.Request)
	}
	resp, err := r.Request.Execute(method, url)

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

type Decode func(res *resty.Response, err error) (*resty.Response, error)

func DecodeJsonData(i interface{}) Decode {
	return func(resp *resty.Response, err error) (*resty.Response, error) {
		if resp.StatusCode() != http.StatusOK {
			return resp, fmt.Errorf("unexpected status code %d", resp.StatusCode())
		}

		if _, ok := i.(*string); ok {
			i = resp.String()
			return resp, nil
		}

		result := gjson.GetBytes(resp.Body(), "code")
		if !result.Exists() {
			err := json.Unmarshal(resp.Body(), i)
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

		dataResult := gjson.GetBytes(resp.Body(), "data")
		err = json.Unmarshal(resp.Body()[dataResult.Index:result.Index+len(result.Raw)], i)
		if err != nil {
			err = errors.Wrap(err, "unmarshal response data")
			return resp, err
		}

		return resp, nil
	}
}
