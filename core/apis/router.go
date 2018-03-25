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
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	core_services "github.com/yroffin/go-boot-sqllite/core/services"
)

// Router internal members
type Router struct {
	// members
	*core_services.SERVICE
	// gin router
	Engine *gin.Engine
	// SwaggerService with injection mecanism
	SwaggerService *SwaggerService `@autowired:"swagger"`
}

// IRouter Test all package methods
type IRouter interface {
	// Bean
	core_bean.IBean
	// Http boot
	HTTP(port int) error
	// Https boot
	HTTPS(port int) error
	// Swagger
	SwaggerModel() func(*gin.Context)
	// HandleFunc
	HandleFunc(path string, f func(c *gin.Context), method string, content string)
	// HandleFuncString declare a string handler
	HandleFuncString(path string, f func() (string, error), method string, content string)
	// HandleFuncStringWithId declare a string handler
	HandleFuncStringWithId(path string, f func(string) (string, error), method string, content string)
}

// New constructor
func (p *Router) New() IRouter {
	bean := Router{SERVICE: &core_services.SERVICE{Bean: &core_bean.Bean{}}}
	return &bean
}

// SetSwagger inject notification
func (p *Router) SetSwagger(value interface{}) {
	if assertion, ok := value.(*SwaggerService); ok {
		p.SwaggerService = assertion
	} else {
		log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
	}
}

// Init Init this API
func (p *Router) Init() error {
	return nil
}

// PostConstruct Init this API
func (p *Router) PostConstruct(name string) error {
	log.Printf("Router::PostConstruct - router creation")
	// define all routes
	p.Engine = gin.Default()
	// Fix default handler
	//p.Engine.HandleMethodNotAllowed = http.HandlerFunc(p.HandlerStaticNotAllowed())
	//p.Engine = http.HandlerFunc(p.HandlerStaticNotFound())

	return nil
}

// Validate Init this API
func (p *Router) Validate(name string) error {
	log.Printf("Router::Validate - router validation")
	p.Engine.Static("/public", "./resources/static")
	return nil
}

// Swagger method
func (p *Router) SwaggerModel() func(*gin.Context) {
	anonymous := func(c *gin.Context) {
		c.IndentedJSON(200, p.SwaggerService.SwaggerModel())
	}
	return anonymous
}

// HTTP boot http service
func (p *Router) HTTP(port int) error {
	gin.SetMode("debug")

	p.Engine.GET("/api/swagger.json", p.SwaggerModel())
	p.Engine.Run(":" + strconv.Itoa(port))
	return nil
}

// HTTP boot http service
func (p *Router) HTTPS(port int) error {
	return nil
}

// HandleFunc declare a handler
func (p *Router) HandleFuncTonic(path string, f func() (interface{}, error), method string, content string) {
	log.Printf("Router::HandleFuncTonic '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticJson(f, 200, ""))
}

// HandleFunc declare a handler
func (p *Router) HandleFunc(path string, f func(c *gin.Context), method string, content string) {
	log.Printf("Router::HandleFunc '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStatic(f, content))
}

// HandleFuncString declare a string handler
func (p *Router) HandleFuncString(path string, f func() (string, error), method string, content string) {
	log.Printf("Router::HandleFuncString '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticString(f, content))
}

// HandleFuncStringWithId declare a string handler
func (p *Router) HandleFuncStringWithId(path string, f func(string) (string, error), method string, content string) {
	log.Printf("Router::HandleFuncStringWithId '%s' with method '%s' with type mime '%s'", path, method, content)
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticStringWithId(f, content))
}

// HandlerStaticNotFound Not found handler
func (p *Router) HandlerStaticNotFound() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		log.Printf("Request request %v", c.Request)
		log.Printf("Request header %v", c.Request.Header)
		log.Printf("Request Encoding %v", c.Request.TransferEncoding)
		// content
		c.Header("Content-type", "text/html")
		c.String(404, "{\"message\":\"Not found\"}")
	}
	return anonymous
}

// HandlerStaticNotAllowed Not found handler
func (p *Router) HandlerStaticNotAllowed() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		log.Printf("Request request %v", c.Request)
		log.Printf("Request header %v", c.Request.Header)
		log.Printf("Request Encoding %v", c.Request.TransferEncoding)
		// content
		c.Header("Content-type", "text/html")
		c.String(405, "{\"message\":\"Not allowed\"}")
	}
	return anonymous
}

// HandlerStaticString render string
func (p *Router) HandlerStaticJson(method func() (interface{}, error), code int, content string) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		// content
		c.Header("Content-type", "text/html")
		data, err := method()
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.JSON(code, data)
		if len(content) > 0 {
			c.Header("Content-Type", content)
		}
	}
	return anonymous
}

// HandlerStatic render string
func (p *Router) HandlerStatic(method func(c *gin.Context), content string) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		method(c)
		if len(content) > 0 {
			c.Header("Content-Type", content)
		}
	}
	return anonymous
}

// HandlerStaticString render string
func (p *Router) HandlerStaticString(method func() (string, error), content string) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		// content
		c.Header("Content-type", "text/html")
		data, err := method()
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.String(200, data)
		if len(content) > 0 {
			c.Header("Content-Type", content)
		}
	}
	return anonymous
}

// HandlerStaticStringWithId render string
func (p *Router) HandlerStaticStringWithId(method func(string) (string, error), content string) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		// content
		c.Header("Content-type", "text/html")
		data, err := method(c.Param("id"))
		if err != nil {
			c.String(400, "{\"message\":\"\"}")
			return
		}
		c.String(200, data)
		if len(content) > 0 {
			c.Header("Content-Type", content)
		}
	}
	return anonymous
}
