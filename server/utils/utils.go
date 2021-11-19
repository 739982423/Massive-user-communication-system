package utils

import (
	"encoding/binary"
	"encoding/json"
	"finalProject/common/message"
	"fmt"
	"net"
)

// 创建一个关于传输的结构体，以后只要遇到传输相关的事，新建一个该结构体的对象，调用其方法即可
type Transfer struct {
	//思考：对于传输，应该需要哪些字段？
	Conn net.Conn
	Buff [2048]byte
}

func (t *Transfer) WritePkg(sendBytes []byte) (err error) {

	fmt.Println("WritePkg: 获取到了新的writePkg请求，下面开始制作包...")
	//要想发送一个包，首先获得其长度，发送该长度变量
	var dataLen uint32
	lenBuff := t.Buff[:4]
	dataLen = uint32(len(sendBytes))
	binary.BigEndian.PutUint32(lenBuff[0:4], dataLen)

	//得到了长度变量lenBuff，并且以[]byte的形式表示了，下面开始发送
	n, err := t.Conn.Write(lenBuff[:])
	if n != 4 || err != nil {
		fmt.Println("WritePkg: conn.Write([len(Buff)]) err, err =", err)
		return
	}
	fmt.Printf("WritePkg: 包的长度是:%v \n", dataLen)

	_, err = t.Conn.Write(sendBytes)
	if err != nil {
		fmt.Println("WritePkg: conn.Write(sendBytes) err, err =", err)
		return
	}
	fmt.Println("WritePkg: 包发送完毕")
	// fmt.Printf("writePkg: 包发送成功，包的内容是:%v \n", string(sendBytes))
	return
}

func (t *Transfer) ReadPkg() (mes message.Message, err error) {
	//读包使用的函数，首先建立一个读取字符的缓存buff
	// buff := make([]byte, 1024)
	//阻塞等待客户端发送包
	fmt.Println("readPkg: 正在等待客户端发送消息...")

	//首先接收到的包是长度包
	//conn.Read在conn未关闭时会阻塞，一直等待客户端发送消息
	//但目前客户端的逻辑是只发送一次，发送完就关闭conn
	//关闭之后，这里的Read就不会阻塞了，会读到EOF，然后报错err = EOF
	_, err = t.Conn.Read(t.Buff[:4])
	if err != nil {
		fmt.Println("readPkg: conn.Read(buff[:4]) err, err =", err)
		return
	}

	//解析该长度是多少
	var pkgLen = binary.BigEndian.Uint32(t.Buff[:4])
	fmt.Printf("readPkg: 读取到客户端发来的长度，%v \n", pkgLen)

	//之后收到的包是消息包，将获取到的字节放入buff缓存的前pkgLen个byte
	_, err = t.Conn.Read(t.Buff[:pkgLen])
	if err != nil {
		fmt.Println("conn.Read(buff[:pkgLen]) err, err =", err)
		return
	}

	//将消息包反序列化，得到message.Message类型的对象mes
	//该变量有两个内置变量，分别为类型和一个序列化后的message.LoginMes对象
	err = json.Unmarshal(t.Buff[:pkgLen], &mes)
	if err != nil {
		fmt.Println("readPkg: buff[:pkgLen] unmarshal err, err =", err)
		return
	}

	fmt.Println("readPkg: 读取到客户端发来的mes结构体，返回mes给上层函数...")
	return
}
