package main

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

const urlBucket = "urlBucket"

var (
	//ErrNotInitialised returned by methods if
	//the URLStore has not been initialised
	ErrNotInitialised = errors.New("Not initialised")
)

/*
URLstore is
*/
type URLstore struct {
	db *bolt.DB
}

/*
Start a URLstore
*/
func (u *URLstore) Start(dbPath string) error {

	var err error
	//OpenDbase
	u.db, err = bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return fmt.Errorf("Open Dbase: %s", err)
	}
	//Create bucket (ok if it already exists)
	if err = u.db.Update(func(tx *bolt.Tx) error {
		_, _err := tx.CreateBucketIfNotExists([]byte(urlBucket))
		if _err != nil {
			return fmt.Errorf("create bucket: %s", _err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

/*
Stop a URLstore
*/
func (u *URLstore) Stop() error {

	if u.db != nil {
		u.db.Close()
		u.db = nil
	}
	return nil
}

/*
Store will store a url mapping in the Dbase
*/
func (u *URLstore) Store(shortName string, fullURL string) error {
	if u.db == nil {
		return ErrNotInitialised
	}
	if err := u.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(urlBucket))
		return b.Put([]byte(shortName), []byte(fullURL))
	}); err != nil {
		return err
	}
	return nil
}

/*
Retrieve the fullUrl for the specified short one
*/
func (u *URLstore) Retrieve(shortName string) (string, error) {
	if u.db == nil {
		return "", ErrNotInitialised
	}

	var value []byte
	if err := u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(urlBucket))
		url := b.Get([]byte(shortName))
		value = make([]byte, len(url))
		copy(value, url)
		return nil
	}); err != nil {
		return "", err
	}

	return string(value), nil
}
