package util_router

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type middleware func(h http.HandlerFunc) http.HandlerFunc
type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}
func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}
func (r *Router) Add(route string, h http.HandlerFunc) {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	r.mux[route] = mergedHandler
}
func (r *Router) AddHandler(route string, h http.Handler) {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler.ServeHTTP)
	}
	r.mux[route] = mergedHandler
}
func (r *Router) Run(addr string) {
	mux := http.NewServeMux()
	for k, v := range r.mux {
		mux.Handle(k, v)
	}
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	logrus.Infoln("start http server at ", addr)
	_ = server.ListenAndServe()
}
