package request

import (
	"context"
	"net/http"

	"github.com/maddsua/eventdb-next/utils"
)

type RequestContext struct {
	Req       *http.Request
	Writer    http.ResponseWriter
	ClientIP  string
	RequestID string
}

type contextKeyType struct{}

func InjectConext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), contextKeyType{}, &RequestContext{
			Req:       req,
			Writer:    writer,
			ClientIP:  utils.ClientIP(req),
			RequestID: utils.RequestID(writer, req),
		})))
	})
}

func From(ctx context.Context) *RequestContext {
	return ctx.Value(contextKeyType{}).(*RequestContext)
}
