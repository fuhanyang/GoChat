package Test

import (
	"common/chain"
	"common/routinePool"
	"errors"
	"log"
	"testing"
	"time"
)

// 更新后的Handler函数，现在返回error
func F1() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		//log.Println("F1: 开始执行")
		time.Sleep(1 * time.Second)
		//log.Println("F1: 执行完成")
		return nil
	}
}

func F2() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		log.Println("F2: 开始执行")
		c.Next()
		log.Println("F2: 执行完成")
		return nil
	}
}

func F3() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		log.Println("F3: 开始执行")
		// 模拟可能失败的情况
		if c.KV["simulateF3Fail"] == true {
			return errors.New("F3执行失败")
		}
		log.Println("F3: 执行完成")
		return nil
	}
}

func F4() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		log.Println("F4: 开始执行")
		// 模拟可能失败的情况
		if c.KV["simulateF4Fail"] == true {
			return errors.New("F4执行失败")
		}
		log.Println("F4: 执行完成")
		return nil
	}
}

func F5() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		log.Println("F5: 开始执行")
		c.Next()
		log.Println("F5: 执行完成")
		return nil
	}
}

func F6() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		//log.Println("F6: 开始执行")
		// 模拟熔断场景
		if c.KV["simulateF6Fail"] == true {
			return errors.New("F6执行失败，触发熔断")
		}
		c.Abort()
		//log.Println("F6: 执行完成并中断链")
		return nil
	}
}

func F7() chain.MyHandlerFunc {
	return func(c *chain.MyContext) error {
		//log.Println("F7: 开始执行")
		//log.Println("F7: 执行完成")
		return nil
	}
}

// 测试基本功能
func TestBasicChain(t *testing.T) {
	log.Println("=== 测试基本链执行 ===")

	c := chain.LoadHandlers(nil, chain.DefaultTimer(), F1(), F2(), F3(), F4(), F5(), F6(), F7())

	// 设置错误处理
	c.SetErrorHandler(func(err error, ctx *chain.MyContext) {
		log.Printf("错误处理被调用: %v", err)
	})

	// 设置自定义日志
	c.SetLogger(func(level string, format string, args ...interface{}) {
		log.Printf("[Test] "+format, args...)
	})

	// 准备测试数据
	c.KV["simulateF3Fail"] = false
	c.KV["simulateF4Fail"] = false
	c.KV["simulateF6Fail"] = false

	err := c.Apply()
	if err != nil {
		t.Logf("链执行失败: %v", err)
	} else {
		t.Logf("链执行成功")
	}
}

// 测试熔断功能
func TestCircuitBreaker(t *testing.T) {
	log.Println("=== 测试熔断功能 ===")

	scenarios := []struct {
		name string
		kv   map[string]interface{}
	}{
		{
			name: "F3失败场景",
			kv: map[string]interface{}{
				"simulateF3Fail": true,
				"simulateF4Fail": false,
				"simulateF6Fail": false,
			},
		},
		{
			name: "F4失败场景",
			kv: map[string]interface{}{
				"simulateF3Fail": false,
				"simulateF4Fail": true,
				"simulateF6Fail": false,
			},
		},
		{
			name: "F6失败场景",
			kv: map[string]interface{}{
				"simulateF3Fail": false,
				"simulateF4Fail": false,
				"simulateF6Fail": true,
			},
		},
	}

	for _, scenario := range scenarios {
		t.Logf("\n--- 测试场景: %s ---", scenario.name)

		c := chain.LoadHandlers(nil, chain.DefaultTimer(), F1(), F2(), F3(), F4(), F5(), F6(), F7())
		c.SetErrorHandler(func(err error, ctx *chain.MyContext) {
			log.Printf("熔断触发: %v", err)
		})

		// 设置测试数据
		for k, v := range scenario.kv {
			c.KV[k] = v
		}

		err := c.Apply()
		if err != nil {
			t.Logf("执行失败: %v", err)
		} else {
			t.Logf("链执行成功")
		}
	}
}

