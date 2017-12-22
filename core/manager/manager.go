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
package manager

import (
	"log"
	"reflect"
	"strings"

	"github.com/yroffin/go-boot-sqllite/core/bean"
)

// Manager interface
type Manager struct {
	Beans map[string]bean.IBean
}

// IManager interface
type IManager interface {
	Init() error
	Register() error
	Boot() error
	// Scan and inject bean in this class
	Inject(interface{}, string, func(interface{})) error
}

// Init a single bean
func (m *Manager) Init() {
	log.Printf("Manager::Init")
	m.Beans = make(map[string]bean.IBean)
}

// Register a single bean
func (m *Manager) Register(name string, b bean.IBean) error {
	m.Beans[name] = b
	b.SetName(name)
	b.Init()
	return nil
}

// Boot Init this manager
func (m *Manager) Boot() error {
	log.Printf("Manager::Boot inject")
	for key, value := range m.Beans {
		m.Inject(key, value)
		log.Printf("Manager::Boot injection sucessfull for %v", key)
	}
	log.Printf("Manager::Boot post-construct")
	for key, value := range m.Beans {
		value.PostConstruct(key)
		log.Printf("Manager::Boot post-construct sucessfull for %v", key)
	}
	log.Printf("Manager::Boot validate")
	for key, value := range m.Beans {
		value.Validate(key)
		log.Printf("Manager::Boot validation sucessfull for %v", key)
	}
	return nil
}

// BeanScanner scan all component on this bean
func (m *Manager) BeanScanner(myBean interface{}) {
	m.reflect(&myBean)
}

// Inject this API
func (m *Manager) Inject(name string, value interface{}) error {
	m.BeanScanner(value)
	return nil
}

// Inject this API
func (m *Manager) reflect(element interface{}) {
	val := reflect.ValueOf(element).Elem()
	m.dump(val)
}

// dumpFields dump all fields
func (m *Manager) isPrivate(val reflect.StructField) bool {
	return strings.ToLower(val.Name[0:1]) == val.Name[0:1]
}

// dumpFields dump all fields
func (m *Manager) dump(val reflect.Value) {
	// Interface case
	if val.Type().Kind() == reflect.Interface {
		m.dump(val.Elem())
		return
	}
	// Pointer case
	if val.Type().Kind() == reflect.Ptr {
		if !val.IsNil() {
			m.dump(val.Elem())
		}
		return
	}
	// Function case
	if val.Type().Kind() == reflect.Func {
		return
	}
	// Ignore primitive types
	if val.Type().Kind() == reflect.Slice {
		return
	}
	// Ignore primitive types
	if val.Type().Kind() == reflect.String {
		return
	}
	log.Printf("**** INJECT/SCAN **** Type: '%v' Kind: '%v'", val.Type(), val.Type().Kind())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if m.isPrivate(typeField) {
			log.Printf("Field  Name: '%s'\tprivate", typeField.Name)
		} else {
			log.Printf("Field  Name: '%s'\tField Value: '%v'\t Tag Value: '%v'", typeField.Name, valueField.Interface(), tag)
			if len(tag.Get("bean")) > 0 {
				var beanName = tag.Get("bean")
				apply, ok := valueField.Interface().(func(interface{}))
				if ok {
					log.Printf("Field  Name: '%s' INJECTION with %v/%v", typeField.Name, m.Beans[beanName], beanName)
					apply(m.Beans[beanName])
				} else {
					log.Printf("Field  Name: '%s' IS NOT COMPATIBLE", typeField.Name)
				}
			}
		}
	}
	// Dump all methods
	for i := 0; i < val.NumMethod(); i++ {
		typeMethod := val.Type().Method(i)
		log.Printf("Method Name: '%s'", typeMethod.Name)
	}
	// Dump all fields for iterate recursively
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		m.dump(valueField)
	}
}
