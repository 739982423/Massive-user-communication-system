package main

import (
	"fmt"
	"finalProject/client/process"
)

func main() {
	var choice int
	for {
		fmt.Println("----------------------欢迎登陆多人聊天系统-----------------------")

		fmt.Println("\t\t\t1. 登录聊天室")
		fmt.Println("\t\t\t2. 注册用户")
		fmt.Println("\t\t\t3. 退出系统")
		fmt.Println("\t\t\t请输入你的选择：")
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			up := &process.UserProcess{}
			up.LoginFunc()
		case 2:
			up := &process.UserProcess{}
			up.RigisterFunc()
		case 3:
			fmt.Println("系统正在退出...")
			return
		default:
			fmt.Println("输入错误，请输入数字(1-3)")
		}

	}
}
