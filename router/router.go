/*
 * @Description:
 * @Date: 2021-10-15 10:16:22
 * @LastEditors: caomao
 * @LastEditTime: 2021-10-18 14:16:37
 */
package router

import (
	"go-gin/wx"
	"net/http"
)

func WxRouter() {
	http.HandleFunc("/", wx.ReceiveMessage)
}
