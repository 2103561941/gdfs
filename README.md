# gdfs
golang分布式文件存储系统

## 功能特性
基础功能：
- 本地上传文件 get 
- 本地下载远程文件
- 删除远程文件
- 查看目录下的文件(元数据)
- 查看指定文件元数据
- 文件/目录重命名

高级功能：

- datanode 支持多机分布解决单点故障。
- 文件上传采用传递方式（client 上传到 datanode 一旦成功，剩下备份的部分又datanode 继续进行。）

## 快速开始

### 依赖检查

`go mod tidy ` 下载相关依赖

### 快速部署

使用 Makefile 脚本生成可执行文件。

``` shell
make cmd.gendocs # 生成命令行文档文件（markdown格式）

make gen.proto # 将 proto 文件转化为 go 文件。

make build.client    #生成 client
make build.namenode [-p=50051] # 生成namenode
make build.datanode [-p=5000]# 生成datanode
```
## 使用指南

#### 配置文件

1. 相关配置，请配置 config 目录下的 yaml 文件。
2. 目前 namenode 仅支持单机模式，在配置 client 和 datanode 配置文件时，需要保证 NamenodeAddr 与 namenode 的 ip 地址一致。
3. datanode 和 namenode 可以修改相关文件存放的路径，但是需要保证路径的文件夹已经存在，除了以 port 命名的文件夹外，其他文件夹不会自动生成，若文件夹不存在会导致启动失败。
4. datanode 和 namenode 可以使用配置文件的 port 地址， 同时可以使用 `-p=[port]` 来指定地址，且命令行指定优先级更高。

#### 本机存储

1. namenode 配置文件中的 stroagePath 目前用来存储目录树的日志文件。（在开放...）
2. datanode 的 storagePath 文件存放 datanode 存放的文件数据块。
3. 每个配置文件里都又日志文件的配置，其中的 LogPath 存放日志文件的存放路径。
#### 项目启动

1. 使用相关 makefile 指令后根据生成的文件依次开启 namenode，datanode 以及 client。
2. client 采用命令行方式进行交互，详细指令可以通过查看 docs/api 目录下的文档查看（若无，可以使用 makefile 自动生成），或使用 ./client.o 进行查看。

*其他细节请查看docs目录。*

## 如何贡献
该项目会不断迭代,欢迎大家贡献代码
## 关于作者
yeebing 2103561941@qq.com
## 许可证
gdfs is licensed under the MIT. See LICENSE for the full license text.