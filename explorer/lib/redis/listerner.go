package redis

import (
	"log"

	src "gopkg.in/redis.v4"
)

//RedisSubChannelUserInfo redis tunnel 名称定义
const (
	RedisSubChannelUserInfo = "RedisSubChannelUserInfo" //用户消息变更
)

//======================================================================

/* MessageListenerInterface
 * 是监听消息时，需要实现的接口
 *
 * MessageListener为默认的消息监听器，输出到标准输出
 */

//MessageListenerInterface 监听接收消息的接口
type MessageListenerInterface interface {
	ReceiveMesage(m *src.Message)
}

//MessageListener 消息监听器
type MessageListener struct {
}

//ReceiveMesage 处理接收消息
func (l *MessageListener) ReceiveMesage(m *src.Message) {
	if nil == m {
		log.Println("receive a nil message.")
	} else {
		log.Printf("receive a message:[%v] from[%v]\n", m.Payload, m.Channel)
	}
}

//================================================================================

//SubMessageKey 消息的唯一标识
type SubMessageKey struct {
	Channel string `json:"tunnel"`
	MsgName string `json:"messageName"`
}

//CheckIsValid 检查消息是否有效
func (s *SubMessageKey) CheckIsValid() bool {
	return len(s.Channel) > 0 && len(s.MsgName) > 0
}

//IsMatch 是否匹配
func (s *SubMessageKey) IsMatch(msgName string) bool {
	return s.MsgName == msgName
}

//IsMatchEx 是否匹配
func (s *SubMessageKey) IsMatchEx(tunnel, msgName string) bool {
	return s.MsgName == msgName && s.Channel == tunnel
}

//SubSimpleMessage 简单的消息
type SubSimpleMessage struct {
	MsgKey  SubMessageKey     `json:"msgKey"`
	MsgData map[string]string `json:"msgData"`
}

//CheckIsValid 检查消息是否有效
func (s *SubSimpleMessage) CheckIsValid() bool {
	return s.MsgKey.CheckIsValid()
}

//SubCommonMessage 通用消息
type SubCommonMessage struct {
	MsgKey  SubMessageKey `json:"msgKey"`
	MsgData []byte        `json:"msgData"`
}

//CheckIsValid 检查消息是否有效
func (s *SubCommonMessage) CheckIsValid() bool {
	return s.MsgKey.CheckIsValid()
}
