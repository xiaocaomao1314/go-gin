package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                    "message": "sfdsfd范德萨范德萨发西方第三方似懂非懂但是范德萨发的反对是多少的",
            })
    })
    r.Run(":9000") 
}
