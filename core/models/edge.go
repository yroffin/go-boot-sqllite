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

	log "github.com/sirupsen/logrus"
)

// EdgeBean simple command model
type EdgeBean struct {
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
	// Source
	Source string `json:"source"`
	// SourceID
	SourceID string `json:"sourceId"`
	// Target
	Target string `json:"target"`
	// TargetID
	TargetID string `json:"targetId"`
	// Link
	Link string `json:"link"`
	// Instance
	Instance string `json:"instance"`
}

// IEdgeBean interface
type IEdgeBean interface {
	// inherit persistent behaviour
	IPersistent
	// inherit ValueBean behaviour
	IValueBean
	// get Type
	GetType() string
	GetSource() string
	GetSourceID() string
	GetTarget() string
	GetTargetID() string
	GetLink() string
	SetInstance(string)
	GetInstance() string
	GetExtend() map[string]interface{}
}

// New constructor
func (p *EdgeBean) New(source string, sourceID string, target string, targetID string, link string) IEdgeBean {
	bean := EdgeBean{}
	bean.Extended = make(map[string]interface{})
	bean.Source = source
	bean.Source = source
	bean.SourceID = sourceID
	bean.Target = target
	bean.TargetID = targetID
	bean.Link = link
	return &bean
}

// SetInstance set instance
func (p *EdgeBean) SetInstance(instance string) {
	p.Instance = instance
}

// GetExtend get extended
func (p *EdgeBean) GetExtend() map[string]interface{} {
	if p.Extended == nil {
		p.Extended = make(map[string]interface{})
	}
	return p.Extended
}

// GetInstance get instance
func (p *EdgeBean) GetInstance() string {
	return p.Instance
}

// GetSource get source name for this edge
func (p *EdgeBean) GetSource() string {
	return p.Source
}

// GetSourceID get source name for this edge
func (p *EdgeBean) GetSourceID() string {
	return p.SourceID
}

// GetTarget get source name for this edge
func (p *EdgeBean) GetTarget() string {
	return p.Target
}

// GetTargetID get source name for this edge
func (p *EdgeBean) GetTargetID() string {
	return p.TargetID
}

// GetLink get source name for this edge
func (p *EdgeBean) GetLink() string {
	return p.Link
}

// GetEntityName get set name
func (p *EdgeBean) GetEntityName() string {
	return "Edge"
}

// GetType get set name
func (p *EdgeBean) GetType() string {
	return p.Type
}

// GetID retrieve ID
func (p *EdgeBean) GetID() string {
	return p.ID
}

// SetID retrieve ID
func (p *EdgeBean) SetID(ID string) {
	p.ID = ID
}

// Set get set name
func (p *EdgeBean) Set(key string, value interface{}) {
}

// Extend vars
func (p *EdgeBean) Extend(e map[string]interface{}) {
	// edge and instance are reserved
	for k, v := range e {
		switch k {
		case "edge":
		case "instance":
			break
		default:
			p.Extended[k] = v
		}
	}
}

// SetString get set name
func (p *EdgeBean) SetString(key string, value string) {
	// Call super method
	IValueBean(p).SetString(key, value)
}

// Get get set name
func (p *EdgeBean) GetAsString(key string) string {
	switch key {
	default:
		// Call super method
		return IValueBean(p).GetAsString(key)
	}
}

// Get get set name
func (p *EdgeBean) GetAsStringArray(key string) []string {
	// Call super method
	return IValueBean(p).GetAsStringArray(key)
}

// ToJSON stringify this commnd
func (p *EdgeBean) ToJSON() string {
	payload, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to marshal")
		return "{}"
	}
	return string(payload)
}

// SetTimestamp set timestamp
func (p *EdgeBean) SetTimestamp(stamp JSONTime) {
	p.Timestamp = stamp
}

// GetTimestamp get timestamp
func (p *EdgeBean) GetTimestamp() JSONTime {
	return p.Timestamp
}

// Copy retrieve ID
func (p *EdgeBean) Copy() IPersistent {
	clone := *p
	return &clone
}

// EdgeBeans simple bean model
type EdgeBeans struct {
	// Collection
	Collection []IPersistent `json:"collections"`
}

// New constructor
func (p *EdgeBeans) New() IPersistents {
	bean := EdgeBeans{Collection: make([]IPersistent, 0)}
	return &bean
}

// Add new bean
func (p *EdgeBeans) Add(bean IPersistent) {
	p.Collection = append(p.Collection, bean)
}

// Get collection of bean
func (p *EdgeBeans) Get() []IPersistent {
	return p.Collection
}

// Index read a single element
func (p *EdgeBeans) Index(index int) INodeBean {
	data, ok := p.Collection[index].(*NodeBean)
	if ok {
		return data
	}
	return nil
}
