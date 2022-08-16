# File Report
文件上报服务

datanode 成功写下文件后，会发发送自己的 address 以及 filekey，表明文件写入成功，同时 namenode 保存到内存中，client在下次获取文件时就可以获取到该 datanode 的 ip 地址。

同时，支持 datanode 关闭重启后，文件上报的服务。
这里我的设计是把同一主机的不同 datanode 的文件按照端口生成文件夹存储，这样下次按照该端口开启 datanode 即可把文件夹里存在的文件进行上传。 