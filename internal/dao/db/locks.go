// +build server

package db

import (
	"sort"
	"sync"
)

type TableName string
type unlockFunc func()

type byName []TableName
func (b byName) Len() int {return len(b)}
func (b byName) Less(i, j int) bool {return b[i]<b[j]}
func (b byName) Swap(i, j int) {b[i],b[j] = b[j],b[i]}


var locklock sync.Mutex
var namedLocks=make(map[TableName]*sync.RWMutex)

func getMutex(tableName TableName) *sync.RWMutex {
	var mutexPtr *sync.RWMutex
	locklock.Lock()
	var ok bool
	if mutexPtr, ok =namedLocks[tableName]; !ok {
		mutexPtr=&sync.RWMutex{}
		namedLocks[tableName]=mutexPtr
	}
	locklock.Unlock()
	return mutexPtr
}

func acquireSingleWrite(tableName TableName) unlockFunc {
	mutexPtr := getMutex(tableName)
	mutexPtr.Lock()
	return func() {
		mutexPtr.Unlock()
	}
}

func acquireWrite(tableNames... TableName) unlockFunc {
	sort.Sort(byName(tableNames))
	var unlockers []unlockFunc
	for _, tableName := range tableNames {
		unlockers=append(unlockers, acquireSingleWrite(tableName))
	}
	return func() {
		for i:=len(unlockers)-1;i>=0;i-- {
			unlockers[i]()
		}
	}
}

func acquireSingleRead(tableName TableName) unlockFunc {
	mutexPtr := getMutex(tableName)
	mutexPtr.RLock()
	return func() {
		mutexPtr.RUnlock()
	}
}

func acquireRead(tableNames... TableName) unlockFunc {
	sort.Sort(byName(tableNames))
	var unlockers []unlockFunc
	for _, tableName := range tableNames {
		unlockers=append(unlockers, acquireSingleRead(tableName))
	}
	return func() {
		for i:=len(unlockers)-1;i>=0;i-- {
			unlockers[i]()
		}
	}
}
