package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"log"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	defer func() {
		// 将消息设置为已读
		msg.AsRead()
	}()

	log.Printf("新消息流向：FromUserName:%s, ToUserName:%s \n",msg.FromUserName, msg.ToUserName)
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信客服，你可以向我提问任何问题。")
			if err != nil {
				log.Fatalf("add friend agree error : %v", err)
				return
			}
		}
	}

	// 私聊
	handlers[UserHandler].handle(msg)
}
