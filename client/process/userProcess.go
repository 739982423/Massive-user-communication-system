package process

import (
	"encoding/json"
	"finalProject/common/message"
	"finalProject/client/utils"
	"fmt"
	"net"
)
type UserProcess struct {
	//字段
}
var ipAddress string = "182.92.234.82"
var port string = "11334"
var des string = ipAddress + ":" + port

func (u *UserProcess) RigisterFunc() (err error) {
	//进入注册函数
	fmt.Println("欢迎您注册多人聊天系统！")

	var userId string
	var passwd string
	var userName string
	//获取用户想注册的账号
	fmt.Print("请输入想注册的账号:")
	fmt.Scanln(&userId)
	fmt.Print("请输入密码:")
	fmt.Scanln(&passwd)
	fmt.Print("请输入昵称:")
	fmt.Scanln(&userName)

	//1. 连接到服务器
	fmt.Println("RigisterFunc: 获取到注册请求，正在连接服务器...")
	// conn, err := net.Dial("tcp", "127.0.0.1:8889")
	conn, err := net.Dial("tcp", des)
	if err != nil {
		fmt.Printf("net.Dial err, err = %v \n", err)
		return
	}
	// fmt.Println("RigisterFunc: 连接成功，现在开始制作序列化后的注册请求字符串...")
	//延时关闭连接，记得及时写上
	defer conn.Close()

	var mes message.Message
	var registerMes message.RigisterMes
	var user message.User

	//2.1 对最内层的User结构体赋值
	user.UserId = userId
	user.Passwd = passwd
	user.UserName = userName

	//2.2 对registerMes的实例赋值
	registerMes.User = user

	//2.3 对外层的传输时的标准信息格式：mes，赋值
	userData, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("userProcess.RigisterFunc(): 序列化rigistermes出错, err =", err)
		return
	}
	mes.Type = message.RigisterMesType
	mes.Data = string(userData)

	//2.4 对外层的标准信息mes序列化
	rigisterData, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("userProcess.RigisterFunc(): 序列化mes出错, err =", err)
		return
	}

	// fmt.Println("userProcess.RigisterFunc(): 序列化后的注册请求字符串制作完成，将其交给writePkg函数...")
	//建立一个可操作读包写包的对象
	tf := &utils.Transfer{
		Conn : conn,
	}
	//2.5 向服务器发送序列化后的字符串（用[]byte形式发送）
	tf.ClientWritePkg([]byte(rigisterData))
	if err != nil {
		fmt.Println("userProcess.RigisterFunc(): 包发送失败 err =", err)
		return
	}

	//3. 等待服务器回复
	// fmt.Println("userProcess.RigisterFunc(): 调用readPkg函数处理服务器将返回的结果...")
	//ClientReadPkg将阻塞等待服务器回复
	mes, err = tf.ClientReadPkg()
	if err != nil {
		fmt.Println("userProcess.RigisterFunc(): 读取服务器返回的包出错, err =", err)
		return
	}
	
	//解析服务器返回的包
	// fmt.Println("userProcess.RigisterFunc(): 解析获取到的包...")
	var rigisterResMes message.RigisterResMes
	err = json.Unmarshal([]byte(mes.Data), &rigisterResMes)
	if err != nil {
		fmt.Println("userProcess.RigisterFunc(): 解析包出错, err =", err)
	}
	
	//看看服务器返回的包内部的Code是多少，以此判断是否注册成功
	switch rigisterResMes.Code {
	case 100:
		fmt.Println(rigisterResMes.Information)
	default:
		fmt.Println(rigisterResMes.Information, "请再次尝试")
	}
	
	return
}

