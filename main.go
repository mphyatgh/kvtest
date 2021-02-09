package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
	_ "crypto/md5"
	"hash/crc32"
	"hash/crc64"
	"encoding/binary"
	"bytes"
	"time"
)

const (
	dbFile = "test.gdbm"
)

var (
	gCrc64tbl = crc64.MakeTable(crc64.ISO)
	db  *DB
)

func help() {
	fmt.Printf(`
    kvdb help                   -- this message 
    kvdb get <key>              -- get a key
    kvdb put <key> <val>        -- set key
    kvdb del <key>              -- delete a key
    kvdb list                   -- list all key in the db

    kvdb ins <num>              -- insert records in batch mode
    kvdb clr                    -- remove all records in the database
    kvdb verify                 -- get all records and verify them

    kvdb rget                   -- get record from the databse randomly
    kvdb tput                   -- put records in multi-threaded mod
    kvdb tget                   -- get records in multi-threaded mod
    kvdb tdel                   -- del records in multi-threaded mod
    kvdb tops                   -- operate records in multi-threaded mod
`)
}

func Marshal(v uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func Unmarshal(b []byte) (v uint64) {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}


func genKey(seq uint64) uint64 {
	b := Marshal(seq)
	h := crc64.New(gCrc64tbl)
	h.Write(b)
	return h.Sum64()
}

func genVal(seq uint64) uint64 {
	b := Marshal(seq)
	h := crc32.NewIEEE()
	h.Write(b)
	return uint64(h.Sum32())
}

func listFunc(key, value uint64) bool {
	fmt.Printf("List: key=%d, value=%d\n", key, value)
	return true
}

func dbIns(num int) {
	st := time.Now()
	for i := 0; i < num; i++ {
		seq := uint64(i)
		k := genKey(seq)
		v := genVal(seq)
		if err := db.Put(k, v); err != nil {
			log.Fatal(err)
		}
	}
	ed := time.Now()
	spt := ed.Sub(st)
	tpr := time.Duration(int(spt)/num)
	fmt.Printf("total time:      %s\n", spt)
	fmt.Printf("time per record: %s\n", tpr)
}

func delFunc(key, value uint64) bool {
	if err := db.Del(key); err != nil {
		log.Fatal(err)
	}
	return true
}

func s2i(s string) int {
	ret, err := strconv.ParseInt(s, 10, 64)
	if err!=nil {
		log.Fatal(err)
	}
	return int(ret)
}

func s2u(s string) uint64 {
	ret, err := strconv.ParseUint(s, 10, 64)
	if err!=nil {
		log.Fatal(err)
	}
	return ret
}

func lenEqs(args []string, num int) {
	if len(args)!=num {
		help()
		os.Exit(1)
	}
}

func main() {
	var (
		cmd string
		key, value uint64
		err error
	)
	if len(os.Args)<2 {
		help()
		return
	}
	if db, err = dbOpen(dbFile, "c"); err!=nil {
		log.Fatal(err)
	}
	defer db.Close()
	cmd = os.Args[1]
	switch cmd {
	case "help":
		help()
	case "del":
		lenEqs(os.Args, 3)
		key = s2u(os.Args[2])
		if err = db.Del(key); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Success\n")
	case "get":
		lenEqs(os.Args, 3)
		key = s2u(os.Args[2])
		if value, err = db.Get(key); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Get: key=%d, value=%d\n", key, value)
	case "put":
		lenEqs(os.Args, 4)
		if len(os.Args)!=4 {
			help()
			return
		}
		key = s2u(os.Args[2])
		value = s2u(os.Args[3])
		if err = db.Put(key, value); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Put: key=%d, value=%d\n", key, value)
	case "list":
		db.List(0, MaxUint, listFunc)
	case "ins":
		lenEqs(os.Args, 3)
		num := s2i(os.Args[2])
		dbIns(num)
	case "clr":
		lenEqs(os.Args, 2)
		db.List(0, MaxUint, delFunc)
	default:
		help()
		return
	}
}



