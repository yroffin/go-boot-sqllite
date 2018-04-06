// Package interfaces for common interfaces
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
	"log"
)

// Bean interface
type Bean struct {
	Name string
}

// BeanInterface interface
type IBean interface {
	Init() error
	GetName() string
	SetName(name string)
	PostConstruct(string) error
	Validate(string) error
}

// Inject Init this bean
func (bean *Bean) Inject(string, map[string]IBean) error {
	log.Printf("Bean::Inject")
	return nil
}

// PostConstruct Init this bean
func (bean *Bean) PostConstruct(string) error {
	log.Printf("Bean::PostConstruct")
	return nil
}

// Validate Init this bean
func (bean *Bean) Validate(string) error {
	log.Printf("Bean::Validate")
	return nil
}

// SetName fix the bean name
func (bean *Bean) SetName(name string) {
	bean.Name = name
}

// GetName get the bean name
func (bean *Bean) GetName() string {
	return bean.Name
}
