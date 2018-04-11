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
	"flag"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/yroffin/go-boot-sqllite/core/winter"
)

func init() {
	winter.Helper.Register("APIManager", (&APIManager{}).New())
}

// APIManager interface
type APIManager struct {
	*winter.Service
	// sync wait group
	wg sync.WaitGroup
	// Properties
	phttp *int
	// Properties
	phttps *int
	// Inject
	Router IRouter `@autowired:"router"`
}

// IAPIManager interface
type IAPIManager interface {
	winter.IService
	// Method
	CommandLine() error
}

// New constructor
func (m *APIManager) New() IAPIManager {
	bean := APIManager{Service: &winter.Service{Bean: &winter.Bean{}}}
	return &bean
}

// Init a single bean
func (m *APIManager) Init() error {
	return nil
}

// CommandLine Init
func (m *APIManager) CommandLine() error {
	// scan flags
	m.phttp = flag.Int("http", -1, "Http port")
	m.phttps = flag.Int("https", -1, "Https port")
	flag.Parse()
	return nil
}

// Validate Init this manager
func (m *APIManager) Validate(name string) error {
	if *m.phttp != -1 {
		// Declarre listener HTTP
		log.WithFields(log.Fields{
			"port": *m.phttp,
		}).Info("Boot listener http")
		m.Router.HTTP(*m.phttp)
	}
	if *m.phttps != -1 {
		// Declarre listener HTTPS
		log.WithFields(log.Fields{
			"port": *m.phttps,
		}).Info("Boot listener https")
		m.Router.HTTPS(*m.phttps)
	}
	m.Wait()
	return nil
}

// HTTP declare http listener
func (m *APIManager) HTTP(port int) error {
	var sport = fmt.Sprintf("%d", port)
	m.wg.Add(1)
	go func(sport string) {
		log.WithFields(log.Fields{
			"port": sport,
		}).Info("Try to serve http proxy")
		// After defining our server, we finally "listen and serve" on port 8080
		err := http.ListenAndServe(":"+sport, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to serve http")
		}
	}(sport)
	return nil
}

// HTTPS declare http listener
func (m *APIManager) HTTPS(port int) error {
	var sport = fmt.Sprintf("%d", port)
	m.wg.Add(1)
	go func(sport string) {
		log.WithFields(log.Fields{
			"port": sport,
		}).Info("Try to serve https proxy")
		// Also serve on https/tls
		err := http.ListenAndServeTLS(":"+sport, ".ssl/hostname.pem", ".ssl/private.key", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to serve https")
		}
	}(sport)
	return nil
}

// Wait for end of all listener
func (m *APIManager) Wait() error {
	m.wg.Wait()
	return nil
}
