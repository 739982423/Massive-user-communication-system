package main

import (
	"finalProject/server/model"
	"fmt"
	"net"
	"time"
)

var listeningPort string = "8889"
var listeningIP string = "10.62.157.41"

func startGoRoutine(conn net.Conn) {
	defer conn.Close()
	fmt.Println("startGoRoutine: 得到了一个连接，开启协程，接下来用mainProcess处理该连接")
	var Processor1 Processor
	Processor1.Conn = conn
	err := Processor1.mainProcess()
	if err != nil {
		fmt.Println("startGoRoutine: Processor1.mainProcess() err, err =", err)
		return
	}
}

func initUserDao() {
	model.MyUserDao = *model.NewUserDao(pool)
}

func main() {
	//在服务器启动时初始化与redis的连接池
	initPool("127.0.0.1:6379", 16, 0, 300 * time.Second)
	//获得连接池之后初始化全局的UserDao
	initUserDao()

	//启动服务程序
	fmt.Println("main: 服务器已启动...")
	listener, err := net.Listen("tcp", listeningIP + ":" + listeningPort)
	if err != nil {
		fmt.Printf("main: net.Listen err, err = %v \n", err)
	} else {
		for {
			fmt.Println("main: 等待客户端连接...")
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("listener.Accept err, err = %v \n", err)
			} else {
				clientId := conn.RemoteAddr()
				fmt.Printf("main: 已建立与%v的连接 \n", clientId)
				go startGoRoutine(conn)
			}
		}
	}
}
