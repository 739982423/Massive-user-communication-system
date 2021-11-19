package model

import (
	"encoding/json"
	"finalProject/common/message"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type UserDao struct {
	pool *redis.Pool
}

//我们在服务器启动时就初始化一个全局的UserDao实例，当任何时候需要与数据库交互时，就可以直接使用该实例
//下面先定义一个空的UserDao，需要在主函数中对该UserDao赋值来完成初始化
var MyUserDao UserDao

//使用工厂模式获得UserDao实例的函数接口
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao {
		pool : pool,
	}
	return
}

func (u *UserDao) Rigister(user *message.User) (err error) {
	//获取与Redis的连接
	conn := u.pool.Get()
	defer conn.Close()

	//在Redis查找Id是否已经存在
	res, err := redis.String(conn.Do("hget","users", user.UserId))
	if err != nil {
		if res == "" {
			//Redis中不存在该id，应该将这个user存入Redis
			userStr, err := json.Marshal(user)
			if err != nil {
				fmt.Println("userDao.Rigister():序列化User结构体出错, err =", err)
				return err
			}
			_, err = conn.Do("hset", "users", user.UserId, userStr)
			if err != nil {
				fmt.Println("userDao.Rigister():向Redis存入用户信息出错， err =", err)
				return err
			}
			return nil
		} else {
			fmt.Println("userDao.Rigister():向redis查询出错, err =", err)
			return err
		}
	} else {
		fmt.Println("userDao.Rigister():UserID查询结果不为空，用户名重复")
		err = ERR_USER_EXISTS
		return err
	}
	return err
}

//根据传来的用户id，检测数据库中是否存在该用户信息，如果存在，以User结构体的形式返回该用户信息，否则返回err
func (u *UserDao) getUserById(conn redis.Conn, Id string) (user *message.User, err error) {

	res, err := redis.String(conn.Do("hget", "users", Id))
	if err != nil {
		fmt.Println("conn.Do(hget) err, err =", err)
		if err == redis.ErrNil {
			err = ERR_USER_NOTEXISTS
		}
		return 
	}
	
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal err, err =", err)
		return
	}

	return
}

//完成登录校验
func (u *UserDao) Login(userId string, passwd string) (user *message.User, err error) { 
	conn := u.pool.Get()
	defer conn.Close()
	user, err = u.getUserById(conn, userId)
	if err != nil {
		if err == ERR_USER_NOTEXISTS {
			fmt.Println("UserDao.Login():u.getUserById err, err =", err)
			return
		}
	}
	//到这里说明已经获取到了对应UserId的user结构体
	if user.Passwd != passwd {
		err = ERR_WRONG_PWD
		return
	}
	return
}