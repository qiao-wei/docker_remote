package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type DockerRemoteContext struct {
	ID   string
	CMD  string
	Args []string
}

func receiveDockerCMD(ctx context.Context, client net.Conn) error {
	buffer := make([]byte, 2048)
	c, err := client.Read(buffer)
	if err != nil {
		return err
	}

	fullCMD := string(buffer[:c])
	fullCMD = strings.Trim(fullCMD, "\n")
	fields := strings.Fields(fullCMD)

	drCtx := ctx.Value("context").(*DockerRemoteContext)
	drCtx.CMD = fields[0]
	if drCtx.CMD == "build" {
		fields[len(fields)-1] = "/tmp/" + drCtx.ID
	}
	drCtx.Args = fields
	client.Write([]byte("ok"))
	return nil
}

func receiveFile(ctx context.Context, client net.Conn) error {
	drCtx := ctx.Value("context").(*DockerRemoteContext)

	rootPath := "/tmp/" + drCtx.ID
	os.Mkdir(rootPath, os.ModePerm)
	buffer := make([]byte, 1024)

	// 读取相对文件路径
	c, err := client.Read(buffer)
	if err != nil {
		return err
	}
	relativeFilePath := string(buffer[:c])
	// 替换掉windows系统传过来的\路径字符为/
	relativeFilePath = strings.ReplaceAll(relativeFilePath, "\\", "/")
	relativePath := path.Dir(relativeFilePath)
	//fmt.Println(pathDir)
	os.MkdirAll(rootPath+relativePath, os.ModePerm)
	client.Write([]byte("ok"))

	// 读取长度
	c, err = client.Read(buffer)
	if err != nil {
		return err
	}
	fileLen, _ := strconv.Atoi(string(buffer[:c]))
	client.Write([]byte("ok"))

	// 读取文件内容并保存
	total := 0
	file, err := os.Create(rootPath + "/" + relativeFilePath)
	if err != nil {
		return err
	}
	for {
		c, err = client.Read(buffer)
		if err != nil {
			return err
		}

		file.Write(buffer[:c])
		total += c
		if total >= fileLen {
			file.Close()
			break
		}

	}
	client.Write([]byte("ok"))
	return nil
}

func executeDockerCMD(ctx context.Context, client net.Conn) {
	// 开始处理
	drCtx := ctx.Value("context").(*DockerRemoteContext)
	fmt.Println(client.RemoteAddr().String(), ":", "docker", drCtx.Args)
	cmd := exec.Command("docker", drCtx.Args...)

	cmd.Stdout = client
	cmd.Stderr = client

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		// client.Write([]byte(err.Error()))
		log.Println("执行docker命令出现错误:", err.Error())
	}

	// 关闭连接
	client.Close()
}

func handleClient(client net.Conn) {

	guid := uuid.New()
	id := guid.String()

	drCtx := &DockerRemoteContext{ID: id}
	ctx := context.WithValue(context.Background(), "context", drCtx)
	buffer := make([]byte, 1024)

	for {
		c, err := client.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println(fmt.Sprintf("client:%v 退出", client.RemoteAddr().String()))
			} else {
				log.Println(fmt.Sprintf("client:%v 读取数据时出现错误", client.RemoteAddr().String()))
			}
			client.Close()
			return
		}

		content := string(buffer[:c])
		if content == "cmd" {
			client.Write([]byte("ok"))
			err = receiveDockerCMD(ctx, client)
			if err != nil {
				client.Write([]byte(err.Error()))
				client.Close()
				return
			}
		} else if content == "file" {
			client.Write([]byte("ok"))
			err = receiveFile(ctx, client)
			if err != nil {
				client.Write([]byte(err.Error()))
				client.Close()
				return
			}
		} else {
			executeDockerCMD(ctx, client)
			// 删除文件夹
			os.RemoveAll("/tmp/" + id)
			return
		}
	}
}

func main() {
	var port uint
	flag.UintVar(&port, "p", 30000, "端口")
	flag.Parse()

	server, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	fmt.Println(fmt.Sprintf("开始监听 port:%v", port))
	if err != nil {
		panic(err.Error())
		return
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatal(err.Error())
			continue
		}
		go handleClient(client)
	}
}
