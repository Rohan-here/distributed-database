package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

// Database is a open bolt databse
type Database struct {
	db *bolt.DB
}

// NewDatabase returns an instance of Database
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {

	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	db = &Database{db: boltDb}
	closeFunc = boltDb.Close

	if err := db.createBucket(); err != nil {
		closeFunc()
		return nil, closeFunc, fmt.Errorf("error creating default bucket: %v", err)
	}

	return db, closeFunc, nil
}

// Create default bucket
func (d *Database) createBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// SetKey sets the key to requested value or returns an error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// GetKey returns the value of requested key or returns an error
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, nil
	}

	return nil, err

}
