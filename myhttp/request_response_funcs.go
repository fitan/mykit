package myhttp

import (
	"context"
	"github.com/gin-gonic/gin"
)

// RequestFunc may take information from an HTTP request and put it into a
// request context. In Servers, RequestFuncs are executed prior to invoking the
// endpoint. In Clients, RequestFuncs are executed after creating the request
// but prior to invoking the HTTP client.
type RequestFunc func(context.Context, *gin.Context) context.Context

// ServerResponseFunc may take information from a request context and use it to
// manipulate a ResponseWriter. ServerResponseFuncs are only executed in
// servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(context.Context, *gin.Context) context.Context

const (
	KeyRequestMethod          = "KeyRequestMethod"
	KeyRequestURI             = "KeyRequestURI"
	KeyRequestPath            = "KeyRequestPath"
	KeyRequestFullPath        = "KeyRequestFullPath"
	KeyRequestProto           = "KeyRequestProto"
	KeyRequestHost            = "KeyRequestHost"
	KeyRequestRemoteAddr      = "KeyRequestRemoteAddr"
	KeyRequestXForwardedFor   = "KeyRequestXForwardedFor"
	KeyRequestXForwardedProto = "KeyRequestXForwardedProto"
	KeyRequestAuthorization   = "KeyRequestAuthorization"
	KeyRequestReferer         = "KeyRequestReferer"
	KeyRequestUserAgent       = "KeyRequestUserAgent"
	KeyRequestXRequestID      = "KeyRequestXRequestID"
	KeyRequestAccept          = "KeyRequestAccept"
)

func PopulateRequestGinKey(ctx context.Context, r *gin.Context) context.Context {
	for k, v := range map[string]string{
		KeyRequestMethod:          r.Request.Method,
		KeyRequestURI:             r.Request.RequestURI,
		KeyRequestPath:            r.Request.URL.Path,
		KeyRequestFullPath:        r.FullPath(),
		KeyRequestProto:           r.Request.Proto,
		KeyRequestHost:            r.Request.Host,
		KeyRequestRemoteAddr:      r.Request.RemoteAddr,
		KeyRequestXForwardedFor:   r.GetHeader("X-Forwarded-For"),
		KeyRequestXForwardedProto: r.GetHeader("X-Forwarded-Proto"),
		KeyRequestAuthorization:   r.GetHeader("Authorization"),
		KeyRequestReferer:         r.GetHeader("Referer"),
		KeyRequestUserAgent:       r.GetHeader("User-Agent"),
		KeyRequestXRequestID:      r.GetHeader("X-Request-Id"),
		KeyRequestAccept:          r.GetHeader("Accept"),
	} {
		r.Set(k, v)
	}
	return ctx
}

// PopulateRequestContext is a RequestFunc that populates several values into
// the context from the HTTP request. Those values may be extracted using the
// corresponding ContextKey type in this package.
func PopulateRequestContext(ctx context.Context, r *gin.Context) context.Context {
	for k, v := range map[contextKey]string{
		ContextKeyRequestMethod:          r.Request.Method,
		ContextKeyRequestURI:             r.Request.RequestURI,
		ContextKeyRequestPath:            r.Request.URL.Path,
		ContextKeyRequestFullPath:        r.FullPath(),
		ContextKeyRequestProto:           r.Request.Proto,
		ContextKeyRequestHost:            r.Request.Host,
		ContextKeyRequestRemoteAddr:      r.Request.RemoteAddr,
		ContextKeyRequestXForwardedFor:   r.GetHeader("X-Forwarded-For"),
		ContextKeyRequestXForwardedProto: r.GetHeader("X-Forwarded-Proto"),
		ContextKeyRequestAuthorization:   r.GetHeader("Authorization"),
		ContextKeyRequestReferer:         r.GetHeader("Referer"),
		ContextKeyRequestUserAgent:       r.GetHeader("User-Agent"),
		ContextKeyRequestXRequestID:      r.GetHeader("X-Request-Id"),
		ContextKeyRequestAccept:          r.GetHeader("Accept"),
	} {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

type contextKey int

const (
	// ContextKeyRequestMethod is populated in the context by
	// PopulateRequestContext. Its value is r.Method.
	ContextKeyRequestMethod contextKey = iota

	ContextKeyRequestDebug

	// ContextKeyRequestURI is populated in the context by
	// PopulateRequestContext. Its value is r.RequestURI.
	ContextKeyRequestURI

	// ContextKeyRequestPath is populated in the context by
	// PopulateRequestContext. Its value is r.URL.Path.
	ContextKeyRequestPath

	ContextKeyRequestFullPath

	// ContextKeyRequestProto is populated in the context by
	// PopulateRequestContext. Its value is r.Proto.
	ContextKeyRequestProto

	// ContextKeyRequestHost is populated in the context by
	// PopulateRequestContext. Its value is r.Host.
	ContextKeyRequestHost

	// ContextKeyRequestRemoteAddr is populated in the context by
	// PopulateRequestContext. Its value is r.RemoteAddr.
	ContextKeyRequestRemoteAddr

	// ContextKeyRequestXForwardedFor is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-For").
	ContextKeyRequestXForwardedFor

	// ContextKeyRequestXForwardedProto is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-Proto").
	ContextKeyRequestXForwardedProto

	// ContextKeyRequestAuthorization is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Authorization").
	ContextKeyRequestAuthorization

	// ContextKeyRequestReferer is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Referer").
	ContextKeyRequestReferer

	// ContextKeyRequestUserAgent is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("User-Agent").
	ContextKeyRequestUserAgent

	// ContextKeyRequestXRequestID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Request-Id").
	ContextKeyRequestXRequestID

	// ContextKeyRequestAccept is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Accept").
	ContextKeyRequestAccept

	// ContextKeyResponseHeaders is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type http.Header, and
	// is captured only once the entire response has been written.
	ContextKeyResponseHeaders

	// ContextKeyResponseSize is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type int64.
	ContextKeyResponseSize
)
