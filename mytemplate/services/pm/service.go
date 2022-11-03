package pm

import (
	"context"
	"github.com/fitan/mykit/mytemplate/repo"
	"github.com/pkg/errors"
)

type Middleware func(service Service) Service

// @tags hello
// @impl
type Service interface {
	// @kit-http /pm GET
	// @kit-http-request ListRequest
	List(ctx context.Context, page, limit int) (list []ListResponse, total int64, err error)
}

type service struct {
	repo *repo.Repo
}

func (s *service) List(ctx context.Context, page, limit int) (list []ListResponse, total int64, err error) {
	res, total, err := s.repo.Pm.List(ctx, page, limit, "", nil)
	if err != nil {
		err = errors.Wrap(err, "repo.Pm.List")
		return
	}
	// @call copy
	list = listDTO(res)

	return
}

func New(repo *repo.Repo) Service {
	return &service{repo}
}
