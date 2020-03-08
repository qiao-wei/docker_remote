# docker_remote

[简体中文](https://github.com/jivi20029/docker_remote/blob/master/README-zh_CN.md) [English](https://github.com/jivi20029/docker_remote/blob/master/README.md)

## 📦 Why do this ?
My program need to docker build first ,and then push to the harbor, last deployment with k8s.

But the hard disk of my local machine is too small. so I don't want to install docker in local, 
so I need to send the program to the remote machine with docker,and then execute the script.
It is so complex. 

So I think that if all this can be done directly locally. Kubectl can operate remote machines directly locally,
But docker cant do that , so I want to implement one .  

###  📦 Why not use the docker api ?   
too complex !

## 📦 Build
### Server
```shell script
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server-liunx server/main.go 
```

### Client
```shell script
go build -o docker client/main.go 
```
Notice : In order to use the same as the docker, the client is named docker.
Of course, you can also name it another , but please change it to the name you named when you encounter the docker command later

## StartUp
### Server
Copy the "server" to the machine with docker
```shell script
nohup server-liunx -p 50000 &
```
### Client
#### \> Set the environment variable : DOCKER_REMOTE_SERVER , point to the IP and port of the server  

* linux mac osx
```shell script
export DOCKER_REMOTE_SERVER=ip:port
```
e.g. 
```shell script
export DOCKER_REMOTE_SERVER=192.168.1.16:50000
```
Write the above script to ~ /.bash_profile file

then 
```shell script
source ~/.bash_profile 
```
notice ：
Write "source ~/.bash_profile" to ~/.zshrc if it doesn't work after Mac OSX restarts

If there is no ~/.zshrc, please create a new one.


* windows

My Computer right button Properties -> Advanced settings -> Environment Variable

#### \> StartUp

Enter the directory where the "docker" is located. 
It is recommended to copy the docker to $path

type
```shell script
docker
```
View all commands .

In fact, it's the docker command, such as viewing all the image lists
```shell script
docker images 
```

If the prompt command not find in work directory, you can type
```shell script
./docker 
```
Or copy docker to $path
