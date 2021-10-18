/*
 * @Description:
 * @Date: 2021-10-15 09:54:06
 * @LastEditors: caomao
 * @LastEditTime: 2021-10-18 15:42:23
 */
package wx

import (
	"crypto/sha1"
	"os"

	"fmt"
	"net/http"
	"sort"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// 验证token 开启服务器
func VerificationToken(w http.ResponseWriter, r *http.Request) {
	// 同步微信公众号的token
	const (
		token = "shuidi"
	)
	fmt.Fprintln(w, "验证token")
	defer r.Body.Close()
	//1.尝试获取4个字段
	nonce := r.URL.Query().Get("nonce")
	timestamp := r.URL.Query().Get("timestamp")
	signature := r.URL.Query().Get("signature")
	echostr := r.URL.Query().Get("echostr")
	// 验证是否有四个字段
	if nonce != "" && timestamp != "" && signature != "" && echostr != "" {
		fmt.Printf("字段提取成功")
	} else {
		fmt.Println("没有数据")
	}

	//2. 赋值一个token

	//3.token，timestamp，nonce按字典排序的字符串list
	strs := sort.StringSlice{token, timestamp, nonce} // 使用本地的token生成校验
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}

	// 4. 哈希算法加密list得到hashcode
	h := sha1.New()
	h.Write([]byte(str))
	hashcode := fmt.Sprintf("%x", h.Sum(nil)) // h.Sum(nil) 做hash  79efadd80a344c0b73b3bd2c403184f7425a5a67

	//5. 判断hashcode是否等于signature
	if hashcode != signature {
		return
	}

	// 6. 如果等于就返回echostr
	_, _ = w.Write([]byte(echostr)) //这个就是往网页端输出的值
}

/**
 * @Description: 接收到通过公众号发送的消息
 */
func ReceiveMessage(rw http.ResponseWriter, req *http.Request) {
	path, _ := os.Executable()
	fmt.Fprintln(rw,path)
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "wx156ac2299ddaba11",
		AppSecret: "35693de7cc39c29f162b1c161c2db775",
		Token:     "shuidi",
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}
