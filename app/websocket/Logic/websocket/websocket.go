package websocket

import (
	"common/chain"
	redis2 "common/redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"rpc/message"
	"rpc/user"
	"sync"
	"time"
	"websocket/DAO/Redis"
	"websocket/Logic"
	"websocket/Models"
	"websocket/rpc/client"
)

var (
	// Upgrade websocket 升级并跨域
	Upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	webSocketMap = make(map[string]*Models.WsConn)
	wsLock       = sync.RWMutex{}

	BroadcastConsumerQueue amqp.Queue
	BroadcastMsg           <-chan amqp.Delivery
	DirectConsumerQueue    amqp.Queue
	DirectMsg              <-chan amqp.Delivery
)

func InitWebsocketMQ() {
	var err error
	Logic.ExchangeInit()
	BroadcastConsumerQueue, err = Logic.QueueDeclare(
		fmt.Sprintf("BroadcastConsumerQueue"+Logic.MachineURL), false, false, false, false, nil,
	)
	if err != nil {
		panic(err)
	}
	err = Logic.QueueBind(BroadcastConsumerQueue.Name, Logic.BExchange, "", false, nil)
	if err != nil {
		panic(err)
	}
	BroadcastMsg = Logic.StartConsume(BroadcastConsumerQueue.Name, true)

	DirectConsumerQueue, err = Logic.QueueDeclare(
		fmt.Sprintf("DirectConsumerQueue"+Logic.MachineURL), false, false, false, false, nil,
	)
	if err != nil {
		panic(err)
	}
	err = Logic.QueueBind(DirectConsumerQueue.Name, Logic.DExchange, Logic.MachineURL, false, nil)
	if err != nil {
		panic(err)
	}

	DirectMsg = Logic.StartConsume(DirectConsumerQueue.Name, true)
	go ConsumeDirectMsg()
	go ConsumeBroadcastMsg()
}

func ConsumeDirectMsg() {
	var (
		err    error
		wsData Models.WebSocketData
	)

	for d := range DirectMsg {
		fmt.Println("Read DirectMsg:", string(d.Body))
		err = json.Unmarshal(d.Body, &wsData)
		if err != nil {
			log.Println(err)
			continue
		}
		exist, _conn := ReadWebSocketMap(wsData.Receiver)
		if !exist {
			log.Println("receiver not exist")
			continue
		}
		// 写入数据
		if err = _conn.InChanWrite(wsData); err != nil {
			_conn.CloseConn()
			continue
		}
	}
}
func ConsumeBroadcastMsg() {
	var (
		err    error
		wsData Models.WebSocketData
	)

	for d := range BroadcastMsg {
		err = json.Unmarshal(d.Body, &wsData)
		if err != nil {
			log.Println(err)
			continue
		}
		exist, _conn := ReadWebSocketMap(wsData.Receiver)
		if !exist {
			log.Println("receiver not exist")
			continue
		}
		// 写入数据
		if err = _conn.InChanWrite(wsData); err != nil {
			_conn.CloseConn()
			continue
		}
	}
}

func KickUser(reason string, accountNum string) error {
	wsData := Models.WebSocketData{
		Data:       reason,
		AccountNum: accountNum,
		Receiver:   accountNum,
		Type:       "kick",
		CreatedAt:  time.Now().String(),
	}
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), checkReceiverExist(), noticedByWebsocket())
	ctx.Set("wsData", wsData)
	return ctx.Apply()
}

// readMsgLoop TODO:读取客户端发送的数据写入到inChan
func readMsgLoop(wsConn *Models.WsConn) {
	// 确定数据结构
	var (
		data []byte
		err  error
	)
	for {
		// 接受数据
		if _, data, err = wsConn.Conn.ReadMessage(); err != nil {
			chain.ZapLogger("error", "ReadMessage error:%s", err.Error())
			goto ERR
		}
		ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), renewalNode(), checkReceiverExist(), checkIsFriend(), persistenceMsg(), noticedByWebsocket())
		ctx.SetWithMap(
			map[string]interface{}{
				"wsConn": wsConn,
				"data":   data,
			})
		ctx.Apply()

	}
ERR:
	wsConn.CloseConn()
}

// writeMsgLoop TODO:读取outChan的数据响应给客户端
func writeMsgLoop(wsConn *Models.WsConn) {
	for {
		var (
			data []byte
			err  error
		)

		// 读取数据
		wsData, err := wsConn.OutChanRead()
		if err != nil {
			fmt.Println(err)
			goto ERR
		}
		if wsData.Type == "init" {
			continue
		}
		if wsData.Type == "error" {
			data, _ = json.Marshal(map[string]interface{}{
				"code": 500,
				"data": wsData,
				"type": "error",
			})
		} else {
			data, _ = json.Marshal(wsData)
		}

		// 发送数据
		if err = wsConn.Conn.WriteMessage(1, data); err != nil {
			goto ERR
		}
	}
ERR:
	wsConn.CloseConn()
}

// InitWebSocket TODO:初始化Websocket
func InitWebSocket(conn *websocket.Conn, accountNum string) (ws *Models.WsConn, err error) {
	ws = Models.NewWsConn(conn, accountNum)
	// 读取客户端数据协程/发送数据协程
	go readMsgLoop(ws)
	go writeMsgLoop(ws)

	//这里更新redis，设置定期时间5分过期
	c := Redis.RedisPoolGet()
	redis2.HmsetWithExpireScript.Do(c, fmt.Sprintf("websocket_%s", accountNum), 360, "url", Logic.MachineURL)
	Redis.RedisPoolPut(c)
	//Redis.RedisDo(Redis2.HmsetWithExpireScriptString, fmt.Sprintf("websocket_%s", AccountNum), 360, "url", Logic.MachineURL)
	//Redis.RedisDo(Redis.HSET, fmt.Sprintf("websocket_%s", AccountNum), Logic.MachineURL)
	//Redis.RedisDo(Redis2.EXPIRE, fmt.Sprintf("websocket_%s", AccountNum), 300)

	data := "websocket Init success"
	var InitData = Models.WebSocketData{
		Data:       data,
		AccountNum: accountNum,
		Receiver:   accountNum,
		Type:       "init",
		CreatedAt:  time.Now().String(),
	}
	err = ws.InChanWrite(InitData)
	if err != nil {
		ws.CloseConn()
	}
	return
}

