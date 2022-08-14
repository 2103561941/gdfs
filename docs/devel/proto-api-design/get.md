# GET 
获取远程文件保存到本地目录

client 获取本地文件路径和远程文件路径, 访问namenode 获取datanode相关节点信息.

返回错误 namenode
- 文件不存在
- 文件是目录文件
- 文件丢失(文件存在目录树里,但是未找到相关的存储节点,或者文件的部分 chunk 丢失)
---

client 根据本地文件路径打开文件, 与datanode进行远程传输.

返回错误 datanode
