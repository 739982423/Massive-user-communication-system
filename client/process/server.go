package process

import (
	"encoding/json"
	"finalProject/client/utils"
	"finalProject/common/message"
	"fmt"
	"net"
	"os"
)
var OnlineUserMap *ClientOnlineUser
type ClientOnlineUser struct {
	OnlineUser map[string]message.User
}
func init() {
	OnlineUserMap = &ClientOnlineUser{
		OnlineUser: make(map[string]message.User, 64),
	}
}

func ShowMenu() {
	for {
		fmt.Println("---------------------多人线上聊天系统-----------------------")
		fmt.Println("\t\t\t1. 显示在线用户列表")
		fmt.Println("\t\t\t2. 发送消息")
		fmt.Println("\t\t\t3. 信息列表")
		fmt.Println("\t\t\t4. 退出系统")
		fmt.Println("\t\t\t请输入选择：")
		var key int
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("\t\t\t1. 显示在线用户列表")
		case 2:
			fmt.Println("\t\t\t2. 发送消息")
		case 3:
			fmt.Println("\t\t\t3. 信息列表")
		default:
			fmt.Println("\t\t\t4. 退出系统")
			os.Exit(0)
		}
	}
}

func ClientProcessMes(mes *message.Message) (err error) {
	fmt.Println("ClientProcessMes: 客户端收到服务器发来的mes...")
	fmt.Println("ClientProcessMes: 正在验证mes类型...")
	switch mes.Type {
	case message.NotifyUserStatusMesType:
		fmt.Println("ClientProcessMes:该mes的类型是用户状态信息，开始处理...")
		err = ClientProcessStateChange(mes)
		if err != nil {
			fmt.Println("ClientProcessMes: 处理状态变化信息失败，err =", err)
		}
	default:
		fmt.Println("ClientProcessMes:未知类型，程序返回")
		// 	break
		// default:
		// 	break
	}
	return
}

func ClientProcessStateChange(mes *message.Message) (err error) {
	fmt.Println("ClientProcessStateChange:开始工作...")
	//因为已知mes的类型是NotifyUserStatusMesType，所以可直接创建一个NotifyUserStatusMes实例，对mes中的Data部分反序列化
	Data := mes.Data
	var stateChangeMes message.NotifyUserStatusMes
	err = json.Unmarshal([]byte(Data), &stateChangeMes)
	if err != nil {
		fmt.Println("ClientProcessStateChange:反序列化mes.Data失败，err =", err)
		return err
	}
	changeId := stateChangeMes.Status
	userId := stateChangeMes.UserId
	userName := stateChangeMes.UserName
	if changeId == 1 {
		fmt.Printf("%v上线了! \n", userName)
	} else if changeId == 0 {
		fmt.Printf("%v下线了! \n", userName)
	}

	//为了修改map的value值，需要先取出来原先的值，然后修改之后再放回去
	var user message.User = OnlineUserMap.OnlineUser[userId]
	user.UserStatus = changeId
	OnlineUserMap.OnlineUser[userId] = user 
	
	return
}

func KeepConnection(conn net.Conn) {
	tf := &utils.Transfer{
		Conn : conn,
	}

	for {
		fmt.Println("KeepConnection: the process is keeping connection with server... waiting for the message...")
		mes, err := tf.ClientReadPkg()
		if err != nil {
			fmt.Println("KeepConnection: Read mes from server err, err =", err)
			return
		}
		//读到了服务器的消息
		fmt.Println("KeepConnection: get messages from the server: message = ", mes)
		ClientProcessMes(&mes)
	}
}