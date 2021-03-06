package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"github.com/cfdrake/go-gdbm"
)

const MaxUint = ^uint64(0)

type Gdbm struct {
	gdbm *gdbm.Database
	name string
	mutex sync.Mutex
}

func gdbmOpen(fname, mode string) (* Gdbm, error) {
	r := new(Gdbm)
	db, err := gdbm.Open(fname, mode)
	r.gdbm = db
	r.name = fname
	return r, err
}

func i2s(v uint64) string {
	return fmt.Sprintf("%016x", v)
}

func (d *Gdbm)Put(key, value uint64) error {
	sk := i2s(key)
	sv := i2s(value)
	d.mutex.Lock()
	err := d.gdbm.Replace(sk, sv)
	d.mutex.Unlock()
	return err
}

func (d *Gdbm)Get(key uint64) (uint64, error) {
	val := uint64(0)
	sk := i2s(key)
	d.mutex.Lock()
	sv, err := d.gdbm.Fetch(sk)
	d.mutex.Unlock()
	if err==nil {
		val, err = strconv.ParseUint(sv, 16, 64)
	}
	return val, err
}


func (d *Gdbm)Del(key uint64) error {
	sk := i2s(key)
	d.mutex.Lock()
	err := d.gdbm.Delete(sk)
	d.mutex.Unlock()
	return err
}

func (d *Gdbm)Close() {
	d.gdbm.Close()
}

func h2u(h string) uint64 {
	v, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func (d *Gdbm)List(k1, k2 uint64, f func (uint64, uint64) bool) error {
	var (
		sk, sv string
		k, v uint64
		err error
	)
	if k1==0 {
		if sk, err = d.gdbm.FirstKey(); err != nil {
			return err
		}
	} else {
		sk = i2s(k)
	}
	for {
		d.mutex.Lock()
		sv, err = d.gdbm.Fetch(sk)
		d.mutex.Unlock()
		if err != nil {
			return err
		}
		v = h2u(sv)
		k = h2u(sk)
		if k>=k2 {
			break
		}

		sk, err = d.gdbm.NextKey(sk)
		cont := f(k, v)
		if !cont {
			break
		}

		if err==gdbm.NoError {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

