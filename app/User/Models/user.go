package Models

import (
	"User/DAO/Mysql"
	"User/DAO/Redis/redisLock"
	"User/Machine_code"
	"context"
	go_redislock "github.com/jefferyjob/go-redislock"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// userPool 用于缓存用户对象,不同的user类型都有对应的缓存池
var userPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

//用户创建后先写入redis再定期销毁，防止频繁上下线导致频繁读取数据库且频繁占用资源

// User 用户超类
type User struct {
	IP             string     `redis:"ip"`
	IsOnline       bool       `redis:"is_online"`
	Type           string     `redis:"type"`
	Name           string     `redis:"name"`
	AccountNum     string     `redis:"account_num"`
	Password       string     `redis:"password"`
	ID             uint       `gorm:"primary_key"`
	CreatedAt      time.Time  `redis:"-"`
	UpdatedAt      time.Time  `redis:"-"`
	DeletedAt      *time.Time `sql:"index" redis:"-"`
	redisLock      go_redislock.RedisLockInter
	lockCancelFunc context.CancelFunc
}

// MatchAUser 匹配一个用户
func MatchAUser() *User {
	// 首先获取记录总数
	var total int64
	if err := Mysql.MysqlDb.Model(&User{}).Count(&total).Error; err != nil {
		log.Printf("Failed to count users: %v", err)
		return nil
	}

	// 如果表中没有记录，直接返回 nil
	if total == 0 {
		return nil
	}

	// 生成一个有效的随机偏移量
	rand.Seed(time.Now().UnixNano())
	randomOffset := rand.Int63n(total)

	// 使用随机偏移量查询一条记录
	var user User
	if err := Mysql.MysqlDb.Offset(int(randomOffset)).Limit(1).Take(&user).Error; err != nil {
		log.Printf("Failed to query user: %v", err)
		return nil
	}

	return &user
}
func GetUserByAccountNum(user *User, accountNum string) {
	//从mysql中查找用户信息
	Mysql.MysqlDb.First(user, "account_num = ?", accountNum)
}
func WriteUser(user *User) {
	Mysql.MysqlDb.Create(user)
}
func NewUser() *User {
	var user = userPool.Get().(*User)
	user.ID = 0
	if user == nil {
		panic("new user failed")
	}
	return user
}
func (u *User) NewLock() {
	ctx := context.Background()
	ctx1, cancel := context.WithCancel(ctx)
	u.lockCancelFunc = cancel
	u.redisLock = redisLock.NewRedisLock(strconv.Itoa(Machine_code.Machine_code)+"_"+u.GetAccountNum(), ctx1)
}

// Repair 补全用户信息
func (u *User) Repair() {
	if u.GetPassword() == "" {
		u.SetPassword("123456")
	}
	if u.GetName() == "" {
		u.SetName("default_name")
	}
	if u.GetAccountNum() == "" {
		panic("account num is empty")
	}
	if u.GetType() == "" {
		u.SetType("User")
	}
	if u.GetIP() == "" {
		u.SetIP("not_set")
	}
	u.NewLock()
}
func (u *User) Release() {
	if u == nil {
		return
	}
	// 释放锁
	u.lockCancelFunc()
	userPool.Put(u)
}

// RedisLock 非阻塞锁
func (u *User) RedisLock() error {
	return u.redisLock.Lock()
}

// RedisBlockLock 阻塞锁
func (u *User) RedisBlockLock() error {
	return u.redisLock.SpinLock(10 * time.Second)
}
func (u *User) RedisUnlock() error {
	return u.redisLock.UnLock()
}

func (u *User) GetIP() string {
	return u.IP
}
func (u *User) GetIsOnline() bool {
	return u.IsOnline
}
func (u *User) GetName() string {
	return u.Name
}
func (u *User) GetAccountNum() string {
	return u.AccountNum
}
func (u *User) GetPassword() string {
	return u.Password
}
func (u *User) GetID() uint {
	return u.ID
}
func (u *User) GetType() string {
	return u.Type
}
func (u *User) SetIP(ip string) {
	u.IP = ip
}
func (u *User) SetAccountNum(accStr string) {
	u.AccountNum = accStr
}
func (u *User) SetIsOnline(isOnline bool) {
	u.IsOnline = isOnline
}
func (u *User) SetName(name string) {
	u.Name = name
}
func (u *User) SetPassword(password string) {
	u.Password = password
}
func (u *User) SetID(id uint) {
	u.ID = id
}
func (u *User) SetType(t string) {
	u.Type = t
}
