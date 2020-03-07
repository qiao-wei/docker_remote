# docker_remote

[English](https://github.com/jivi20029/docker_remote/blob/master/README.md) [简体中文](https://github.com/jivi20029/docker_remote/blob/master/README-zh_CN.md) 

## 📦 BUILD
### 服务端
```shell script
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server-linux server/main.go 
```

### 客户端
```shell script
go build -o docker client/main.go 
```
注 : 为了使用起来和docker命令一样，所以将客户端命名为docker。
当然你也可以将他命名为其它名称，只是后面遇到docker命令时请改成你所命名的名称


## 运行
### 服务端
将server拷贝到装有docker的机器上
```shell script
nohup server-linux -p 50000 &
```
### 客户端
#### \> 设置环境变量 DOCKER_REMOTE_SERVER 指向服务端所在的IP和端口

* linux mac osx 设置如下 
```shell script
export DOCKER_REMOTE_SERVER=ip:port
```
例如
```shell script
export DOCKER_REMOTE_SERVER=192.168.1.16:50000
```
将以上脚本写入~/.bash_profile  
然后
```shell script
source ~/.bash_profile 
```
注：mac osx 重启后如果不生效 需要再往 ~/.zshrc  写入 

``
source ~/.bash_profile
``
如果没有~/.zshrc请新建 。  

* windows设置如下

我的电脑 右键 属性 -> 高级系统设置 -> 环境变量 


#### \> 启动

进入到docker所在目录 , 建议将docker拷贝到 $PATH 指向的目录 
```shell script
docker
```
查看所有命令 
其实就是docker命令 例如 查看所有镜象列表
```shell script
docker images 
```

如果在docker所在目录提示命令不存在 可使用 
```shell script
./docker 
```
或者将docker拷贝到 $PATH 指向的目录