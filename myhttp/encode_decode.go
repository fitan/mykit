package myhttp

import (
	"encoding/json"
	"net/http"
)

func ResponseJsonEncode(w http.ResponseWriter, v interface{}) {
	res, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
