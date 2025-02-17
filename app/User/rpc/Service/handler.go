package Service

import (
	"User/Logic/MatchAUser"
	"User/Logic/Search/SearchUser"
	"User/Logic/UserLogIn"
	"User/Logic/UserLogOff"
	"User/Logic/UserSignUp"
	"User/StatusCode"
	"context"
	"errors"
	"fmt"
	"rpc/User"
)

func (s *server) Register(ctx context.Context, req *User.RegisterRequest) (*User.RegisterResponse, error) {
	return s.register(ctx, req)
}
func (s *server) register(ctx context.Context, req *User.RegisterRequest) (*User.RegisterResponse, error) {
	resp := &User.RegisterResponse{}
	resp.HandlerName = "Register"
	if req.GetPassword() != req.GetPasswordConfirm() {
		err := errors.New("Password and PasswordConfirm are not equal")
		resp.Msg = err.Error()
		resp.Code = 501
		return resp, nil
	}
	acc, err := UserSignUp.UserSignUp(req.GetUsername(), req.GetPassword(), req.GetIp())
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 502
		return resp, nil
	}
	resp.Code = StatusCode.StatusOK
	resp.Msg = "Register success"
	resp.AccountNum = acc

	return resp, nil
}
func (s *server) Login(ctx context.Context, req *User.LoginRequest) (*User.LoginResponse, error) {
	resp := &User.LoginResponse{}
	resp.HandlerName = "Login"

	err := UserLogIn.UserLogIn(req.GetAccountNum(), req.GetPassword(), req.GetIp())
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 501
		return resp, nil
	}
	resp.Code = StatusCode.StatusOK
	resp.Msg = "Login success"
	return resp, nil
}
func (s *server) LogOff(ctx context.Context, req *User.LogOffRequest) (*User.LogOffResponse, error) {
	resp := &User.LogOffResponse{}
	resp.HandlerName = "LogOff"

	err := UserLogOff.UserLogOff(req.GetAccountNum(), req.GetPassword(), req.GetIp())
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 503
		return resp, nil
	}
	resp.Code = StatusCode.StatusOK
	resp.Msg = "LogOff success"
	return resp, nil
}
func (s *server) FindUser(ctx context.Context, req *User.FindUserRequest) (*User.FindUserResponse, error) {
	resp := &User.FindUserResponse{}
	resp.HandlerName = "FindUser"

	user, err := SearchUser.SearchUserByAccountNum(req.GetAccountNum())
	if user == nil {
		err = errors.New("user is nil !!!!")
		goto ERR
	}
	if err != nil {
		goto ERR
	}
	resp.Username = user.Name
	resp.Code = StatusCode.StatusOK
	resp.Msg = "FindUser success"
	return resp, nil
ERR:
	fmt.Println(err)
	resp.Msg = err.Error()
	resp.Code = 504
	return resp, nil

}
func (s *server) MatchAUser(ctx context.Context, req *User.MatchAUserRequest) (*User.MatchAUserResponse, error) {
	resp := &User.MatchAUserResponse{}
	resp.HandlerName = "MatchAUser"
	resp.UserAccountNum = req.GetAccountNum()

	user, err := MatchAUser.MatchAUser(req.GetAccountNum())
	if err != nil {
		fmt.Println(err)
		resp.Msg = err.Error()
		resp.Code = 505
		return resp, nil
	}
	resp.Code = StatusCode.StatusOK
	resp.Msg = "Match a user success"
	resp.Username = user.Name
	resp.UserAccountNum = user.AccountNum

	return resp, nil
}
