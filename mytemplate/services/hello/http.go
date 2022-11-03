package hello

import (
	"context"
	"github.com/fitan/mykit/myhttp"
	"github.com/go-kit/kit/endpoint"
	http1 "net/http"

	govalidator "github.com/asaskevich/govalidator"
	http "github.com/go-kit/kit/transport/http"
	mux "github.com/gorilla/mux"
	errors "github.com/pkg/errors"
)

type Handler struct{}

func MakeHTTPHandler(r *mux.Router, s Service, mws Mws, ops Ops) Handler {
	eps := NewEndpoint(s, mws)
	r.Handle("/hello/{id}", http.NewServer(eps.HelloEndpoint, decodeHelloRequest, http.EncodeJSONResponse, ops[HelloMethodName]...)).Methods("GET").Name("@kit-http /hello/{id} GET")

	return Handler{}
}

type Ops map[string][]http.ServerOption

func AllMethodAddOps(options map[string][]http.ServerOption, option http.ServerOption) {
	methods := []string{HelloMethodName}
	for _, v := range methods {
		options[v] = append(options[v], option)
	}
}

type HttpKit struct {
	service service
}

func (h *HttpKit) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HelloRequest)
		res, err := h.service.Hello(ctx, req.ID, req.Query)
		return myhttp.WrapResponse(res, err)
	}
}

func (h *HttpKit) DecodeRequest(ctx context.Context, r *http1.Request) (res interface{}, err error) {
	req := HelloRequest{}

	var age string

	var email string

	var id string

	vars := mux.Vars(r)

	id = vars["id"]

	age = r.URL.Query().Get("age")

	email = r.URL.Query().Get("email")

	req.ID = id

	req.Query.Age = age

	req.Query.Email = email

	validReq, err := govalidator.ValidateStruct(req)

	if err != nil {
		err = errors.Wrap(err, "govalidator.ValidateStruct")
		return
	}

	if !validReq {
		err = errors.Wrap(err, "valid false")
		return
	}

	return req, err
}

/*


Hello
@Summary @kit-http /hello/{id} GET
@Description @kit-http /hello/{id} GET

@Accept json
@Produce json
@Param id path string true
@Param age query string false
@Param email query string false
@Success 200 {object} encode.Response{data=HelloResponse}
@Router /hello/{id} [GET]
*/
func decodeHelloRequest(ctx context.Context, r *http1.Request) (res interface{}, err error) {

	req := HelloRequest{}

	var age string

	var email string

	var id string

	vars := mux.Vars(r)

	id = vars["id"]

	age = r.URL.Query().Get("age")

	email = r.URL.Query().Get("email")

	req.ID = id

	req.Query.Age = age

	req.Query.Email = email

	validReq, err := govalidator.ValidateStruct(req)

	if err != nil {
		err = errors.Wrap(err, "govalidator.ValidateStruct")
		return
	}

	if !validReq {
		err = errors.Wrap(err, "valid false")
		return
	}

	return req, err
}
