// package main

// import (
// 	"fmt"
// 	wechat "github.com/silenceper/wechat/v2"
// 	"github.com/silenceper/wechat/v2/cache"
// 	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
// 	"github.com/silenceper/wechat/v2/officialaccount/message"
// 	"net/http"
// )

// func main() {
// 	// 配置好路由
// 	http.HandleFunc("/", serveWeachat)
// 	fmt.Println("开始")
// 	err := http.ListenAndServe(":8090", nil)
// 	fmt.Println("wechat server listener at", ":8090")
// 	if err != nil {
// 		fmt.Printf("start server error , err=%v", err)
// 	}
// }

// // 路由调用的函数
// func serveWeachat(rw http.ResponseWriter, req *http.Request) {
// 	fmt.Println("更新我看")
// 	fmt.Fprintln(rw, "我怎么看这是什么")
// 	wc := wechat.NewWechat()
// 	// 本地内存保存access_token
// 	memory := cache.NewMemory()
// 	cfg := &offConfig.Config{
// 		AppID: "wx156ac2299ddaba11",
// 		AppSecret: "35693de7cc39c29f162b1c161c2db775",
// 		Token:     "shuidi",
// 		EncodingAESKey:"l1MxeqQfFOylpQQojS1XEB94VaUfUz3KUpjqwD0CJtn",
// 		Cache:     memory,
// 	}
// 	officialAccount := wc.GetOfficialAccount(cfg)
// 	server := officialAccount.GetServer(req, rw)
// 	//  设置接收消息的处理方法
// 	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

// 		//回复消息：演示回复用户发送的消息
// 		text := message.NewText(msg.Content)
// 		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
// 	})

// 	err := server.Serve()
// 	if err != nil {

// 		fmt.Println("错误是:", err)
// 		return
// 	}
// 	server.Send()
// }

// gbk转为utf-8
package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"
)

const (
	token = "shuidi" //跟微信公众平台的token一样即可
)

//http.ResponseWriter : 回复http的对应
//http.Request ： http请求的对象
func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("http url ", r) //打印http的请求url
	defer r.Body.Close()

	//1.尝试获取4个字段
	nonce := r.URL.Query().Get("nonce")
	timestamp := r.URL.Query().Get("timestamp")
	signature := r.URL.Query().Get("signature")
	echostr := r.URL.Query().Get("echostr")

	// if nonce != nil && timestamp != nil && signature != nil && echostr != nil {
	// 	fmt.Printf("字段提取成功")
	// } else {
	// 	return
	// }

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
	fmt.Println("url once my_signature signature", nonce, hashcode, signature)
	if hashcode != signature {
		return
	}

	// 6. 如果等于就返回echostr
	_, _ = w.Write([]byte(echostr)) //这个就是往网页端输出的值
}

func main() {
	fmt.Println("服务器程序")

	http.HandleFunc("/", sayHello)   // 匹配url的/就会调用sayHello
	http.HandleFunc("/wx", sayHello) // 匹配url的/wx就会调用sayHello

	err := http.ListenAndServe(":8090", nil) //这个就是绑定服务器的80端口
	if err != nil {
		fmt.Println("ListenAndServer  error", err)
	}

}
