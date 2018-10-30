// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jgin

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/go-xorm/core"
)

// NETCacher implments cache object facilities
type NETCacher struct {
	idList         *list.List
	sqlList        *list.List
	idIndex        map[string]map[string]*list.Element
	sqlIndex       map[string]map[string]*list.Element
	store          core.CacheStore
	mutex          sync.Mutex
	Expired        time.Duration
	GcInterval     time.Duration
}

// XromNetCacher creates a cacher
func XromNetCacher(store core.CacheStore) *NETCacher {
	return XromNetCacher2(store)
}

// XromNetCacher2 creates a cache include different params
func XromNetCacher2(store core.CacheStore) *NETCacher {
	cacher := &NETCacher{store: store}
	cacher.RunGC()
	return cacher
}

// RunGC run once every m.GcInterval
func (m *NETCacher) RunGC() {
	//time.AfterFunc(m.GcInterval, func() {
	//	m.RunGC()
	//	m.GC()
	//})
}

// GC check ids lit and sql list to remove all element expired
func (m *NETCacher) GC() {
}

// GetIds returns all bean's ids according to sql and parameter from cache
func (m *NETCacher) GetIds(tableName, sql string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, err := m.store.Get(sql); err == nil {
		return v
	}
	return nil
}

// GetBean returns bean according tableName and id from cache
func (m *NETCacher) GetBean(tableName string, id string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tid := genID(tableName, id)
	if v, err := m.store.Get(tid); err == nil {
		return v
	}
	return nil
}


// ClearIds clears all sql-ids mapping on table tableName from cache
func (m *NETCacher) ClearIds(tableName string) {
}

// ClearBeans clears all beans in some table
func (m *NETCacher) ClearBeans(tableName string) {
}

// PutIds pus ids into table
func (m *NETCacher) PutIds(tableName, sql string, ids interface{}) {
	m.mutex.Lock()
	m.store.Put(sql, ids)
	m.mutex.Unlock()
}

// PutBean puts beans into table
func (m *NETCacher) PutBean(tableName string, id string, obj interface{}) {
	m.mutex.Lock()
	m.store.Put(genID(tableName, id), obj)
	m.mutex.Unlock()
}

// DelIds deletes ids
func (m *NETCacher) DelIds(tableName, sql string) {
}

// DelBean deletes beans in some table
func (m *NETCacher) DelBean(tableName string, id string) {
}

type idNode struct {
	tbName    string
	id        string
	lastVisit time.Time
}

type sqlNode struct {
	tbName    string
	sql       string
	lastVisit time.Time
}

func genSQLKey(sql string, args interface{}) string {
	return fmt.Sprintf("%v-%v", sql, args)
}

func genID(prefix string, id string) string {
	return fmt.Sprintf("%v-%v", prefix, id)
}

func newIDNode(tbName string, id string) *idNode {
	return &idNode{tbName, id, time.Now()}
}

func newSQLNode(tbName, sql string) *sqlNode {
	return &sqlNode{tbName, sql, time.Now()}
}
