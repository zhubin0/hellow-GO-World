package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhubin/zbService/controller"
)

func main(){
	router := gin.New()

	router.GET("/", HelloGin)
	router.GET("/zzz", HHH)
	router.POST("/loginJSON", controller.AA)

	// 绑定普通表单的例子 (user=manu&password=123)
	router.POST("/loginForm", controller.BB)

	// 绑定多媒体表单的例子 (user=manu&password=123)
	router.POST("/login", controller.CC)



	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}

func HelloGin(c *gin.Context){
	c.String(http.StatusOK, "Hello World")
}

func HHH(c *gin.Context){
	c.String(http.StatusOK, "Hello HHH")
}
