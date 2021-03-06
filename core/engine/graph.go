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
package engine

import (
	"github.com/yroffin/go-boot-sqllite/core/models"
	"github.com/yroffin/go-boot-sqllite/core/winter"
)

// GraphCrudBusiness internal members
type GraphCrudBusiness struct {
	*winter.Service
	// Store with injection mecanism
	Store IGraphStore `@autowired:"cayley-manager"`
}

// New constructor
func (p *GraphCrudBusiness) New() ILinkBusiness {
	bean := GraphCrudBusiness{Service: &winter.Service{Bean: &winter.Bean{}}}
	return &bean
}

// Init this bean
func (p *GraphCrudBusiness) Init() error {
	return nil
}

// Clear this bean
func (p *GraphCrudBusiness) Clear() error {
	p.Store.Clear()
	return nil
}

// Statistics some statistics
func (p *GraphCrudBusiness) All() ([]IQuad, error) {
	return p.Store.All()
}

// Statistics some statistics
func (p *GraphCrudBusiness) Statistics() ([]IStats, error) {
	return p.Store.Statistics()
}

// Export some values
func (p *GraphCrudBusiness) Export() (map[string][]map[string]interface{}, error) {
	return p.Store.Export()
}

// PostConstruct this bean
func (p *GraphCrudBusiness) PostConstruct(name string) error {
	return nil
}

// Validate this bean
func (p *GraphCrudBusiness) Validate(name string) error {
	return nil
}

// CreateLink retrieve this link
func (p *GraphCrudBusiness) CreateLink(toCreate models.IEdgeBean) (models.IEdgeBean, error) {
	return toCreate, p.Store.CreateLink(toCreate)
}

// UpdateLink retrieve this link
func (p *GraphCrudBusiness) UpdateLink(toUpdate models.IEdgeBean) (models.IEdgeBean, error) {
	return toUpdate, p.Store.UpdateLink(toUpdate)
}

// GetAllLink retrieve this bean by its id
func (p *GraphCrudBusiness) GetAllLink(model string, id string, toGets []models.IEdgeBean, targetType string) ([]models.IEdgeBean, error) {
	p.Store.GetAllLink(model, id, &toGets, targetType)
	return toGets, nil
}

// DeleteLink a bean
func (p *GraphCrudBusiness) DeleteLink(toDelete models.IEdgeBean) (models.IEdgeBean, error) {
	return toDelete, p.Store.DeleteLink(toDelete)
}

// TruncateLink a bean
func (p *GraphCrudBusiness) TruncateLink(toTruncate models.IPersistent) (models.IPersistent, error) {
	p.Store.TruncateLink(toTruncate)
	return toTruncate, nil
}

// PatchLink a bean
func (p *GraphCrudBusiness) PatchLink(toPatch models.IEdgeBean) (models.IEdgeBean, error) {
	return toPatch, p.Store.CreateLink(toPatch)
}
