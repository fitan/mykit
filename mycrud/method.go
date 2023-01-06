package mycrud

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

// method path is /{tableName}/{id}/action/{methodName}
// relation method is /{tableName}/{id}/{relationTableName}/action/{methodName}

func (c *Core) MethodHandler() {

	c.Handler("methodName", http.MethodPost, " /{tableName}/{id}/action/{methodName}")

}

type MethodRequest struct {
	HttpMethod string `json:"httpMethod"`
	TableName  string `json:"tableName"`
	Id         string `json:"id"`
	MethodName string `json:"methodName"`
	Body       []byte `json:"body"`
}

func (c *Core) MethodDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := MethodRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.MethodName = r.Method
		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

	}
}

func (c *Core) MethodEndpoint() {

}
