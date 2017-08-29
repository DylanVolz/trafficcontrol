package db

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"errors"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type Query string
type Queryx struct {
	Query  Query
	Params []interface{}
}
type DB struct {
	*sqlx.DB
}
type Tx struct {
	*sqlx.Tx
}

var (
	//ErrNoGetterFound   = errors.New("No getter found")
	ErrNoDeleterFound  = errors.New("No deleter found")
	//ErrNoSelecterFound = errors.New("No getter found")
	ErrNoUpdaterFound  = errors.New("No updater found")
	ErrNoInserterFound = errors.New("No inserter found")
)

func Limit(offset, count int) selectOption {
	return &limitOption{offset, count}
}

type limitOption struct {
	offset int
	count  int
}

func (o *limitOption) Wrap(query string, params []interface{}) (string, []interface{}) {
	query = fmt.Sprintf("SELECT a.* FROM (%s) a LIMIT ?, ?", query)
	params = append(params, o.offset)
	params = append(params, o.count)
	return query, params
}




type selectOption interface {
	Wrap(string, []interface{}) (string, []interface{})
}






func (tx *Tx) Selectx(o interface{}, qx Queryx, options ...selectOption) error {
	q := string(qx.Query)
	params := qx.Params

	log.Debugln(q)
	for _, option := range options {
		q, params = option.Wrap(q, params)
	}
	if u, ok := o.(Selecter); ok {
		return u.Select(tx.Tx, Query(q), params...)
	}
	stmt, err := tx.Preparex(q)
	if err != nil {
		return err
	}
	return stmt.Select(o, params...)
}






func (tx *Tx) Countx(qx Queryx) (int, error) {
	stmt, err := tx.Preparex(fmt.Sprintf("SELECT COUNT(*) FROM (%s) q", string(qx.Query)))
	if err != nil {
		return 0, err
	}
	count := 0
	err = stmt.Get(&count, qx.Params...)
	return count, err
}






func (tx *Tx) Getx(o interface{}, qx Queryx) error {
	if u, ok := o.(Getter); ok {
		return u.Get(tx.Tx, qx.Query, qx.Params...)
	}
	stmt, err := tx.Preparex(string(qx.Query))
	if err != nil {
		return err
	}
	return stmt.Get(o, qx.Params...)
}




func (tx *Tx) Get(o interface{}, query Query, params ...interface{}) error {
	if u, ok := o.(Getter); ok {
		return u.Get(tx.Tx, query, params...)
	}
	stmt, err := tx.Preparex(string(query))
	if err != nil {
		return err
	}
	return stmt.Get(o, params...)
}





func (tx *Tx) Update(o interface{}) error {
	if u, ok := o.(Updater); ok {
		return u.Update(tx.Tx)
	}
	log.Debugln("No updater found for object: %s", reflect.TypeOf(o))
	return ErrNoUpdaterFound
}





func (tx *Tx) Delete(o interface{}) error {
	if u, ok := o.(Deleter); ok {
		return u.Delete(tx.Tx)
	}
	log.Debugln("No deleter found for object: %s", reflect.TypeOf(o))
	return ErrNoDeleterFound
}





func (tx *Tx) Insert(o interface{}) error {
	if u, ok := o.(Inserter); ok {
		err := u.Insert(tx.Tx)
		if err != nil {
			log.Errorln(err.Error())
		}
		return err
	}
	log.Debugln("No inserter found for object: %s", reflect.TypeOf(o))
	return ErrNoInserterFound
}






func (db *DB) Begin() *Tx {
	tx := db.MustBegin()
	return &Tx{tx}
}








type Updater interface {
	Update(*sqlx.Tx) error
}
type Inserter interface {
	Insert(*sqlx.Tx) error
}
type Selecter interface {
	Select(*sqlx.Tx, Query, ...interface{}) error
}
type Getter interface {
	Get(*sqlx.Tx, Query, ...interface{}) error
}
type Deleter interface {
	Delete(*sqlx.Tx) error
}
