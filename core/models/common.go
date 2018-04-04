// Package models for all models
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
package models

import (
	"encoding/json"
	"log"
)

// ValueBean simple value model
type ValueBean struct {
	// Extended internal store
	Extended map[string]interface{} `json:"extended"`
}

// IValueSetter interface
type IValueSetter interface {
	Set(string, interface{})
}

// IValueBean simple interface
type IValueBean interface {
	IValueSetter
	SetString(string, string)
	GetAsString(string) string
	GetAsStringArray(string) []string
	ToString() string
	ToJSON() string
}

// New constructor
func (p *ValueBean) New() IValueBean {
	bean := ValueBean{}
	bean.Extended = make(map[string]interface{})
	return &bean
}

// ToString stringify this commnd
func (p *ValueBean) ToString() string {
	payload, err := json.Marshal(p.Extended)
	if err != nil {
		log.Println("Unable to marshal:", err)
		return "{}"
	}
	return string(payload)
}

// ToJSON return o json formated value (in pretty format)
func (p *ValueBean) ToJSON() string {
	payload, err := json.MarshalIndent(p.Extended, "", "\t")
	if err != nil {
		log.Println("Unable to marshal:", err)
		return "{}"
	}
	return string(payload)
}

// Set a value for a key
func (p *ValueBean) Set(key string, value interface{}) {
	p.Extended[key] = value
}

// SetString a value for a key
func (p *ValueBean) SetString(key string, value string) {
	p.Extended[key] = value
}

// GetAsString field value
func (p *ValueBean) GetAsString(key string) string {
	if assertion, ok := p.Extended[key].(string); ok {
		return assertion
	}
	log.Fatalf("Unable to render key %v for string type", key)
	return ""
}

// GetAsStringArray field value
func (p *ValueBean) GetAsStringArray(key string) []string {
	if assertion, ok := p.Extended[key].([]string); ok {
		return assertion
	}
	log.Fatalf("Unable to render key %v for string type", key)
	return make([]string, 0)
}
