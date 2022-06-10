package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", rootHandler)
	setupRoutes(r)
	r.Run(":8989")
}
