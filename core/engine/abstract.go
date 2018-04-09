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
package engine

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yroffin/go-boot-sqllite/core/models"
	"github.com/yroffin/go-boot-sqllite/core/winter"
)

var (
	// Tables Fix tables names
	Tables = []string{"Node"}
)

func init() {
	winter.Helper.Register("graph-crud-business", (&GraphCrudBusiness{}).New())
	winter.Helper.Register("sql-crud-business", (&SqlCrudBusiness{}).New())
	winter.Helper.Register("cayley-manager", (&Graph{}).New("./cayley.db"))
	winter.Helper.Register("sqllite-manager", (&Store{}).New("./sqllite.db"))
}

// API base class
type API struct {
	*winter.Bean
	// all mthods to declare
	methods []APIMethod
	// Router with injection mecanism
	Router IRouter `@autowired:"router"`
	// SqlCrudBusiness with injection mecanism
	SQLCrudBusiness ICrudBusiness `@autowired:"sql-crud-business"`
	// GraphBusiness with injection mecanism
	GraphBusiness ILinkBusiness `@autowired:"graph-crud-business"`
	// Factory
	Factory          func() models.IPersistent
	Factories        func() models.IPersistents
	HandlerTasks     func(name string, body string) (interface{}, int, error)
	HandlerTasksByID func(id string, name string, body string) (interface{}, int, error)
}

// APIMethod single structure to modelise api declaration
type APIMethod struct {
	path    string
	handler string
	method  string
	target  IAPI
	// mime type
	typeMime string
	addr     reflect.Value
	// Fields
	summary string
	desc    string
	// method handler
	in    []interface{}
	out   map[string]interface{}
	args  map[string]interface{}
	query map[string]interface{}
}

// CrudHandler single structure to modelise api declaration
type CrudHandler interface {
	HandlerPost(body string) (interface{}, error)
}

// IAPI all package methods
type IAPI interface {
	winter.IBean
	Declare(APIMethod, interface{}) error
	// Data handled by this API
	GetFactory() models.IPersistent
	GetFactories() models.IPersistents
	// All
	GetAll() ([]models.IPersistent, error)
	// Links
	GetAllLinks(id string, targetType IAPI) ([]models.IPersistent, error)
	LoadAllLinks(name string, factory func() models.IPersistent, targetType IAPI) (interface{}, int, error)
}

// GetFactory return on new bean
func (p *API) GetFactory() models.IPersistent {
	if p.Factory != nil {
		return p.Factory()
	}
	log.Printf("Factory is nil for %v\n", p.GetName())
	return nil
}

// GetFactories return a bean list
func (p *API) GetFactories() models.IPersistents {
	if p.Factories != nil {
		return p.Factories()
	}
	return nil
}

// Call params
func Call(params ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return in
}

