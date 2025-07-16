package service

import (
	"context"
	"errors"
	"fmt"
	"rpc/user"
	"user/Const"
	"user/DAO/Mysql"
	"user/Logic/Search/SearchUser"
	"user/Logic/authentication"
	"user/Logic/forcedOffline"
	"user/Logic/getUserInfo"
	"user/Logic/matchUser"
)

func (s *server) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	return s.register(ctx, req)
}
func (s *server) register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	resp := &user.RegisterResponse{}
	resp.HandlerName = "Register"
	if req.GetPassword() != req.GetPasswordConfirm() {
		err := errors.New("Password and PasswordConfirm are not equal")
		resp.Msg = err.Error()
		resp.Code = 501
		return resp, err
	}
	acc, err := authentication.UserSignUp(req.GetUsername(), req.GetPassword(), req.GetIp(), Mysql.Db)
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 502
		return resp, err
	}
	resp.Code = Const.StatusOK
	resp.Msg = "Register success"
	resp.AccountNum = acc

	return resp, nil
}
func (s *server) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	resp := &user.LoginResponse{}
	resp.HandlerName = "Login"

	err := authentication.UserLogIn(req.GetAccountNum(), req.GetPassword(), req.GetIp(), Mysql.Db)
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 501
		return resp, err
	}
	resp.Code = Const.StatusOK
	resp.Msg = "Login success"
	return resp, nil
}
func (s *server) LogOff(ctx context.Context, req *user.LogOffRequest) (*user.LogOffResponse, error) {
	resp := &user.LogOffResponse{}
	resp.HandlerName = "LogOff"

	err := authentication.UserLogOff(req.GetAccountNum(), req.GetPassword(), req.GetIp(), Mysql.Db)
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 503
		return resp, err
	}
	resp.Code = Const.StatusOK
	resp.Msg = "LogOff success"
	return resp, nil
}
func (s *server) FindUser(ctx context.Context, req *user.FindUserRequest) (*user.FindUserResponse, error) {
	resp := &user.FindUserResponse{}
	resp.HandlerName = "FindUser"

	user, err := SearchUser.SearchUserByAccountNum(req.GetAccountNum(), Mysql.Db)
	if user == nil {
		err = errors.New("user is nil !!!!")
		goto ERR
	}
	if err != nil {
		goto ERR
	}
	resp.Username = user.Name
	resp.Code = Const.StatusOK
	resp.Msg = "FindUser success"
	return resp, nil
ERR:
	fmt.Println(err)
	resp.Msg = err.Error()
	resp.Code = 504
	return resp, err

}
func (s *server) MatchAUser(ctx context.Context, req *user.MatchAUserRequest) (*user.MatchAUserResponse, error) {
	resp := &user.MatchAUserResponse{}
	resp.HandlerName = "MatchAUser"
	resp.UserAccountNum = req.GetAccountNum()

	_user, err := MatchAUser.MatchAUser(req.GetAccountNum(), req.Choice, Mysql.Db)
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 505
		return resp, err
	}
	if _user == nil {
		err = errors.New("_user is nil !!!!")
		resp.Msg = err.Error()
		resp.Code = 505
		return resp, err
	}
	resp.Code = Const.StatusOK
	resp.Msg = "Match a _user success"
	resp.Username = _user.Name
	resp.UserAccountNum = _user.AccountNum

	return resp, nil
}
func (s *server) ForcedOffline(ctx context.Context, req *user.ForcedOfflineRequest) (*user.ForcedOfflineResponse, error) {
	err := forcedOffline.ForcedOffline(req.GetAccountNum())
	resp := &user.ForcedOfflineResponse{
		AccountNum: req.GetAccountNum(),
	}
	return resp, err
}
func (s *server) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	resp := &user.GetUserInfoResponse{}
	resp.HandlerName = "GetUserInfo"
	resp.AccountNum = req.GetAccountNum()
	infoData, err := getUserInfo.GetUserInfo(req.GetAccountNum(), Mysql.Db)
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 506
		return resp, err
	}
	resp.Code = Const.StatusOK
	resp.Msg = "Get user info success"
	resp.AccountNum = infoData.AccountNum
	resp.Username = infoData.Name
	resp.Email = infoData.Email
	resp.Ip = infoData.IP
	resp.CreateAt = infoData.CreateAt

	return resp, nil
}
