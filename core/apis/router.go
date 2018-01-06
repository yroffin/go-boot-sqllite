// Package apis for all exposed api
// MIT License
//
// Copyright (c) 2017 yroffin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package apis

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yroffin/go-boot-sqllite/core/bean"
)

// Router internal members
type Router struct {
	// Base component
	*bean.Bean
	// mux router
	Router *mux.Router
}

// IRouter Test all package methods
type IRouter interface {
	bean.IBean
}

// Init Init this API
func (p *Router) Init() error {
	return nil
}

// PostConstruct Init this API
func (p *Router) PostConstruct(name string) error {
	log.Printf("Router::PostConstruct - router creation")
	// define all routes
	p.Router = mux.NewRouter()
	// Fix default handler
	p.Router.MethodNotAllowedHandler = http.HandlerFunc(p.HandlerStaticNotAllowed())
	p.Router.NotFoundHandler = http.HandlerFunc(p.HandlerStaticNotFound())
	return nil
}

// Validate Init this API
func (p *Router) Validate(name string) error {
	log.Printf("Router::Validate - router validation")
	// handle now all requests
	http.Handle("/", p.Router)
	return nil
}

// HandleFunc declare a handler
func (p *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), method string, content string) {
	log.Printf("Router::HandleFunc '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	var res = p.Router.HandleFunc(path, p.HandlerStatic(f)).Methods(method)
	if len(content) > 0 {
		res.Headers("Content-Type", content)
	}
}

// HandleFuncString declare a string handler
func (p *Router) HandleFuncString(path string, f func() (string, error), method string, content string) {
	log.Printf("Router::HandleFunc '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	var res = p.Router.HandleFunc(path, p.HandlerStaticString(f)).Methods(method)
	if len(content) > 0 {
		res.Headers("Content-Type", content)
	}
}

// HandlerStaticNotFound Not found handler
func (p *Router) HandlerStaticNotFound() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Url %v", r.URL)
		log.Printf("Request Headers %v", r.Header)
		log.Printf("Request Encoding %v", r.TransferEncoding)
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Not found")
	}
	return anonymous
}

// HandlerStaticNotFound Not found handler
func (p *Router) HandlerStaticNotAllowed() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Url %v", r.URL)
		log.Printf("Request Headers %v", r.Header)
		log.Printf("Request Encoding %v", r.TransferEncoding)
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(405)
		fmt.Fprintf(w, "Not allowed")
	}
	return anonymous
}

// HandlerStatic render string
func (p *Router) HandlerStatic(method func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		// security header
		w.Header().Set("Strict-Transport-Security", "")
		w.Header().Set("Content-Security-Policy", "")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		method(w, r)
	}
	return anonymous
}

// HandlerStaticString render string
func (p *Router) HandlerStaticString(method func() (string, error)) func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		// security header
		w.Header().Set("Strict-Transport-Security", "")
		w.Header().Set("Content-Security-Policy", "")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		// content
		w.Header().Set("Content-type", "text/html")
		data, err := method()
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(data))
	}
	return anonymous
}
