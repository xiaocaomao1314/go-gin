package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                    "message": "sfdsfd范德萨范德萨发西方的反对是多少的",
            })
    })
    r.Run(":9000") 
}
