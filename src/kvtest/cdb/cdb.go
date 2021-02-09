// Packge gdbm implements a wrapper around libgdbm, the GNU DataBase Manager
// library, for Go.
package cdb
/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -lkvdb -L.
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include "kvdb.h"

int test(void) { return 0; }
 
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Kvdb struct {
	fd	C.kvdb_t
}

/*
Simple function to open a database file with default parameters (block size
is default for the filesystem and file permissions are set to 0666).

mode is one of:
  "r" - reader
  "w" - writer
  "c" - rw / create
  "n" - new db
*/
func Open(filename string, mode string) (db *Kvdb, err error) {
	db = new(Kvdb)
	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	db.fd = C.kvdb_open(cs)
	return db, nil
}

func (db *Kvdb)Get(key uint64) (uint64, error) {
	var cv C.uint64_t
	ret := C.kvdb_get(db.fd, C.uint64_t(key), &cv)
	if ret!=0 {
		return 0, errors.New("not found")
	}
	return uint64(cv), nil
}

func (db *Kvdb)Put(key, value uint64) error {
	ret := C.kvdb_put(db.fd, C.uint64_t(key), C.uint64_t(value))
	if ret!=0 {
		return errors.New("error")
	}
	return nil
}

func (db *Kvdb)Del(key uint64) error {
	ret := C.kvdb_del(db.fd, C.uint64_t(key))
	if ret!=0 {
		return errors.New("error")
	}
	return nil
}

func (db *Kvdb)Next(sk uint64) (uint64, uint64, error) {
	var ck, cv C.uint64_t
	ret := C.kvdb_next(db.fd, C.uint64_t(sk), &ck, &cv)
	if ret!=0 {
		return 0, 0, errors.New("error")
	}
	return uint64(ck), uint64(cv), nil
}

