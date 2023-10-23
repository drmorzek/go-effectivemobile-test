package framework

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	routes           map[string]map[string]HandlerFunc
	middlewares      map[string][]MiddlewareFunc
	globalMiddleware []MiddlewareFunc
	TemplateDir      string
}

func (r *Router) SetTemplateDir(dir string) {
	r.TemplateDir = dir
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func NewRouter() *Router {

	return &Router{
		routes:           make(map[string]map[string]HandlerFunc),
		middlewares:      make(map[string][]MiddlewareFunc),
		globalMiddleware: make([]MiddlewareFunc, 0),
	}
}

func (r *Router) Run(path string) {

	log.Fatal(http.ListenAndServe(path, r))

}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.globalMiddleware = append(r.globalMiddleware, middleware)
}

func (r *Router) UseForRoute(route string, middleware MiddlewareFunc) {
	r.middlewares[route] = append(r.middlewares[route], middleware)
}

func (r *Router) ServeFiles(path string, root http.FileSystem) {
	fileServer := http.FileServer(root)

	r.AddRoute("GET", path+"/*filepath", func(c *Context) {
		reqPath := c.Request.URL.Path
		if strings.Contains(reqPath, "..") {
			http.NotFound(c.Writer, c.Request)
			return
		}
		reqPath = strings.TrimPrefix(reqPath, path)
		c.Request.URL.Path = reqPath
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (r *Router) AddRoute(method string, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	if _, exists := r.routes[method]; !exists {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = ApplyMiddleware(handler, middlewares)
}

func (r *Router) GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.AddRoute("GET", path, handler, middlewares...)
}
func (r *Router) POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.AddRoute("POST", path, handler, middlewares...)
}
func (r *Router) PUT(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.AddRoute("PUT", path, handler, middlewares...)
}
func (r *Router) DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.AddRoute("DELETE", path, handler, middlewares...)
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var handler HandlerFunc
	params := make(map[string]string)
	// Проверка на наличие точных совпадений маршрутов
	if routes, ok := r.routes[req.Method]; ok {
		if h, ok := routes[req.URL.Path]; ok {
			handler = h
		} else {
			for path, h := range routes {
				routeParts := strings.Split(path, "/")
				requestParts := strings.Split(req.URL.Path, "/")
				if len(routeParts) != len(requestParts) && !strings.HasSuffix(path, "/*") {
					continue
				}
				matches := true
				for i := range routeParts {
					if routeParts[i] == requestParts[i] || (len(routeParts[i]) > 0 && routeParts[i][0] == ':') {
						if len(routeParts[i]) > 0 && routeParts[i][0] == ':' {
							params[routeParts[i][1:]] = requestParts[i]
						}
						continue
					} else if routeParts[i] == "*" {
						params["*"] = strings.Join(requestParts[i:], "/")
						break
					}
					matches = false
					break
				}
				if matches {
					handler = h
					break
				}
			}
		}
	}

	if handler != nil {
		ctx := NewContext(w, req, &r)
		ctx.Params = params
		handler = ApplyMiddleware(handler, r.middlewares[req.URL.Path])
		handler = ApplyMiddleware(handler, r.globalMiddleware)
		handler(ctx)
	} else {
		http.NotFound(w, req)
	}
}
