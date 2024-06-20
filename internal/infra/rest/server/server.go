package server

import (
	"fmt"
	"net/http"
	"slices"
)

type middleware func(http.Handler) http.Handler

type Server struct {
	*http.ServeMux
	chain []middleware
}

func New(midds ...middleware) *Server {
	return &Server{
		ServeMux: &http.ServeMux{},
		chain:    midds,
	}
}

func (s *Server) Use(midds ...middleware) {
	s.chain = append(s.chain, midds...)
}

func (s *Server) Group(fn func(s *Server)) {
	fn(&Server{ServeMux: s.ServeMux, chain: slices.Clone(s.chain)})
}

func (s *Server) Post(path string, fn http.HandlerFunc, midds ...middleware) {
	s.handle(http.MethodPost, path, fn, midds)
}
func (s *Server) Get(path string, fn http.HandlerFunc, midds ...middleware) {
	s.handle(http.MethodGet, path, fn, midds)
}
func (s *Server) Put(path string, fn http.HandlerFunc, midds ...middleware) {
	s.handle(http.MethodPut, path, fn, midds)
}
func (s *Server) Patch(path string, fn http.HandlerFunc, midds ...middleware) {
	s.handle(http.MethodPatch, path, fn, midds)
}
func (s *Server) Delete(path string, fn http.HandlerFunc, midds ...middleware) {
	s.handle(http.MethodDelete, path, fn, midds)
}

func (s *Server) handle(method, path string, fn http.HandlerFunc, midds []middleware) {
	s.Handle(fmt.Sprintf("%s %s", method, path), s.wrapper(fn, midds))
}

func (s *Server) wrapper(fn http.Handler, midds []middleware) http.Handler {
	output := fn
	midds = append(slices.Clone(s.chain), midds...)

	slices.Reverse(midds)

	for _, m := range midds {
		output = m(output)
	}
	return output
}
