package db

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Pointer ...
type Pointer *r.Session

// ConnectDB ...
func ConnectDB(address string) (*r.Session, error) {
	db, err := r.Connect(r.ConnectOpts{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
