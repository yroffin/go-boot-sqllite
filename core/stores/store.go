// Package stores for all sgbd operation
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
package stores

import (

	// for import driver
	_ "github.com/mattn/go-sqlite3"

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
)

// IStats stats
type IStats interface {
	GetKey() string
	GetValue() string
}

// StoreStats statss
type StoreStats struct {
	Key   string
	Value string
}

// GetKey some statistics
func (p *StoreStats) GetKey() string {
	return p.Key
}

// GetValue some statistics
func (p *StoreStats) GetValue() string {
	return p.Value
}

// IDataStore interface
type IDataStore interface {
	core_bean.IBean
	Create(entity models.IPersistent) error
	Update(id string, entity models.IPersistent) error
	Delete(id string, entity models.IPersistent) error
	Truncate(entity models.IPersistent) error
	Get(id string, entity models.IPersistent) error
	GetAll(entity models.IPersistent, array models.IPersistents) error
	Clear([]string) error
	Statistics() ([]IStats, error)
}

// IGraphStore interface
type IGraphStore interface {
	core_bean.IBean
	CreateLink(data models.IEdgeBean) error
	DeleteLink(entity models.IEdgeBean) error
	TruncateLink(entity models.IPersistent) error
	GetLink(entity models.IEdgeBean) error
	GetAllLink(model string, id string, collection *[]models.IEdgeBean) error
	Clear() error
	Statistics() ([]IStats, error)
}
