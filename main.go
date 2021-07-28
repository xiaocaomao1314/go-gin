package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                    "message": "sfdsfdsfdsfds大大说的",
            })
    })
    r.Run(":9000") 
}
