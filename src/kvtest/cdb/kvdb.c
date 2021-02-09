#include <stdint.h>
#include <stddef.h>

#include "kvdb.h"

kvdb_t kvdb_open(char *name)
{
	return NULL;
}

int kvdb_close(kvdb_t db){ return 0; }
int kvdb_get(kvdb_t db, uint64_t k, uint64_t *v){ return 0; }
int kvdb_put(kvdb_t db, uint64_t k, uint64_t v){ return 0; }
int kvdb_del(kvdb_t db, uint64_t k){ return 0; }
int kvdb_next(kvdb_t db, uint64_t sk, uint64_t *k, uint64_t *v){ return 0; }


