package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"net/http"
)

type UpdateMany struct {
	Repo *Repo
	*KitHttpConfig
}

type UpdateManyRequest struct {
	Body interface{} `json:"body"`
}

func (u *UpdateMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := UpdateManyRequest{}

		body := u.Repo.TableMsg.manyObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder.Decode")
			return
		}

		req.Body = body
		return req, nil
	}
}

func (u *UpdateMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateManyRequest)
		err = u.Repo.UpdateMany(ctx, req.Body)
		return nil, err
	}
}
