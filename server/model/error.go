package model

import (
	"errors"
)

var (
	ERR_USER_EXISTS = errors.New("用户名已存在")
	ERR_USER_NOTEXISTS = errors.New("用户名不存在")
	ERR_WRONG_PWD = errors.New("密码错误")
)