func (u *UserProcess) LoginFunc() (err error) {

	//进入登录函数
	var userId string
	var passwd string
	//获取用户的输入
	fmt.Print("请输入账号:")
	fmt.Scanln(&userId)
	fmt.Print("请输入密码:")
	fmt.Scanln(&passwd)
	
	originUserId := userId
	//1. 连接到服务器
	fmt.Println("LoginFunc: 获取到登录请求，正在连接服务器...")

	conn, err := net.Dial("tcp", des)
	if err != nil {
		fmt.Printf("net.Dial err, err = %v \n", err)
		return
	}
	// fmt.Println("LoginFunc: 连接成功，现在开始制作序列化后的登录请求字符串...")
	//fmt.Println("当前连接是：", conn)
	//延时关闭连接，记得及时写上
	defer conn.Close()

	//2. 连接成功，首先向服务器发送message长度
	//2.1 首先建立发送消息的两种结构体
	var mes message.Message
	var loginmes message.LoginMes

	//2.2 给内层结构体赋值
	loginmes.UserId = userId
	loginmes.Passwd = passwd

	//2.3 对内层结构体序列化
	userdata, err := json.Marshal(loginmes)
	if err != nil {
		fmt.Println("userdata: json.Marshal err, err = ", err)
		return
	}

	//2.4 内层结构体组装完毕，将其赋值给外层结构体，并为外层结构体的种类赋值
	mes.Type = message.LoginMesType
	mes.Data = string(userdata)

	//2.5 对外层结构体序列化
	logindata, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("logindata: json.Marshal err, err = ", err)
		return
	}
	// fmt.Println("logindata:", string(logindata))
	//2.6 外层结构体组装完毕，将序列化后的mes传给 专门发送包的函数
	// fmt.Println("LoginFunc: 序列化后的登录请求字符串制作完成，将其交给writePkg函数...")
	tf := &utils.Transfer{
		Conn : conn,
	}
	err = tf.ClientWritePkg(logindata)
	if err != nil {
		fmt.Println("LoginFunc: writePkg(logindata) err, err =", err)
		return
	}

	//3 等待服务器的返回
	//这里readPkg返回的是message.Message类型的对象，我们要获取服务器返回的信息，需要对其data部分反序列化

	// fmt.Println("LoginFunc: 调用readPkg函数处理服务器将返回的结果...")
	mes, err = tf.ClientReadPkg()
	if err != nil {
		fmt.Println("LoginFunc: readPkg(conn) err, err =", err)
		return
	}

	// fmt.Println("LoginFunc: 解析获取到的包...")
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("LoginFunc: mes.Data unmarshal err, err =", err)
		return
	}

	if loginResMes.Code == 100 {
		// fmt.Println("LoginFunc:", loginResMes.Information)
		//登录成功，首先显示一下当前在线用户列表
		
		// 创建一个客户端维护的当前在线用户列表
		// var OnlineUserMap ClientOnlineUser
		// OnlineUserMap.OnlineUser = make(map[string]message.User)

		//将服务器传来的在线用户ID列表的值以User结构体的形式存入客户端维护的在线用户列表
		//即初始化客户端的在线用户列表：
		for _, val := range(loginResMes.OnlineUsers) {
			//这里放入客户端维护在线列表的是user结构体的复制，因为如果放指向下一行定义的user的指针，在当前循环结束后，user会被清理
			//指针就变成了野指针，再访问就会出错。所以在定义OnlineUserMap.OnlineUser这个map的时候，它的value必须是值而非指针形式
			//当然也可以改变这个用户在线列表的初始化方式，使得定义OnlineUserMap.OnlineUser时，map的value可以以指针的形式存在而不出错。
			userId = val.UserId
			OnlineUserMap.OnlineUser[userId] = val
		}
		// fmt.Println("OnlineUserMap.OnlineUser", OnlineUserMap.OnlineUser)
		//展示当前在线用户
		fmt.Println("当前在线的用户有：")

		var myself message.User
		for _, user := range OnlineUserMap.OnlineUser {	//遍历map时如果只填一个值，则对应的是遍历key
			if user.UserId == originUserId {
				myself = user
			}
			fmt.Println(user.UserName)
		}
		fmt.Println()
		// fmt.Println("---------------------------------")
		// fmt.Println(CurUser)
		//登录成功，初始化客户端存储当前登录用户信息的结构体
		CurUser.Conn = conn
		// fmt.Println("---------------------------------")
		// fmt.Println("CurUser.Conn = ", CurUser.Conn)
		CurUser.User = myself
		// fmt.Println("CurUser.User =", CurUser.User)
		// fmt.Println("---------------------------------")
		//开启协程保持与服务器的连接
		go KeepConnection(conn)
		ShowMenu()
	} else {
		fmt.Println("LoginFunc: Login failed, err code =", loginResMes.Code, "err information =", loginResMes.Information)
	}

	return
}
