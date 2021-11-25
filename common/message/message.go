package message

const (
	LoginMesType    		= 	"LoginMes"
	LoginMesResType			= 	"LoginResMes"
	RigisterMesType			= 	"RigisterMes"
	RigisterResMesType 		= 	"RigisterResMes"
	NotifyUserStatusMesType = 	"NotifyUserStatusMes"
	SmsMesType 				= 	"SmsMes"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	UserId   string `json:"userId"`
	Passwd   string `json:"passwd"`
	// UserName string `json:"userName"`
}

type LoginResMes struct {
	Code        int    `json:"code"`
	Information string `json:"information"`
	//增加一个服务器返回的保存当前在线用户的切片
	OnlineUsers []User `json:"onlineUsers"`
}

type RigisterMes struct {
	User User
}

type RigisterResMes struct {
	Code        int    `json:"code"`
	Information string `json:"information"`
}

//用户状态变化通知所使用的mes类型
type NotifyUserStatusMes struct {
	UserName string		`json:"userName"`
	UserId string 		`json:"userId"`
	Status int			`json:"status"`
}

type SmsMes struct {
	User User			`json:"user"`
	Content string		`json:"content"`
}