package service

import (
	"context"
	"errors"
	"fmt"
	"friend/DAO/Mysql"
	"friend/Logic"
	"friend/Models"
	"friend/rpc/client"
	"rpc/friend"
	"rpc/user"
	"time"
)

func (s *server) AddFriend(ctx context.Context, req *friend.AddFriendRequest) (*friend.AddFriendResponse, error) {
	var (
		resp       friend.AddFriendResponse
		err        error
		matchResp  *user.MatchAUserResponse
		timer      *time.Timer
		timeLimit  = 2 * time.Second
		oneWait    = time.Millisecond * 500
		retry      int
		choiceType = false
	)
	resp.HandlerName = "AddFriend"
	resp.AccountNum = req.GetAccountNum()
	// 先查找发起请求的用户存不存在
	_user, err := client.UserServiceClient.Client.FindUser(ctx, &user.FindUserRequest{AccountNum: req.GetAccountNum(), HandlerName: "FindUser"})
	if err != nil {
		goto ERR
	}
	if _user == nil {
		err = errors.New("_user not found")
		goto ERR
	}
	timer = time.NewTimer(timeLimit)
	defer timer.Stop()
	for {
		retry++
		if retry >= int(timeLimit/oneWait)-1 {
			// 采用mysql匹配
			choiceType = true
		}
		select {
		case <-timer.C:
			err = errors.New("not match a _user in expected time")
			goto ERR
		default:
			matchResp, err = client.UserServiceClient.Client.MatchAUser(ctx, &user.MatchAUserRequest{
				AccountNum:  req.GetAccountNum(),
				HandlerName: "MatchAUser",
				Choice:      choiceType,
			})
			if err != nil {
				goto ERR
			}
			err = Logic.AddFriend(Mysql.MysqlDb, req.GetAccountNum(), matchResp.GetUserAccountNum(), _user.GetUsername(), matchResp.GetUsername())
			if err != nil {
				fmt.Println(err)
				if errors.Is(err, Models.ErrDuplicateFriend) {
					time.Sleep(oneWait)
					continue
				}
				goto ERR
			}
			goto Correction
		}

	}
Correction:
	resp.Friend = &friend.Friend{
		AccountNum: matchResp.GetUserAccountNum(),
		Name:       matchResp.GetUsername(),
	}
	fmt.Println("friend:", resp.Friend.String())
	resp.Msg = "success"
	resp.Code = 200
	return &resp, nil
ERR:
	fmt.Println(err)
	resp.Msg = err.Error()
	resp.Code = 500
	return &resp, nil
}
func (s *server) GetFriends(ctx context.Context, req *friend.GetFriendsRequest) (*friend.GetFriendsResponse, error) {
	var (
		resp friend.GetFriendsResponse
	)
	resp.Friends = make([]*friend.Friend, 0)
	resp.HandlerName = "GetFriends"
	resp.AccountNum = req.GetAccountNum()

	friends := Logic.GetFriends(Mysql.MysqlDb, req.GetAccountNum())
	for _, _friend := range friends {
		resp.Friends = append(resp.Friends, _friend)
	}
	resp.Msg = "success"
	resp.Code = 200

	return &resp, nil
}
func (s *server) DeleteFriend(ctx context.Context, req *friend.DeleteFriendRequest) (*friend.DeleteFriendResponse, error) {

	return &friend.DeleteFriendResponse{}, nil
}
func (s *server) CheckFriend(ctx context.Context, req *friend.CheckFriendRequest) (*friend.CheckFriendResponse, error) {
	return &friend.CheckFriendResponse{}, nil
}
func (s *server) AddFriendWithAccountNum(ctx context.Context, req *friend.AddFriendWithAccountNumRequest) (*friend.AddFriendResponse, error) {
	var (
		resp    friend.AddFriendResponse
		err     error
		_user   *user.FindUserResponse
		_friend *user.FindUserResponse
	)
	resp.HandlerName = "AddFriendWithAccountNum"
	resp.AccountNum = req.GetAccountNum()
	// 先查找发起请求的用户存不存在
	_user, err = client.UserServiceClient.Client.FindUser(ctx, &user.FindUserRequest{AccountNum: req.GetAccountNum(), HandlerName: "FindUser"})
	if err != nil {
		goto ERR
	}
	if _user == nil {
		err = errors.New("_user not found")
		goto ERR
	}
	// 再查找被添加的用户存不存在
	_friend, err = client.UserServiceClient.Client.FindUser(ctx, &user.FindUserRequest{AccountNum: req.GetTargetAccountNum(), HandlerName: "FindUser"})
	if err != nil {
		goto ERR
	}
	// 进行添加好友
	err = Logic.AddFriend(Mysql.MysqlDb, req.GetAccountNum(), req.GetTargetAccountNum(), _user.GetUsername(), _friend.GetUsername())
	if err != nil {
		goto ERR
	}
	resp.Friend = &friend.Friend{
		AccountNum: req.GetTargetAccountNum(),
		Name:       _friend.GetUsername(),
	}
	fmt.Println("friend:", resp.Friend.String())
	resp.Msg = "success"
	resp.Code = 200
	return &resp, nil
ERR:
	fmt.Println(err)
	resp.Msg = err.Error()
	resp.Code = 500
	return &resp, nil
}
