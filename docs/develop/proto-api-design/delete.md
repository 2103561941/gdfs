# Delete
删除 gdfs 中指定目录的文件。
需要删除目录树上的文件夹，需要删除datanode节点中的指定缓存，需要删除cache中的filekey.


## 个人思考
需要删除目录树中的文件，通过该文件找到file chunk（其中还可能有文件夹里面的文件，需要递归删除），在cache中使用filechunk 找到对应的 datanode，对其进行删除操作。 如果还有 datanode 对应 filechunk 的结构，也要里面的filechunk进行删除。
考虑的点，如果删除失败了，怎么处理，我的理解就是目录树删了和map删了其实就可以，这样client 访问不到，剩下的残党就算了，留着也行。
datanode 重新连接上namenode时会进行文件上报，当遇到自己存有但namenode中不存在的文件，就将本地的文件删除。
## 思考题
只要目录树里面的删了就行，这样对用户来说是黑盒的，datanode里面是否存有这个 file chunk 问题不大。


## 扩展
hdfs 有个回收站功能.