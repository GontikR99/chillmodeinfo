// +build server

package db

import (
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

var database *bolthold.Store

func init() {
	var err error
	database, err = bolthold.Open("chillmodeinfo.db", 0640, nil)
	if err!=nil {
		panic(err)
	}
}

//func Insert(key interface{}, value interface{}) error {
//	return database.Insert(key, value)
//}

func TxInsert(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxInsert(tx, key, value)
}

//func Update(key interface{}, value interface{}) error {
//	return database.Update(key, value)
//}

func TxUpdate(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxUpdate(tx, key, value)
}

//func Upsert(key interface{}, value interface{}) error {
//	return database.Upsert(key, value)
//}

func TxUpsert(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxUpsert(tx, key, value)
}

//func Delete(key interface{}, dataType interface{}) error {
//	return database.Delete(key, dataType)
//}

func TxDelete(tx *bbolt.Tx, key interface{}, dataType interface{}) error {
	return database.TxDelete(tx, key, dataType)
}

//func Get(key interface{}, result interface{}) error {
//	return database.Get(key, result)
//}

func TxGet(tx *bbolt.Tx, key interface{}, result interface{}) error {
	return database.TxGet(tx, key, result)
}

//func Find(result interface{}, query *bolthold.Query) error {
//	return database.Find(result, query)
//}

func TxFind(tx *bbolt.Tx, result interface{}, query *bolthold.Query) error {
	return database.TxFind(tx, result, query)
}

//func ForEach(query *bolthold.Query, callback interface{}) error {
//	return database.ForEach(query, callback)
//}

func TxForEach(tx *bbolt.Tx, query *bolthold.Query, callback interface{}) error {
	return database.TxForEach(tx, query, callback)
}

//func DeleteMatching(dataType interface{}, query *bolthold.Query) error {
//	return database.DeleteMatching(dataType, query)
//}

func TxDeleteMatching(tx *bbolt.Tx, dataType interface{}, query *bolthold.Query) error {
	return database.TxDeleteMatching(tx, dataType, query)
}

func MakeView(tables []TableName, body func(tx *bbolt.Tx)error)error {
	unlocker := acquireRead(tables...)
	defer unlocker()
	return database.Bolt().View(body)
}

func MakeUpdate(tables []TableName, body func(tx *bbolt.Tx)error)error {
	unlocker := acquireWrite(tables...)
	defer unlocker()
	return database.Bolt().Update(body)
}