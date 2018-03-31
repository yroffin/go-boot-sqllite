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
	"fmt"
	"log"
	"strings"
	"time"
)

// JSONTime for json datetime handle
type JSONTime time.Time

// MarshalJSON for json datetime handle
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339Nano))
	return []byte(stamp), nil
}

// UnmarshalJSON for json datetime handle
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	var ns = strings.Replace(string(data), "\"", "", 1000)
	v, err := time.Parse(time.RFC3339Nano, ns)
	if err != nil {
		log.Printf("Error while decoding '%s' -> '%s'", ns, err)
	}
	*t = JSONTime(v)
	return nil
}

// IPersistent interface
type IPersistent interface {
	GetID() string
	GetName() string
	SetID(string)
	GetTimestamp() JSONTime
	SetTimestamp(JSONTime)
	Copy() IPersistent
}

// IPersistents interface
type IPersistents interface {
	Add(IPersistent)
	Get() []IPersistent
}
