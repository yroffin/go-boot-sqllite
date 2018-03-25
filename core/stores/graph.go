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
	"time"

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
func (p *Graph) New(tables []string, dbpath string) IStore {
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

// uuid generates a random UUID according to RFC 4122
func (p *Graph) uuid(entity interface{}) (string, error) {
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

// Create this persistent bean n store
func (p *Graph) Create(entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// Fix timestamp
	entity.SetTimestamp(models.JSONTime(time.Now()))
	// fix UUID
	uuid, _ := p.uuid(entity)
	entity.SetID(uuid)
	// insert
	data, _ := json.Marshal(entity)
	p.store.AddQuad(quad.Make("/"+entityName+"/"+uuid, "is", entityName, nil))
	p.store.AddQuad(quad.Make("/"+entityName+"/"+uuid, "uuid", uuid, nil))
	p.store.AddQuad(quad.Make("/"+entityName+"/"+uuid, "class", "XXXX", nil))
	p.store.AddQuad(quad.Make("/"+entityName+"/"+uuid, "link", string(data), nil))
	return nil
}

// Update this persistent bean
func (p *Graph) Update(id string, entity models.IPersistent) error {
	// get entity name
	//var entityName = entity.SetName()
	// Fix timestamp
	entity.SetTimestamp(models.JSONTime(time.Now()))
	// Fix ID
	entity.SetID(id)
	// prepare statement
	//node := cayley.StartPath(p.store).Has(id)
	/*
		data, _ := json.Marshal(entity)
		res, _ := statement.Exec(string(data), id)
		rowAffected, _ := res.RowsAffected()
		if rowAffected == 0 {
			log.Printf("'%s' with id '%v' affected %d row(s)", "UPDATE "+entityName+" SET json = ? WHERE id = ?", id, rowAffected)
		}
	*/
	return nil
}

// Delete this persistent bean
func (p *Graph) Delete(id string, entity models.IPersistent) error {
	// get entity name
	//var entityName = entity.SetName()
	// Fix ID
	entity.SetID(id)
	// prepare statement
	/*
		statement, _ := p.database.Prepare("DELETE FROM " + entityName + " WHERE id = ?")
		res, _ := statement.Exec(id)
		rowAffected, _ := res.RowsAffected()
		if rowAffected == 0 {
			log.Printf("'%s' with id '%v' affected %d row(s)", "DELETE FROM "+entityName+" WHERE id = ?", id, rowAffected)
		}
	*/
	return nil
}

// Truncate method
func (p *Graph) Truncate(entity models.IPersistent) error {
	// get entity name
	//var entityName = entity.SetName()
	// prepare statement
	/*
		statement, _ := p.database.Prepare("DELETE FROM " + entityName)
		res, _ := statement.Exec()
		rowAffected, _ := res.RowsAffected()
		if rowAffected == 0 {
			log.Printf("'%s' affected %d row(s)", "DELETE FROM "+entityName, rowAffected)
		}
	*/
	return nil
}

// Get this persistent bean
func (p *Graph) Get(id string, entity models.IPersistent) error {
	// get entity name
	//var entityName = entity.SetName()
	// prepare statement
	/*
		rows, _ := p.database.Query("SELECT id, json FROM "+entityName+" WHERE id = ?", &id)
		var data string
		for rows.Next() {
			rows.Scan(&id, &data)
			var bin = []byte(data)
			entity.SetID(id)
			json.Unmarshal(bin, entity)
		}
	*/
	return nil
}

// GetAll this persistent bean
func (p *Graph) GetAll(entity models.IPersistent, array models.IPersistents) error {
	// get entity name
	//var entityName = entity.SetName()

	log.Println("Size:", p.store.Size())

	var query = `
var output = g.V().As('source').Has('is', 'Node').Out(null, 'edge').As('target').All()
g.Emit(output)
`
	p.QueryGizmo(query, "")

	query = `
var output = g.V('/Node/fd30437c-db55-456e-b8b9-dc6401260f32').As('source').Has('is', 'Node').Out(null, 'edge').As('target').All()
g.Emit(output)
`
	p.QueryGizmo(query, "")

	// get entity name
	//var entityName = entity.SetName()
	// prepare statement
	/*
		rows, _ := p.database.Query("SELECT id, json FROM " + entityName)
		var id string
		var data string
		for rows.Next() {
			rows.Scan(&id, &data)
			var bin = []byte(data)
			copy := entity.Copy()
			json.Unmarshal(bin, &copy)
			copy.SetID(id)
			array.Add(copy)
		}
	*/
	return nil
}

func quadValueToString(v quad.Value) string {
	if s, ok := v.(quad.String); ok {
		return string(s)
	}
	return quad.StringOf(v)
}

// Query query gizmo
func (p *Graph) QueryGizmo(text string, tag string) ([]string, error) {
	session := gizmo.NewSession(p.store)
	c := make(chan query.Result, 1)
	go func() {
		session.Execute(context.TODO(), text, c, -1)
	}()

	// resultSet := make(map[string]map[string]string)

	var results []string
	for res := range c {
		if err := res.Err(); err != nil {
			return results, err
		}
		switch result := res.(type) {
		case *gizmo.Result:
			//log.Println("Result/Meta:", result.Meta)
			//log.Println("Result/Tags:", result.Tags)
			log.Println("Source:", p.store.NameOf(result.Tags["source"]).String(), result.Tags["source"].Key(), "Edge:", p.store.NameOf(result.Tags["edge"]).String(), "Target:", p.store.NameOf(result.Tags["target"]).String())
			break
		default:
			log.Println("Unknown:", res)
		}
	}
	return results, nil
}
