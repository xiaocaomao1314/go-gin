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
    r.Run(":9000") 
}
