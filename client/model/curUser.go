package model

import (
	"finalProject/common/message"
	"net"
)

//CurUser保存当前客户端登录的用户的信息和连接
//因为在客户端很多地方都会用到CurUser的实例，所以可以直接初始化为一个全局的实例
type CurUser struct {
	User message.User
	Conn net.Conn
}