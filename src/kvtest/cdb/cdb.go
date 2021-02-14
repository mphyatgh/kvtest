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
	"kv/db"
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
func KvdbOpen(filename string) (db *Kvdb, err error) {
	db = new(Kvdb)
	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	db.fd = C.kvdb_open(cs)
	return db, nil
}

func (db *Kvdb)Close() {
	C.kvdb_close(db.fd)
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
	return 0, 0, nil
}

func (db *Kvdb)List(k1, k2 uint64, f func (uint64, uint64) bool) error {
	var (
		cont bool
		cur C.cursor_t
		k, v C.uint64_t
	)

	cont = true
	cur = C.kvdb_open_cursor(db.fd, C.uint64_t(k1), C.uint64_t(k2))
	for ;cont; {
		ret := C.kvdb_get_next(db.fd, cur, &k, &v)
		if ret<0 {
			break
		}
		cont = f(uint64(k), uint64(v))
	}
	C.kvdb_close_cursor(db.fd, cur)
	return nil
}

func cdbTest() {
	var d db.DB
	d, _ = KvdbOpen("test.kvdb")
	d.Close()
}

