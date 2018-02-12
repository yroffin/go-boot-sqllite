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
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	// for import driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/yroffin/go-boot-sqllite/core/bean"
	"github.com/yroffin/go-boot-sqllite/core/models"
)

// Store internal members
type Store struct {
	// Base component
	*bean.Bean
	// Store SQL lite
	database *sql.DB
	// Tables
	Tables []string
	// Db path
	DbPath string
}

// IStore interface
type IStore interface {
	bean.IBean
}

// Init Init this bean
func (p *Store) Init() error {
	return nil
}

// PostConstruct this bean
func (p *Store) PostConstruct(name string) error {
	// Create database
	database, _ := sql.Open("sqlite3", p.DbPath)
	p.database = database

	// create all tables
	for i := 0; i < len(p.Tables); i++ {
		// prepare statement
		statement, _ := p.database.Prepare("CREATE TABLE IF NOT EXISTS " + p.Tables[i] + " (id TEXT NOT NULL PRIMARY KEY, json TEXT)")
		statement.Exec()
	}

	return nil
}

// Validate Init this bean
func (p *Store) Validate(name string) error {
	return nil
}

// uuid generates a random UUID according to RFC 4122
func (p *Store) uuid(entity interface{}) (string, error) {
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
func (p *Store) Create(entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// Fix timestamp
	entity.SetTimestamp(models.JSONTime(time.Now()))
	// fix UUID
	uuid, _ := p.uuid(entity)
	entity.SetID(uuid)
	// insert
	statement, _ := p.database.Prepare("INSERT INTO " + entityName + " (id, json) VALUES (?,?)")
	data, _ := json.Marshal(entity)
	statement.Exec(uuid, string(data))
	return nil
}

// Update this persistent bean
func (p *Store) Update(id string, entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// Fix timestamp
	entity.SetTimestamp(models.JSONTime(time.Now()))
	// Fix ID
	entity.SetID(id)
	// prepare statement
	statement, _ := p.database.Prepare("UPDATE " + entityName + " SET json = ? WHERE id = ?")
	data, _ := json.Marshal(entity)
	res, _ := statement.Exec(string(data), id)
	rowAffected, _ := res.RowsAffected()
	if rowAffected == 0 {
		log.Printf("'%s' with id '%v' affected %d row(s)", "UPDATE "+entityName+" SET json = ? WHERE id = ?", id, rowAffected)
	}
	return nil
}

// Delete this persistent bean
func (p *Store) Delete(id string, entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// Fix ID
	entity.SetID(id)
	// prepare statement
	statement, _ := p.database.Prepare("DELETE FROM " + entityName + " WHERE id = ?")
	res, _ := statement.Exec(id)
	rowAffected, _ := res.RowsAffected()
	if rowAffected == 0 {
		log.Printf("'%s' with id '%v' affected %d row(s)", "DELETE FROM "+entityName+" WHERE id = ?", id, rowAffected)
	}
	return nil
}

// Truncate method
func (p *Store) Truncate(entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// prepare statement
	statement, _ := p.database.Prepare("DELETE FROM " + entityName)
	res, _ := statement.Exec()
	rowAffected, _ := res.RowsAffected()
	if rowAffected == 0 {
		log.Printf("'%s' affected %d row(s)", "DELETE FROM "+entityName, rowAffected)
	}
	return nil
}

// Get this persistent bean
func (p *Store) Get(id string, entity models.IPersistent) error {
	// get entity name
	var entityName = entity.SetName()
	// prepare statement
	rows, _ := p.database.Query("SELECT id, json FROM "+entityName+" WHERE id = ?", &id)
	var data string
	for rows.Next() {
		rows.Scan(&id, &data)
		var bin = []byte(data)
		entity.SetID(id)
		json.Unmarshal(bin, entity)
	}
	return nil
}

// GetAll this persistent bean
func (p *Store) GetAll(entity models.IPersistent, array models.IPersistents) error {
	// get entity name
	var entityName = entity.SetName()
	// prepare statement
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
	return nil
}
