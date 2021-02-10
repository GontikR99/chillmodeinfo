// +build server

package db

import "github.com/timshannon/bolthold"

var database *bolthold.Store

func init() {
	var err error
	database, err = bolthold.Open("chillmodeinfo.db", 0600, nil)
	if err!=nil {
		panic(err)
	}
}

func Insert(key interface{}, value interface{}) error {
	return database.Insert(key, value)
}

func Update(key interface{}, value interface{}) error {
	return database.Update(key, value)
}

func Upsert(key interface{}, value interface{}) error {
	return database.Upsert(key, value)
}

func Delete(key interface{}, dataType interface{}) error {
	return database.Delete(key, dataType)
}

func Get(key interface{}, result interface{}) error {
	return database.Get(key, result)
}

func Find(result interface{}, query *bolthold.Query) error {
	return database.Find(result, query)
}

func ForEach(query *bolthold.Query, callback interface{}) error {
	return database.ForEach(query, callback)
}
