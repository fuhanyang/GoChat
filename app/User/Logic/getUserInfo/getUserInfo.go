package getUserInfo

import (
	"common/chain"
	"errors"
	"github.com/jinzhu/gorm"
	"user/Models"
)

type UserInfo struct {
	Name       string
	AccountNum string
	Email      string
	IP         string
	CreateAt   string
}

func GetUserInfo(accountNum string, db *gorm.DB) (UserInfo, error) {
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), getInfoFromMySQL())
	ctx.Set("accountNum", accountNum)
	ctx.Set("db", db)

	err := ctx.Apply()

	if err != nil {
		return UserInfo{}, err
	}
	_userInfo := ctx.Get("userInfo")
	if _userInfo == nil {
		return UserInfo{}, errors.New("user not found")
	}
	return _userInfo.(UserInfo), nil
}

func getInfoFromMySQL() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		user := Models.NewUser()
		//如果有函数结束要释放的资源，一般不要直接return nil,而是先ctx.Next()再释放资源，否则可能导致错误
		defer Models.ReleaseUser(user)
		userInfo := UserInfo{}
		_db := ctx.Get("db")
		if _db == nil {
			return errors.New("db is nil")
		}
		db, _ := _db.(*gorm.DB)
		_accountNum := ctx.Get("accountNum")
		if _accountNum == nil {
			return errors.New("accountNum is nil")
		}
		accountNum, _ := _accountNum.(string)

		Models.GetUserByAccountNum(db, user, accountNum)
		if user.ID == 0 {
			return errors.New("user not found")
		}
		userInfo.Name = user.Name
		userInfo.AccountNum = user.AccountNum
		userInfo.Email = "此功能暂未开放"
		userInfo.IP = user.IP
		userInfo.CreateAt = user.CreatedAt.Format("2006-01-02 15:04:05")
		ctx.Set("userInfo", userInfo)
		ctx.Next()
		return nil
	}
}
