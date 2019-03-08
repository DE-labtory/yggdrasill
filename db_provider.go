package yggdrasill

import (
	"sync"

	"github.com/DE-labtory/leveldb-wrapper/key_value_db"
)

type DBHandle struct {
	dbName string
	db     key_value_db.KeyValueDB
}

type DBProvider struct {
	db        key_value_db.KeyValueDB
	mux       sync.Mutex
	dbHandles map[string]*DBHandle
}

func CreateNewDBProvider(keyValueDB key_value_db.KeyValueDB) *DBProvider {
	keyValueDB.Open()
	return &DBProvider{keyValueDB, sync.Mutex{}, make(map[string]*DBHandle)}
}

func (p *DBProvider) Close() {
	p.db.Close()
}

func (p *DBProvider) GetDBHandle(dbName string) *DBHandle {
	p.mux.Lock()
	defer p.mux.Unlock()

	dbHandle := p.dbHandles[dbName]
	if dbHandle == nil {
		dbHandle = &DBHandle{dbName, p.db}
		p.dbHandles[dbName] = dbHandle
	}

	return dbHandle
}

func (h *DBHandle) Get(key []byte) ([]byte, error) {
	return h.db.Get(dbKey(h.dbName, key))
}

func (h *DBHandle) Put(key []byte, value []byte, sync bool) error {
	return h.db.Put(dbKey(h.dbName, key), value, sync)
}

func (h *DBHandle) Delete(key []byte, sync bool) error {
	return h.db.Delete(dbKey(h.dbName, key), sync)
}

func (h *DBHandle) WriteBatch(KVs map[string][]byte, sync bool) error {
	return h.db.WriteBatch(KVs, sync)
}

func (h *DBHandle) GetIteratorWithPrefix() key_value_db.KeyValueDBIterator {
	return h.db.GetIteratorWithPrefix([]byte(h.dbName + "_"))
}

func (h *DBHandle) Snapshot() (map[string][]byte, error) {
	return h.db.Snapshot()
}

func dbKey(dbName string, key []byte) []byte {
	dbName = dbName + "_"
	return append([]byte(dbName), key...)
}