// ScanHandler this API
func (p *API) ScanHandler(swagger ISwaggerService, ptr interface{}) {
	// define all methods
	types := reflect.TypeOf(ptr).Elem()
	values := reflect.ValueOf(ptr).Elem()
	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		value := values.Field(i)
		// declare a standard mux handler
		if len(field.Tag.Get("@handler")) > 0 {
			log.Println("Info:", field.Name, "handler")
			p.add(ptr, field.Tag.Get("path"), field.Tag.Get("@handler"), field.Tag.Get("method"), field.Tag.Get("mime-type"), "", "", map[string]interface{}{}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{})
		}
		// declare a crud handler
		if len(field.Tag.Get("@crud")) > 0 {
			assert, conv := ptr.(IAPI)
			if conv {
				log.Println("Info:", field.Name, "is API")
				p.add(ptr, field.Tag.Get("@crud"), "HandlerStaticGetAll", "GET", "application/json", "Get all", "Get all resources", map[string]interface{}{}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactories()})
				p.add(ptr, field.Tag.Get("@crud"), "HandlerStaticPost", "POST", "application/json", "Execute a task or create", "Execute a task on all resources", map[string]interface{}{}, map[string]interface{}{"task": "params"}, []interface{}{assert.GetFactory()}, map[string]interface{}{"200": assert.GetFactory()})
				p.add(ptr, field.Tag.Get("@crud")+"/:id", "HandlerStaticGetByID", "GET", "application/json", "Get by id", "Get a resource by its id", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactory()})
				p.add(ptr, field.Tag.Get("@crud")+"/:id", "HandlerStaticPutByID", "PUT", "application/json", "Update by id", "Update a resource by its id", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{assert.GetFactory()}, map[string]interface{}{"200": assert.GetFactory()})
				p.add(ptr, field.Tag.Get("@crud")+"/:id", "HandlerStaticDeleteByID", "DELETE", "application/json", "Delete by id", "Delete a resource by its id", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactory()})
				p.add(ptr, field.Tag.Get("@crud")+"/:id", "HandlerStaticPatchByID", "PATCH", "application/json", "Patch by id", "Patch a resource by its id", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{assert.GetFactory()}, map[string]interface{}{"200": assert.GetFactory()})
				p.add(ptr, field.Tag.Get("@crud")+"/:id", "HandlerStaticPostByID", "POST", "application/json", "Execute a task", "Execute a new task on resource", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{assert.GetFactory()}, map[string]interface{}{"201": assert.GetFactory()})
			} else {
				log.Println("Warning:", field.Name, "is not API", reflect.TypeOf(value))
			}
		}
		// declare a link handler
		if len(field.Tag.Get("@link")) > 0 {
			assert, conv := value.Interface().(IAPI)
			if conv {
				var linkName = field.Tag.Get("@href")
				log.Println("Info:", field.Name, "is API/HREF", assert.GetName())
				p.addLink(ptr, field.Tag.Get("@link")+"/:id/"+linkName, "HandlerLinkStaticGetAll", "GET", "application/json", "Get all", "Get all resources", map[string]interface{}{"id": "Id"}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactories()}, assert)
				p.addLink(ptr, field.Tag.Get("@link")+"/:id/"+linkName+"/:link", "HandlerLinkStaticGetByID", "GET", "application/json", "Get by id", "Get a resource by its id", map[string]interface{}{"id": "Id", "link": "Link"}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactory()}, assert)
				p.addLink(ptr, field.Tag.Get("@link")+"/:id/"+linkName+"/:link", "HandlerLinkStaticPostByID", "POST", "application/json", "Update by id", "Update a resource by its id", map[string]interface{}{"id": "Id", "link": "Link"}, map[string]interface{}{}, []interface{}{assert.GetFactory()}, map[string]interface{}{"200": assert.GetFactory()}, assert)
				p.addLink(ptr, field.Tag.Get("@link")+"/:id/"+linkName+"/:link", "HandlerLinkStaticPutByID", "PUT", "application/json", "Update by id", "Update a resource by its id", map[string]interface{}{"id": "Id", "link": "Link"}, map[string]interface{}{}, []interface{}{assert.GetFactory()}, map[string]interface{}{"200": assert.GetFactory()}, assert)
				p.addLink(ptr, field.Tag.Get("@link")+"/:id/"+linkName+"/:link", "HandlerLinkStaticDeleteByID", "DELETE", "application/json", "Delete by id", "Delete a resource by its id", map[string]interface{}{"id": "Id", "link": "Link"}, map[string]interface{}{}, []interface{}{}, map[string]interface{}{"200": assert.GetFactory()}, assert)
			} else {
				log.Println("Warning:", field.Name, "is not API/HREF", reflect.TypeOf(value.Interface()))
			}
		}
	}
	// Add method to swagger
	for k := range p.methods {
		method := p.methods[k]
		var typ = strings.SplitAfter(reflect.TypeOf(ptr).String(), ".")
		swagger.AddPaths(typ[1], p.methods[k].path, p.methods[k].method, p.methods[k].summary, p.methods[k].desc, method.args, method.query, method.in, method.out)
	}
	// call bean init
	p.Init()
}

