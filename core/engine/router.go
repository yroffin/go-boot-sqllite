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
package engine

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
	"github.com/yroffin/go-boot-sqllite/core/winter"
)

func init() {
	winter.Helper.Register("router", (&service{}).New())
}

// service internal members
type service struct {
	*winter.Service
	// gin router
	Engine *gin.Engine
	// Packr
	box packr.Box
	// SwaggerService with injection mecanism
	Swagger ISwaggerService `@autowired:"swagger"`
}

// IRouter Test all package methods
type IRouter interface {
	winter.IService
	// Http boot
	HTTP(port int) error
	// Https boot
	HTTPS(port int) error
	// Swagger
	SwaggerModel() func(*gin.Context)
	// HandleFunc
	HandleFunc(path string, f func(c *gin.Context), method string, content string)
	// HandleFunc
	HandleFuncLink(path string, f func(c *gin.Context, target IAPI), method string, content string, target IAPI)
	// HandleFuncString declare a string handler
	HandleFuncString(path string, f func() (string, error), method string, content string)
	// HandleFuncStringWithId declare a string handler
	HandleFuncStringWithId(path string, f func(string) (string, error), method string, content string)
}

// New constructor
func (p *service) New() IRouter {
	bean := service{Service: &winter.Service{Bean: &winter.Bean{}}}
	return &bean
}

// Init Init this API
func (p *service) Init() error {
	return nil
}

// PostConstruct Init this API
func (p *service) PostConstruct(name string) error {
	// define all routes
	p.Engine = gin.Default()

	return nil
}

// Resources Init this API
func (p *service) Resources(name string, box packr.Box) error {
	// Packr
	p.box = box
	for _, resource := range box.List() {
		var content = "text/html"
		if strings.HasSuffix(resource, ".html") {
			content = "text/html"
		}
		if strings.HasSuffix(resource, ".js") {
			content = "application/javascript"
		}
		if strings.HasSuffix(resource, ".json") {
			content = "application/json"
		}
		if strings.HasSuffix(resource, ".css") {
			content = "text/css"
		}
		log.WithFields(log.Fields{
			"resource": resource,
			"path":     "/public/" + resource,
			"content":  content,
		}).Info("Resources")
		p.Engine.GET("/public/"+resource, p.HandlerStaticFile(resource, content))
	}
	return nil
}

// Validate Init this API
func (p *service) Validate(name string) error {
	return nil
}

// Swagger method
func (p *service) SwaggerModel() func(*gin.Context) {
	anonymous := func(c *gin.Context) {
		c.IndentedJSON(200, p.Swagger.SwaggerModel())
	}
	return anonymous
}

// HTTP boot http service
func (p *service) HTTP(port int) error {
	gin.SetMode("debug")

	p.Engine.GET("/api/swagger.json", p.SwaggerModel())
	p.Engine.Run(":" + strconv.Itoa(port))
	return nil
}

// HTTP boot http service
func (p *service) HTTPS(port int) error {
	return nil
}

// HandleFunc declare a handler
func (p *service) HandleFuncTonic(path string, f func() (interface{}, error), method string, content string) {
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticJson(f, 200, ""))
}

// HandleFunc declare a handler
func (p *service) HandleFunc(path string, f func(c *gin.Context), method string, content string) {
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStatic(f, content))
}

// HandleFuncLink declare a handler
func (p *service) HandleFuncLink(path string, f func(c *gin.Context, target IAPI), method string, content string, target IAPI) {
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticLink(f, content, target))
}

// HandleFuncString declare a string handler
func (p *service) HandleFuncString(path string, f func() (string, error), method string, content string) {
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticString(f, content))
}

// HandleFuncStringWithId declare a string handler
func (p *service) HandleFuncStringWithId(path string, f func(string) (string, error), method string, content string) {
	// declare it to the router
	p.Engine.Handle(method, path, p.HandlerStaticStringWithId(f, content))
}

// HandlerStaticNotFound Not found handler
func (p *service) HandlerStaticNotFound() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		log.WithFields(log.Fields{
			"request":  c.Request,
			"header":   c.Request.Header,
			"encoding": c.Request.TransferEncoding,
		}).Warn("While retrrieve row(s)")
		// content
		c.Header("Content-type", "text/html")
		c.String(404, "{\"message\":\"Not found\"}")
	}
	return anonymous
}

// HandlerStaticNotAllowed Not found handler
func (p *service) HandlerStaticNotAllowed() func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		log.WithFields(log.Fields{
			"request":  c.Request,
			"header":   c.Request.Header,
			"encoding": c.Request.TransferEncoding,
		}).Warn("While retrrieve row(s)")
		// content
		c.Header("Content-type", "text/html")
		c.String(405, "{\"message\":\"Not allowed\"}")
	}
	return anonymous
}

// HandlerStaticString render string
func (p *service) HandlerStaticFile(resource string, content string) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		// content
		c.Header("Content-type", content)
		data, _ := p.box.MustString(resource)
		c.String(200, data)
	}
	return anonymous
}

// HandlerStaticString render string
func (p *service) HandlerStaticJson(method func() (interface{}, error), code int, content string) func(c *gin.Context) {
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
func (p *service) HandlerStatic(method func(c *gin.Context), content string) func(c *gin.Context) {
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

// HandlerStaticLink render static handler
func (p *service) HandlerStaticLink(method func(c *gin.Context, target IAPI), content string, target IAPI) func(c *gin.Context) {
	anonymous := func(c *gin.Context) {
		// security header
		c.Header("Strict-Transport-Security", "")
		c.Header("Content-Security-Policy", "")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "same-origin")
		method(c, target)
		if len(content) > 0 {
			c.Header("Content-Type", content)
		}
	}
	return anonymous
}

// HandlerStaticString render string
func (p *service) HandlerStaticString(method func() (string, error), content string) func(c *gin.Context) {
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
func (p *service) HandlerStaticStringWithId(method func(string) (string, error), content string) func(c *gin.Context) {
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
