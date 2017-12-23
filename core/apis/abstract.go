// Package apis for common interfaces
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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/business"
	"github.com/yroffin/go-boot-sqllite/core/models"
)

// API base class
type API struct {
	// members
	*bean.Bean
	// mux router
	Router *mux.Router
	// all mthods to declare
	methods []APIMethod
	// Router with injection mecanism
	SetRouterBean func(interface{}) `bean:"router"`
	RouterBean    *Router
	// Router with injection mecanism
	SetCrudBusiness func(interface{}) `bean:"crud-business"`
	CrudBusiness    *business.CrudBusiness
	// Crud
	HandlerGetByID    func(id string) (string, error)
	HandlerGetAll     func() (string, error)
	HandlerPost       func(body string) (string, error)
	HandlerPutByID    func(id string, body string) (string, error)
	HandlerDeleteByID func(id string) (string, error)
	HandlerPatchByID  func(id string, body string) (string, error)
}

// APIMethod single structure to modelise api declaration
type APIMethod struct {
	path    string
	handler string
	method  string
	// mime type
	typeMime string
	addr     reflect.Value
}

// APIInterface all package methods
type APIInterface interface {
	bean.IBean
	Declare(APIMethod, interface{})
	HandlerStatic() func(w http.ResponseWriter, r *http.Request)
	HandlerGetByID(id string) (string, error)
}

// ScanHandler this API
func (p *API) ScanHandler(ptr interface{}) {
	// define all methods
	types := reflect.TypeOf(ptr).Elem()
	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		// declare a standard mux handler
		if len(field.Tag.Get("handler")) > 0 {
			p.append(ptr, field.Tag.Get("path"), field.Tag.Get("handler"), field.Tag.Get("method"), field.Tag.Get("mime-type"))
		}
		// declare a crud handler
		if strings.Contains(field.Name, "crud") {
			p.append(ptr, field.Tag.Get("path")+"/{id:[0-9a-zA-Z-_]*}", "HandlerStaticGetByID", "GET", "application/json")
			p.append(ptr, field.Tag.Get("path"), "HandlerStaticGetAll", "GET", "application/json")
			p.append(ptr, field.Tag.Get("path"), "HandlerStaticPost", "POST", "application/json")
			p.append(ptr, field.Tag.Get("path")+"/{id:[0-9a-zA-Z-_]*}", "HandlerStaticPutByID", "PUT", "application/json")
			p.append(ptr, field.Tag.Get("path")+"/{id:[0-9a-zA-Z-_]*}", "HandlerStaticDeleteByID", "DELETE", "application/json")
			p.append(ptr, field.Tag.Get("path")+"/{id:[0-9a-zA-Z-_]*}", "HandlerStaticPatchByID", "PATCH", "application/json")
		}
	}
	// call bean init
	p.Init()
}

// Init initialize the API
func (p *API) append(ptr interface{}, path string, handler string, method string, mime string) {
	addr := reflect.ValueOf(ptr).MethodByName(handler)
	if addr.IsNil() {
		log.Fatalf("Unable to find any method called '%v'", handler)
	} else {
		log.Printf("Successfully mounted method called '%v' on path '%s' with method '%s' - '%s'", handler, path, method, mime)
	}
	p.methods = append(p.methods, APIMethod{path: path, handler: handler, method: method, addr: addr, typeMime: mime})
}

// Init initialize the APIf
func (p *API) Init() error {
	// inject SlideBusiness
	p.SetCrudBusiness = func(value interface{}) {
		if assertion, ok := value.(*business.CrudBusiness); ok {
			p.CrudBusiness = assertion
		} else {
			log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
		}
	}
	// inject RouterBean
	p.SetRouterBean = func(value interface{}) {
		if assertion, ok := value.(*Router); ok {
			p.RouterBean = assertion
		} else {
			log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
		}
	}
	// build arguments
	arr := [1]reflect.Value{reflect.ValueOf(p)}
	var arguments = arr[1:1]
	// build all static acess to low level function (private)
	for i := 0; i < len(p.methods); i++ {
		// compute rvalue
		var rvalue = p.methods[i].addr.Call(arguments)[0]
		// declare this new method
		p.Declare(p.methods[i], rvalue.Interface())
	}
	return nil
}

