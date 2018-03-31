// Package apis for common apis
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

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
)

// Node internal members
type Node struct {
	// Base component
	*API
	// internal members
	Name string
	// mounts
	Crud interface{} `@crud:"/api/nodes"`
	Link INode       `@autowired:"node-api" @link:"/api/nodes" @href:"nodes"`
	Node INode       `@autowired:"node-api"`
	// SwaggerService with injection mecanism
	Swagger ISwaggerService `@autowired:"swagger"`
}

// INode implements IBean
type INode interface {
	APIInterface
}

// New constructor
func (p *Node) New() INode {
	bean := &Node{API: &API{Bean: &core_bean.Bean{}}}
	return bean
}

// SetNode injection
func (p *Node) SetNode(value interface{}) {
	if assertion, ok := value.(INode); ok {
		p.Node = assertion
	} else {
		log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
	}
}

// SetLink injection
func (p *Node) SetLink(value interface{}) {
	if assertion, ok := value.(INode); ok {
		p.Link = assertion
	} else {
		log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
	}
}

// SetSwagger injection
func (p *Node) SetSwagger(value interface{}) {
	if assertion, ok := value.(ISwaggerService); ok {
		p.Swagger = assertion
	} else {
		log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
	}
}

// Init this API
func (p *Node) Init() error {
	// Crud
	p.Factory = func() models.IPersistent {
		return (&models.NodeBean{}).New()
	}
	p.Factories = func() models.IPersistents {
		return (&models.NodeBeans{}).New()
	}
	return p.API.Init()
}

// PostConstruct this API
func (p *Node) PostConstruct(name string) error {
	// Scan struct and init all handler
	p.ScanHandler(p.Swagger, p)
	return nil
}

// Validate this API
func (p *Node) Validate(name string) error {
	return nil
}
