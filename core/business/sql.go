// Package business for business interface
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
package business

import (
	"log"
	"reflect"

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
	core_services "github.com/yroffin/go-boot-sqllite/core/services"
	"github.com/yroffin/go-boot-sqllite/core/stores"
)

// SqlCrudBusiness internal members
type SqlCrudBusiness struct {
	// members
	*core_services.SERVICE
	// Store with injection mecanism
	Store stores.IStore `@autowired:"sqllite-manager"`
}

// New constructor
func (p *SqlCrudBusiness) New() ICrudBusiness {
	bean := SqlCrudBusiness{SERVICE: &core_services.SERVICE{Bean: &core_bean.Bean{}}}
	return &bean
}

// SetSqlliteManager injection
func (p *SqlCrudBusiness) SetSqlliteManager(value interface{}) {
	if assertion, ok := value.(stores.IStore); ok {
		p.Store = assertion
	} else {
		log.Fatalf("Unable to validate injection with %v type is %v", value, reflect.TypeOf(value))
	}
}

// Init this bean
func (p *SqlCrudBusiness) Init() error {
	return nil
}

// PostConstruct this bean
func (p *SqlCrudBusiness) PostConstruct(name string) error {
	return nil
}

// Validate this bean
func (p *SqlCrudBusiness) Validate(name string) error {
	return nil
}

// GetAll retrieve this bean by its id
func (p *SqlCrudBusiness) GetAll(toGet models.IPersistent, toGets models.IPersistents) (models.IPersistents, error) {
	p.Store.GetAll(toGet, toGets)
	return toGets, nil
}

// Get retrieve this bean by its id
func (p *SqlCrudBusiness) Get(toGet models.IPersistent) (models.IPersistent, error) {
	p.Store.Get(toGet.GetID(), toGet)
	return toGet, nil
}

// Create create a new persistent bean
func (p *SqlCrudBusiness) Create(toCreate models.IPersistent) (models.IPersistent, error) {
	p.Store.Create(toCreate)
	return toCreate, nil
}

// Update an existing bean
func (p *SqlCrudBusiness) Update(toUpdate models.IPersistent) (models.IPersistent, error) {
	p.Store.Update(toUpdate.GetID(), toUpdate)
	return toUpdate, nil
}

// Delete a bean
func (p *SqlCrudBusiness) Delete(toDelete models.IPersistent) (models.IPersistent, error) {
	p.Store.Delete(toDelete.GetID(), toDelete)
	return toDelete, nil
}

// Delete a bean
func (p *SqlCrudBusiness) Truncate(toTruncate models.IPersistent) (models.IPersistent, error) {
	p.Store.Truncate(toTruncate)
	return toTruncate, nil
}

// Patch a bean
func (p *SqlCrudBusiness) Patch(toPatch models.IPersistent) (models.IPersistent, error) {
	p.Store.Update(toPatch.GetID(), toPatch)
	return toPatch, nil
}
