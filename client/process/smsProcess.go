package process

import (
	"encoding/json"
	"finalProject/client/utils"
	"finalProject/common/message"
	"fmt"
)

type SmsProcess struct {
}

func (s *SmsProcess) ClientSendGroupMes(content string) (err error) {
	var sendMes message.SmsMes
	sendMes.Content = content
	sendMes.User= CurUser.User
	// sendMes.User.UserId = CurUser.UserName
	// fmt.Println("*************")
	// fmt.Println("发送给服务器的sendMes：", sendMes)
	// fmt.Println("*************")
	sendMesData, err := json.Marshal(sendMes)
	if err != nil {
		fmt.Println("smsProcess.SendGroupMes(): 序列化内层消息失败，err =", err)
	}

	var smsMes message.Message
	smsMes.Type = message.SmsMesType
	smsMes.Data = string(sendMesData)

	
	smsMesData, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("smsProcess.SendGroupMes(): 序列化外层消息失败，err =", err)
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	tf.ClientWritePkg(smsMesData)
	if err != nil {
		fmt.Println("smsProcess.SendGroupMes(): 包发送失败 err =", err)
		return
	}
	return
}

func (s *SmsProcess) ClientReceiveGroupMes(mes *message.Message) (err error) {
	//取出mes中的数据部分
	data := mes.Data

	//反序列化Data
	var groupMes message.SmsMes
	err = json.Unmarshal([]byte(data), &groupMes)
	if err != nil {
		fmt.Println("smsProcess.ClientReceiveGroupMes(): 反序列化外层消息错误，err =", err)
	} 
	// senderID := groupMes.User.UserId
	senderName := groupMes.User.UserName
	content := groupMes.Content

	fmt.Printf("%v：%v \n" , senderName, content)
	return 
}