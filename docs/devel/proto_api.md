# protobuf API 划分
目前需要实现的功能：

## datanode 需要的接口
### 1. Get 

描述: 将文件数据拉到本地
- 输入 
    local_file_path 写入的本地路径
    remote_file_path 
- 输出 


## NameNode 接口配置
基本的接口都需要访问namenode来获取到datanode对应的节点数据，通过节点数据访问datanode来获取最终的数据或修改。

List、Stat和Rename 三个接口的返回可以直接通过Namenode内存中维护的文件目录树来返回，而部分数据的修改可以通过异步去datanode中修改
