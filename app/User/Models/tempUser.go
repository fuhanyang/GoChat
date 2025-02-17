package Models

import "sync"

type TempUser struct {
	IP         string `redis:"ip"`
	Type       string `redis:"type"`
	Name       string `redis:"name"`
	AccountNum string `redis:"account_num"`
	IsOnline   bool   `redis:"is_online"`
}

// tempUserPool 用于缓存用户对象,不同的user类型都有对应的缓存池
var tempUserPool = sync.Pool{
	New: func() interface{} {
		return &TempUser{}
	},
}

func NewTempUser() *TempUser {
	var user = tempUserPool.Get().(*TempUser)
	return user
}
func (u *TempUser) Release() {
	if u == nil {
		return
	}
	tempUserPool.Put(u)
}
func (u *TempUser) GetIP() string {
	return u.IP
}
func (u *TempUser) GetIsOnline() bool {
	return u.IsOnline
}
func (u *TempUser) GetName() string {
	return u.Name
}
func (u *TempUser) GetAccountNum() string {
	return u.AccountNum
}
func (u *TempUser) GetType() string {
	return u.Type
}
func (u *TempUser) SetIP(ip string) {
	u.IP = ip
}
func (u *TempUser) SetAccountNum(accStr string) {
	u.AccountNum = accStr
}
func (u *TempUser) SetIsOnline(isOnline bool) {
	u.IsOnline = isOnline
}
func (u *TempUser) SetName(name string) {
	u.Name = name
}
func (u *TempUser) SetType(t string) {
	u.Type = t
}
