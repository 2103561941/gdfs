// namenode 中维护目录树结构
// 目录树结构采用前缀树来保存在内存中，定期向磁盘写入目录树
// 目录树同时保存有一份操作日志文件，类似 redis 的日志文件，用于重启 namenode 时恢复部分未写入磁盘的文件信息。

package directorytree