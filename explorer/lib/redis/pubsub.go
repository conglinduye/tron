package redis

import (
	"errors"
	"log"

	"sync"

	"fmt"

	"time"

	src "gopkg.in/redis.v4"
)

/* PubSubServer
 * 基于redis的pub sub指令进行消息的通知
 *
 * 作为sub方：
 * 使用时，需要创建一个实现了MessageListenerInterface的listener
 * 并将该listener注册到PubSubServer中
 *
 *　作为pub方：
 *　使用时，使用PubSubServer广播一个消息
 */

//PubSubServer redis订阅消息的
type PubSubServer struct {
	client   *TronRedis                            //连接客户端
	listener map[string][]MessageListenerInterface //监听器
	locker   sync.RWMutex                          //监听器锁对象
	pubsub   *src.PubSub                           //消息订阅器
}

//NewPubSubServer 创建一个新的PubSubServer
func NewPubSubServer(addr, password string, db int) *PubSubServer {
	var server PubSubServer
	var client = NewClient(addr, password, db, 10)
	if nil != client {
		server.client = client
	} else {
		return nil
	}
	return &server
}

//IsValid 是否可用，如果可用返回nil, 否则返回错误值
func (s *PubSubServer) IsValid() error {
	if s.client == nil {
		return errors.New("PubSubServer not initialize or can not connect to server")
	}
	return nil
}

//GetAddress get服务器地址
func (s *PubSubServer) GetAddress() string {
	if nil != s.client {
		return s.client.Addr
	}
	return ""
}

//SubScriber 订阅消息，如果成功返回nil, 否则返回错误值
func (s *PubSubServer) SubScriber(channel string, listener MessageListenerInterface) error {
	var err error
	//check conn
	if err = s.IsValid(); err != nil {
		return err
	}

	//check parameter
	if len(channel) == 0 || listener == nil {
		return errors.New("channel or listener is nil")
	}

	//如果pubsub为空，需要构造这个对象
	if s.pubsub == nil {
		s.pubsub, err = s.client.Subscribe(channel)
		if err == nil || s.pubsub != nil {
			//启动监听
			go s.processQueryMessage()
		} else {
			return err
		}
	}

	//向redis订阅消息
	err = s.pubsub.Subscribe(channel)
	if err == nil {
		err = s.registeListener(channel, listener)
	}
	return err
}

//Unsubscriber 取消监听
func (s *PubSubServer) Unsubscriber(channel string, listener MessageListenerInterface) error {
	if s.pubsub == nil {
		return nil
	}

	var err error
	//check conn
	if err = s.IsValid(); err != nil {
		return err
	}

	//check parameter
	if len(channel) == 0 {
		return errors.New("channel is nil")
	}
	//如果监听器为空，则取消全部的监听
	if listener == nil {
		if err = s.pubsub.Unsubscribe(channel); err != nil {
			return err
		}
	}
	s.unregisteListener(channel, listener)
	return err
}

//Exit 退出
func (s *PubSubServer) Exit() error {
	//check conn
	if err := s.IsValid(); err != nil {
		return err
	}

	return s.client.Close()
}

//PublishMessage 发布消息
func (s *PubSubServer) PublishMessage(channel, message string) error {
	//check conn
	if err := s.IsValid(); err != nil {
		return err
	}
	if len(channel) == 0 || len(message) == 0 {
		return errors.New("channel or message invalid")
	}

	if cmd := s.client.Publish(channel, message); cmd != nil {
		return cmd.Err()
	}
	return nil
}

//registeListener 注册监听器
func (s *PubSubServer) registeListener(channel string, listener MessageListenerInterface) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	if nil == s.listener {
		s.listener = make(map[string][]MessageListenerInterface, 0)
	}

	v, ok := s.listener[channel]
	if ok == false || nil == v {
		v = make([]MessageListenerInterface, 0)
	}
	v = append(v, listener)
	s.listener[channel] = v

	return nil
}

//unregisteListener 注销监听器
func (s *PubSubServer) unregisteListener(channel string, listener MessageListenerInterface) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if nil == s.listener {
		return
	}

	v, ok := s.listener[channel]
	if ok == false || nil == v {
		return
	}

	for i, item := range v {
		if item == listener {
			v = append(v[:i], v[i+1:]...)
			s.listener[channel] = v
			break
		}
	}

	return
}

//processQueryMessage 处理监听消息
func (s *PubSubServer) processQueryMessage() {
	fmt.Println("the PubSubServer, processQueryMessage is starting ... ...")
	for {
		if nil == s.pubsub {
			fmt.Println("the redis pubsub object is nil, please check the code or configuration.")
			time.Sleep(1 * time.Second)
			continue
		}

		//尝试请求获取消息
		msg, err := s.pubsub.ReceiveMessage()
		if err != nil {
			log.Println("ReceiveMessage error :", err)
			continue
		}

		//如果没有收到消息，则需要休息一会儿，防止过快的访问redis
		if nil == msg {
			if err := s.pubsub.Ping("just a ping message."); err != nil {
				fmt.Printf("the redis ping error :[%v]\n", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		s.deliveryMessage(msg)
	}
}

//deliveryMessage 分发消息
func (s *PubSubServer) deliveryMessage(msg *src.Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("deliveryMessage recover a error : ", err)
		}
	}()

	s.locker.RLock()
	defer s.locker.RUnlock()

	listeners, ok := s.listener[msg.Channel]
	if ok == false {
		return
	}

	for _, item := range listeners {
		if item == nil {
			continue
		}
		item.ReceiveMesage(msg)
	}
}
