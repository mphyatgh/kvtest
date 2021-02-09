package db

type DB interface {
	Put(key, value uint64) error
	Get(key uint64) (uint64, error)
	Del(key uint64) error
	Close()
	List(k1, k2 uint64, f func (uint64, uint64) bool) error
}

