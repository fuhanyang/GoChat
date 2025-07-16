package main

import (
	"testing"
	"websocket/Logic"
)

func TestBroadCast(t *testing.T) {
	err = Logic.BroadcastMessage([]byte("it is a test message"), false, false)
	if err != nil {
		t.Error(err)
	}
}
