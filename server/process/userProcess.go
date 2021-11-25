package process

import (
	"encoding/json"
	"finalProject/common/message"
	"finalProject/server/model"
	"finalProject/server/utils"
	"fmt"
	"net"
)

type UserProcess struct {
	User *message.User
	Conn net.Conn
}


func (u *UserProcess) NotifyOtherStateChange(state int) (err error) {
	fmt.Println("userProcess.NotifyOtherStateChange(): state =", state)
	for _, v := range UserMgrPtr.OnlineUser {
		// if v.UserId == u.UserId {
		// 	continue
		// }
		conn := v.Conn
		var NotifyMes message.NotifyUserStatusMes
		NotifyMes.UserName = u.User.UserName
		NotifyMes.UserId = u.User.UserId
		NotifyMes.Status = state

		//序列化该notify结构体
		notifyMesData, err := json.Marshal(NotifyMes)
		if err != nil {
			fmt.Println("userProcess.NotifyOtherOlineUser(): Notify结构体序列化失败, err =", err)
			return err
		}
		
		//制作外层传输mes的标准形式
		var NotifyResMes message.Message
		NotifyResMes.Type = message.NotifyUserStatusMesType
		NotifyResMes.Data = string(notifyMesData)

		NotifyResMesData, err := json.Marshal(NotifyResMes)
		if err != nil {
			fmt.Println("userProcess.NotifyOtherOlineUser(): NotifyResMes结构体序列化失败, err =", err)
			return err
		}
		
		//交给向客户端写包的函数WritePkg
		tf := &utils.Transfer{
			Conn : conn,
		}
		err = tf.WritePkg(NotifyResMesData)
		if err != nil {
			fmt.Println("userProcess.NotifyOtherOlineUser(): notifyMes发送失败, err =", err)
		}
	}
	return
}

func (u *UserProcess) ServerProcessRigister(mes *message.Message) (err error) {
	//取出mes中的data部分，首先需要反序列化
	data := mes.Data
	var rigistermes message.RigisterMes
	err = json.Unmarshal([]byte(data), &rigistermes)
	if err != nil {
		fmt.Println("userProcess.serverProcessRigister():message.RigisterMes unmarshal err, err =", err)
		return
	}

	//得到了message.RigisterMes类型的结构体
	//取出其中的User结构体
	user := &rigistermes.User
	user.UserStatus = 0

	//去Redis里查询该UserID是否存在，如果err = nil，则说明不存在，可以注册
	//先声明之后要返回给客户端的消息结构体
	var rigitserResMes message.RigisterResMes

	err = model.MyUserDao.Rigister(user)
	if err != nil {
		fmt.Println("userProcess.serverProcessRigister(): model.MyUserDao.Rigister(user) err, err =", err)
		rigitserResMes.Code = 402
		rigitserResMes.Information = "用户名已存在/其他错误"
	} else {
		//到这里说明Redis中没有该UserID，则注册成功
		//开始制作返回给客户端的消息结构体
		rigitserResMes.Code = 100
		rigitserResMes.Information = "注册成功"
	}

	//对message.RigisterResMes类型的结构体进行序列化
	rigitserResMesData, err := json.Marshal(rigitserResMes)
	if err != nil {
		fmt.Println("userProcess.serverProcessRigister():rigitserResMes marshal err, err =", err)
		return
	}

	//制作最外层的message.Message结构体对象（客户端和服务器之间传输的标准包）
	var resMes message.Message
	resMes.Type = message.RigisterResMesType
	resMes.Data = string(rigitserResMesData)

	resMesData, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("userProcess.serverProcessRigister():resMes marshal err, err =", err)
		return
	}

	//得到了最终序列化后的message.Message结构体对象，将其交给包传送函数WritePkg
	//首先创建一个可操作WritePkg函数的操作对象（transfer类型）
	tf := &utils.Transfer{
		Conn : u.Conn,
	}

	err = tf.WritePkg([]byte(resMesData))
	if err != nil {
		fmt.Println("userProcess.serverProcessRigister():tf.WritePkg err, err =", err)
		return
	}
	return
}

func (u *UserProcess) ServerProcessLogin(mes *message.Message) (user *message.User, err error) {
	//取出mes中的data部分，首先需要反序列化
	data := mes.Data
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(data), &loginMes)
	if err != nil {
		fmt.Println("message.LoginMes unmarshal err, err =", err)
		return
	}

	//此时得到了LoginMes，即带有登录账号信息的结构体对象
	userId:= loginMes.UserId
	passwd := loginMes.Passwd
	// userName := loginMes.UserName
	fmt.Println("ServerProcessLogin: 获取到了用户名和密码，正在验证...")

	//判断账号密码是否正确匹配的逻辑……
	var loginResMes message.LoginResMes

	user, err = model.MyUserDao.Login(userId, passwd)
	if err != nil {
		fmt.Println("userProcess.ServerProcessLogin():model.MyUserDao.Login err, err =", err)
		switch err {
		case model.ERR_USER_NOTEXISTS:
			loginResMes.Code = 401
			loginResMes.Information = "用户名不存在!"
		case model.ERR_WRONG_PWD:
			loginResMes.Code = 404
			loginResMes.Information = "密码错误!"
		default:
			loginResMes.Code = 001
			loginResMes.Information = "未知错误!"
		}

	} else {
		loginResMes.Code = 100
		loginResMes.Information = "登录成功!"
		//登录成功，说明提供的账号密码正确，那么给这个UserProcess的实例绑定上该用户的userId和userName
		u.User.UserId = userId
		u.User.UserName = user.UserName

		//登录成功，首先通知其他人该用户上线了（向客户端发送一个message.NotifyUserStatusMes类型的包）
		u.NotifyOtherStateChange(1)

		//登录成功，我们将该用户放入服务器维护的“用户在线列表中”
		UserMgrPtr.AddOnlineUser(u)
		//为了让客户端收到服务器此时的“用户在线列表”，需要把这个列表放入传回客户端的mes中(传回用户名而不是用户id)
		for _, val := range UserMgrPtr.OnlineUser {
			var tmpUser message.User = *val.User
			tmpUser.Passwd = ""
			tmpUser.UserStatus = 1
			loginResMes.OnlineUsers = append(loginResMes.OnlineUsers, tmpUser)
		}

		fmt.Printf("userProcess.ServerProcessLogin(): %v 登录成功！\n", user.UserName)
	}
	
	//制作返回给客户端的消息结构体
	loginResMesData, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("loginResMes marshal err, err =", err)
		return
	}

	//序列化完成后，将该结果放入message.Message对象中（客户端和服务器之间传输的标准包），然后再序列化该resMes对象
	var resMes message.Message
	resMes.Type = message.LoginMesResType
	resMes.Data = string(loginResMesData)

	resMesData, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("resMes marshal err, err =", err)
		return
	}

	fmt.Println("serverProcessLogin: 序列化字符串制作完成，调用writePkg函数...")
	//将该序列化完成的byte切片交给统一的 包发送函数writePkg

	tf := &utils.Transfer{
		Conn : u.Conn,
	}

	err = tf.WritePkg(resMesData)
	if err != nil {
		fmt.Println("ResMes marshal err, err =", err)
		return
	}
	return
}