// Init initialize the API
func (p *API) add(ptr interface{}, path string, handler string, method string, mime string, sum string, desc string, args map[string]interface{}, query map[string]interface{}, in []interface{}, out map[string]interface{}) {
	addr := reflect.ValueOf(ptr).MethodByName(handler)
	if !addr.IsValid() || addr.IsNil() {
		log.Fatalf("Unable to find any method called '%v'", handler)
	} else {
		log.Printf("Successfully mounted method called '%v' on path '%s' with method '%s' - '%s'", handler, path, method, mime)
	}
	p.methods = append(p.methods, APIMethod{path: path, handler: handler, method: method, addr: addr, typeMime: mime, summary: sum, desc: desc, args: args, query: query, in: in, out: out})
}

// Init initialize the API
func (p *API) addLink(ptr interface{}, path string, handler string, method string, mime string, sum string, desc string, args map[string]interface{}, query map[string]interface{}, in []interface{}, out map[string]interface{}, target IAPI) {
	addr := reflect.ValueOf(ptr).MethodByName(handler)
	if !addr.IsValid() || addr.IsNil() {
		log.Fatalf("Unable to find any method called '%v'", handler)
	} else {
		log.Printf("Successfully mounted method called '%v' on path '%s' with method '%s' - '%s'", handler, path, method, mime)
	}
	p.methods = append(p.methods, APIMethod{path: path, handler: handler, method: method, addr: addr, typeMime: mime, summary: sum, desc: desc, args: args, query: query, in: in, out: out, target: target})
}

