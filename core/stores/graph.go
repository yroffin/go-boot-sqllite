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
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/query"
	"github.com/cayleygraph/cayley/query/gizmo"
	// bolt
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/quad"

	core_bean "github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
	core_services "github.com/yroffin/go-boot-sqllite/core/services"
)

// Graph internal members
type Graph struct {
	// members
	*core_services.SERVICE
	// Store SQL lite
	store *graph.Handle
	// Tables
	Tables []string
	// Db path
	DbPath string
}

// New constructor
func (p *Graph) New(tables []string, dbpath string) IGraphStore {
	bean := Graph{SERVICE: &core_services.SERVICE{Bean: &core_bean.Bean{}}, Tables: tables, DbPath: dbpath}
	return &bean
}

// Init Init this bean
func (p *Graph) Init() error {
	return nil
}

// PostConstruct this bean
func (p *Graph) PostConstruct(name string) error {
	// Initialize the database
	graph.InitQuadStore("bolt", p.DbPath, nil)

	// Open and use the database
	database, err := cayley.NewGraph("bolt", p.DbPath, nil)
	if err != nil {
		log.Fatalln(err)
	}
	p.store = database

	return nil
}

// Validate Init this bean
func (p *Graph) Validate(name string) error {
	return nil
}

// Clear Init this bean
func (p *Graph) Clear() error {
	it := p.store.QuadsAllIterator()
	for it.Next(context.Background()) {
		qu := p.store.Quad(it.Result())
		tx := cayley.NewTransaction()
		tx.RemoveQuad(qu)
		p.store.ApplyTransaction(tx)
	}

	return nil
}

// Statistics some statistics
func (p *Graph) Statistics() ([]IStats, error) {
	stats := make([]IStats, 0)
	it := p.store.QuadsAllIterator()
	for it.Next(context.Background()) {
		qu := p.store.Quad(it.Result())
		stat := StoreStats{}
		stat.Key = qu.Subject.Native().(string)
		stat.Value = qu.Predicate.Native().(string) + " " + qu.Object.Native().(string) + " " + qu.Label.Native().(string)
		stats = append(stats, &stat)
	}
	return stats, nil
}

// uuid generates a random UUID according to RFC 4122
func (p *Graph) uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	var text = fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return text, nil
}

// CreateLink in graph db
func (p *Graph) CreateLink(data models.IEdgeBean) error {
	// fix UUID
	uuid, _ := p.uuid()
	data.SetID(uuid)
	// insert
	jsonData, _ := json.Marshal(data)
	quad := quad.Make("/"+data.GetSource()+"/"+data.GetSourceID(), data.GetLink()+":"+uuid, "/"+data.GetTarget()+"/"+data.GetTargetID(), string(jsonData))
	log.Println("Create:", quad, string(jsonData))
	p.store.AddQuad(quad)
	return nil
}

// DeleteLink this persistent bean
func (p *Graph) DeleteLink(toDelete models.IEdgeBean) error {
	it := p.store.QuadsAllIterator()
	for it.Next(context.Background()) {
		qu := p.store.Quad(it.Result())
		if qu.Predicate.Native().(string) == toDelete.GetLink()+":"+toDelete.GetID() {
			log.Println("Remove:", qu.Subject.Native(), qu.Predicate.Native(), qu.Object.Native())
			tx := cayley.NewTransaction()
			tx.RemoveQuad(qu)
			p.store.ApplyTransaction(tx)
		}
	}

	return nil
}

// TruncateLink method
func (p *Graph) TruncateLink(entity models.IPersistent) error {
	return nil
}

// GetLink this persistent bean
func (p *Graph) GetLink(entity models.IEdgeBean) error {
	return nil
}

// GetAllLink this persistent bean
func (p *Graph) GetAllLink(model string, id string, array *[]models.IEdgeBean) error {
	var query = `g.V('/` + model + `/` + id + `').As('source').Out(null, 'edge').As('target').Labels().As('label').All()`
	results, _ := p.QueryGizmo(query, "")
	for _, v := range results {
		if id == v.GetSourceID() {
			*array = append(*array, v)
		}
	}

	return nil
}

// QueryGizmo query gizmo
func (p *Graph) QueryGizmo(text string, tag string) ([]models.IEdgeBean, error) {
	session := gizmo.NewSession(p.store)
	c := make(chan query.Result, 1)
	go func() {
		session.Execute(context.TODO(), text, c, -1)
	}()

	resultSet := make([]models.IEdgeBean, 0)

	for res := range c {
		if err := res.Err(); err != nil {
			return nil, err
		}
		switch result := res.(type) {
		case *gizmo.Result:
			// Tags are source, target, edge
			// they also stored in native labels
			data := models.EdgeBean{}
			json.Unmarshal([]byte(p.store.NameOf(result.Tags["label"]).Native().(string)), &data)
			resultSet = append(resultSet, &data)
			break
		default:
			log.Println("Unknown:", res)
		}
	}
	return resultSet, nil
}
