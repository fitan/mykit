package mycrud

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type KitHttpImpl interface {
	KitHttpBaseImpl
	GetDecode() kithttp.DecodeRequestFunc
	GetEndpoint() endpoint.Endpoint
}

type KitHttpBaseImpl interface {
	GetName() string
	GetHttpMethod() string
	GetEncode() kithttp.EncodeResponseFunc
	GetOptions() []kithttp.ServerOption
	GetHttpPath() string
	GetEndpointMid() []endpoint.Middleware
}

type KitHttpConfig struct {
	Name        string
	HttpMethod  string
	HttpPath    string
	Encode      kithttp.EncodeResponseFunc
	Options     []kithttp.ServerOption
	EndpointMid []endpoint.Middleware
}

func (c *KitHttpConfig) GetHttpPath() string {
	return c.HttpPath
}

func (c *KitHttpConfig) GetName() string {
	return c.Name
}

func (c *KitHttpConfig) GetHttpMethod() string {
	return c.HttpMethod
}

func (c *KitHttpConfig) GetEncode() kithttp.EncodeResponseFunc {
	return c.Encode
}

func (c *KitHttpConfig) GetOptions() []kithttp.ServerOption {
	return c.Options
}

func (c *KitHttpConfig) GetEndpointMid() []endpoint.Middleware {
	return c.EndpointMid
}

//type CrudService struct {
//	GetOne     KitHttpImpl
//	GetMany    KitHttpImpl
//	CreateOne  KitHttpImpl
//	CreateMany KitHttpImpl
//	UpdateOne  KitHttpImpl
//	UpdateMany KitHttpImpl
//	DeleteOne  KitHttpImpl
//	DeleteMany KitHttpImpl
//
//	GetRelationOne  KitHttpImpl
//	GetRelationMany KitHttpImpl
//
//	CreateRelationOne  KitHttpImpl
//	CreateRelationMany KitHttpImpl
//}
//
//func (c *CrudService) Handler() {
//	if c.GetOne != nil {
//		c.GetOne.Handler()
//	}
//	if c.GetMany != nil {
//		c.GetMany.Handler()
//	}
//
//	if c.CreateOne != nil {
//		c.CreateOne.Handler()
//	}
//
//	if c.CreateMany != nil {
//		c.CreateMany.Handler()
//	}
//
//	if c.UpdateOne != nil {
//		c.UpdateOne.Handler()
//	}
//
//	if c.UpdateMany != nil {
//		c.UpdateMany.Handler()
//	}
//
//	if c.DeleteOne != nil {
//		c.DeleteOne.DeleteOneHandler()
//	}
//
//	if c.DeleteMany != nil {
//		c.DeleteMany.Handler()
//	}
//
//	if c.GetRelationOne != nil {
//		c.GetRelationOne.Handler()
//	}
//
//	if c.GetRelationMany != nil {
//		c.GetRelationMany.Handler()
//	}
//
//	if c.CreateRelationOne != nil {
//		c.CreateRelationOne.Handler()
//	}
//
//	if c.CreateRelationMany != nil {
//		c.CreateRelationMany.Handler()
//	}
//
//}
