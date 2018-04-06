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

// ICrudBusiness interface
type ICrudBusiness interface {
	winter.IBean
	// Relationnal data
	GetAll(models.IPersistent, models.IPersistents) (models.IPersistents, error)
	Get(models.IPersistent) (models.IPersistent, error)
	Create(models.IPersistent) (models.IPersistent, error)
	Update(models.IPersistent) (models.IPersistent, error)
	Delete(models.IPersistent) (models.IPersistent, error)
	Patch(models.IPersistent) (models.IPersistent, error)
	Clear([]string) error
	Statistics() ([]IStats, error)
}

// ILinkBusiness interface
type ILinkBusiness interface {
	winter.IBean
	// Linked ones
	CreateLink(toCreate models.IEdgeBean) (models.IEdgeBean, error)
	UpdateLink(toUpdate models.IEdgeBean) (models.IEdgeBean, error)
	DeleteLink(toCreate models.IEdgeBean) (models.IEdgeBean, error)
	PatchLink(toPatch models.IEdgeBean) (models.IEdgeBean, error)
	GetAllLink(model string, id string, toGets []models.IEdgeBean, targetType string) ([]models.IEdgeBean, error)
	Clear() error
	Statistics() ([]IStats, error)
	Export() (map[string][]map[string]interface{}, error)
}