// Init initialize the APIf
func (p *API) Init() error {
	// build arguments
	arr := [1]reflect.Value{reflect.ValueOf(p)}
	var arguments = arr[1:1]
	// build all static acess to low level function (private)
	for i := 0; i < len(p.methods); i++ {
		log.Printf("Build static link for %v with %v", p.methods[i].path, arguments)
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

// Validate this API
func (p *API) Validate(name string) error {
	return p.Bean.Validate(name)
}

// Declare a new interface
func (p *API) Declare(data APIMethod, intf interface{}) error {
	// verify type
	if value, ok := intf.(func(c *gin.Context, target IAPI)); ok {
		log.Printf("Declare handler() '%s' on '%s' with method '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (p.Router).GetName(), data.typeMime)
		// declare it to the router
		(p.Router).HandleFuncLink(data.path, value, data.method, data.typeMime, data.target)
		return nil
	}
	// verify type
	if value, ok := intf.(func(c *gin.Context)); ok {
		log.Printf("Declare handler() '%s' on '%s' with method '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (p.Router).GetName(), data.typeMime)
		// declare it to the router
		(p.Router).HandleFunc(data.path, value, data.method, data.typeMime)
		return nil
	}
	// verify type
	if value, ok := intf.(func() (string, error)); ok {
		log.Printf("Declare function() '%s' on '%s' with method '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (p.Router).GetName(), data.typeMime)
		// declare it to the router
		(p.Router).HandleFuncString(data.path, value, data.method, data.typeMime)
		return nil
	}
	// verify type
	if value, ok := intf.(func(string) (string, error)); ok {
		log.Printf("Declare function() '%s' on '%s' with method with id '%s' ('%s') with type '%s'", data.handler, data.path, data.method, (p.Router).GetName(), data.typeMime)
		// declare it to the router
		(p.Router).HandleFuncStringWithId(data.path, value, data.method, data.typeMime)
		return nil
	}
	// Error case
	return errors.New("Unable to find any type for " + data.handler)
}

// HandlerStaticGetAll is the GET by ID handler
func (p *API) HandlerStaticGetAll() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		data, err := p.GetAll()
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerStaticGetByID is the GET by ID handler
func (p *API) HandlerStaticGetByID() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		data, err := p.GetByID(c.Param("id"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerStaticPost is the POST handler
func (p *API) HandlerStaticPost() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		if len(c.Query("task")) > 0 {
			data, count, err := p.HandlerTasks(c.Query("task"), string(body))
			if err != nil {
				c.String(400, "{\"message\":\"\"}")
				return
			}
			p.XTotalCount(c, count)
			c.IndentedJSON(202, data)
		} else {
			_, ok := c.Request.URL.Query()["filter"]
			if ok {
				var objmap map[string]string
				err := json.Unmarshal(body, &objmap)
				data, err := p.HandlerFilter(objmap)
				if err != nil {
					c.String(400, "{\"message\":\"\"}")
					return
				}
				c.IndentedJSON(200, data)
			} else {
				data, err := p.HandlerPost(string(body))
				if err != nil {
					c.String(400, "{\"message\":\"\"}")
					return
				}
				c.IndentedJSON(201, data)
			}
		}
	}
	return anonymous
}

// XTotalCount handle X-Total-Count
func (p *API) XTotalCount(c *gin.Context, count int) {
	// handle X-total-count
	if count >= 0 {
		c.Header("X-total-count", strconv.Itoa(count))
	}
}

// HandlerStaticPostByID is the POST handler
func (p *API) HandlerStaticPostByID() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		if len(c.Query("task")) > 0 {
			if c.Param("id") == "*" {
				// id with * is like post on all resources
				data, count, err := p.HandlerTasks(c.Query("task"), string(body))
				if err != nil {
					c.String(400, "{\"message\":\"\"}")
					return
				}
				p.XTotalCount(c, count)
				c.IndentedJSON(202, data)
			} else {
				data, count, err := p.HandlerTasksByID(c.Param("id"), c.Query("task"), string(body))
				if err != nil {
					c.String(400, "{\"message\":\"\"}")
					return
				}
				p.XTotalCount(c, count)
				c.IndentedJSON(202, data)
			}
		} else {
			data, err := p.HandlerPost(string(body))
			if err != nil {
				c.String(400, "{\"message\":\"\"}")
				return
			}
			c.IndentedJSON(201, data)
		}
	}
	return anonymous
}

// HandlerStaticPutByID is the PUT by ID handler
func (p *API) HandlerStaticPutByID() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		data, err := p.HandlerPutByID(c.Param("id"), string(body))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerStaticDeleteByID is the DELETE by ID handler
func (p *API) HandlerStaticDeleteByID() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		data, err := p.HandlerDeleteByID(c.Param("id"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerStaticPatchByID is the PATCH by ID handler
func (p *API) HandlerStaticPatchByID() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		data, err := p.HandlerPatchByID(c.Param("id"), string(body))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerLinkStaticGetAll is the GET by ID handler
func (p *API) HandlerLinkStaticGetAll() func(c *gin.Context, targetType IAPI) {
	anonymous := func(c *gin.Context, targetType IAPI) {
		c.Header("Content-type", "application/json")
		data, err := p.GetAllLinks(c.Param("id"), targetType)
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerLinkStaticGetByID is the GET by ID handler
func (p *API) HandlerLinkStaticGetByID() func(c *gin.Context, targetType string) {
	anonymous := func(c *gin.Context, targetType string) {
		c.Header("Content-type", "application/json")
		data, err := p.GetByID(c.Param("id"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerLinkStaticPostByID is the PUT by ID handler
func (p *API) HandlerLinkStaticPostByID() func(c *gin.Context, targetType IAPI) {
	anonymous := func(c *gin.Context, targetType IAPI) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		data, err := p.HandlerLinkPostByID(c.Param("id"), c.Param("link"), string(body), targetType)
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerLinkStaticPutByID is the PUT by ID handler
func (p *API) HandlerLinkStaticPutByID() func(c *gin.Context, targetType IAPI) {
	anonymous := func(c *gin.Context, targetType IAPI) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		data, err := p.HandlerLinkPutByID(c.Param("id"), c.Param("link"), string(body), targetType, c.Query("instance"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// HandlerLinkStaticDeleteByID is the DELETE by ID handler
func (p *API) HandlerLinkStaticDeleteByID() func(c *gin.Context, targetType IAPI) {
	anonymous := func(c *gin.Context, targetType IAPI) {
		c.Header("Content-type", "application/json")
		body, _ := ioutil.ReadAll(c.Request.Body)
		data, err := p.HandlerLinkDeleteByID(c.Param("id"), c.Param("link"), string(body), targetType, c.Query("instance"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.IndentedJSON(200, data)
	}
	return anonymous
}

// GetAll get all
func (p *API) GetAll() ([]models.IPersistent, error) {
	return p.GenericGetAll(p.Factory(), p.Factories())
}

// GetByID get by id
func (p *API) GetByID(id string) (models.IPersistent, error) {
	return p.GenericGetByID(id, p.Factory())
}

// HandlerPost create handler
func (p *API) HandlerPost(body string) (interface{}, error) {
	return p.GenericPost(body, p.Factory())
}

// HandlerFilter task handler for filter
func (p *API) HandlerFilter(body map[string]string) (models.IPersistents, error) {
	return nil, nil
}

// HandlerPutByID update by id
func (p *API) HandlerPutByID(id string, body string) (interface{}, error) {
	return p.GenericPutByID(id, body, p.Factory())
}

// HandlerDeleteByID delete by id
func (p *API) HandlerDeleteByID(id string) (interface{}, error) {
	return p.GenericDeleteByID(id, p.Factory())
}

// HandlerPatchByID pach by id
func (p *API) HandlerPatchByID(id string, body string) (interface{}, error) {
	return p.GenericPatchByID(id, body, p.Factory())
}

// HandlerLinkPostByID update by id
func (p *API) HandlerLinkPostByID(src string, dst string, body string, targetType IAPI) (models.IPersistent, error) {
	source := p.Factory()
	p.GenericGetByID(src, source)
	target := targetType.GetFactory()
	p.GenericGetByID(dst, target)
	log.Println("output", target)
	toCreate := (&models.EdgeBean{}).New(source.GetEntityName(), source.GetID(), targetType.GetName(), target.GetID(), "HREF")
	// add edge extended data
	var ext = make(map[string]interface{})
	json.Unmarshal([]byte(body), &ext)
	toCreate.Extend(ext)
	_, err := p.GenericLinkPostByID(toCreate)
	// edge is reserved keyword
	delete(ext, "edge")
	ext["instance"] = toCreate.GetID()
	target.Extend(ext)
	return target, err
}

// HandlerLinkPutByID update by id
func (p *API) HandlerLinkPutByID(src string, dst string, body string, targetType IAPI, instance string) (models.IPersistent, error) {
	source := p.Factory()
	p.GenericGetByID(src, source)
	target := targetType.GetFactory()
	p.GenericGetByID(dst, target)
	toUpdate := (&models.EdgeBean{}).New(source.GetEntityName(), source.GetID(), targetType.GetName(), target.GetID(), "HREF")
	// add edge extended data, edge and instance are reserved keyword
	var ext = make(map[string]interface{})
	json.Unmarshal([]byte(body), &ext)
	toUpdate.Extend(ext)
	toUpdate.SetInstance(instance)
	_, err := p.GenericLinkPutByID(toUpdate)
	// edge is reserved keyword
	delete(ext, "edge")
	target.Extend(ext)
	ext["instance"] = toUpdate.GetID()
	return target, err
}

// HandlerLinkDeleteByID update by id
func (p *API) HandlerLinkDeleteByID(src string, dst string, body string, targetType IAPI, instance string) (interface{}, error) {
	source := p.Factory()
	p.GenericGetByID(src, source)
	target := p.Factory()
	p.GenericGetByID(dst, target)
	toDelete := &models.EdgeBean{}
	json.Unmarshal([]byte(body), toDelete)
	toDelete.SetInstance(instance)
	return p.GenericLinkDeleteByID(toDelete)
}

// GetAllLinks get all
func (p *API) GetAllLinks(id string, targetType IAPI) ([]models.IPersistent, error) {
	return p.GenericLinkGetAll(id, make([]models.IEdgeBean, 0), targetType)
}

// GenericGetAll default method
func (p *API) GenericGetAll(toGet models.IPersistent, toGets models.IPersistents) ([]models.IPersistent, error) {
	p.SQLCrudBusiness.GetAll(toGet, toGets)
	return toGets.Get(), nil
}

// GenericGetByID default method
func (p *API) GenericGetByID(id string, toGet models.IPersistent) (models.IPersistent, error) {
	toGet.SetID(id)
	p.SQLCrudBusiness.Get(toGet)
	return toGet, nil
}

// GenericPost adefault method
func (p *API) GenericPost(body string, toCreate models.IPersistent) (interface{}, error) {
	var bin = []byte(body)
	result := json.Unmarshal(bin, &toCreate)
	log.Println("JSON", body, models.ToJSON(toCreate))
	// check unmashal errors
	if result != nil {
		log.Printf("Error, while Unmarshaling body %v - %v", body, result)
		return body, result
	}
	bean, _ := p.SQLCrudBusiness.Create(toCreate)
	return bean, nil
}

// GenericPutByID default method
func (p *API) GenericPutByID(id string, body string, toUpdate models.IPersistent) (interface{}, error) {
	toUpdate.SetID(id)
	var bin = []byte(body)
	json.Unmarshal(bin, &toUpdate)
	bean, _ := p.SQLCrudBusiness.Update(toUpdate)
	return bean, nil
}

// GenericPatchByID default method
func (p *API) GenericPatchByID(id string, body string, toPatch models.IPersistent) (interface{}, error) {
	toPatch.SetID(id)
	var bin = []byte(body)
	json.Unmarshal(bin, &toPatch)
	bean, _ := p.SQLCrudBusiness.Patch(toPatch)
	return bean, nil
}

// GenericDeleteByID default method
func (p *API) GenericDeleteByID(id string, toDelete models.IPersistent) (interface{}, error) {
	toDelete.SetID(id)
	p.SQLCrudBusiness.Get(toDelete)
	old, _ := p.SQLCrudBusiness.Delete(toDelete)
	return old, nil
}

// GenericLinkPutByID default method
func (p *API) GenericLinkPostByID(assoc models.IEdgeBean) (interface{}, error) {
	bean, _ := p.GraphBusiness.CreateLink(assoc)
	return bean, nil
}

// GenericLinkPutByID default method
func (p *API) GenericLinkPutByID(assoc models.IEdgeBean) (interface{}, error) {
	bean, _ := p.GraphBusiness.UpdateLink(assoc)
	return bean, nil
}

// GenericLinkDeleteByID default method
func (p *API) GenericLinkDeleteByID(assoc models.IEdgeBean) (interface{}, error) {
	bean, _ := p.GraphBusiness.DeleteLink(assoc)
	return bean, nil
}

// GenericLinkGetAll default method
func (p *API) GenericLinkGetAll(id string, links []models.IEdgeBean, targetType IAPI) ([]models.IPersistent, error) {
	// Retrieve all links
	edges, _ := p.GraphBusiness.GetAllLink(p.GetFactory().GetEntityName(), id, links, targetType.GetName())
	// Build output
	output := make([]models.IPersistent, 0)
	for _, edge := range edges {
		// Retrive bean
		t := targetType.GetFactory()
		// Filter by type
		if edge.GetTarget() == t.GetEntityName() {
			t.SetID(edge.GetTargetID())
			p.SQLCrudBusiness.Get(t)
			ex := make(map[string]interface{})
			ex["instance"] = edge.GetInstance()
			ex["edge"] = edge
			t.Extend(edge.GetExtend())
			t.Extend(ex)
			output = append(output, t)
		}
	}
	return output, nil
}

// LoadAllLinks read all views and all data
func (p *API) LoadAllLinks(name string, factory func() models.IPersistent, targetType IAPI) (interface{}, int, error) {
	// Read all rows
	all, _ := p.GetAll()
	for _, element := range all {
		// Retrieve all links
		targets := make([]models.IPersistent, 0)
		edges, _ := p.GenericLinkGetAll(element.GetID(), make([]models.IEdgeBean, 0), targetType)
		for _, edge := range edges {
			targets = append(targets, edge)
		}
		element.(models.IValueSetter).Set(name, targets)
	}
	return all, len(all), nil
}
