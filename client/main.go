package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	PS = string(os.PathSeparator)
)

func checkResponse(server net.Conn) {
	buffer := make([]byte, 1024)
	c, err := server.Read(buffer)
	if err != nil {
		fmt.Println("服务器返回错误", err.Error())
		server.Close()
		os.Exit(0)
	}
	r := string(buffer[:c])
	if r != "ok" {
		fmt.Println("服务器返回错误", r)
		os.Exit(0)
	}
}

func sendFile(server net.Conn, rootPath, relativeFilePath string) {
	// 发送 file 标志位
	//fmt.Println("发送file标志位")
	server.Write([]byte("file"))
	checkResponse(server)

	filePath := rootPath + relativeFilePath
	fi, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("获取文件信息出错")
		server.Close()
	}

	// 发送文件名
	//fmt.Println("发送文件名", relativeFilePath)
	server.Write([]byte(relativeFilePath))
	checkResponse(server)

	// 发送文件长度
	fileLen := strconv.FormatInt(fi.Size(), 10)
	//fmt.Println("发送文件长度", fileLen)
	server.Write([]byte(fileLen))
	checkResponse(server)

	// 发送文件内容
	//fmt.Println("发送文件内容")
	file, err := os.Open(filePath)
	fBuffer := make([]byte, 1024)
	for {
		c, err := file.Read(fBuffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("读取文件失败", err)
			}
			file.Close()
			break
		}
		//fmt.Println(c)
		server.Write(fBuffer[:c])
	}
	checkResponse(server)
}

func sendDirectory(server net.Conn, rootPath, relativeFilePath string) {
	files, _ := ioutil.ReadDir(rootPath + relativeFilePath)

	for _, f := range files {
		newRelativeFilePath := relativeFilePath + PS + f.Name()
		if f.IsDir() {
			sendDirectory(server, rootPath, newRelativeFilePath)
		} else {
			sendFile(server, rootPath, newRelativeFilePath)
		}
	}
}

func sendDockerCMD(server net.Conn, fullCMD string) {
	// 发送cmd标志位
	//fmt.Println("发送cmd标志位")
	server.Write([]byte("cmd"))
	checkResponse(server)

	// 发送cmd内容
	//fmt.Println("发送cmd内容")
	server.Write([]byte(fullCMD))
	checkResponse(server)
}

func main() {
	srvAddress := os.Getenv("DOCKER_REMOTE_SERVER")
	if srvAddress == "" {
		fmt.Println("请先设置DOCKER_REMOTE_SERVER环境变量")
		fmt.Println("* linux mac osx ")
		fmt.Println("  vim ~/.bash_profile")
		fmt.Println("  增加 export DOCKER_REMOTE_SERVER=IP:PORT 例如 127.0.0.1:30000")
		fmt.Println("* windows ")
		fmt.Println("  右键 我的电脑")
		fmt.Println("  属性 -> 高级系统设置 -> 环境变量 增加相应的设置")
		return
	}
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	cmd := os.Args[1]
	dir := ""
	if cmd == "build" {
		if len(os.Args) < 5 {
			fmt.Println("参数有误")
			return
		}

		dir = os.Args[len(os.Args)-1]
	}
	fullCMD := strings.Join(os.Args[1:], " ")

	server, err := net.Dial("tcp", srvAddress)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sendDockerCMD(server, fullCMD)

	if cmd == "build" {
		//fmt.Println("发送目录")
		sendDirectory(server, dir, "")
		//sendFile(server, "test")
	}

	//fmt.Println("发送execute标志位")
	server.Write([]byte("execute"))

	// 显示服务端返回
	buffer := make([]byte, 1024)
	for {
		c, err := server.Read(buffer)
		if err != nil {
			if err != io.EOF {
				os.Stdout.Write([]byte(err.Error()))
			}
			//os.Stdout.Write([]byte("\n"))
			server.Close()
			return
		}
		os.Stdout.Write(buffer[:c])
	}
}
