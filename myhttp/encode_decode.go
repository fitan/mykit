package myhttp

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

// 默认返回对象
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	TraceID string      `json:"trace_id"`
}

func WrapResponse(data interface{}, err error) (interface{}, error) {
	res := Response{Data: data, Code: 200}

	if err != nil {
		res.Error = err.Error()
		res.Code = 500
	}

	return res, nil
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := Response{}
	res.Data = response
	res.Code = 200

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	_, _ = w.Write(b)
	return nil
}

func KitErrorEncoder() kithttp.ServerOption {
	return kithttp.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
		res := Response{Error: err.Error(), Code: 500}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(b)
	})
}

func ResponseJsonEncode(w http.ResponseWriter, v interface{}) {
	res, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
