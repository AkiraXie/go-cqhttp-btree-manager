# go-cqhttp-btree-manager

这是一个粗糙且简单的管理 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp) 内置的缓存数据库的小工具

`cache` 和 `btree` 目录均来自于 [go-cqhttp/internal](https://github.com/Mrs4s/go-cqhttp/tree/master/internal) 并做了微调
## 编译
1. [下载并安装 go](https://go.dev/dl/) (version >= 1.17)
2. 在当前目录使用 `go build`

## 使用
**以数据库 image.db 为例**
1. 插入缓存数据
 
``./gocq-http-manager insert -f xxx.image -o image.db``

2. 查看缓存数据

``./gocq-http-manager showimg -f xxx.image``

3. 从数据库中查询缓存数据, 查询的 `key` 为缓存的 `md5`

``./gocq-http-manager select -i {md5} -o image.db``

4. 查询数据库中所有缓存数据

``./gocq-http-manager showall -o image.db``

5. 从数据库导出指定缓存数据到指定文件, `-f` 参数为导出文件名，可省略, 省略则为 `{md5}.image`

``./gocq-http-manager export -f xxx.image -i {md5} -o image.db``

6. 从数据库导出所有数据到另一数据库

``./gocq-http-manager dump -s image.db -d xxx.db``