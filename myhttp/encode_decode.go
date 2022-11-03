package myhttp

import (
	"encoding/json"
	"net/http"
)

// 默认返回对象
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   error       `json:"message"`
	TraceID string      `json:"trace_id"`
}

func WrapResponse(data interface{}, err error) (interface{}, error) {
	return Response{
		Data:  data,
		Error: err,
	}, err
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
