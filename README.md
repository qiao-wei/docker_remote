# docker_remote

## 📦 BUILD
### Server
```shell script
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server server/main.go 
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
nohup server -p 50000 &
```
### Client
#### \> Set the environment variable : DOCKER_REMOTE_SERVER , point to the IP and port where the server is located

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
