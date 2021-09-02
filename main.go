package main

import (
"github.com/gin-gonic/gin"
"net/http"
)
func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                    "message": "s878",
            })
    })
    r.GET("/list",func (c *gin.Context)  {
        c.String(http.StatusOK, "Hello World")
    })
    r.GET("/listView",func (c *gin.Context)  {
        c.String(http.StatusNotFound, "未发现页面")
    })
    r.GET("/yaml",func (c *gin.Context)  {
        c.YAML(http.StatusOK,gin.H{"name":"patch 请求 yaml 格式化","age":18})
      })
      r.GET("/xml",func (c *gin.Context)  {
        c.XML(http.StatusOK,gin.H{"name":"delete 请求 xml 格式化","age":18})
      })
      //get请求 html界面显示   http://localhost:8080/html
  r.GET("/html",func (c *gin.Context)  {
        // router.LoadHTMLGlob("../view/tem/index/*")  //这是前台的index
        // router.LoadHTMLGlob("../view/tem/admin/*")  //这是后台的index
        r.LoadHTMLFiles("./index.html")  //指定加载某些文件
        c.HTML(http.StatusOK,"index.html",nil)
      })
    r.Run(":900
0") 
}
