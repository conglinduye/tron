/**
 * @author [yanzheng]
 * @email [yan_zheng2018@163.com@mail.com]
 * @create date 2018-09-16 16:57:14
 * @modify date 2018-09-16 16:57:14
 * @desc [基于gorilla实现的websocket，读队列写队列，应答式]
 */
package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//MessageInfo xiaoxi
type MessageInfo struct {
	Address   string `json:"address"`
	ToAddress string `json:"to_address"`
	Amount    string `json:"amount"`
}

// 客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *wsMessage // 读队列
	outChan  chan *wsMessage // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

func (wsConn *wsConnection) wsReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &wsMessage{
			msgType,
			data,
		}

		// 放入请求队列
		h.broadcast <- req
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if msg != nil {
				if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
					goto error
				} else {
					fmt.Printf("send user:[%v],message:[%v]-[%v]\n", wsConn.wsSocket.RemoteAddr, msg.messageType, string(msg.data))
				}
			}

		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection) procLoop() {
	address := "TLJBkA2po1DgxJjALYnc3Vkqg333TW6s31"
	toAddress := "TBt6G8aXxDAYhnrjhoSSJS7wRBBxDbhfV8"
	amount := "199 TRX"
	// 启动一个gouroutine发送心跳
	//message := &MessageInfo{Address: "TLJBkA2po1DgxJjALYnc3Vkqg333TW6s31",
	//	ToAddress: "TBt6G8aXxDAYhnrjhoSSJS7wRBBxDbhfV8", Amount: "199 TRX"}
	message := fmt.Sprintf("%v, %v,%v", address, toAddress, amount)
	//msgByte := utils.HexDecode(message)
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if err := wsConn.wsWrite(websocket.TextMessage, []byte(message)); err != nil {
				fmt.Println("heartbeat fail")
				//当检测到某连接发送心跳失败，则服务端关闭链接，并从广播列表中删除
				wsConn.wsClose()
				h.unregister <- wsConn
				break
			}
		}
	}()
}

//WsHandler 启动websocket
func WsHandler(resp gin.ResponseWriter, req *http.Request) {
	go h.run()
	// 应答客户端告知升级连接为websocket
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	h.register <- wsConn

	// 心跳处理器
	go wsConn.procLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	h.broadcast <- &wsMessage{messageType, data}
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *wsConnection) wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
