# kvtest

这是一个测试kv数据库的工具。

## 命令说明

现在支持如下命令：

    root@03fa02a224fc:~/kvtest/src/kvtest# ./kv 
    kv help                   -- this message 
    kv get <key>              -- get a key
    kv put <key> <val>        -- set key
    kv del <key>              -- delete a key
    kv list                   -- list all key in the db
    
    kv ins <num>              -- insert records in batch mode
    kv clr                    -- remove all records in the database
    kv verify                 -- get all records and verify them
* kv get <key>  --  获取一个key的值
* kv put <key> <val> -- 如果key存在，则把key的值改为val。如果不存在，则创建可一个key，并设置值为val。
* kv del <key> -- 删除key/value对。
* kv list -- 列出数据库中所有kv对
* kv ins <num> -- 批量插入key/value对，其中num是key/value对的数量。
* kv clr -- 删除数据库中所有key/value对
* kvdb verify -- 验证用kv ins命令批量插入记录的值是否正确。

## 移植说明

默认情况下，kvtest支持gdbm数据库。但是可以方便的移植到其他kv数据库上。

* src/kvtest/db.go 是针对gdbm接口的实现，只需要更改这个文件，实现Get/Put/Del/List等几个函数，就可以测试其他数据库。
* src/kvtest/cdb/ 是一个用于测试C语言实现的数据库例子。

