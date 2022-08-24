# 功能测试

## 功能1
1.Put 文件到远程目录
./client.o put ./storage/local/hello /hello

2. 使用stat查看文件的相关信息, 比对信息是否正确。
./client.o stat /hello

3. Get 文件并验证数据
./client.o get ./storage/local/hello.cp /hello
diff hello hello.cp

4. delete 文件，并使用stat 查看文件是否真的删除了
./client.o delete /hello
./client.o stat /hello


## 功能2
1. 创建文件夹
查看文件夹是否已经存在
./client.o list /
创建文件夹
./client.o mkdir /dir

2. put 文件到文件夹中
./client.o put ./storage/local/hello /dir/hello
./client.o stat /dir/hello

3. 将刚刚put的文件进行重命名
./client.o rename /dir/hello /dir/hello.txt
./client.o list /dir

4. 删除目录
./client.o delete /dir
./client.o list /

## 展示负载均衡，多次上传文件
## 展示读取文件时datanode挂掉后选择新的datanode

## 展示重启 namenode 时目录树的读取及重启datanode时文件上报 
./client.o put ./storage/local/big  /big
./client.o stat /big
重启namenode和datanode
查看文件是否存在
./client.o list /
./client.o get ./storage/local/big.cp /big
diff big.cp big

