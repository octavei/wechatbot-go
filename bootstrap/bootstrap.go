package bootstrap

import (
	"fmt"
	"github.com/869413421/wechatbot/handlers"
	"github.com/eatmoreapple/openwechat"
	"log"
)



func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}
	self, _:=bot.GetCurrentUser()
	friend, _:=self.Friends()

	fmt.Println("你的朋友列表：============", self.ID())
	for i:=0; i< friend.Count(); i++ {
		fmt.Println("序号：", (i+1), friend[i].City,
			"\nID:",friend[i].ID(),
			"\nID:",friend[i].User.ID(),
			"\nUserName:",friend[i].User.UserName,
			"\nNickName:",friend[i].User.NickName,
			//"\nMemberList:",friend[i].User.MemberList,
			//"\nSignature:",friend[i].User.Signature,
			//"\nDisplayName:",friend[i].User.DisplayName,
			"\nHeadImgUrl:",friend[i].User.HeadImgUrl)

		fmt.Println("-----------------------------------end")
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
