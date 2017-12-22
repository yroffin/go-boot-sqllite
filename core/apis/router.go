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
	log.Printf("Router::HandleFunc %s with method %s", path, method)
	// declare it to the router
	p.Router.HandleFunc(path, f).Methods(method).Headers("Content-Type", content)
}
