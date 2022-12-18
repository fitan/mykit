package hello

import (
	"context"
	http1 "net/http"
	"strings"

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

func MethodAddOps(options map[string][]http.ServerOption, option http.ServerOption, methods []string) {
	for _, v := range methods {
		options[v] = append(options[v], option)
	}
}

/*


Hello
@Summary @kit-http /hello/{id} GET
@Description @kit-http /hello/{id} GET

@Accept json
@Produce json
@Param id path string true
@Param age query string false
@Param between_time query string false
@Param email query string false
@Param idIn query string false
@Success 200 {object} encode.Response{data=HelloResponse}
@Router /hello/{id} [GET]
*/
func decodeHelloRequest(ctx context.Context, r *http1.Request) (res interface{}, err error) {

	req := HelloRequest{}

	var _age string

	var _email string

	var _idIn []string

	var _between_time []string

	var _id string

	vars := mux.Vars(r)

	_id = vars["id"]

	_age = r.URL.Query().Get("age")

	_email = r.URL.Query().Get("email")

	_idIn = strings.Split(r.URL.Query().Get("idIn"), ",")

	_between_time = strings.Split(r.URL.Query().Get("between_time"), ",")

	req.ID = _id

	req.Query.IDIn = _idIn

	req.Query.BetweenTime = _between_time

	req.Query.Age = _age

	req.Query.Email = _email

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
