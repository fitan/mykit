package handlers

import (
	"github.com/fitan/mykit/myhttp"
	"github.com/gorilla/mux"
	"net/http"
)

func Handlers(r *mux.Router) {
	r.HandleFunc("/echo/{id}", Echo).Methods(http.MethodGet).Name("测试")
}

type EchoRequest struct {
	ID string `json:"id" path:"id"`
}

type EchoResponse struct {
	ID string `json:"id"`
}

func Echo(w http.ResponseWriter, request *http.Request) {
	req := EchoRequest{}

	err := myhttp.BindAndValidate(&req, request)
	if err != nil {
		myhttp.ResponseJsonEncode(w, err.Error())
		return
	}

	myhttp.ResponseJsonEncode(w, EchoResponse{req.ID})
}
