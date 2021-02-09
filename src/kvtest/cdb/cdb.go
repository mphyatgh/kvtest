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
	_ "errors"
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

