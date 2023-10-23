package framework

import "net/http"

func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		f(ctx.Writer, ctx.Request)
	}
}

func WrapH(h http.Handler) HandlerFunc {
	return func(ctx *Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func ApplyMiddleware(handler HandlerFunc, middlewares []MiddlewareFunc) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
