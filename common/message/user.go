package message


//用户信息的结构体，为了与redis交互

type User struct { 
	UserId string 	 `json:"userId"`
	Passwd string    `json:"passwd"`
	UserName string  `json:"userName"`
	UserStatus int	 `json:"userStatus"`
}