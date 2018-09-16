/**
 * @author [yanzheng]
 * @email [yan_zheng2018@163.com@mail.com]
 * @create date 2018-09-16 16:58:32
 * @modify date 2018-09-16 16:58:32
 * @desc [基于websocket实现消息广播，注册，注销链接功能]
 */
package websocket

type wsPool struct {
	// 注册了的连接器
	connections map[*wsConnection]bool

	// 从连接器中发入的信息
	broadcast chan *wsMessage

	// 从连接器中注册请求
	register chan *wsConnection

	// 从连接器中注销请求
	unregister chan *wsConnection
}

var h = wsPool{
	broadcast:   make(chan *wsMessage),
	register:    make(chan *wsConnection),
	unregister:  make(chan *wsConnection),
	connections: make(map[*wsConnection]bool),
}

func (h *wsPool) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.outChan)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.outChan <- m:
				default:
					delete(h.connections, c)
					close(c.outChan)
				}

			}
		}
	}
}
