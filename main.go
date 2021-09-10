package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"encoding/json"
	"io"
	"io/ioutil"
	mathRand "math/rand"
	"strconv"
	"sync"
	"time"
)
type WxAccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}
type WxJsApiTicket struct {
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
type WxSignature struct {
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
	AppID     string `json:"appId"`
}
 
type WxSignRtn struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data WxSignature `json:"data"`
}
var (
	MemoryCacheVar  *MemoryCache
	AppID           string = "wx0d45a180607ace86"
	AppSecret       string = "8485007d3272349315455cbbeadad445"
	AccessTokenHost string = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppID + "&secret=" + AppSecret
	JsAPITicketHost string = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)
func main()  {
	// 绑定路由
	// http.HandleFunc("/",checkout）
	http.HandleFunc("/wx", getWxSign)
	// 启动监听=j
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
	 fmt.Println("服务器启动失败！")
	}
}
func checkout(response http.ResponseWriter, request *http.Request)  {
	 

    response.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域

    response.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

    response.Header().Set("content-type", "application/json")   
	
  
	//解析URL参数
	err := request.ParseForm()
	if err != nil {
		fmt.Println("URL解析失败！")
		return
	}
	// token
	var token string = "kkcm"
	// 获取参数
	signature := request.FormValue("signature")
	timestamp := request.FormValue("timestamp")
	nonce := request.FormValue("nonce")
	echostr := request.FormValue("echostr")
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray  = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		_, err := response.Write([]byte(echostr))
		if err != nil {
			fmt.Println("响应失败。。。")
		}
	} else {
		fmt.Println("验证失败")
	}
}
func getWxSign(w http.ResponseWriter, r *http.Request) {
	 

    w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域

    w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

    w.Header().Set("content-type", "application/json")   
  
	fmt.Println(3242)

	query := r.URL.Query()
var	url string= query.Get("url")
	var (
		noncestr, jsapi_ticket, timestamp, signature, signatureStr, access_token string
		wxAccessToken WxAccessToken
		wxJsApiTicket WxJsApiTicket
		wxSignature WxSignature
		wxSignRtn WxSignRtn
	)
 
	
	
	noncestr = RandStringBytes(16)
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	//获取access_token，如果缓存中有，则直接取出数据使用；否则重新调用微信端接口获取
	client := &http.Client{}
	if MemoryCacheVar.Get("access_token") == nil {
		request, _ := http.NewRequest("GET", AccessTokenHost, nil)
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		err = json.Unmarshal(body, &wxAccessToken)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		if wxAccessToken.Errcode == 0 {
			access_token = wxAccessToken.Access_token
		} else {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = wxAccessToken.Errmsg
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		MemoryCacheVar.Put("access_token", access_token, time.Duration(wxAccessToken.Expires_in)*time.Second)
 
		//获取 jsapi_ticket
		requestJs, _ := http.NewRequest("GET", JsAPITicketHost+"?access_token="+access_token+"&type=jsapi", nil)
		responseJs, _ := client.Do(requestJs)
		defer responseJs.Body.Close()
		bodyJs, err := ioutil.ReadAll(responseJs.Body)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		err = json.Unmarshal(bodyJs, &wxJsApiTicket)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		if wxJsApiTicket.Errcode == 0 {
			jsapi_ticket = wxJsApiTicket.Ticket
		} else {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = wxJsApiTicket.Errmsg
			fmt.Fprintln(w,wxSignRtn)
			return
		}
		MemoryCacheVar.Put("jsapi_ticket", jsapi_ticket, time.Duration(wxJsApiTicket.Expires_in)*time.Second)
	} else {
		//缓存中存在access_token，直接读取
		access_token = MemoryCacheVar.Get("access_token").(*Item).Value
		jsapi_ticket = MemoryCacheVar.Get("jsapi_ticket").(*Item).Value
	}
    fmt.Println("access_token:",access_token)
	fmt.Println("jsapi_ticket:",jsapi_ticket)
 
	// 获取 signature
	signatureStr = "jsapi_ticket=" + jsapi_ticket + "&noncestr=" + noncestr + "&timestamp=" + timestamp + "&url=" + url
	signature = GetSha1(signatureStr)
 
	wxSignature.Url = url
	wxSignature.Noncestr = noncestr
	wxSignature.Timestamp = timestamp
	wxSignature.Signature = signature
	wxSignature.AppID = AppID
 
	// 返回前端需要的数据
	wxSignRtn.Code = 0
	wxSignRtn.Msg = "success"
	wxSignRtn.Data = wxSignature
	fmt.Fprintln(w,wxSignRtn)
}
 
//生成指定长度的字符串
func RandStringBytes(n int) string {
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathRand.Intn(len(letterBytes))]
	}
	return string(b)
}
 
//SHA1加密
func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
 

 
// 数据缓存处理
type Item struct {
	Value      string
	CreateTime time.Time
	LifeTime   time.Duration
}
 
type MemoryCache struct {
	sync.RWMutex
	Items map[string]*Item
}
 
func (mc *MemoryCache) Put(key string, value string, lifeTime time.Duration) {
	mc.Lock()
	defer mc.Unlock()
	mc.Items[key] = &Item{
		LifeTime:   lifeTime,
		Value:      value,
		CreateTime: time.Now(),
	}
}
 
func (mc *MemoryCache) Get(key string) interface{} {
	mc.RLock()
	defer mc.RUnlock()
	if e, ok := mc.Items[key]; ok {
		if !e.isExpire() {
			return e
		} else {
			delete(mc.Items, key)
		}
	}
	return nil
}
 
func (e *Item) isExpire() bool {
	if e.LifeTime == 0 {
		return false
	}
	//根据创建时间和生命周期判断元素是否失效
	return time.Now().Sub(e.CreateTime) > e.LifeTime
}