func ReadWebSocketMap(key string) (exist bool, ws *Models.WsConn) {
	wsLock.RLock()
	ws, exist = webSocketMap[key]
	wsLock.RUnlock()
	if exist && ws.IsClose() {
		//第一次关闭的时候是软删除，这里再次判断是否关闭，如果已经关闭，则删除
		deleteWebSocketMap(key)
		exist = false
	}
	return
}

func WriteWebSocketMap(accountNum string, ws *Models.WsConn) {
	wsLock.Lock()
	fmt.Println("writeWebSocketMap", accountNum, ws)
	webSocketMap[accountNum] = ws
	wsLock.Unlock()
}

func deleteWebSocketMap(accountNum string) {
	wsLock.Lock()
	delete(webSocketMap, accountNum)
	wsLock.Unlock()
}

// renewalNode 节点续期
func renewalNode() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {

		wsConn, err := chain.GetToType[*Models.WsConn](ctx, "wsConn")
		if err != nil {
			return err
		}
		data, err := chain.GetToType[[]byte](ctx, "data")
		if err != nil {
			return err
		}

		// 活跃节点续期
		c := Redis.RedisPoolGet()
		_, err = redis2.HmsetWithExpireScript.Do(c, fmt.Sprintf("websocket_%s", wsConn.AccountNum), 360, "url", Logic.MachineURL)
		Redis.RedisPoolPut(c)
		//_, err := Redis.RedisDo(Redis2.HmsetWithExpireScriptString, fmt.Sprintf("websocket_%s", wsConn.AccountNum), 360, "url", Logic.MachineURL)
		if err != nil {
			log.Println(err)
		}
		//Redis.RedisDo(Redis2.EXPIRE, fmt.Sprintf("websocket_%s", wsConn.AccountNum), 360)
		var wsData = Models.WebSocketData{}
		err = json.Unmarshal(data, &wsData)
		if err != nil {
			return err
		}
		wsData.CreatedAt = time.Now().Format(time.RFC3339)
		ctx.Set("wsData", wsData)
		return nil
	}
}
func checkReceiverExist() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		wsData, err := chain.GetToType[Models.WebSocketData](ctx, "wsData")
		if err != nil {
			return err
		}

		// 检验接收者是否存在
		_user, err := client.UserServiceClient.Client.FindUser(
			context.Background(),
			&user.FindUserRequest{AccountNum: wsData.Receiver})
		if err != nil {
			return err
		}
		if _user == nil {
			return errors.New("用户不存在")
		}
		return nil
	}
}

func checkIsFriend() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		//// 检查是否为好友关系
		//resp, err := client.FriendServiceClient.Client.CheckFriend(
		//	context.Background(),
		//	&friend.CheckFriendRequest{
		//		Sender:      wsData.AccountNum,
		//		Receiver:    wsData.Receiver,
		//		HandlerName: "CheckFriend",
		//	})
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		//if !resp.IsFriend {
		//	fmt.Println("非好友无法发送消息")
		//	continue
		//}
		return nil
	}
}

func persistenceMsg() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		wsData, err := chain.GetToType[Models.WebSocketData](ctx, "wsData")
		if err != nil {
			return err
		}
		// 调用微服务发送消息进行持久化
		//TODO 改成向消息队列发送消息，消息队列处理这个持久化
		resp, err := client.MessageServiceClient.Client.SendText(context.Background(), &message.SendTextRequest{
			SenderAccountNum:   wsData.AccountNum,
			ReceiverAccountNum: wsData.Receiver,
			HandlerName:        wsData.Type,
			Content:            wsData.Data,
			Time:               []byte(wsData.CreatedAt),
		})
		if err != nil {
			return err
		}
		ctx.Logger("info", "SendText response %v", resp)
		return nil
	}
}

// noticedByWebsocket 通过websocket通知用户
func noticedByWebsocket() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		wsData, err := chain.GetToType[Models.WebSocketData](ctx, "wsData")
		if err != nil {
			return err
		}

		//通过websocket发送消息
		//TODO检查是否在线,在线才通过websocket发送消息

		//这里先查本地
		exist, _conn := ReadWebSocketMap(wsData.Receiver)
		fmt.Println("get wsconn from local", exist, wsData.Receiver)
		if exist {
			// 写入数据
			if err = _conn.InChanWrite(wsData); err != nil {
				// 写入失败，关闭连接
				_conn.CloseConn()
				return err
			}
			return nil
		}
		data, _ := json.Marshal(wsData)
		// 本地没有则查redis然后通过消息队列传给其他服务器
		reply, err := Redis.RedisDo(redis2.HMGET, fmt.Sprintf("websocket_%s", wsData.Receiver), "url")
		url, err := redis.Strings(reply, err)
		if errors.Is(err, redis.ErrNil) || len(url) == 0 {
			goto Broadcast
		}
		if err != nil {
			return err
		}
		err = Logic.DirectionMessage(data, url[0], false, false)
		if err != nil {
			return err
		}
		return nil
		// redis没有则通过消息队列广播
	Broadcast:
		err = Logic.BroadcastMessage(data, false, false)
		return err
	}
}
