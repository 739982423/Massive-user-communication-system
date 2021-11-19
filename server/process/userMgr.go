package process

import "errors"

var UserMgrPtr *UserMgr

type UserMgr struct {
	OnlineUser map[string]*UserProcess
}

func init() {
	UserMgrPtr = &UserMgr{
		OnlineUser: make(map[string]*UserProcess, 1024),
	}
}

func (u *UserMgr) AddOnlineUser(up *UserProcess) {
	u.OnlineUser[up.User.UserId] = up
}

func (u *UserMgr) DelOnlineUser(up *UserProcess) {
	delete(u.OnlineUser, up.User.UserId)
}

func (u *UserMgr) GetAllOnlineUsers() map[string]*UserProcess {
	return u.OnlineUser
}

func (u *UserMgr) GetOneUser(userId string) (up *UserProcess, err error) {
	res, ok := u.OnlineUser[userId]
	if ok {
		return res, nil
	} else {
		err = errors.New("该用户不在线")
		return nil, err
	}
}