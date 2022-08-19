# TODO LIST


2. namenode 若在put时没有找到合适的datanode，应该将其从目录树中删除。 ok

4. 其他 crud 操作。 ok
    mkdir ok
    rename ok
    stat ok
    list ok
 
    delete: ok 


10. grpc 中间件，实现错误自动恢复（namenode和datanode节点），同时进行日志打印。 ok

11. 使用 metadata 代替结构体里的address ok

12. 给文件加上上次访问时间（修改时间） ok

1. namenode 目录树备份，

5. datanode 负载均衡

7. 对受影响的副本，做副本数补全，即要确认副本数不能少。
可以在节点挂掉后做副本数补全操作。（alive进行监控，需要保存每个节点存了哪些文件，到指定的文件夹去删除，同时需要完成对应的重新 put 操作。）
datanode 挂了之后，需要在namenode中进行日志打印，并把 cache 中 datanode 相关的节点删了，以免下次客户端去调用

11. config 业务代码之间解耦，即把相应要用到的config 写到结构体里面，而不是直接在业务代码里使用，直接在业务代码里使用的话不利于做单元测试和mock。

13. 美化client的返回信息，而不是统一返回log信息。

---

9. namenode 一致性算法，解决单点故障问题。 raft 或 主从复制

8. 实现http监控指标大盘 普罗米修斯

测试：
1. 轮询是否达到随机 ok

2. 大文件分块存储  ok 
这里使用日志文件，大小是1m，按照 100k进行划分， 

3. 多个 datanode 测试能否选三个出来备份  ok

4. 文件完整性测试 使用 diff 来查看文件的复制情况。 ok