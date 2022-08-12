# protobuf API 划分
目前需要实现的功能：

以Get为例
要求输入两个参数：local_file_path 和 remote_file_path,
实现目标是将 datanode 中的存储的文件内容拷贝到本地路径的文件中。
返回 result 表示操作结果。

1. 先向 client 输入本地文件路径和远程文件路径
2. client 通过 rpc 服务，将 remote_file_path 传输给 namenode，获取对应 datanode 的信息。
3. namenode 接受该 remote_file_path, 先再目录树中查询是否存在该文件，若存在，则获取到该文件对应的 key，查询 cache 获取到对应 datanode 的响应信息，返回给 client。
4. client 获取到 datanode 数据后，通过其访问该服务器。通过 local_file_path 打开本地文件，通过 grpc 的 stream 传输文件数据，保存到本地文件。 