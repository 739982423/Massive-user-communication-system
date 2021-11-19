package main

import (
	"finalProject/common/message"
	"finalProject/server/process"
	"finalProject/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
	userProcessPtr *process.UserProcess
}

func (p *Processor) ServerProcessMes(mes *message.Message) (err error) {
	fmt.Println("ServerProcessMes: mes处理函数开始处理传来的mes结构体...")
	switch mes.Type {
	case message.LoginMesType:
		fmt.Println("ServerProcessMes: 该mes的类型是登录信息，开始处理...")
		//首先建立一个up(userProcess)的实例，将conn赋值给up，准备用up调用与用户消息处理相关的函数
		up := &process.UserProcess{
			Conn: p.Conn,
			User: &message.User{},
		}
		//这个分支是处理登录信息的，所以用up调用处理登录信息的的函数ServerProcessLogin，如果登录成功，会返回一个message.User类型的结构体实例
		//用该结构体实例存储该用户的所有信息
		user, err := up.ServerProcessLogin(mes)
		if err != nil {
			fmt.Println("ServerProcessMes: 处理登录mes错误，err =", err)
			return err
		}
		//如果没有出错，表示登录成功，将存有该用户信息的user放入up中，这样就将该conn与用户的信息放在了一起管理
		up.User = user

		//因为每个客户端的访问都对应着一个连接，而每个连接都会用一个协程处理，所以一个协程对应一个用户
		//现在用p(processor)表示总的与用户接头的实例，该实例可以调用解析mes的函数等等
		//因为每个p对应一个用户，所以将用户处理的接口(up)也放在processor结构体里面，方便统一管理
		p.userProcessPtr = up

	case message.RigisterMesType:
		fmt.Println("ServerProcessMes: 该mes的类型是注册信息，开始处理...")
		tmpUP := &process.UserProcess{
			Conn: p.Conn,
		}
		err = tmpUP.ServerProcessRigister(mes)
		if err != nil {
			fmt.Println("ServerProcessMes: 处理注册mes错误，err =", err)
		}
	default:
		fmt.Println("ServerProcessMes: 未知类型，程序返回")
		// 	break
		// default:
		// 	break
	}
	return
}

func (p *Processor) mainProcess() (err error) {
	for {
		fmt.Println("mainProcess: mainProcess开始了新一轮循环！")
		fmt.Println("mainProcess: 正在等待读包...")
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err != io.EOF {
				up := p.userProcessPtr
				up.NotifyOtherStateChange(0)
			}
			fmt.Println("mainProcess: readPkg(conn) err, err =", err)
			return err
		}
		fmt.Println("mainProcess: 获得了ReadPkg返回的消息结构体，下面将其交给消息处理函数ServerProcessMes...")
		p.ServerProcessMes(&mes)
	}
}
