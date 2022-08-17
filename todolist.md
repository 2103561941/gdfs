# TODO LIST


2. namenode 若在put时没有找到合适的datanode，应该将其从目录树中删除。 ok

4. 其他 crud 操作。 ok
    mkdir ok
    rename ok
    stat ok
    list ok
 
    delete: ok 
    需要删除目录树中的文件，通过该文件找到file chunk（其中还可能有文件夹里面的文件，需要递归删除），在cache中使用filechunk 找到对应的 datanode，对其进行删除操作。 如果还有 datanode 对应 filechunk 的结构，也要里面的filechunk进行删除。
    考虑的点，如果删除失败了，怎么处理，我的理解就是目录树删了和map删了其实就可以，这样client 访问不到，剩下的残党就算了，留着也行。
    如果残党太多，可以手动做一个上报，datanode 跟 namenode 里面的cache 比较文件是否存在，不存在就删了。
    
    思考题，只要目录树里面的删了就行，这样对用户来说是黑盒的，datanode里面是否存有这个 file chunk 问题不大。
    
    hdfs 有个回收站功能，这里

5. datanode 负载均衡

10. grpc 中间件，实现错误自动恢复（namenode和datanode节点），同时进行日志打印。 ok

11. 使用 metadata 代替结构体里的address ok

12. 给文件加上上次访问时间（修改时间） ok

---

1. namenode 目录树备份， 

7. 对受影响的副本，做副本数补全，即要确认副本数不能少。
可以在节点挂掉后做副本数补全操作。（alive进行监控，需要保存每个节点存了哪些文件，到指定的文件夹去删除，同时需要完成对应的重新 put 操作。）
datanode 挂了之后，需要在namenode中进行日志打印，并把 cache 中 datanode 相关的节点删了，以免下次客户端去调用

8. 实现http监控指标大盘 普罗米修斯

9. namenode 一致性算法，解决单点故障问题。 raft 或 主从复制 
    

测试：
1. 轮询是否达到随机 ok

2. 大文件分块存储  ok 
这里使用日志文件，大小是1m，按照 100k进行划分， 

3. 多个 datanode 测试能否选三个出来备份  ok

4. 文件完整性测试 使用 diff 来查看文件的复制情况。 ok