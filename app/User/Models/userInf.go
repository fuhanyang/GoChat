package Models

type UserInf interface {
	GetIP() string
	GetIsOnline() bool
	GetAccountNum() string
	GetName() string
	GetType() string
	SetIP(string)
	SetAccountNum(string)
	SetIsOnline(bool)
	SetName(string)
	SetType(string)
	Release()
}
