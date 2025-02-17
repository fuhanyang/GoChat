package TypeDefine

import "sync"

const (
	Ch         = "ch"
	Name       = "name"
	AccountNum = "account_num"
	Password   = "password"
	IP         = "ip"
	ID         = "id"
	Opt        = "opt"
)

// UserMap 后续可以使用redis代替
var UserMap map[int64]UserObj
var Mu sync.Mutex

type UserOptChan chan map[interface{}]interface{}
type UserViewChan chan ViewData

type ViewData struct {
	Type interface{}
	Data interface{}
}

func init() {
	UserMap = make(map[int64]UserObj)
}

type UserOpt struct {
	Ch UserOptChan
	ID uint
}

type UserViewData struct {
	Ch UserViewChan
	ID uint
}

// UserObj 用户交互数据的实例
type UserObj struct {
	Opt      UserOpt      //用户操作的通道
	ViewData UserViewData //与前端交互的数据
}
