package myhttp

import (
	"context"
	"encoding/json"
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
	res, ok := response.(Response)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("response type error"))
		return nil
	}

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
