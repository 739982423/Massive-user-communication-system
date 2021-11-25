package process

import (
	"encoding/json"
	"finalProject/common/message"
	"finalProject/server/utils"
	"fmt"
	"net"
)

type SmsProcess struct {
	User message.User
	Conn net.Conn
}

func (s *SmsProcess) ServerReceiveGroupMes(mes *message.Message) (err error) {
	//取出mes中的数据部分
	data := mes.Data

	//反序列化Data
	var groupMes message.SmsMes
	err = json.Unmarshal([]byte(data), &groupMes)
	if err != nil {
		fmt.Println("smsProcess.ServerReceiveGroupMes(): 反序列化外层消息错误，err =", err)
	} 

	//得到发送来的消息的content和发送者的user结构体，需要把该结构体中的passwd部分隐藏掉，然后将该包转发给各在线用户
	groupMes.User.Passwd = ""
	// fmt.Println("——————————————————")
	// fmt.Println("groupMes:", groupMes)
	// fmt.Println("——————————————————")
	
	//序列化该包
	groupMesData, err := json.Marshal(groupMes)
	if err != nil {
		fmt.Println("smsProcess.ServerReceiveGroupMes(): 序列化内层消息错误，err =", err)
	} 
	
	//创建一个外层的message.Message来传输
	var serverSendMes message.Message
	serverSendMes.Type = message.SmsMesType
	serverSendMes.Data = string(groupMesData)

	serverSendMesData, err := json.Marshal(serverSendMes)
	if err != nil {
		fmt.Println("smsProcess.ServerReceiveGroupMes(): 序列化内层消息错误，err =", err)
	} 

	//创建一个发送对象
	tf := &utils.Transfer{}
	//找到需要发送的对象并将序列化好的结果交给发送函数
	for _, up := range UserMgrPtr.OnlineUser {
		tf.Conn = up.Conn
		tf.WritePkg(serverSendMesData)
		if err != nil {
			fmt.Println("smsProcess.ServerReceiveGroupMes(): 群发消息出现错误，up =", up, "err =", err)
		} 
	}
	return
}