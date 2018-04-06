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
package winter

import (
	"log"
	"reflect"
	"strings"
)

var (
	// Helper manage beans
	Helper = (&Manager{}).New("manager")
)

// Manager interface
type Manager struct {
	*Service
	// Bean registry
	ArrayOfBeans     []interface{}
	ArrayOfBeanNames []string
	// Bean registry
	MapOfBeans map[string]interface{}
}

// IManager interface
type IManager interface {
	IService
	// Method
	Register(name string, b IBean) error
	Boot() error
	GetBean(name string) interface{}
	GetBeanNames() []string
	ForEach(func(interface{}))
}

// New constructor
func (m *Manager) New(name string) IManager {
	bean := Manager{Service: &Service{Bean: &Bean{}}}
	bean.ArrayOfBeans = make([]interface{}, 0)
	bean.ArrayOfBeanNames = make([]string, 0)
	bean.MapOfBeans = make(map[string]interface{})
	bean.Register(name, &bean)
	return &bean
}

// Init a single bean
func (m *Manager) Init() error {
	log.Printf("Manager::Init")
	return nil
}

// ForEach interate on beans
func (m *Manager) ForEach(iter func(interface{})) {
	for i := 0; i < len(m.GetBeanNames()); i++ {
		bean := m.GetBean(m.GetBeanNames()[i])
		iter(bean)
	}
}

// Register a single bean
func (m *Manager) Register(name string, b IBean) error {
	log.Println("Manager::Register", name)
	m.ArrayOfBeans = append(m.ArrayOfBeans, b)
	m.ArrayOfBeanNames = append(m.ArrayOfBeanNames, name)
	m.MapOfBeans[name] = b
	b.SetName(name)
	b.Init()
	return nil
}

// Boot Init this manager
func (m *Manager) Boot() error {
	log.Printf("Manager::Boot inject")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.Inject(m.ArrayOfBeanNames[index], m.ArrayOfBeans[index])
		log.Printf("Manager::Boot injection sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot post-construct")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		log.Printf("Manager::Boot post-construct execute for %v", m.ArrayOfBeanNames[index])
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "PostConstruct")
		log.Printf("Manager::Boot post-construct sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot validate")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "Validate")
		log.Printf("Manager::Boot validation sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	return nil
}

// GetBean get bean
func (m *Manager) GetBean(name string) interface{} {
	return m.MapOfBeans[name]
}

// GetBeanNames get bean
func (m *Manager) GetBeanNames() []string {
	return m.ArrayOfBeanNames
}

// Inject this API
func (m *Manager) Inject(name string, intf interface{}) error {
	m.autowire(false, 0, name, intf, reflect.ValueOf(intf))
	return nil
}

// dumpFields dump all fields
func (m *Manager) isPrivate(val reflect.StructField) bool {
	return strings.ToLower(val.Name[0:1]) == val.Name[0:1]
}

// dumpFields dump all fields
func (m *Manager) autowire(debug bool, level int, name string, intf interface{}, val reflect.Value) {
	if debug {
		log.Printf("%02d: **** METHOD ***", level)
		for i := 0; i < val.NumMethod(); i++ {
			typeMethod := val.Type().Method(i)
			log.Printf("Method Name: '%s'", typeMethod.Name)
		}
	}

	var kind = val.Type().Kind()

	// Interface case and Pointer case
	if kind == reflect.Interface || kind == reflect.Ptr {
		if !val.IsNil() {
			m.autowire(debug, level+1, name, intf, val.Elem())
		}
		return
	}

	// Function case
	if kind == reflect.Func {
		return
	}
	// Ignore primitive types
	if kind == reflect.Slice {
		return
	}
	// Ignore primitive types
	if kind == reflect.String {
		return
	}
	// Ignore primitive types
	if kind == reflect.Map {
		return
	}

	if debug {
		log.Printf("%02d: **** FIELDS ****", level)
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			log.Printf("Field  Name: '%s'\tField Value: '%v'\t Tag Value: '%v'", typeField.Name, valueField.Interface(), tag)
		}
	}

	// Dump all methods
	if debug {
		log.Printf("**** INJECT/SCAN **** Type: '%v' Kind: '%v'", val.Type(), val.Type().Kind())
	}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if !m.isPrivate(typeField) {
			if len(tag.Get("@autowired")) > 0 {
				if valueField.IsNil() {
					var beanName = tag.Get("@autowired")
					// myBean contain the target bean to inject
					var myBean = m.MapOfBeans[beanName]
					var camelCaseOperation = "Set" + typeField.Name
					var setter = reflect.ValueOf(intf).MethodByName(camelCaseOperation)
					arr := [1]reflect.Value{reflect.ValueOf(myBean)}
					var arguments = arr[:1]
					log.Printf("Apply: '%s' on %v with %v(%v) - %v", camelCaseOperation, beanName, setter, arguments, name)
					//setter.Call(arguments)
					valueField.Set(reflect.ValueOf(myBean))
				}
			}
		}
	}
	// Dump all fields for iterate recursively
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if !m.isPrivate(typeField) {
			// autowired fields are excluded for recursive init
			if !(len(tag.Get("@autowired")) > 0) {
				m.autowire(debug, level+1, name, valueField.Interface(), valueField)
			}
		}
	}
}

// dumpFields dump all fields
func (m *Manager) execute(debug bool, beanName string, intf interface{}, handler string) {
	val := reflect.ValueOf(intf)
	for i := 0; i < val.NumMethod(); i++ {
		typeMethod := val.Type().Method(i)
		if typeMethod.Name == handler {
			var setter = reflect.ValueOf(intf).MethodByName(handler)
			arr := [1]reflect.Value{reflect.ValueOf(beanName)}
			var arguments = arr[:1]
			if debug {
				log.Printf("Apply: '%s' on %v with %v(%v)", handler, beanName, setter, beanName)
			}
			setter.Call(arguments)
		}
	}
}
