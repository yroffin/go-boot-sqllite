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
	"reflect"
	"strings"

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
	core_services "github.com/yroffin/go-boot-sqllite/core/services"
)

// SwaggerService internal members
type SwaggerService struct {
	// members
	*core_services.SERVICE
	// swagger model
	Swagger *models.SwaggerModel
}

// ISwaggerService Test all package methods
type ISwaggerService interface {
	// Bean
	core_bean.IBean
	// Swagger
	SwaggerModel() *models.SwaggerModel
	Version(string) string
	BasePath(string) string
	// swagger method
	AddPaths(tags string, route string, method string, sumary string, description string, args map[string]interface{}, params map[string]interface{}, in []interface{}, out map[string]interface{})
}

// New constructor
func (p *SwaggerService) New() ISwaggerService {
	bean := SwaggerService{SERVICE: &core_services.SERVICE{Bean: &core_bean.Bean{}}}
	return &bean
}

// Init Init this API
func (p *SwaggerService) Init() error {
	// root model
	p.Swagger = &models.SwaggerModel{
		Swagger:  "2.0",
		Host:     "localhost:3000",
		BasePath: "/",
		Info: models.SwaggerInfo{
			Description:    "Todo.",
			Version:        "1.x",
			Title:          "Swagger",
			TermsOfService: "Todo.",
			Contact: models.SwaggerContact{
				Email: "yroffin@gmail.com",
			},
			License: models.SwaggerLicense{
				Name: "MIT",
				URL:  "https://github.com/yroffin/go-jarvis/blob/master/LICENSE",
			},
		},
		Tags:                make([]models.SwaggerTags, 0),
		Schemes:             make([]string, 0),
		Paths:               make(map[string]models.SwaggerRoute),
		SecurityDefinitions: make(map[string]models.SwaggerSecurityDefinitions),
		Definitions:         make(map[string]models.SwaggerDefinitions),
		ExternalDocs: models.SwaggerExternalDocs{
			Description: "Todo.",
			URL:         "Todo.",
		},
	}
	return nil
}

// PostConstruct Init this API
func (p *SwaggerService) PostConstruct(name string) error {
	return nil
}

// Validate Init this API
func (p *SwaggerService) Validate(name string) error {
	return nil
}

// SwaggerModel method
func (p *SwaggerService) SwaggerModel() *models.SwaggerModel {
	return p.Swagger
}

// Version method
func (p *SwaggerService) Version(vers string) string {
	return vers
}

// BasePath method
func (p *SwaggerService) BasePath(base string) string {
	return base
}

// AddTags method
func (p *SwaggerService) AddTags() {
	// tags content
	p.Swagger.Tags = make([]models.SwaggerTags, 0)
	p.Swagger.Tags = append(p.Swagger.Tags, models.SwaggerTags{
		Name:        "xxxx",
		Description: "xxxx",
		ExternalDocs: models.SwaggerExternalDocs{
			Description: "xxxx",
			URL:         "xxxx",
		},
	})
}

// AddSchemes method
func (p *SwaggerService) AddSchemes() {
	// tags content
	p.Swagger.Schemes = append(p.Swagger.Schemes, "http")
}

// AddPaths method
func (p *SwaggerService) AddPaths(tags string, route string, method string, summary string, description string, args map[string]interface{}, query map[string]interface{}, in []interface{}, out map[string]interface{}) {
	// Parameter path
	for k := range args {
		route = strings.Replace(route, ":"+k, "{"+k+"}", -1)
	}
	// Method body
	detail := &models.SwaggerMethodBody{
		Tags:        make([]string, 0),
		Summary:     summary,
		Description: description,
		Consumes:    make([]string, 0),
		Produces:    make([]string, 0),
		Parameters:  make([]models.SwaggerMethodParamBody, 0),
		Responses:   make(map[string]models.SwaggerMethodResp),
	}
	detail.Tags = append(detail.Tags, tags)
	detail.Consumes = append(detail.Consumes, "application/json")
	detail.Produces = append(detail.Produces, "application/json")
	// route content
	var met = strings.ToLower(method)
	if p.Swagger.Paths[route] == nil {
		p.Swagger.Paths[route] = make(map[string]models.SwaggerMethodBody)
	}
	// Parameters in
	for index := 0; index < len(in); index++ {
		prm := models.SwaggerMethodParamBody{
			In:          "body",
			Name:        "Name",
			Description: "Desc",
			Required:    false,
			Type:        "string",
		}
		detail.Parameters = append(detail.Parameters, prm)
	}
	// Parameter path
	for k, v := range args {
		prm := models.SwaggerMethodParamBody{
			In:          "path",
			Name:        k,
			Description: reflect.TypeOf(v).String(),
			Required:    false,
			Type:        "string",
		}
		detail.Parameters = append(detail.Parameters, prm)
	}
	// Parameter query
	for k, v := range query {
		prm := models.SwaggerMethodParamBody{
			In:          "query",
			Name:        k,
			Description: reflect.TypeOf(v).String(),
			Required:    false,
			Type:        "string",
		}
		detail.Parameters = append(detail.Parameters, prm)
	}
	// Response
	for index, outv := range out {
		body := models.SwaggerMethodResp{Schema: make(map[string]string)}
		var key = reflect.TypeOf(outv).String()
		var unary = strings.SplitAfter(key, ".")[1]

		if strings.Contains(key, "[]") {
			body.Schema["type"] = "array"
			body.Description = "Array of " + unary
		} else {
			body.Description = "Instance of " + unary
			// Add definitions
			p.AddDefinition(unary, outv)
		}

		body.Schema["$ref"] = "#/definitions/" + unary
		detail.Responses[index] = body
	}
	p.Swagger.Paths[route][met] = *detail
}

// AddDefinition method
func (p *SwaggerService) getType(typ string) (string, string) {
	switch typ {
	case "string":
		return "string", "string"
	case "JSONTime":
		return "string", "date-time"
	default:
		return "string", "string"
	}
}

// AddDefinition method
func (p *SwaggerService) AddDefinition(name string, ptr interface{}) string {
	if _, ok := p.Swagger.Definitions[name]; ok {
		return name
	}

	p.Swagger.Definitions[name] = models.SwaggerDefinitions{
		Type:       "Object",
		Properties: make(map[string]models.SwaggerFormat),
	}

	var fields = reflect.TypeOf(reflect.ValueOf(ptr))
	for index := 0; index < fields.NumField(); index++ {
		var fieldName = fields.Field(index).Name
		var a, b = p.getType(fields.Field(index).Type.Name())
		p.Swagger.Definitions[name].Properties[fieldName] = models.SwaggerFormat{
			Type:   a,
			Format: b,
		}
	}
	return name
}