// 测试带名称的Handler
func TestNamedHandlers(t *testing.T) {
	log.Println("=== 测试带名称的Handler ===")

	handlers := map[string]chain.MyHandlerFunc{
		"参数校验": func(c *chain.MyContext) error {
			if c.KV["accountNum"] == nil || c.KV["accountNum"].(string) == "" {
				return errors.New("accountNum is empty")
			}
			accountNum := c.KV["accountNum"].(string)
			log.Println("参数校验通过 accountNum: ", accountNum)
			return nil
		},
		"用户加锁": func(c *chain.MyContext) error {
			log.Println("用户加锁成功")
			return nil
		},
		"查询用户": func(c *chain.MyContext) error {
			if c.KV["simulateQueryFail"] == true {
				return errors.New("用户不存在")
			}
			log.Println("查询用户成功")
			return nil
		},
		"更新状态": func(c *chain.MyContext) error {
			log.Println("更新状态成功")
			return nil
		},
	}

	c := chain.LoadHandlersWithNames(nil, chain.DefaultTimer(), handlers)
	c.SetErrorHandler(func(err error, ctx *chain.MyContext) {
		log.Printf("错误处理: %v", err)
	})

	// 测试成功场景
	c.KV["accountNum"] = "123456"
	c.KV["simulateQueryFail"] = false

	err := c.Apply()
	if err != nil {
		t.Logf("执行失败: %v", err)
	} else {
		t.Logf("链执行成功")
	}

	// 测试失败场景
	c.KV["accountNum"] = "123456"
	c.KV["simulateQueryFail"] = true
	err = c.Apply()
	if err != nil {
		t.Logf("执行失败: %v", err)
	} else {
		t.Logf("执行成功，耗时: %v", c.GetExecutionTime())
	}
}

// 测试性能（更新后的版本）
func TestPerformance(t *testing.T) {
	log.Println("=== 测试性能 ===")

	c := chain.LoadHandlers(nil, nil, F1(), F2(), F3(), F4(), F5(), F6(), F7())
	c.KV["simulateF3Fail"] = false
	c.KV["simulateF4Fail"] = false
	c.KV["simulateF6Fail"] = false

	al := routinePool.GetAllLeaders()

	var t1, t2 time.Time
	go func() {
		t1 = time.Now()
		for i := 0; i < 100000; i++ {
			_c := *c
			_c.Reset()
			task := chain.CreateTask(&_c)
			err := al.TaskToLeader(task)
			if err != nil {
				log.Println(err)
			}
		}
		t2 = time.Now()
	}()

	// 监控一段时间
	for j := 0; j < 3; j++ {

		for i := 0; i < 5; i++ {
			log.Printf("leader %d running num: %d", i%routinePool.MaxLeaderNum, al.GetLeader(i%routinePool.MaxLeaderNum).GetLeaderRunning())
		}
		time.Sleep(2 * time.Second)
	}

	log.Printf("性能测试完成，耗时: %v", t2.Sub(t1))
}

// 测试错误恢复
func TestErrorRecovery(t *testing.T) {
	log.Println("=== 测试错误恢复 ===")

	c := chain.LoadHandlers(nil, F1(), F2(), F3(), F4(), F5(), F6(), F7())

	// 设置错误处理，但不中断执行
	c.SetErrorHandler(func(err error, ctx *chain.MyContext) {
		log.Printf("错误被捕获: %v", err)
		// 可以选择继续执行或中断
		c.Next()
	})

	c.KV["simulateF3Fail"] = true
	c.KV["simulateF4Fail"] = false
	c.KV["simulateF6Fail"] = false

	err := c.Apply()
	if err != nil {
		t.Logf("链执行失败: %v", err)
	} else {
		t.Logf("链执行成功")
	}
}
