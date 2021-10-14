package main

import (
	"fmt"
	"net/http"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func main() {
	// 配置好路由
	http.HandleFunc("/", serveWeachat)
	err := http.ListenAndServe(":8090", nil)
	fmt.Println("wechat server listener at", ":8090")
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}

// 路由调用的函数
func serveWeachat(rw http.ResponseWriter, req *http.Request) {
 fmt.Println("请求一次我做了什么")
	// fmt.Fprintln(rw,"我怎么看")
	wc := wechat.NewWechat()
	// 本地内存保存access_token
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "wx156ac2299ddaba11",
		AppSecret: "35693de7cc39c29f162b1c161c2db775",
		Token:     "shuidi",
		Cache:     memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	server := officialAccount.GetServer(req, rw)
	//  设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	err := server.Serve()
	if err != nil {
		
		fmt.Println("错误是:",err)
		return
	}
	server.Send()
}
// gbk转为utf-8