// PostConstruct this API
func (p *API) PostConstruct(name string) error {
	return p.Bean.PostConstruct(name)
}

// Declare a new interface
func (p *API) Declare(data APIMethod, intf interface{}) error {
	// verify type
	if value, ok := intf.(func(http.ResponseWriter, *http.Request)); ok {
		log.Printf("Declare handler() '%s' on '%s' with method '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (*p.RouterBean).GetName(), data.typeMime)
		// declare it to the router
		(*p.RouterBean).HandleFunc(data.path, value, data.method, data.typeMime)
		return nil
	}
	// verify type
	if value, ok := intf.(func() (string, error)); ok {
		log.Printf("Declare function() '%s' on '%s' with method '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (*p.RouterBean).GetName(), data.typeMime)
		// declare it to the router
		(*p.RouterBean).HandleFuncString(data.path, value, data.method, data.typeMime)
		return nil
	}
	// Error case
	return errors.New("Unable to find any type for " + data.handler)
}

// HandlerStaticGetAll is the GET by ID handler
func (p *API) HandlerStaticGetAll() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		data, err := p.HandlerGetAll()
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// HandlerStaticGetByID is the GET by ID handler
func (p *API) HandlerStaticGetByID() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		data, err := p.HandlerGetByID(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// HandlerStaticPost is the POST handler
func (p *API) HandlerStaticPost() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		body, _ := ioutil.ReadAll(r.Body)
		data, err := p.HandlerPost(string(body))
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(201)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// HandlerStaticPutByID is the PUT by ID handler
func (p *API) HandlerStaticPutByID() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		body, _ := ioutil.ReadAll(r.Body)
		data, err := p.HandlerPutByID(mux.Vars(r)["id"], string(body))
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// HandlerStaticDeleteByID is the DELETE by ID handler
func (p *API) HandlerStaticDeleteByID() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		data, err := p.HandlerDeleteByID(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// HandlerStaticPatchByID is the PATCH by ID handler
func (p *API) HandlerStaticPatchByID() func(w http.ResponseWriter, r *http.Request) {
	anonymous := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		body, _ := ioutil.ReadAll(r.Body)
		data, err := p.HandlerPatchByID(mux.Vars(r)["id"], string(body))
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"message\":\"\"}")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, data)
	}
	return anonymous
}

// GenericGetAll default method
func (p *API) GenericGetAll(toGet models.IPersistent, toGets models.IPersistents) (string, error) {
	p.CrudBusiness.GetAll(toGet, toGets)
	var arr = toGets.Get()
	data, _ := json.Marshal(&arr)
	return string(data), nil
}

// GenericGetByID default method
func (p *API) GenericGetByID(id string, toGet models.IPersistent) (string, error) {
	toGet.SetID(id)
	p.CrudBusiness.Get(toGet)
	data, _ := json.Marshal(&toGet)
	return string(data), nil
}

// GenericPost adefault method
func (p *API) GenericPost(body string, toCreate models.IPersistent) (string, error) {
	var bin = []byte(body)
	json.Unmarshal(bin, &toCreate)
	bean, _ := p.CrudBusiness.Create(toCreate)
	data, _ := json.Marshal(&bean)
	return string(data), nil
}

// GenericPutByID default method
func (p *API) GenericPutByID(id string, body string, toUpdate models.IPersistent) (string, error) {
	toUpdate.SetID(id)
	var bin = []byte(body)
	json.Unmarshal(bin, &toUpdate)
	bean, _ := p.CrudBusiness.Update(toUpdate)
	data, _ := json.Marshal(&bean)
	return string(data), nil
}

// GenericPatchByID default method
func (p *API) GenericPatchByID(id string, body string, toPatch models.IPersistent) (string, error) {
	toPatch.SetID(id)
	var bin = []byte(body)
	json.Unmarshal(bin, &toPatch)
	bean, _ := p.CrudBusiness.Patch(toPatch)
	data, _ := json.Marshal(&bean)
	return string(data), nil
}

// GenericDeleteByID default method
func (p *API) GenericDeleteByID(id string, toDelete models.IPersistent) (string, error) {
	toDelete.SetID(id)
	old, _ := p.CrudBusiness.Delete(toDelete)
	data, _ := json.Marshal(&old)
	return string(data), nil
}
