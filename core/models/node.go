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

// NodeBean simple command model
type NodeBean struct {
	// Id
	ID string `json:"id"`
	// Timestamp
	Timestamp JSONTime `json:"timestamp"`
	// Name
	Name string `json:"name"`
	// Type
	Type string `json:"type"`
	// Extended internal store
	Extended map[string]interface{} `json:"extended"`
}

// INodeBean interface
type INodeBean interface {
	// inherit persistent behaviour
	IPersistent
	// inherit ValueBean behaviour
	IValueBean
	// command
	GetType() string
}

// New constructor
func (p *NodeBean) New() INodeBean {
	bean := NodeBean{}
	bean.Extended = make(map[string]interface{})
	return &bean
}

// GetName get set name
func (p *NodeBean) GetName() string {
	return "NodeBean"
}

// Extend vars
func (p *NodeBean) Extend(e map[string]interface{}) {
	for k, v := range e {
		p.Extended[k] = v
	}
}

// GetType get set name
func (p *NodeBean) GetType() string {
	return p.Type
}

// GetID retrieve ID
func (p *NodeBean) GetID() string {
	return p.ID
}

// SetID retrieve ID
func (p *NodeBean) SetID(ID string) {
	p.ID = ID
}

// Set get set name
func (p *NodeBean) Set(key string, value interface{}) {
}

// SetString get set name
func (p *NodeBean) SetString(key string, value string) {
	// Call super method
	IValueBean(p).SetString(key, value)
}

// Get get set name
func (p *NodeBean) GetAsString(key string) string {
	switch key {
	default:
		// Call super method
		return IValueBean(p).GetAsString(key)
	}
}

// Get get set name
func (p *NodeBean) GetAsStringArray(key string) []string {
	// Call super method
	return IValueBean(p).GetAsStringArray(key)
}

// ToString stringify this commnd
func (p *NodeBean) ToString() string {
	// Call super method
	return IValueBean(p).ToString()
}

// ToJSON stringify this commnd
func (p *NodeBean) ToJSON() string {
	// Call super method
	return IValueBean(p).ToJSON()
}

// SetTimestamp set timestamp
func (p *NodeBean) SetTimestamp(stamp JSONTime) {
	p.Timestamp = stamp
}

// GetTimestamp get timestamp
func (p *NodeBean) GetTimestamp() JSONTime {
	return p.Timestamp
}

// Copy retrieve ID
func (p *NodeBean) Copy() IPersistent {
	clone := *p
	return &clone
}

// NodeBeans simple bean model
type NodeBeans struct {
	// Collection
	Collection []IPersistent `json:"collections"`
}

// New constructor
func (p *NodeBeans) New() IPersistents {
	bean := NodeBeans{Collection: make([]IPersistent, 0)}
	return &bean
}

// Add new bean
func (p *NodeBeans) Add(bean IPersistent) {
	p.Collection = append(p.Collection, bean)
}

// Get collection of bean
func (p *NodeBeans) Get() []IPersistent {
	return p.Collection
}

// Index read a single element
func (p *NodeBeans) Index(index int) INodeBean {
	data, ok := p.Collection[index].(*NodeBean)
	if ok {
		return data
	}
	return nil
}
