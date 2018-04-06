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
package engine

import (
	"github.com/yroffin/go-boot-sqllite/core/models"
	"github.com/yroffin/go-boot-sqllite/core/winter"
)

// Node internal members
type Node struct {
	// Base component
	*API
	// internal members
	Name string
	// mounts
	Crud interface{} `@crud:"/api/nodes"`
	Link INode       `@autowired:"NodeBean" @link:"/api/nodes" @href:"nodes"`
	Node INode       `@autowired:"NodeBean"`
	// SwaggerService with injection mecanism
	Swagger ISwaggerService `@autowired:"swagger"`
}

// INode implements IBean
type INode interface {
	IAPI
}

// New constructor
func (p *Node) New() INode {
	bean := &Node{API: &API{Bean: &winter.Bean{}}}
	return bean
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
