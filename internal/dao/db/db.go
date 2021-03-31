// +build server

package db

import (
	"compress/gzip"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"os"
	"time"
)

var database *bolthold.Store

func init() {
	var err error
	database, err = bolthold.Open("chillmodeinfo.db", 0640, nil)
	if err!=nil {
		panic(err)
	}
	go func() {
		for {
			<-time.After(100*time.Millisecond)
			if time.Now().Minute()==59 && time.Now().Hour()%12==11 {
				<-time.After(1*time.Minute)
			} else {
				continue
			}
			database.Bolt().View(func(tx *bbolt.Tx) error {
				out, err := os.Create("chillmodeinfo-"+time.Now().Format("2006-01-02T15:04")+".db.gz")
				if err!=nil {
					return err
				}
				defer out.Close()
				outzip := gzip.NewWriter(out)
				defer outzip.Close()
				_, err = tx.WriteTo(outzip)
				return err
			})
		}
	}()
}

func TxInsert(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxInsert(tx, key, value)
}

func TxUpdate(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxUpdate(tx, key, value)
}

func TxUpsert(tx *bbolt.Tx, key interface{}, value interface{}) error {
	return database.TxUpsert(tx, key, value)
}

func TxDelete(tx *bbolt.Tx, key interface{}, dataType interface{}) error {
	return database.TxDelete(tx, key, dataType)
}

func TxGet(tx *bbolt.Tx, key interface{}, result interface{}) error {
	return database.TxGet(tx, key, result)
}

func TxFind(tx *bbolt.Tx, result interface{}, query *bolthold.Query) error {
	return database.TxFind(tx, result, query)
}

func TxForEach(tx *bbolt.Tx, query *bolthold.Query, callback interface{}) error {
	return database.TxForEach(tx, query, callback)
}

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