package process

import (
	"encoding/json"
	"finalProject/client/model"
	"finalProject/client/utils"
	"finalProject/common/message"
	"fmt"
	"net"
	_"os"
)

var OnlineUserMap *ClientOnlineUser
var CurUser model.CurUser //声明一个存储当前客户端登录的用户信息的结构体，我们将在用户登录成功后初始化它
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
			ShowOnlineUsers()
		case 2:
			fmt.Println("\t\t\t2. 发送消息")
			var content string
			fmt.Println("请输入：")
			fmt.Scanln(&content)
			sp := SmsProcess{}
			sp.ClientSendGroupMes(content)
		case 3:
			fmt.Println("\t\t\t3. 信息列表")
			fmt.Println("---------------------------------")
			fmt.Println(CurUser)
			fmt.Println("---------------------------------")
		default:
			fmt.Println("\t\t\t4. 退出系统")
			// os.Exit(0)
		}
	}
}

func ShowOnlineUsers() (err error) {
	onlineList := make([]string, 0)
	for _, v := range OnlineUserMap.OnlineUser {
		if v.UserStatus == 1 {
			onlineList = append(onlineList, v.UserName)
		}
	}
	if len(onlineList) > 1 {
		fmt.Println("当前在线用户有：")
		for i := 0; i < len(onlineList); i++ {
			fmt.Println(onlineList[i])
		}
	} else {
		fmt.Println("当前无其他在线用户！")
	}
	return
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

	case message.SmsMesType:
		fmt.Println("ClientProcessMes:该mes的类型是群发消息，开始处理...")
		sp := SmsProcess{}
		err = sp.ClientReceiveGroupMes(mes)
		if err != nil {
			fmt.Println("ClientProcessMes: 解析接收到的群发信息失败，err =", err)
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
	v, ok := OnlineUserMap.OnlineUser[userId]
	if ok {
		fmt.Println("map中有user，修改状态前", v)
		v.UserStatus = changeId
		OnlineUserMap.OnlineUser[userId] = v
		fmt.Println("map中有user，修改状态后", v)
	} else {
		var user message.User = message.User{
			UserId:     userId,
			UserName:   userName,
			UserStatus: changeId,
		}
		OnlineUserMap.OnlineUser[userId] = user
		fmt.Println("map中无user，添加user", user)
	}
	fmt.Println("此时的在线列表:", OnlineUserMap.OnlineUser)

	return
}

func KeepConnection(conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
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
