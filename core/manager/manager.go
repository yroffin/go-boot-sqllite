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
	"flag"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/yroffin/go-boot-sqllite/core/apis"
	"github.com/yroffin/go-boot-sqllite/core/bean"
)

// Manager interface
type Manager struct {
	// Bean registry
	ArrayOfBeans     []interface{}
	ArrayOfBeanNames []string
	// Bean registry
	MapOfBeans map[string]interface{}
	// sync wait group
	wg sync.WaitGroup
	// Properties
	phttp *int
	// Properties
	phttps *int
}

// IManager interface
type IManager interface {
	Init() error
	Register() error
	CommandLine() error
	Boot() error
	// Scan and inject bean in this class
	Inject(interface{}, string, func(interface{})) error
}

// Init a single bean
func (m *Manager) Init() {
	log.Printf("Manager::Init")
	m.ArrayOfBeans = make([]interface{}, 0)
	m.ArrayOfBeanNames = make([]string, 0)
	m.MapOfBeans = make(map[string]interface{})
}

// Register a single bean
func (m *Manager) Register(name string, b bean.IBean) error {
	m.ArrayOfBeans = append(m.ArrayOfBeans, b)
	m.ArrayOfBeanNames = append(m.ArrayOfBeanNames, name)
	m.MapOfBeans[name] = b
	b.SetName(name)
	b.Init()
	return nil
}

// CommandLine Init
func (m *Manager) CommandLine() error {
	// scan flags
	m.phttp = flag.Int("http", -1, "Http port")
	m.phttps = flag.Int("https", -1, "Https port")
	flag.Parse()
	return nil
}

// Boot Init this manager
func (m *Manager) Boot(routerBeanName string) error {
	log.Printf("Manager::Boot inject")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.Inject(m.ArrayOfBeanNames[index], m.ArrayOfBeans[index])
		log.Printf("Manager::Boot injection sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot post-construct")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "PostConstruct")
		log.Printf("Manager::Boot post-construct sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot validate")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "Validate")
		log.Printf("Manager::Boot validation sucessfull for %v", m.ArrayOfBeanNames[index])
	}

	var router apis.IRouter
	if assertion, ok := m.MapOfBeans[routerBeanName].(apis.IRouter); ok {
		router = assertion
	} else {
		log.Fatalf("Unable to validate bean %s", routerBeanName)
	}

	if *m.phttp != -1 {
		// Declarre listener HTTP
		log.Printf("Manager::Boot listen on %v", *m.phttp)
		router.HTTP(*m.phttp)
	}
	if *m.phttps != -1 {
		// Declarre listener HTTPS
		log.Printf("Manager::Boot listen on %v", *m.phttps)
		router.HTTPS(*m.phttps)
	}
	m.Wait()
	return nil
}

// HTTP declare http listener
func (m *Manager) HTTP(port int) error {
	var sport = fmt.Sprintf("%d", port)
	m.wg.Add(1)
	go func(sport string) {
		log.Printf("Try to serve HTTP proxy on %s", sport)
		// After defining our server, we finally "listen and serve" on port 8080
		err := http.ListenAndServe(":"+sport, nil)
		if err != nil {
			log.Fatalf("Unable to serve HTTP %v", err)
		}
	}(sport)
	return nil
}

// HTTPS declare http listener
func (m *Manager) HTTPS(port int) error {
	var sport = fmt.Sprintf("%d", port)
	m.wg.Add(1)
	go func(sport string) {
		log.Printf("Try to serve HTTPS proxy on %s", sport)
		// Also serve on https/tls
		err := http.ListenAndServeTLS(":"+sport, ".ssl/hostname.pem", ".ssl/private.key", nil)
		if err != nil {
			log.Fatalf("Unable to serve HTTPS %v", err)
		}
	}(sport)
	return nil
}

// Wait for end of all listener
func (m *Manager) Wait() error {
	m.wg.Wait()
	return nil
}

// Inject this API
func (m *Manager) Inject(name string, intf interface{}) error {
	m.autowire(false, 0, intf, reflect.ValueOf(intf))
	return nil
}

// dumpFields dump all fields
func (m *Manager) isPrivate(val reflect.StructField) bool {
	return strings.ToLower(val.Name[0:1]) == val.Name[0:1]
}

// dumpFields dump all fields
func (m *Manager) autowire(debug bool, level int, intf interface{}, val reflect.Value) {
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
			m.autowire(debug, level+1, intf, val.Elem())
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
					var myBean = m.MapOfBeans[beanName]
					var camelCase = ""
					for _, word := range strings.Split(beanName, "-") {
						camelCase += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
					}
					var camelCaseOperation = "Set" + camelCase
					var setter = reflect.ValueOf(intf).MethodByName(camelCaseOperation)
					arr := [1]reflect.Value{reflect.ValueOf(myBean)}
					var arguments = arr[:1]
					log.Printf("Apply: '%s' on %v with %v(%v)", camelCaseOperation, beanName, setter, arguments)
					setter.Call(arguments)
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
				m.autowire(debug, level+1, valueField.Interface(), valueField)
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
