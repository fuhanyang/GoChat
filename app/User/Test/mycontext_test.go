package Test

import (
	"User/GoRoutinePool"
	"User/MyHandlerFunc"
	"fmt"
	"testing"
	"time"
)

func F1() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		//fmt.Println("F1")
		time.Sleep(5 * time.Second)
	}
}
func F2() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		c.Next()
		//fmt.Println("F2")
	}
}
func F3() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		//fmt.Println("F3")
	}
}
func F4() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		//fmt.Println("F4")
	}
}
func F5() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		c.Next()
		//fmt.Println("F5")
	}
}
func F6() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		c.Abort()
		//fmt.Println("F6")
	}
}
func F7() MyHandlerFunc.MyHandlerFunc {
	return func(c *MyHandlerFunc.MyContext) {
		//fmt.Println("F7")
	}
}
func Test1(t *testing.T) {
	c := MyHandlerFunc.LoadHandlers(F1(), F2(), F3(), F4(), F5(), F6(), F7())
	al := GoRoutinePool.GetAllLeaders()
	var t1, t2 time.Time
	go func() {
		t1 = time.Now()
		for i := 0; i < 1000000; i++ {
			_c := *c
			err := al.TaskToLeader(MyHandlerFunc.CreateTask(&_c))
			if err != nil {
				fmt.Println(err)
			}
		}
		t2 = time.Now()

	}()
	for j := 0; j < 10; j++ {
		time.Sleep(3 * time.Second)
		for i := 0; i < 10; i++ {
			fmt.Println("leader ", i%GoRoutinePool.MaxLeaderNum, " running num:", al.GetLeader(i%GoRoutinePool.MaxLeaderNum).GetLeaderRunning())
		}
	}
	fmt.Println(t2.Sub(t1))
}
