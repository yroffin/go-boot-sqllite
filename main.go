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
package main

import (
	"github.com/yroffin/go-boot-sqllite/core/apis"
	"github.com/yroffin/go-boot-sqllite/core/business"
	"github.com/yroffin/go-boot-sqllite/core/manager"
	"github.com/yroffin/go-boot-sqllite/core/stores"
)

// Main
func main() {
	// declare manager and boot it
	m := (&manager.Manager{}).New("manager")
	// Command Line
	m.CommandLine()
	// Core beans
	m.Register("swagger", (&apis.SwaggerService{}).New())
	m.Register("router", (&apis.Router{}).New())
	m.Register("sql-crud-business", (&business.SqlCrudBusiness{}).New())
	m.Register("graph-crud-business", (&business.GraphCrudBusiness{}).New())
	m.Register("sqllite-manager", (&stores.Store{}).New([]string{"Node"}, "./sqllite.db"))
	m.Register("cayley-manager", (&stores.Graph{}).New([]string{"Node"}, "./cayley.db"))
	// API beans
	m.Register("node-api", (&apis.Node{}).New())
	m.Boot()
}
