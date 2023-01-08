package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type UpdateOneImpl interface {
	UpdateOneHandler()
	UpdateOneDecode() kithttp.DecodeRequestFunc
	UpdateOneEndpoint() endpoint.Endpoint
}

type UpdateOne struct {
	Repo *Repo
	*KitHttpConfig
}

type UpdateOneRequest struct {
	Id   string      `json:"id"`
	Body interface{} `json:"body"`
}

func (u *UpdateOne) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := UpdateOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]

		body := u.Repo.TableMsg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder(r.Body).Decode(body)")
			return
		}

		req.Body = body
		return req, err

	}
}

func (u *UpdateOne) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateOneRequest)
		err = u.Repo.UpdateOne(ctx, req.Id, req.Body)
		return nil, err
	}
}
