package Models

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// WsConn TODO:封装的基本结构体
type WsConn struct {
	inChan     chan WebSocketData
	outChan    chan WebSocketData
	closeChan  chan []byte
	isClose    bool // 通道closeChan是否已经关闭
	closeMutex sync.RWMutex
	mutex      sync.Mutex
	Conn       *websocket.Conn
	AccountNum string
	once       sync.Once
}

func NewWsConn(conn *websocket.Conn, accountNum string) *WsConn {
	ws := &WsConn{
		inChan:     make(chan WebSocketData, 10),
		outChan:    make(chan WebSocketData, 10),
		closeChan:  make(chan []byte, 1024),
		Conn:       conn,
		AccountNum: accountNum,
	}
	return ws
}

// InChanRead TODO:读取inChan的数据
func (conn *WsConn) InChanRead() (data WebSocketData, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// InChanWrite TODO:inChan写入数据
func (conn *WsConn) InChanWrite(data WebSocketData) (err error) {
	select {
	case conn.inChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanRead TODO:读取inChan的数据
func (conn *WsConn) OutChanRead() (data WebSocketData, err error) {
	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanWrite TODO:inChan写入数据
func (conn *WsConn) OutChanWrite(data WebSocketData) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// CloseConn TODO:关闭WebSocket连接
func (conn *WsConn) CloseConn() {
	// 关闭closeChan以控制inChan/outChan策略,仅此一次
	conn.once.Do(func() {
		fmt.Println("---------------close conn accountNum:", conn.AccountNum)
		conn.closeMutex.Lock()
		if !conn.isClose {
			close(conn.closeChan)
			conn.isClose = true
		}
		conn.closeMutex.Unlock()
		//关闭WebSocket的连接,Conn.Close()是并发安全可以多次关闭
		_ = conn.Conn.Close()

	})

}
func (conn *WsConn) IsClose() bool {
	conn.closeMutex.RLock()
	defer conn.closeMutex.RUnlock()
	return conn.isClose
}
