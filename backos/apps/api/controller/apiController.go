package controller

import "github.com/gin-gonic/gin"

func CheckStat (c *gin.Context) {
	c.String(200, "ok")
}


func Ping (c *gin.Context) {
	c.String(200, "pong")
}