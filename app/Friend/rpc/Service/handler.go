package Service

import (
	"Friend/Logic"
	"Friend/Models"
	"Friend/rpc/client"
	"context"
	"errors"
	"fmt"
	"rpc/Friend"
	"rpc/User"
	"time"
)

func (s *server) AddFriend(ctx context.Context, req *Friend.AddFriendRequest) (*Friend.AddFriendResponse, error) {
	var (
		resp       Friend.AddFriendResponse
		err        error
		match_resp *User.MatchAUserResponse
		timer      *time.Timer
	)
	resp.HandlerName = "AddFriend"
	resp.AccountNum = req.GetAccountNum()
	user, err := client.UserServiceClient.Client.FindUser(ctx, &User.FindUserRequest{AccountNum: req.GetAccountNum(), HandlerName: "FindUser"})
	if err != nil {
		goto ERR
	}
	if user == nil {
		err = errors.New("user not found")
		goto ERR
	}
	timer = time.NewTimer(time.Second * 2)
	for {
		select {
		case <-timer.C:
			err = errors.New("not match a user in expected time")
			goto ERR
		default:
			match_resp, err = client.UserServiceClient.Client.MatchAUser(ctx, &User.MatchAUserRequest{
				AccountNum:  req.GetAccountNum(),
				HandlerName: "MatchAUser",
			})
			if err != nil {
				goto ERR
			}
			err = Logic.AddFriend(req.GetAccountNum(), match_resp.GetUserAccountNum(), user.GetUsername(), match_resp.GetUsername())
			if err != nil {
				fmt.Println(err)
				if errors.Is(err, Models.ErrDuplicateFriend) {
					time.Sleep(time.Millisecond * 500)
					continue
				}
				goto ERR
			}
			goto Correction
		}

	}
Correction:
	resp.Friend = &Friend.Friend{
		AccountNum: match_resp.GetUserAccountNum(),
		Name:       match_resp.GetUsername(),
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
func (s *server) GetFriends(ctx context.Context, req *Friend.GetFriendsRequest) (*Friend.GetFriendsResponse, error) {
	var (
		resp Friend.GetFriendsResponse
	)
	resp.Friends = make([]*Friend.Friend, 0)
	resp.HandlerName = "GetFriends"
	resp.AccountNum = req.GetAccountNum()

	friends := Logic.GetFriends(req.GetAccountNum())
	for _, friend := range friends {
		resp.Friends = append(resp.Friends, friend)
	}
	resp.Msg = "success"
	resp.Code = 200

	return &resp, nil
}
func (s *server) DeleteFriend(ctx context.Context, req *Friend.DeleteFriendRequest) (*Friend.DeleteFriendResponse, error) {

	return &Friend.DeleteFriendResponse{}, nil
}
