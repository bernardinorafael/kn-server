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

func New(m ...middleware) *Server {
	return &Server{
		ServeMux: &http.ServeMux{},
		chain:    m,
	}
}

func (s *Server) Use(m ...middleware) {
	s.chain = append(s.chain, m...)
}

func (s *Server) Group(fn func(s *Server)) {
	fn(&Server{ServeMux: s.ServeMux, chain: slices.Clone(s.chain)})
}

func (s *Server) Post(path string, fn http.HandlerFunc, m ...middleware) {
	s.handle(http.MethodPost, path, fn, m)
}
func (s *Server) Get(path string, fn http.HandlerFunc, m ...middleware) {
	s.handle(http.MethodGet, path, fn, m)
}
func (s *Server) Put(path string, fn http.HandlerFunc, m ...middleware) {
	s.handle(http.MethodPut, path, fn, m)
}
func (s *Server) Patch(path string, fn http.HandlerFunc, m ...middleware) {
	s.handle(http.MethodPatch, path, fn, m)
}
func (s *Server) Delete(path string, fn http.HandlerFunc, m ...middleware) {
	s.handle(http.MethodDelete, path, fn, m)
}

func (s *Server) handle(method, path string, fn http.HandlerFunc, m []middleware) {
	s.Handle(fmt.Sprintf("%s %s", method, path), s.wrapper(fn, m))
}

func (s *Server) wrapper(fn http.Handler, m []middleware) http.Handler {
	output := fn
	m = append(slices.Clone(s.chain), m...)

	slices.Reverse(m)

	for _, m := range m {
		output = m(output)
	}
	return output
}
