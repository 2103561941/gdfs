# Put
将本地文件存储到分布式存储系统上.

client 获取本地文件路径和远程文件路径, 访问namenode 获取datanode相关节点信息.同时会生成文件的每个 chunck 分块的专属 uuid 用于保存在 datanode 中的文件名. (一个 datanode 只会保存一个文件备份这里可以对文件内容进行抽样扫描, 生成 md5 字符串码用于鉴别)

返回错误 namenode
- 文件已经存在在目录里面了(可能是路径冲突,但是这里不做覆盖)
- 目前没有 datanode 支持存储文件. (这个可以返回空的response)
---

client 根据本地文件路径打开文件, 与datanode进行远程传输. 将文件分块存储到对应的存储节点上. 

返回错误 datanode

---

节点在接受成功新的数据后进行文件上报, namenode 将新的文件存储到 cache 中.
