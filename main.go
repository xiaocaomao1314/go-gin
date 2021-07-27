package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                    "message": "随到随装粉色的范德萨发",
            })
    })
    r.Run(":9000") 
}
