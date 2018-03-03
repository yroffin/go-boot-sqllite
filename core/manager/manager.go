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

	"github.com/yroffin/go-boot-sqllite/core/bean"
)

// Manager interface
type Manager struct {
	// Bean registry
	ArrayOfBeans     []bean.IBean
	ArrayOfBeanNames []string
	// Bean registry
	MapOfBeans map[string]bean.IBean
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
	m.ArrayOfBeans = make([]bean.IBean, 0)
	m.ArrayOfBeanNames = make([]string, 0)
	m.MapOfBeans = make(map[string]bean.IBean)
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
func (m *Manager) Boot() error {
	log.Printf("Manager::Boot inject")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.Inject(m.ArrayOfBeanNames[index], m.ArrayOfBeans[index])
		log.Printf("Manager::Boot injection sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot post-construct")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.ArrayOfBeans[index].PostConstruct(m.ArrayOfBeanNames[index])
		log.Printf("Manager::Boot post-construct sucessfull for %v", m.ArrayOfBeanNames[index])
	}
	log.Printf("Manager::Boot validate")
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.ArrayOfBeans[index].Validate(m.ArrayOfBeanNames[index])
		log.Printf("Manager::Boot validation sucessfull for %v", m.ArrayOfBeanNames[index])
	}

	if *m.phttp != -1 {
		// Declarre listener HTTP
		log.Printf("Manager::Boot listen on %v", *m.phttp)
		m.HTTP(*m.phttp)
	}
	if *m.phttps != -1 {
		// Declarre listener HTTPS
		log.Printf("Manager::Boot listen on %v", *m.phttps)
		m.HTTPS(*m.phttps)
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
		if val.IsNil() {
			return
		}
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
					log.Printf("Field  Name: '%s' INJECTION with %v/%v", typeField.Name, m.MapOfBeans[beanName], beanName)
					apply(m.MapOfBeans[beanName])
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
