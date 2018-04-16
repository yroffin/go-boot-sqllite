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
	"os"
	"reflect"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	// Helper manage beans
	Helper = (&Manager{}).New("manager")
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

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
	m.ArrayOfBeans = append(m.ArrayOfBeans, b)
	m.ArrayOfBeanNames = append(m.ArrayOfBeanNames, name)
	m.MapOfBeans[name] = b
	b.SetName(name)
	b.Init()
	return nil
}

// Boot Init this manager
func (m *Manager) Boot() error {
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		m.Inject(m.ArrayOfBeanNames[index], m.ArrayOfBeans[index])
		log.WithFields(log.Fields{
			"name": m.ArrayOfBeanNames[index],
		}).Info("Boot injection sucessfull")
	}
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		log.WithFields(log.Fields{
			"index": index,
			"count": len(m.ArrayOfBeans),
			"name":  m.ArrayOfBeanNames[index],
		}).Info("Boot post-construct execute")
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "PostConstruct")
		log.WithFields(log.Fields{
			"index": index,
			"count": len(m.ArrayOfBeans),
			"name":  m.ArrayOfBeanNames[index],
		}).Info("Boot post-construct execute sucessfull")
	}
	for index := 0; index < len(m.ArrayOfBeans); index++ {
		log.WithFields(log.Fields{
			"index": index,
			"count": len(m.ArrayOfBeans),
			"name":  m.ArrayOfBeanNames[index],
		}).Info("Boot validate execute")
		m.execute(false, m.ArrayOfBeanNames[index], m.ArrayOfBeans[index], "Validate")
		log.WithFields(log.Fields{
			"index": index,
			"count": len(m.ArrayOfBeans),
			"name":  m.ArrayOfBeanNames[index],
		}).Info("Boot validate execute sucessfull")
	}
	// Wait infinite
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
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
		log.WithFields(log.Fields{
			"level": level,
		}).Debug("Method")
		for i := 0; i < val.NumMethod(); i++ {
			typeMethod := val.Type().Method(i)
			log.WithFields(log.Fields{
				"name": typeMethod.Name,
			}).Debug("Method name")
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
		log.WithFields(log.Fields{
			"level": level,
		}).Debug("Field")
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			log.WithFields(log.Fields{
				"name":      typeField.Name,
				"interface": valueField.Interface(),
				"tag":       tag,
			}).Debug("Field")
		}
	}

	// Dump all methods
	if debug {
		log.WithFields(log.Fields{
			"type": val.Type(),
			"kind": val.Type().Kind(),
		}).Debug("Injection/scan")
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
					log.WithFields(log.Fields{
						"bean": myBean,
						"name": typeField.Name,
					}).Debug("Set")
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
				log.WithFields(log.Fields{
					"handler": handler,
					"name":    beanName,
				}).Debug("Execute")
			}
			setter.Call(arguments)
		}
	}
}
