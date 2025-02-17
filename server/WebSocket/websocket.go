package WebSocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/Jwt"
	"sync"
	"time"
)

var webSocketMap = make(map[string]*wsConn)
var wsLock = sync.RWMutex{}

func ReadWebSocketMap(key string) (exist bool, ws *wsConn) {
	wsLock.RLock()
	ws, exist = webSocketMap[key]
	wsLock.RUnlock()
	return
}
func writeWebSocketMap(accountNum string, ws *wsConn) {
	wsLock.Lock()
	webSocketMap[accountNum] = ws
	wsLock.Unlock()
}
func deleteWebSocketMap(accountNum string) {
	wsLock.Lock()
	delete(webSocketMap, accountNum)
	wsLock.Unlock()
}

// wsConn TODO:封装的基本结构体
type wsConn struct {
	inChan     chan WebSocketData
	outChan    chan WebSocketData
	closeChan  chan []byte
	isClose    bool // 通道closeChan是否已经关闭
	mutex      sync.Mutex
	conn       *websocket.Conn
	accountNum string
}

// InitWebSocket TODO:初始化Websocket
func InitWebSocket(conn *websocket.Conn, accountNum string) (ws *wsConn, err error) {
	ws = &wsConn{
		inChan:     make(chan WebSocketData, 10),
		outChan:    make(chan WebSocketData, 10),
		closeChan:  make(chan []byte, 1024),
		conn:       conn,
		accountNum: accountNum,
	}
	writeWebSocketMap(accountNum, ws)
	// 读取客户端数据协程/发送数据协程
	go ws.readMsgLoop()
	go ws.writeMsgLoop()
	data := "Websocket Init success"
	var InitData = WebSocketData{
		Data:       data,
		AccountNum: accountNum,
		Receiver:   accountNum,
		Type:       "Init",
		CreatedAt:  time.Now().String(),
	}
	fmt.Println(InitData)
	err = ws.InChanWrite(InitData)
	return
}

// InChanRead TODO:读取inChan的数据
func (conn *wsConn) InChanRead() (data WebSocketData, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// InChanWrite TODO:inChan写入数据
func (conn *wsConn) InChanWrite(data WebSocketData) (err error) {
	select {
	case conn.inChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanRead TODO:读取inChan的数据
func (conn *wsConn) OutChanRead() (data WebSocketData, err error) {
	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanWrite TODO:inChan写入数据
func (conn *wsConn) OutChanWrite(data WebSocketData) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// CloseConn TODO:关闭WebSocket连接
func (conn *wsConn) CloseConn() {
	// 关闭closeChan以控制inChan/outChan策略,仅此一次
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
	//关闭WebSocket的连接,conn.Close()是并发安全可以多次关闭
	_ = conn.conn.Close()
	deleteWebSocketMap(conn.accountNum)
}

// readMsgLoop TODO:读取客户端发送的数据写入到inChan
func (conn *wsConn) readMsgLoop() {
	for {
		// 确定数据结构
		var (
			data []byte
			err  error
		)
		// 接受数据
		if _, data, err = conn.conn.ReadMessage(); err != nil {
			goto ERR
		}
		var wsData = WebSocketData{}
		err = json.Unmarshal(data, &wsData)
		if err != nil {
			fmt.Println(err, "data", string(data))
			continue
		}
		//TODO 布隆过滤器检测接收方用户是否存在

		wsData.CreatedAt = time.Now().Format(time.RFC3339)
		//发给消息队列进行信息处理,如果错误则熔断并返回错误消息给用户，不进行后续操作
		err = SendMessage(wsData)
		if err != nil {
			errData := FormErrWebSocketData(err, wsData.AccountNum, wsData.Receiver)
			if err = conn.InChanWrite(errData); err != nil {
				goto ERR
			}
			continue
		}
		//通过websocket直接传递给客户端
		exist, _conn := ReadWebSocketMap(wsData.Receiver)
		if !exist {
			continue
		}

		// 写入数据
		if err = _conn.InChanWrite(wsData); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

// writeMsgLoop TODO:读取outChan的数据响应给客户端
func (conn *wsConn) writeMsgLoop() {
	for {
		var (
			data []byte
			err  error
		)
		// 读取数据
		if wsData, err := conn.OutChanRead(); err != nil {
			fmt.Println(err)
			goto ERR
		} else {
			if wsData.Type == "Init" {
				continue
			}
			if wsData.Type == "error" {
				data, _ = json.Marshal(map[string]interface{}{
					"code":        500,
					"content":     wsData.Data,
					"handlerName": "error",
				})
			} else {
				data = []byte(wsData.Data)
			}
			fmt.Println("write message", wsData.Data)
		}
		// 发送数据

		if err = conn.conn.WriteMessage(1, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

// websocket 升级并跨域
var (
	upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// WebSocketMiddleware
// websocket中间件，每次请求时判断是否建立了websocket连接,如果没有则创建连接并返回,如果已经建立连接则直接返回
func WebSocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			conn *websocket.Conn
			ws   *wsConn
			mc   *Jwt.CustomClaims
		)
		fmt.Println(" start websocket middleware")
		token := c.Query("token")
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err = Jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		sender := mc.AccountNum
		fmt.Println("sender:", sender, "password:", mc.Password)
		if sender == "" {
			c.Abort()
			fmt.Println("account_num is empty")
			return
		}
		// TODO 布隆过滤器检测账号存不存在，不存在则熔断并返回错误消息给用户，不进行后续操作

		// 已经建立连接则直接返回
		if exist, _ := ReadWebSocketMap(sender); exist {
			c.Abort()
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "account_num is exist",
			})
			return
		}

		if conn, err = upgrade.Upgrade(c.Writer, c.Request, nil); err != nil {
			fmt.Println(err)
			return
		}
		if ws, err = InitWebSocket(conn, sender); err != nil {
			fmt.Println(err)
			return
		}
		// 使得inChan和outChan耦合起来
		// 管道关闭则关闭WebSocket连接
		go func() {
			for {
				var data WebSocketData
				if data, err = ws.InChanRead(); err != nil {
					fmt.Println(err)
					goto ERR
				}
				if err = ws.OutChanWrite(data); err != nil {
					fmt.Println(err)
					goto ERR
				}
			}
		ERR:
			ws.CloseConn()
		}()
		return
	}
}
